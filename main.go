package main

import (
	"fmt"
	"os"

	"github.com/nce/ics2mattermost/icsparser"
	"github.com/nce/ics2mattermost/logger"
	"github.com/nce/ics2mattermost/mattermost"

	"strings"
  _ "embed"
)

type DailyIngest struct {
  Daily icsparser.Event
  TravellingPersons string
  AbsentPersons string
}

//go:generate bash setVersion.sh
//go:embed Version
var Version string

func main() {
  logger.SetupLogging(strings.ToLower("debug"))

  logger.Info(fmt.Sprintf("Application version %s", Version))

  var err error
  var icsUrl, icsUser, icsToken string
  var ok bool

  icsUrl, ok = os.LookupEnv("ICS_URL")
  if !ok {
    logger.Fatal("ENV Var ICS_URL not set")
  }
  icsUser, ok = os.LookupEnv("ICS_USER")
  if !ok {
    logger.Fatal("ENV Var ICS_USER not set")
  }
  icsToken, ok = os.LookupEnv("ICS_TOKEN")
  if !ok {
    logger.Fatal("ENV Var ICS_TOKEN not set")
  }

  cal := icsparser.Setup(
      icsUrl,
      icsUser,
      icsToken)

  webhook := mattermost.Setup(os.Getenv("MATTERMOST_URL"))

  cal.GetTodaysEvents()

  logger.Info(fmt.Sprintf("Meetings: %d", len(cal.Events)))
  for _, foo := range cal.Events {
    logger.Info(foo.Summary)
  }

  ingest := DailyIngest{Daily: icsparser.Event{}, TravellingPersons: "*no one*", AbsentPersons: "*no one*"}

  for _, event := range cal.Events {
    travelers, err := event.GetPersonsByCategory("travel")

    if err == nil {
      ingest.TravellingPersons = travelers
    }

    absents, err := event.GetPersonsByCategory("leaves")
    if err == nil {
      ingest.AbsentPersons = absents
    }
  }

  ingest.Daily, err = cal.GetEventByName("DAILY (ALL)")
  logger.Info(err.Error())
  if err == nil {

    dailyMessage := map[string]string{
      "name": "Foobar",
      "text": "#### Welcome to today's daily ingest\n " +
      ":calendar: " + ingest.Daily.Summary + " -- " + ingest.Daily.Start.Format("15:04 MST") +
      " - " + ingest.Daily.End.Format("15:04 MST") + "\n" +
      ":link: *Daily* ➞ [Microsoft Teams](" + ingest.Daily.Location + ") \n" +
      ":airplane: " + ingest.TravellingPersons + "\n" +
      ":palm_tree: " + ingest.AbsentPersons,
    }
    webhook.Send(dailyMessage)

  }

}
