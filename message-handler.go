package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type PubSubMessage struct {
	Id           string    `json:"id"`
	Subscription string    `json:"subscription"`
	Message      Message   `json:"message"`
	RawMessage   string    `json:"rawMessage"`
	ReceivedAt   time.Time `json:"receivedAt"`
}

type Message struct {
	PublishTime   string            `json:"publishTime"`
	Data          string            `json:"data"`
	MessageId     string            `json:"messageId"`
	Attributes    map[string]string `json:"attributes"`
	ExtractedData string
}

func LaunchHandler(statusChannel chan PubSubMessage) {
	subpath := "/messages"
	port := 2999
	http.HandleFunc(subpath, messageHandler(statusChannel))

	fmt.Printf("Listening on port %d at %s\n", port, subpath)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}

func messageHandler(statusChannel chan PubSubMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("could not read body: %s\n", err)
		}

		msg := PubSubMessage{
			RawMessage: string(body),
			ReceivedAt: time.Now(),
		}
		slog.Info(fmt.Sprintf("Received message: %+v", msg))

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
