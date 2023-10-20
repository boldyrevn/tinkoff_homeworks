package model

type Device struct {
	SerialNum string `json:"serialNum"`
	Model     string `json:"model"`
	IP        string `json:"IP"`
}

type Message struct {
	Text string `json:"message"`
}
