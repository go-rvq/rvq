package gitedit

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
)

// ErrPrivateAddress is a download refused: the address is of the server's
// own network — itself, a private network, a link-local address — that a URL
// given by a user must not reach through the server.
var ErrPrivateAddress = errors.New("the address is not public")

// SafeHTTPClient is the client that downloads the URLs a user imports: it
// connects only to public addresses — checked on the address it connects to,
// after any redirect and DNS answer —, gives up after a minute, and follows
// at most 5 redirects.
func SafeHTTPClient() *http.Client {
	dialer := &net.Dialer{
		Timeout: 15 * time.Second,
		Control: func(_, address string, _ syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil {
				return err
			}
			if ip := net.ParseIP(host); ip == nil || !PublicIP(ip) {
				return fmt.Errorf("%s: %w", host, ErrPrivateAddress)
			}
			return nil
		},
	}
	return &http.Client{
		Timeout: time.Minute,
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return dialer.DialContext(ctx, network, addr)
			},
			TLSHandshakeTimeout: 15 * time.Second,
		},
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return errors.New("too many redirects")
			}
			return nil
		},
	}
}

// PublicIP says whether ip is an address of the internet: not the server
// itself (loopback, unspecified), a private network, link-local, multicast,
// nor the shared address space of carriers (100.64.0.0/10).
func PublicIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() || ip.IsMulticast() {
		return false
	}
	if v4 := ip.To4(); v4 != nil && v4[0] == 100 && v4[1]&0xc0 == 64 {
		return false
	}
	return true
}
