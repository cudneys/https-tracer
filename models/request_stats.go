package models

import (
	"crypto/tls"
	"encoding/json"
	"net"
	"time"
)

type RequestStats struct {
	RemoteHost            string              `json:"remote_host"`
	Start                 time.Time           `json:"connection_start"`
	DnsStart              time.Time           `json:"dns_start"`
	DnsDone               time.Time           `json:"dns_done"`
	DnsTimeDelta          time.Duration       `json:"dns_time_delta"`
	DnsAddrs              []net.IPAddr        `json:"dns_resolved_addresses"`
	DnsError              string              `json:"dns_error,omitempty"`
	DnsCoalesced          bool                `json:"dns_coalesced"`
	ConnectStart          time.Time           `json:"connection_start_time"`
	ConnectDone           time.Time           `json:"connection_done_time"`
	ConnectTimeDelta      time.Duration       `json:"connection_time_delta"`
	TLSHandshakeStart     time.Time           `json:"tls_handshake_start_time"`
	TLSHandshakeDone      time.Time           `json:"tls_handshake_done_time"`
	TLSHandshakeTimeDelta time.Duration       `json:"tls_handshake_time_delta"`
	TLSConnectionState    tls.ConnectionState `json:"tls_connection_state,omitempty"`
}

func (s *RequestStats) AsJSON() (string, error) {
	val, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
