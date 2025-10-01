package icsparser

import (
  "fmt"
  "time"
  "errors"
  "strings"

	"github.com/nce/ics2mattermost/logger"
)

type DailyIngest struct {
  EventsToday []Event
  TravellingPersons []string
  AbsentPersons []string
}

func (c *Calendar) gatherRelevantEvents() DailyIngest {

  var ingest DailyIngest

  for _, event := range c.Events {
    travelers, err := event.GetPersonsByCategory("travel")

    if err == nil {
      ingest.TravellingPersons = travelers
      continue
    }

    absents, err := event.GetPersonsByCategory("leaves")
    if err == nil {
      ingest.AbsentPersons = absents
      continue
    }

    ingest.EventsToday = append(ingest.EventsToday, event)
  }

  return ingest
}

func (c *Calendar) PrepareDailyIngest() (map[string]string, error) {

  ingest := c.gatherRelevantEvents()
  if len(ingest.EventsToday) == 0 {

    return nil, errors.New("no events today")

  } else  {
    logger.Info(fmt.Sprintf("amount of meetings: %d", len(ingest.EventsToday)))

    loc := time.Local
    var formattedEvents, travellers, absentees string

    for _, e := range ingest.EventsToday {
      formattedEvents = formattedEvents + ":calendar: " + e.Start.In(loc).Format("15:04") + " - " +
        e.End.In(loc).Format("15:04 MST") + " :fire: [" + e.Summary + "](" + e.Location + ")\n"
    }

    if len(ingest.AbsentPersons) == 0 {
      absentees = "*no one*"
    } else {
      absentees = strings.Join(ingest.AbsentPersons, ", ")
    }

    if len(ingest.TravellingPersons) == 0 {
      travellers = "*no one*"
    } else {
      travellers = strings.Join(ingest.TravellingPersons, ", ")
    }

    dailyMessage := map[string]string{
      "name": "Foobar",
      "text": "#### Welcome to VRP's daily ingest\n" +
      formattedEvents +
      ":airplane: " + travellers + "\n" +
      ":palm_tree: " + absentees,
    }

//    logger.Info(
//        fmt.Sprintf("Sent out daily digest with %d persons travelling " +
//                    "and %d persons absent", strings.Count(ingest.TravellingPersons, ","),
//                    strings.Count(ingest.AbsentPersons, ",")))

    return dailyMessage, nil

  }
}
