package messenger

import "github.com/kevingdc/pulley/pkg/user"

type Message struct {
	User    user.User
	Content string
}

type Messenger interface {
	CanSend(m Message) bool
	Send(m Message)
}

var messageHandler messenger

type messenger struct {
	handlers []Messenger
}

func Register(handler Messenger) {
	messageHandler.handlers = append(messageHandler.handlers, handler)
}

func Send(m Message) {
	for _, handler := range messageHandler.handlers {
		if handler.CanSend(m) {
			go handler.Send(m)
		}
	}
}
