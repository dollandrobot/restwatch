package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	simulationMode bool
	statusChannel  chan PubSubMessage
	messages       []PubSubMessage
}

// NewApp creates a new App application struct
func NewApp(ch chan PubSubMessage) *App {
	return &App{
		statusChannel:  ch,
		simulationMode: true,
		messages:       []PubSubMessage{{RawMessage: "{}"}},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	if a.simulationMode {
		slog.Info("Running in offline mode")
		go a.runSimulationMode()
	} else {
		go a.processingIncoming()
	}
}

func (a *App) GetMessages() []PubSubMessage {
	return a.messages
}

func (a *App) runSimulationMode() {
	cnt := 0
	for {
		cnt += 1
		val := fmt.Sprintf(`{"name":"event-%d"}`, cnt)
		msg := PubSubMessage{
			RawMessage: val,
			ReceivedAt: time.Now(),
		}
		runtime.LogPrintf(a.ctx, "Received message: %s", msg)
		a.messages = append(a.messages, msg)
		runtime.EventsEmit(a.ctx, "messageReceived", msg)

		channel := make(chan bool)
		// this is a goroutine which executes asynchronously
		go func() {
			time.Sleep(5 * time.Second)
			// send a message to the channel
			channel <- true
		}()

		// setup a channel listener
		<-channel
		runtime.LogPrintf(a.ctx, "Received value from channel: %+v", val)
	}
}

func (a *App) processingIncoming() {
	slog.Info("Processing incoming messages...")
	for {
		msg := <-a.statusChannel
		slog.Info("Received message", "message", msg)
		a.messages = append(a.messages, msg)
		runtime.EventsEmit(a.ctx, "messageReceived", msg)
	}
}
