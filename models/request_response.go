package models

type TraceRequest struct {
	Url           string `json:"url"`
	LogTLSDetails bool   `json:"log_tls_details"`
}

type TraceResponse struct {
	Request TraceRequest `json:"trace_request"`
	Trace   RequestStats `json:"trace_response"`
}
