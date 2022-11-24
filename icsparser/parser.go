package icsparser

import (
	"bytes"
	"io"

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
   t gocal.Event
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
  start, end := time.Now(), time.Now().Add(2*24*time.Hour).Truncate(24*time.Hour)
  c.cal.Start, c.cal.End = &start, &end

  err := c.cal.Parse()
  if err != nil {
    logger.Fatal(fmt.Sprintf("could not parse calendar: %s", err.Error()))
  }

  //c.Events = c.cal.Events
  for _, foo := range c.cal.Events {
    c.Events = append(c.Events, foo)
  }
  //c.Events = Event{t: c.cal.Events}
}

func (e *Event) GetTravellingPersons() string {
  //return e.Attendees[0].Cn
}

func (c *Calendar) GetDate() {
  logger.Error(c.cal.Start.String())
}
