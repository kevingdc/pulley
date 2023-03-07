package messenger

import (
	"fmt"

	"github.com/kevingdc/pulley/pkg/app"
	"golang.org/x/sync/errgroup"
)

type Messenger interface {
	CanSend(m *app.Message) bool
	Send(m *app.Message) error
}

var messageHandler messenger

type messenger struct {
	handlers []Messenger
}

func Register(handler Messenger) {
	messageHandler.handlers = append(messageHandler.handlers, handler)
}

func Send(m *app.Message) error {
	if m.User == nil {
		return fmt.Errorf("user is nil")
	}

	if m.Content == nil {
		return fmt.Errorf("content is nil")
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

func SendToUsers(u []*app.User, content *app.MessageContent) error {
	g := new(errgroup.Group)

	for _, user := range u {
		if user == nil {
			continue
		}

		user := user
		g.Go(func() error {
			return SendToUser(user, content)
		})
	}

	return g.Wait()
}

func SendToUser(u *app.User, content *app.MessageContent) error {
	return Send(&app.Message{User: u, Content: content})
}
