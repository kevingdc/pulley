package eventhandler

import (
	"github.com/google/go-github/github"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type InstallationEventHandler struct {
	event *event.Payload
}

func (h *InstallationEventHandler) Handle() (event.HandlerResponse, error) {
	userService := h.event.App.UserService

	installationEvent := h.event.Payload.(*github.InstallationEvent)

	action := event.Action(installationEvent.GetAction())
	if action == event.ActionUninstalled {
		return nil, nil
	}

	id := idconv.ToRepoID(installationEvent.GetSender().GetID())
	user, err := userService.FindOneByRepositoryIDAndType(id, app.RepoGitHub)
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
