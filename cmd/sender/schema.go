package sender

import "context"

type message struct {
	ID      uint
	To      string
	Subject string
	Content string
	Ctx     context.Context
}

type result struct {
	ID        uint
	MessageId string
	Status    int
	Error     string
}

var reqChan = make(chan message, 1000)
var resultChan = make(chan result)
