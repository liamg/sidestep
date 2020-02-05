package sidestep

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/miekg/dns"
)

func TestCustomResolver(t *testing.T) {

	dns.HandleFunc("service.", func(w dns.ResponseWriter, r *dns.Msg) {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Compress = false

		switch r.Opcode {
		case dns.OpcodeQuery:
			for _, q := range m.Question {
				rr, err := dns.NewRR(fmt.Sprintf("%s TXT HELLO!", q.Name))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}

			}
		}

		_ = w.WriteMsg(m)
	})

	port := 53535
	server := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}

	go func() {
		_ = server.ListenAndServe()
	}()

	time.Sleep(time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resolver := newResolver("127.0.0.1:"+strconv.Itoa(port), "udp")
	results, err := resolver.LookupTXT(ctx, "test.service")
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "HELLO!", results[0])

	_ = server.Shutdown()

}
