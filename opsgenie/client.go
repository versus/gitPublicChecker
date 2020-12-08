package opsgenie

import (
	"github.com/opsgenie/opsgenie-go-sdk-v2/client"
)

type Opsgenie struct {
	api    string
	config *client.Config
}

func New(apiKey string) Opsgenie {
	config := &client.Config{
		ApiKey: apiKey,
	}

	return Opsgenie{
		api:    apiKey,
		config: config,
	}
}
