package model

type Device struct {
	SerialNum string `json:"serialNum,omitempty"`
	Model     string `json:"model,omitempty"`
	IP        string `json:"IP,omitempty"`
}

type Message struct {
	Text string `json:"message"`
}
