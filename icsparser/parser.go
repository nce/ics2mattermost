package icsparser

import (
	"bytes"
	"io"
	"strings"
  "errors"

	"github.com/apognu/gocal"
	"github.com/nce/ics2mattermost/logger"

	"fmt"
	"net/http"
	"time"
)

type ics struct {
  icsUrl string
  authEmail string
  authToken string
}

type Calendar struct {
  cal *gocal.Gocal
  Events []Event
}

type Event struct {
   gocal.Event
}

func (i *ics) queryCalendar() *gocal.Gocal {

  client := &http.Client{}

  req, err := http.NewRequest("GET", i.icsUrl, nil)

  if err != nil {
    panic(fmt.Errorf("Got error %s", err.Error()))
  }

  req.SetBasicAuth(i.authEmail, i.authToken)

  response, err := client.Do(req)
  if err != nil {
    panic(fmt.Errorf("Got error %s", err.Error()))
  }
  defer response.Body.Close()

  calendar, err := io.ReadAll(response.Body)
  if err !=  nil {
    logger.Fatal(fmt.Sprintf("could not read from the http request: %s", err.Error()))
  }

  c := gocal.NewParser(bytes.NewReader(calendar))

  return c
}

func Setup(icsUrl string, authEmail string, authToken string) *Calendar {

  var confluence = ics{
    authToken: authToken,
    authEmail: authEmail,
    icsUrl: icsUrl}

  var cal = &Calendar{confluence.queryCalendar(), nil}

  return cal
}

func (c *Calendar) GetTodaysEvents() {
  //loc, _ := time.LoadLocation("Europe/Berlin")

  // truncate a day at 23:59 to filter only TODAYS events
  start, end := time.Now(), time.Now().Add(1*24*time.Hour).Truncate(24*time.Hour)
  c.cal.Start, c.cal.End = &start, &end

  err := c.cal.Parse()
  if err != nil {
    logger.Fatal(fmt.Sprintf("could not parse calendar: %s", err.Error()))
  }

  // a bit messy, but we need to "cast" gocal.Event to Event
  // to add new Methods on the outside package
  for _, e := range c.cal.Events {
    c.Events = append(c.Events, Event{e})
  }
}

func (c *Calendar) GetEventByName(eventName string) (Event, error) {

  for _, event := range c.Events {
    if event.Summary == eventName {
      return event, nil
    }
  }

  return Event{gocal.Event{Summary: "No event found"}}, errors.New("no event found")
}

func (e *Event) GetPersonsByCategory(calendarCategory string) (string, error) {
  var attendees []string

  for _, category := range e.Categories {
    if category == calendarCategory {
      for _, name := range e.Attendees {
        if name.Cn == "" {
          continue
        }
        attendees = append(attendees, name.Cn)
      }
    }
  }

  if len(attendees) == 0 {
    return "", errors.New("No attendees")
  }

  return strings.Join(attendees, ", "), nil
}

