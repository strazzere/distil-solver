// RedNaga / Tim Strazzere (c) 2018-*

package main

// Statistics is the structure used to generate a JSON structure of stats used by the 'dump' handler
type Statistics struct {
	RequestsHandled int64         `json:",omitempty"`
	StartTime       string        `json:",omitempty"`
	SessionData     []SessionData `json:",omitempty"`
	NormalSessions  int64         `json:",omitempty"`
	NoHandle        int64         `json:",omitempty"`
	Errors          int64         `json:",omitempty"`
	BadAuth         int64         `json:",omitempty"`
	BannedAuth      int64         `json:",omitempty"`
}

// SessionData is the structure used to submit to the handleSession function from the clients perspective
type SessionData struct {
	JsSha1 string `json:"JsSha1"`
	JsURI  string `json:"JsUri"`
	JsData string `json:"JsData"`
}

// Exception is a structure for returning errors in JSON format
type Exception struct {
	Code    int64
	Message string
}

// Task is the structure which gives the client the ability to create multiple requests on their end
// while ensuring the integrity of the session (assuming all rules are followed)
type Task struct {
	URI      string   `json:"uri"`
	Method   string   `json:"method"`
	Headers  []string `json:"headers"`
	Data     string   `json:"data"`
	Interval int64    `json:"interval"`
}

// SessionDataResponse is the structure which Tasks and Headers are stored and returned
// as part of a "session" creation with this API
type SessionDataResponse struct {
	Tasks   []Task   `json:"tasks"`
	Headers []string `json:"headers"`
}

// DistilConfig is the parsed structure for how the distil script is configured
type DistilConfig struct {
	Path              string `json:"uri"`
	XDistilAjax       string `json:"X-Distil-Ajax"`
	HeartbeatInterval int64  `json:"heartbeat-timer"`
}

// Key structure for capped-usage based keys
type Key struct {
	Key   string `json:"key"`
	Limit int64  `json:"daily_limit"`
}
