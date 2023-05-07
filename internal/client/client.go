package client

import (
	"context"

	"github.com/nbd-wtf/go-nostr"

	"distrise/internal/configs"
)

func GetClient(ctx context.Context) (*nostr.Relay, error) {
	config := configs.GetConfig()
	relay, err := nostr.RelayConnect(ctx, config.RelayUrl)
	if err != nil {
		return nil, err
	}

	return relay, nil
}
