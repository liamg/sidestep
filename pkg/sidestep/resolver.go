package sidestep

import (
	"context"
	"net"
)

func newResolver(nameserverAddr string, networkProtocol string) *net.Resolver {
	return &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, networkProtocol, nameserverAddr)
		},
	}
}
