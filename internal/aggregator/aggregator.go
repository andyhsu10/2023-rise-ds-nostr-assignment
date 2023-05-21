package aggregator

import (
	"context"
	"distrise/internal/client"
	"fmt"

	"github.com/nbd-wtf/go-nostr"
)

func Aggregator() {
	ctx := context.Background()

	relay, err := client.GetClient(ctx, "")
	if err != nil {
		panic(err)
	}

	// Create filters
	var filters = []nostr.Filter{{
		Kinds: []int{1}, // type 1 event (note)
	}}

	sub, err := relay.Subscribe(ctx, filters)
	if err != nil {
		panic(err)
	}

	for ev := range sub.Events {
		// handle returned event.
		// channel will stay open until the ctx is cancelled
		fmt.Println(relay.URL, ev)
	}
}
