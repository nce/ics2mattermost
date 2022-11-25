package mattermost

import (
	//"github.com/nce/ics2mattermost/logger"

	"bytes"
	"encoding/json"
	"net/http"
)

type webhook struct {
  url string
}

func (w *webhook) Send(message map[string]string) {

    json_data, err := json.Marshal(message)
    if err != nil {
        panic(err)
    }

    _, err = http.Post(w.url, "application/json",
        bytes.NewBuffer(json_data))

    if err != nil {
        panic(err)
    }
}


func Setup(webhookUrl string) *webhook {
    w := webhook{url: webhookUrl }

  return &w
}
