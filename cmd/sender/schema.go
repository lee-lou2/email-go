package sender

import "context"

type request struct {
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

var reqChan = make(chan request, 1000)
var resultChan = make(chan result)
