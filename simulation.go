package main

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) runSimulationMode(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// Exit the loop when the context is canceled
			return
		default:
			val := fmt.Sprintf(`{"field1":"%s", "field2":"%s", "field3":"%s", "field4":"%s", "field5":%d, "field6":%d}`,
				gofakeit.Noun(),
				gofakeit.Noun(),
				gofakeit.Noun(),
				gofakeit.Noun(),
				gofakeit.Int(),
				gofakeit.Int())

			id, err := uuid.NewV7()
			if err != nil {
				runtime.LogErrorf(a.ctx, "could not generate id: %s", err)
			}

			a.messageCount++

			msg := Message{
				Id:            id.String(),
				Number:        a.messageCount,
				ReceivedAt:    time.Now(),
				Method:        gofakeit.HTTPMethod(),
				Body:          val,
				FormattedBody: a.wrapBodyInMarkdown([]byte(val)),
				ContentLength: int64(len(val)),
				RemoteAddr:    gofakeit.IPv4Address(),
				Header:        fakeHeaders(),
			}
			a.receiveNewMessage(msg)

			channel := make(chan struct{})
			// this is a goroutine which executes asynchronously
			go func() {
				time.Sleep(time.Duration(gofakeit.IntRange(0, 5)) * time.Second)
				// send a message to the channel
				channel <- struct{}{}
			}()

			// setup a channel listener
			select {
			case <-channel:
				// success
			case <-time.After(10 * time.Second):
				// timeout handling
			}
		}
	}
}

func fakeHeaders() map[string][]string {
	headers := make(map[string][]string)
	headers["Content-Type"] = []string{"application/json"}
	headers["Authorization"] = []string{"Bearer " + gofakeit.UUID()}
	headers["User-Agent"] = []string{gofakeit.UserAgent()}
	headers["Accept"] = []string{"*/*"}
	return headers
}
