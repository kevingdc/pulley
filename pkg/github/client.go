package github

import (
	"encoding/json"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v50/github"
	"github.com/kevingdc/pulley/pkg/config"
)

type payloadInstallation struct {
	Installation *github.Installation `json:"installation"`
}

func NewClientFromPayload(config *config.Config, payloadBytes []byte) (*github.Client, error) {
	installationID, err := ParseInstallationID(payloadBytes)
	if err != nil {
		return nil, err
	}

	client, err := NewClient(config, installationID)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewClient(config *config.Config, installationID int64) (*github.Client, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, config.GithubAppID, installationID, config.GithubAppPrivateKeyPath)
	if err != nil {
		return nil, err
	}

	client := github.NewClient(&http.Client{Transport: itr})

	return client, nil
}

func ParseInstallationID(payloadBytes []byte) (int64, error) {
	payload := &payloadInstallation{}
	err := json.Unmarshal(payloadBytes, payload)
	if err != nil {
		return 0, err
	}

	installationID := payload.Installation.GetID()

	return installationID, nil
}
