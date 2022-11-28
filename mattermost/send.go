package mattermost

import (
	//"github.com/nce/ics2mattermost/logger"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nce/ics2mattermost/logger"
)

type webhook struct {
  url string
}

func (w *webhook) Send(message map[string]string) {

    json_data, err := json.Marshal(message)
    if err != nil {
        panic(err)
    }

    resp, err := http.Post(w.url, "application/json",
        bytes.NewBuffer(json_data))
    if err != nil {
        panic(err)
    }

    body, _ := io.ReadAll(resp.Body)
    logger.Debug(fmt.Sprintf("Mattermost HTTP status: %s, body: %s", resp.Status, string(body)))
}

func Setup(webhookUrl string) *webhook {
    w := webhook{url: webhookUrl }

  return &w
}
