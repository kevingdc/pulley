package messenger

import (
	"fmt"

	"github.com/kevingdc/pulley/pkg/app"
	"golang.org/x/sync/errgroup"
)

type Message struct {
	User    *app.User
	Content string
}

type Messenger interface {
	CanSend(m Message) bool
	Send(m Message) error
}

var messageHandler messenger

type messenger struct {
	handlers []Messenger
}

func Register(handler Messenger) {
	messageHandler.handlers = append(messageHandler.handlers, handler)
}

func Send(m Message) error {
	if m.User == nil {
		return fmt.Errorf("user is nil")
	}

	g := new(errgroup.Group)

	for _, handler := range messageHandler.handlers {
		if handler.CanSend(m) {
			handler := handler
			g.Go(func() error {
				return handler.Send(m)
			})
		}
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}
