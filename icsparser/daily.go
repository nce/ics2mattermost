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


  var err error

  ingest := c.gatherRelevantEvents()
  if len(ingest.EventsToday) == 0 {

    logger.Error(err.Error())
    return nil, errors.New("no events today")

  } else  {
    logger.Info(fmt.Sprintf("amount of meetings: %d", len(ingest.EventsToday)))

    loc := time.Local
    var formattedEvents string

    for _, e := range ingest.EventsToday {
      formattedEvents = formattedEvents + ":calendar: " + e.Start.In(loc).Format("15:04") + " - " +
        e.End.In(loc).Format("15:04 MST") + " :fire: [" + e.Summary + "](" + e.Location + ")\n"
    }

    dailyMessage := map[string]string{
      "name": "Foobar",
      "text": "#### Welcome to VRP's daily ingest\n" +
      formattedEvents +
      ":airplane: " + strings.Join(ingest.TravellingPersons, ", ") + "\n" +
      ":palm_tree: " + strings.Join(ingest.AbsentPersons, ", "), 
    }

//    logger.Info(
//        fmt.Sprintf("Sent out daily digest with %d persons travelling " +
//                    "and %d persons absent", strings.Count(ingest.TravellingPersons, ","),
//                    strings.Count(ingest.AbsentPersons, ",")))

    return dailyMessage, nil

  }
}
