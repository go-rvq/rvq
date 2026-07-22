package login_session

import (
	"crypto/sha256"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetStringHash(v string, len int) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(v)))[:len]
}

func GetIP(r *http.Request) string {
	if r == nil {
		return ""
	}

	ips := GetProxy(r)
	if len(ips) > 0 && ips[0] != "" {
		rip, _, err := net.SplitHostPort(ips[0])
		if err != nil {
			rip = ips[0]
		}
		return rip
	}

	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}

func GetProxy(r *http.Request) []string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}

	return nil
}

func IsTokenValid(v LoginSession) bool {
	return time.Now().Sub(v.ExpiredAt) > 0
}
