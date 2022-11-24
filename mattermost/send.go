package mattermost

// http://localhost:8065/hooks/bma5si7katyembq4c945jeuego

import (
  //"github.com/nce/ics2mattermost/logger"

	"github.com/apognu/gocal"

    "bytes"
    "encoding/json"
    "net/http"
  
 )

type webhook struct {
  url string
}

func (w *webhook) Send(event gocal.Event) {

  values := map[string]string{
    "name": "John Doe",
    "text": "#### Welcome to today's daily\n " +
    ":calendar: " + event.Summary + "\n" +
    ":link: " + event.URL + "\n" +
    "" + event.Categories[0] + "\n" +
    "" + event.Attendees[0].Cn,
  }

    json_data, err := json.Marshal(values)
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
