package client

import (
	"encoding/json"
)

type Request struct {
	RequestInfo *RequestInfo `json:"RequestInfo"`
	Body        interface{}  `json:"Body"`
}

type RequestInfo struct {
	Context string `json:"context"`
	Query   string `json:"query,omitempty"`
}

// String permit to get request object as Json string
func (r *Request) String() string {
	json, _ := json.Marshal(r)
	return string(json)
}
