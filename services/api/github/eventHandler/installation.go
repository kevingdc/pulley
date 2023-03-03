package eventhandler

import (
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/app"
	"github.com/kevingdc/pulley/pkg/idconv"
	"github.com/kevingdc/pulley/pkg/messenger"
	"github.com/kevingdc/pulley/services/api/github/event"
)

type InstallationEventHandler struct {
	event *event.Payload
}

func NewInstallationEventHandler(e *event.Payload) *InstallationEventHandler {
	return &InstallationEventHandler{event: e}
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

	message := app.NewSimpleMessage(user, "Awesome, I'm connected to your GitHub account! "+
		"I'll notify you when there's an update on your pull requests or you get one to review.")

	err = messenger.Send(message)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
