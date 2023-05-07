package main

import (
	"bufio"
	"context"
	"distrise/internal/client"
	"distrise/internal/configs"
	"distrise/internal/libs/keylib"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip19"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config := configs.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	relay, err := client.GetClient(ctx)
	if err != nil {
		panic(err)
	}

	// Subscribe to an event, by entering the user's public key
	reader := bufio.NewReader(os.Stdin)
	var pk string
	var b [65]byte
	fmt.Printf("Using relay: %s\n--------\nPlease enter you public key to subscribe to type 1 event: ", config.RelayUrl)
	if n, err := reader.Read(b[:]); err == nil {
		pk = strings.TrimSpace(string(b[:n]))
	} else {
		panic(err)
	}
	npub := keylib.GetNPub(pk)

	// Create filters
	var filters nostr.Filters
	if _, v, err := nip19.Decode(npub); err == nil {
		t := make(map[string][]string)
		t["p"] = []string{v.(string)}
		filters = []nostr.Filter{{
			Kinds: []int{1}, // type 1 event (note)
			Tags:  t,
			Limit: 3, // Get the three most recent notes
		}}
	} else {
		panic("Not a valid npub!")
	}

	// Submit to relay. Results will be returned on the sub.Events channel
	sub, _ := relay.Subscribe(ctx, filters)

	for event := range sub.Events {
		fmt.Println("Event:", event)
	}

	fmt.Printf("--------\nPlease enter your private key to publish a note: ")
	var sk string
	if n, err := reader.Read(b[:]); err == nil {
		s := string(b[:n])
		sk = strings.TrimSpace(s)
	} else {
		panic(err)
	}
	nsec := keylib.GetNSec(sk)

	if _, s, e := nip19.Decode(nsec); e == nil {
		sk = s.(string)
	} else {
		panic("Not a valid nsec!")
	}

	fmt.Println("Please enter your one line note:")
	content, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	event := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      1,
		Content:   strings.TrimRight(strings.TrimSpace(content), "\r\n"),
	}
	event.Sign(sk)
	relay.Publish(ctx, event)

	fmt.Printf("Publishing event to %s:\n%s\n", config.RelayUrl, event)

	go func() {
		<-sub.EndOfStoredEvents
		fmt.Println("Canceling subscription")
		// The subscription is closed when context ctx is cancelled ("CLOSE" in NIP-01).
		cancel()
	}()
}
