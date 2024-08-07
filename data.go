package main

import "time"

type DataChunk struct {
	Index     int       `json:"index"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

var Data map[string]string

func init() {
	// implemented but not consumed
	Data = make(map[string]string)
	Data["what is the model name?"] = "provider Bot"
	Data["what is the model version?"] = "1.0"
}
