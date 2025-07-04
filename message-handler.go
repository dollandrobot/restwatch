package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Message struct {
	Id            string              `json:"id"`
	ReceivedAt    time.Time           `json:"receivedAt"`
	Method        string              `json:"method"`
	Body          string              `json:"body"`
	BodyMarkdown  string              `json:"bodyMarkdown"`
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

func (a *App) launchHandler(statusChannel chan Message) {
	subpath := a.userOptions.DefaultEndpoint
	port := a.userOptions.Port
	http.HandleFunc(subpath, a.messageHandler(statusChannel))

	runtime.LogInfof(a.ctx, "Listening on port %d at %s", port, subpath)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func (a *App) wrapBodyInMarkdown(body []byte) string {
	var jsonMap map[string]interface{}
	err := json.Unmarshal(body, &jsonMap)
	if err != nil {
		return fmt.Sprintf("```\n%s\n```", body)
	}

	formatted, err := json.MarshalIndent(jsonMap, "", "  ")
	if err != nil {
		return fmt.Sprintf("```\n%s\n```", body)
	}
	return fmt.Sprintf("```js\n%s\n```", string(formatted))
}

func (a *App) messageHandler(statusChannel chan Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			runtime.LogErrorf(a.ctx, "could not read body: %s", err)
		}

		id, err := uuid.NewV7()
		if err != nil {
			runtime.LogErrorf(a.ctx, "could not generate id: %s", err)
		}
		msg := Message{
			Id:            id.String(),
			ReceivedAt:    time.Now(),
			Method:        r.Method,
			Body:          string(body),
			BodyMarkdown:  a.wrapBodyInMarkdown(body),
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

		statusChannel <- msg
	}
}
