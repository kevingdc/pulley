package event

import (
	"strconv"

	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/pkg/user"
)

type InstallationEventHandler struct {
	event *EventPayload
}

const (
	ActionInstalled   Action = "created"
	ActionUninstalled Action = "deleted"
)

func (h *InstallationEventHandler) Handle() (EventHandlerResponse, error) {
	event := h.event.Payload.(*github.InstallationEvent)

	action := Action(event.GetAction())
	if action == ActionUninstalled {
		return nil, nil
	}

	id := strconv.FormatInt(event.GetSender().GetID(), 10)
	user, err := user.FindOneByRepositoryIDAndType(id, user.RepoGitHub)
	if err != nil {
		return nil, err
	}

	err = messenger.Send(messenger.Message{
		User: user,
		Content: "Awesome, I'm connected to your GitHub account! " +
			"I'll notify you when there's an update on your pull requests or you get one to review.",
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
