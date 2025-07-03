package main

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx            context.Context
	simulationMode bool
	statusChannel  chan SimpleMessage
	messages       []SimpleMessage
	userOptions    UserOptions
}

// NewApp creates a new App application struct
func NewApp(ch chan SimpleMessage) *App {
	options, err := loadUserOptions()
	if err != nil {
		slog.Error("Cannot determine config directory. Options will not be saved")
	}
	err = saveUserOptions(options)
	if err != nil {
		slog.Error("Cannot save options")
	}

	return &App{
		statusChannel:  ch,
		simulationMode: true,
		messages:       []SimpleMessage{},
		userOptions:    options,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	runtime.LogSetLogLevel(a.ctx, logger.DEBUG)
	if a.simulationMode {
		slog.Info("Running in offline mode")
		go a.runSimulationMode()
	} else {
		go a.processingIncoming()
		go a.launchHandler(a.statusChannel)
	}
}

func (a *App) GetMessages() []SimpleMessage {
	return a.messages
}

func (a *App) GetUserOptions() UserOptions {
	return a.userOptions
}

func (a *App) SaveUserOptions(opts UserOptions) error {
	a.userOptions = opts
	return saveUserOptions(opts)
}

func (a *App) runSimulationMode() {
	cnt := 0
	for {

		for range 4 {
			cnt += 1
			val := fmt.Sprintf(`{"name":"simulated-event-%d"}`, cnt)

			id, err := uuid.NewV7()
			if err != nil {
				runtime.LogErrorf(a.ctx, "could not generate id: %s", err)
			}

			msg := SimpleMessage{
				Id:         id.String(),
				Content:    val,
				ReceivedAt: time.Now(),
			}
			a.receiveNewMessage(msg)

		}

		channel := make(chan struct{})
		// this is a goroutine which executes asynchronously
		go func() {
			time.Sleep(5 * time.Second)
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

func (a *App) receiveNewMessage(msg SimpleMessage) {
	slog.Info("Received message", "message", msg)

	a.messages = append(a.messages, msg)
	// Limit messages to MaxMessagesToKeep
	slog.Info(fmt.Sprintf("MessagesToKeep %d of %d", a.userOptions.MaxMessagesToKeep, len(a.messages)))
	if len(a.messages) > a.userOptions.MaxMessagesToKeep && a.userOptions.MaxMessagesToKeep > 0 {
		excess := len(a.messages) - a.userOptions.MaxMessagesToKeep
		slog.Info(fmt.Sprintf("Trimming %d messages from store", excess))
		a.messages = a.messages[excess:]
	}

	runtime.EventsEmit(a.ctx, "messageReceived", msg)
}

func (a *App) processingIncoming() {
	slog.Info("Processing incoming messages...")
	for {
		msg := <-a.statusChannel
		a.receiveNewMessage(msg)
	}
}
