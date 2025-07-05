package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Message struct {
	Id            string              `json:"id"`
	Number        int                 `json:"number"`
	ReceivedAt    time.Time           `json:"receivedAt"`
	Method        string              `json:"method"`
	Body          string              `json:"body"`
	FormattedBody string              `json:"formattedBody"`
	ContentLength int64               `json:"contentLength"`
	RemoteAddr    string              `json:"remoteAddr"`
	Header        map[string][]string `json:"header"`
}

// type PubSubMessage struct {
// 	PublishTime   string            `json:"publishTime"`
// 	Data          string            `json:"data"`
// 	MessageId     string            `json:"messageId"`
// 	Attributes    map[string]string `json:"attributes"`
// 	ExtractedData string
// }

func (a *App) launchHandler() {
	subpath := a.userOptions.DefaultEndpoint
	port := a.userOptions.Port

	mux := http.NewServeMux()
	mux.HandleFunc(subpath, a.messageHandler())
	a.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		runtime.LogInfof(a.ctx, "Listening on port %d at %s", port, subpath)
		if err := a.srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
		a.serverWaitGroup.Done()
	}()
}

func (a *App) startServer() {
	a.serverWaitGroup = &sync.WaitGroup{}

	a.serverWaitGroup.Add(1)
	a.launchHandler()
}

func (a *App) restartServer() {
	if err := a.srv.Shutdown(context.TODO()); err != nil {
		runtime.LogErrorf(a.ctx, "could not shutdown server: %s", err)
	}

	a.serverWaitGroup.Wait()
	a.startServer()
}

func (a *App) wrapBodyInMarkdown(body []byte) string {
	var jsonMap map[string]any
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		return string(body)
	}

	formatted, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		return string(body)
	}
	return string(formatted)
}

func (a *App) messageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			runtime.LogErrorf(a.ctx, "could not read body: %s", err)
		}

		id, err := uuid.NewV7()
		if err != nil {
			runtime.LogErrorf(a.ctx, "could not generate id: %s", err)
		}

		a.messageCount++

		msg := Message{
			Id:            id.String(),
			Number:        a.messageCount,
			ReceivedAt:    time.Now(),
			Method:        r.Method,
			Body:          string(body),
			FormattedBody: a.wrapBodyInMarkdown(body),
			ContentLength: r.ContentLength,
			RemoteAddr:    r.RemoteAddr,
			Header:        r.Header,
		}

		//err = json.Unmarshal(body, &msg)
		//if err != nil {
		//	fmt.Printf("could not unmarshal body")
		//}
		//
		//decodedData, err := base64.StdEncoding.DecodeString(msg.Message.Data)
		//if err != nil {
		//	fmt.Printf("could not decode data: %s\n", err)
		//	return
		//}
		//
		//if len(decodedData) != 0 {
		//	var data map[string]interface{}
		//	err = json.Unmarshal(decodedData, &data)
		//	if err != nil {
		//		fmt.Printf("could not unmarshal body")
		//	}
		//	msg.Message.ExtractedData = string(decodedData)
		//}

		a.statusChannel <- msg
	}
}
