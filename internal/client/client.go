package client

import (
	"context"

	"github.com/nbd-wtf/go-nostr"

	"distrise/internal/configs"
)

func GetClient(ctx context.Context, relayUrl string) (*nostr.Relay, error) {
	config := configs.GetConfig()
	var url string
	if relayUrl != "" {
		url = relayUrl
	} else {
		url = config.RelayUrl
	}

	relay, err := nostr.RelayConnect(ctx, url)
	if err != nil {
		return nil, err
	}

	return relay, nil
}
