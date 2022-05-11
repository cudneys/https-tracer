package tracer

import (
	"crypto/tls"
	"fmt"
	"github.com/cudneys/http-tracer/models"
	"net/http"
	"net/http/httptrace"
	"time"
)

type Transport struct {
	current *http.Request
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	return http.DefaultTransport.RoundTrip(req)
}

func (t *Transport) GotConn(info httptrace.GotConnInfo) {
	fmt.Printf("Connection reused for %v? %v\n", t.current.URL, info.Reused)
}

func Trace(requestURL string, logTLSInfo bool) (models.RequestStats, error) {
	stats := models.RequestStats{}
	t := &Transport{}

	req, _ := http.NewRequest("GET", requestURL, nil)
	trace := &httptrace.ClientTrace{
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			stats.DnsStart = time.Now()
			stats.RemoteHost = dnsInfo.Host
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			stats.DnsDone = time.Now()
			stats.DnsTimeDelta = stats.DnsDone.Sub(stats.DnsStart)
			stats.DnsAddrs = dnsInfo.Addrs

			if dnsInfo.Err != nil {
				stats.DnsError = dnsInfo.Err.Error()
			} else {
				stats.DnsError = ""
			}

			stats.DnsCoalesced = dnsInfo.Coalesced
		},
		ConnectStart: func(network, addr string) {
			//fmt.Printf("NETWORK: %s || ADDR: %s\n", network, addr)
			stats.ConnectStart = time.Now()
		},
		ConnectDone: func(network, addr string, err error) {
			stats.ConnectDone = time.Now()
			stats.ConnectTimeDelta = stats.ConnectDone.Sub(stats.ConnectStart)
		},
		TLSHandshakeStart: func() {
			stats.TLSHandshakeStart = time.Now()
		},
		TLSHandshakeDone: func(s tls.ConnectionState, err error) {
			if logTLSInfo == true {
				stats.TLSConnectionState = s
			}
			stats.TLSHandshakeDone = time.Now()
			stats.TLSHandshakeTimeDelta = stats.TLSHandshakeDone.Sub(stats.TLSHandshakeStart)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{Transport: t}
	if _, err := client.Do(req); err != nil {
		return stats, err
	}
	return stats, nil
}
