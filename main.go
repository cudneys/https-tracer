package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"time"
)

type RequestStats struct {
	RemoteHost            string              `json:"remote_host"`
	Start                 time.Time           `json:"connection_start"`
	DnsStart              time.Time           `json:"dns_start"`
	DnsDone               time.Time           `json:"dns_done"`
	DnsTimeDelta          time.Duration       `json:"dns_time_delta"`
	DnsAddrs              []net.IPAddr        `json:"dns_resolved_addresses"`
	DnsError              error               `json:"dns_error,omitempty"`
	DnsCoalesced          bool                `json:"dns_coalesced"`
	ConnectStart          time.Time           `json:"connection_start_time"`
	ConnectDone           time.Time           `json:"connection_done_time"`
	ConnectTimeDelta      time.Duration       `json:"connection_time_delta"`
	TLSHandshakeStart     time.Time           `json:"tls_handshake_start_time"`
	TLSHandshakeDone      time.Time           `json:"tls_handshake_done_time"`
	TLSHandshakeTimeDelta time.Duration       `json:"tls_handshake_time_delta"`
	TLSConnectionState    tls.ConnectionState `json:"tls_connection_state,omitempty"`
}

type transport struct {
	current *http.Request
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.current = req
	return http.DefaultTransport.RoundTrip(req)
}

func (t *transport) GotConn(info httptrace.GotConnInfo) {
	fmt.Printf("Connection reused for %v? %v\n", t.current.URL, info.Reused)
}

func main() {

	requestURL := flag.String("url", "", "Request URL to test")
	logTLSInfo := flag.Bool("log-tls-info", false, "Log TLS info (Including certificate info)")
	flag.Parse()

	fmt.Printf("LOG TLS INFO: %s\n", *logTLSInfo)

	if *requestURL == "" {
		fmt.Println("ERROR: Missing --url flag")
		flag.Usage()
		return
	}

	stats := RequestStats{}
	t := &transport{}

	req, _ := http.NewRequest("GET", *requestURL, nil)
	trace := &httptrace.ClientTrace{
		DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
			stats.DnsStart = time.Now()
			stats.RemoteHost = dnsInfo.Host
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			stats.DnsDone = time.Now()
			stats.DnsTimeDelta = stats.DnsDone.Sub(stats.DnsStart)
			stats.DnsAddrs = dnsInfo.Addrs
			stats.DnsError = dnsInfo.Err
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
			if *logTLSInfo == true {
				stats.TLSConnectionState = s
			}
			stats.TLSHandshakeDone = time.Now()
			stats.TLSHandshakeTimeDelta = stats.TLSHandshakeDone.Sub(stats.TLSHandshakeStart)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := &http.Client{Transport: t}
	if _, err := client.Do(req); err != nil {
		log.Fatal(err)
	}
	dumpStats(stats)
}

func dumpStats(s RequestStats) {
	val, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		fmt.Printf("JSON ERROR: %s\n", err.Error())
		return
	}

	//_ = val
	fmt.Println(string(val))
}
