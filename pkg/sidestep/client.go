package sidestep

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type client struct {
	host         string
	port         uint16
	domain       string
	nameServer   string
	random       *rand.Rand
	inBuffer     []byte
	bufLock      sync.Mutex
	transmission uint8
	resolver     *net.Resolver
	ctx          context.Context
	cancel       context.CancelFunc
}

type ConnectionOption func(c *client)

func Connect(host string, port uint16, options ...ConnectionOption) (Connection, error) {

	ctx, cancel := context.WithCancel(context.Background())

	client := &client{
		host:     host,
		port:     port,
		random:   rand.New(rand.NewSource(time.Now().UnixNano())),
		resolver: net.DefaultResolver,
		ctx:      ctx,
		cancel:   cancel,
	}

	for _, option := range options {
		option(client)
	}

	if client.nameServer != "" {
		client.resolver = newResolver(client.nameServer, "udp") // TODO make protocol configurable
	}

	if client.domain == "" {
		target := fmt.Sprintf("%06d.net", client.random.Int63())
		client.domain = target[len(target)-10:] //todo can this safely be shorter?
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func Domain(domain string) func(c *client) {
	return func(c *client) {
		c.domain = domain
	}
}

func NameServer(nameServer string) func(c *client) {
	return func(c *client) {
		c.nameServer = nameServer
	}
}

func (c *client) connect() error {
	packet := clientPacket{
		operation:    OpOpen,
		transmission: 0,
		sequence:     0,
		data:         []byte(fmt.Sprintf("%s:%d", c.host, c.port)),
	}
	return c.sendViaDNS(packet)
}

func (c *client) sendViaDNS(packet clientPacket) error {

	// TODO mutex lock for dns resolver?

	packet.baseSize = uint8(len(c.domain))

	// TODO: handle connecting to DNS nameserver directly
	results, err := c.resolver.LookupTXT(c.ctx, packet.ToDNS(c.domain))
	if err != nil {
		return err
	}
	switch packet.operation {
	case OpOpen:
		if len(results) != 1 || results[0] != "OK" {
			return fmt.Errorf("connection failed") // TODO: better info here - output result text?
		}
	case OpSend, OpReceive:
		c.bufLock.Lock()
		defer c.bufLock.Unlock()
		for _, result := range results {
			c.inBuffer = append(c.inBuffer, []byte(result)...)
		}
	}
	return nil
}

func (c *client) Write(data []byte) (n int, err error) {

	target := c.domain

	//------------
	var packetData []byte
	var sent int

	// 10 because 1 octet prefix, 1 octet suffix, 1 octet subdomain prefix, 7 bytes header
	remainingSpace := 0xff - (len(target) + 10)

	c.transmission++
	var sequence uint8

	for _, byt := range data {
		if remainingSpace == 0 {
			if err := c.sendViaDNS(clientPacket{
				operation:    OpSend,
				transmission: c.transmission,
				sequence:     sequence,
				data:         packetData,
			}); err != nil {
				return 0, err
			}
			sequence++
			sent += len(packetData)
			packetData = nil
		}
		packetData = append(packetData, byt)
	}

	if len(packetData) > 0 {
		if err := c.sendViaDNS(clientPacket{
			operation:    OpSend,
			transmission: c.transmission,
			sequence:     sequence,
			data:         packetData,
		}); err != nil {
			return 0, err
		}
		sent += len(packetData)
	}

	//------------

	return sent, nil
}

func (c *client) Read(data []byte) (n int, err error) {
	c.bufLock.Lock()
	if len(c.inBuffer) == 0 {
		c.bufLock.Unlock()
		c.transmission++
		if err := c.sendViaDNS(clientPacket{
			operation:    OpReceive,
			transmission: c.transmission,
		}); err != nil {
			return 0, err
		}
	}

	c.bufLock.Lock()
	defer c.bufLock.Unlock()

	for i, b := range c.inBuffer {
		if i >= len(data) {
			c.inBuffer = c.inBuffer[i:]
			return i, nil
		}
		data[i] = b
	}
	received := len(c.inBuffer)
	c.inBuffer = nil
	return received, nil
}

func (c *client) Close() error {
	c.cancel()
	return nil
}
