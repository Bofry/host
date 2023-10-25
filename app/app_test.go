package app_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Bofry/host/app"
)

func TestBuild(t *testing.T) {
	ap, err := app.Build("demo",
		app.WithDefaultMessageHandler(func(ctx *app.Context, message *app.Message) {
			ctx.Send(message.Format, message.Protocol, message.Body)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	alice := &MockMessageClient{
		In:  make(chan []byte),
		Out: make(chan []byte),
	}
	bob := &MockMessageClient{
		In:  make(chan []byte),
		Out: make(chan []byte),
	}
	ap.MessageClientManager().Join(alice)
	ap.MessageClientManager().Join(bob)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ap.Start(ctx)

	go func() {
		var words = strings.Split("a quick brown fox jumps over the lazy dog", " ")
		for _, word := range words {
			alice.In <- []byte(fmt.Sprintf("alice: %s", word))
			bob.In <- []byte(fmt.Sprintf("bob: %s", word))
		}
	}()

	go func() {
		for {
			select {
			case b, ok := <-alice.Out:
				if ok {
					t.Logf("alice:: %s", string(b))
				}
			case b, ok := <-bob.Out:
				if ok {
					t.Logf("bob:: %s", string(b))
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		ap.Stop(ctx)
	}
}

func TestInit(t *testing.T) {
	ap := app.Init(&MockModule,
		app.BindServiceProvider(&MockServiceProvider{ID: "service_provider"}),
		app.BindConfig(&MockConfig{Env: "local"}),
		app.BindEventClient(app.MultiEventClient{
			"foo_topic":        nil,
			"bar_topic":        nil,
			app.InvalidChannel: nil,
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ap.Start(ctx)

	client := &MockMessageClient{
		In:                make(chan []byte),
		Out:               make(chan []byte),
		MessageClientInfo: app.NewMessageClientInfo(),
	}
	client.RegisterCloseHandler(func(mc app.MessageClient) {
		t.Logf("ID:: %s", mc.ID())
		t.Logf("StartAt:: %v", mc.StartAt())
	})
	ap.MessageClientManager().Join(client)

	go func() {
		client.In <- []byte(fmt.Sprintf("foo$%s", "foo_proto"))
		client.In <- []byte(fmt.Sprintf("bar$%s", "bar_proto"))
		client.In <- []byte(fmt.Sprintf("nop$%s", "nop_proto"))
	}()

	go func() {
		for {
			select {
			case b, ok := <-client.Out:
				if ok {
					t.Logf("%s", string(b))
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		ap.Stop(ctx)
	}
}
