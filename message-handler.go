package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type SimpleMessage struct {
	Id         string    `json:"id"`
	Content    string    `json:"content"`
	ReceivedAt time.Time `json:"receivedAt"`
}

// type PubSubMessage struct {
// 	PublishTime   string            `json:"publishTime"`
// 	Data          string            `json:"data"`
// 	MessageId     string            `json:"messageId"`
// 	Attributes    map[string]string `json:"attributes"`
// 	ExtractedData string
// }

func (a *App) launchHandler(statusChannel chan SimpleMessage) {
	subpath := a.userOptions.DefaultEndpoint
	port := a.userOptions.Port
	http.HandleFunc(subpath, a.messageHandler(statusChannel))

	runtime.LogInfof(a.ctx, "Listening on port %d at %s", port, subpath)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func (a *App) messageHandler(statusChannel chan SimpleMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			runtime.LogErrorf(a.ctx, "could not read body: %s", err)
		}

		id, err := uuid.NewV7()
		if err != nil {
			runtime.LogErrorf(a.ctx, "could not generate id: %s", err)
		}
		msg := SimpleMessage{
			Id:         id.String(),
			Content:    string(body),
			ReceivedAt: time.Now(),
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
