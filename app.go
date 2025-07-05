package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx             context.Context
	simulationMode  bool
	statusChannel   chan Message
	messages        []Message
	userOptions     UserOptions
	messageCount    int
	srv             *http.Server
	serverWaitGroup *sync.WaitGroup
}

// NewApp creates a new App application struct
func NewApp(ch chan Message) *App {
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
		simulationMode: false,
		messages:       []Message{},
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
		go a.runSimulationMode(a.ctx)
	} else {
		go a.processingIncoming(a.ctx)
		go a.startServer()
	}
}

func (a *App) GetMessages() []Message {
	return a.messages
}

func (a *App) GetUserOptions() UserOptions {
	return a.userOptions
}

func (a *App) SaveUserOptions(opts UserOptions) error {
	var restartService = false
	if a.userOptions.Port != opts.Port || a.userOptions.DefaultEndpoint != opts.DefaultEndpoint {
		restartService = true
	}
	a.userOptions = opts
	saveUserOptions(opts)
	if restartService {
		a.restartServer()
	}
	return nil
}

func (a *App) receiveNewMessage(msg Message) {
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

func (a *App) processingIncoming(ctx context.Context) {
	slog.Info("Processing incoming messages...")
	for {
		select {
		case <-ctx.Done():
		case msg := <-a.statusChannel:
			a.receiveNewMessage(msg)
		}
	}
}
