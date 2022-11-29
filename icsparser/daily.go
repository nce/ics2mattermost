package icsparser

import (
  "fmt"
  "time"
  "errors"

	"github.com/nce/ics2mattermost/logger"
)

type DailyIngest struct {
  EventsToday []Event
  TravellingPersons string
  AbsentPersons string
}

func (c *Calendar) PrepareDailyIngest() (map[string]string, error) {

  logger.Info(fmt.Sprintf("amount of meetings: %d", len(c.Events)))

  ingest := DailyIngest{
      EventsToday: []Event{},
      TravellingPersons: "*no one*",
      AbsentPersons: "*no one*",
  }

  for _, event := range c.Events {
    travelers, err := event.GetPersonsByCategory("travel")

    if err == nil {
      ingest.TravellingPersons = travelers
    }

    absents, err := event.GetPersonsByCategory("leaves")
    if err == nil {
      ingest.AbsentPersons = absents
    }
  }

  var err error

  //ingest.Daily, err = c.GetEventByName("DAILY (ALL)")
  if len(c.Events) == 0 {

    logger.Error(err.Error())
    return nil, errors.New("no events today")

  } else  {

    loc := time.Local
    var formattedEvents string

    for _, e := range c.Events {
      formattedEvents = formattedEvents + ":calendar: " + e.Start.In(loc).Format("15:04") + " - " +
        e.End.In(loc).Format("15:04 MST") + " :fire: [" + e.Summary + "](" + e.Location + ")\n"
    }

    dailyMessage := map[string]string{
      "name": "Foobar",
      "text": "#### Welcome to today's daily ingest\n" +
      formattedEvents +
      ":airplane: " + ingest.TravellingPersons + "\n" +
      ":palm_tree: " + ingest.AbsentPersons,
    }

//    logger.Info(
//        fmt.Sprintf("Sent out daily digest with %d persons travelling " +
//                    "and %d persons absent", strings.Count(ingest.TravellingPersons, ","),
//                    strings.Count(ingest.AbsentPersons, ",")))

    return dailyMessage, nil

  }
}
