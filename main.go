package main

import (
  "fmt"
  "os"

  "github.com/nce/ics2mattermost/icsparser"
  "github.com/nce/ics2mattermost/logger"
  "github.com/nce/ics2mattermost/mattermost"

  "strings"
  "time"
  _ "embed"
)

type DailyIngest struct {
  Daily icsparser.Event
  TravellingPersons string
  AbsentPersons string
}

//go:generate sh setVersion.sh
//go:embed version
var Version string

func checkIfEmpty(env string) string {
  ret, ok := os.LookupEnv(env)
  if !ok {
    logger.Fatal(fmt.Sprintf("ENV Var '%s' not set; Check help", env))
  }
  return ret
}

func main() {
  logger.SetupLogging(strings.ToLower("debug"))

  logger.Info(fmt.Sprintf("Application version %s", Version))

  var err error
  var icsUrl, icsUser, icsToken, mattermostUrl string

  icsUrl = checkIfEmpty("ICS_URL")
  icsUser = checkIfEmpty("ICS_USER")
  icsToken = checkIfEmpty("ICS_TOKEN")
  mattermostUrl = checkIfEmpty("MATTERMOST_URL")

  cal := icsparser.Setup(
      icsUrl,
      icsUser,
      icsToken)

  webhook := mattermost.Setup(mattermostUrl)

  cal.GetTodaysEvents()

  logger.Info(fmt.Sprintf("Meetings: %d", len(cal.Events)))
  for _, foo := range cal.Events {
    logger.Info(foo.Summary)
  }

  ingest := DailyIngest{
      Daily: icsparser.Event{},
      TravellingPersons: "*no one*", 
      AbsentPersons: "*no one*",
  }

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
  if err != nil {
    logger.Error(err.Error())
  } else  {

    loc, _ := time.LoadLocation("Europe/Berlin")
    dailyMessage := map[string]string{
      "name": "Foobar",
      "text": "#### Welcome to today's daily ingest\n " +
      ":calendar: " + ingest.Daily.Summary + " -- " + ingest.Daily.Start.In(loc).Format("15:04 MST") +
      " - " + ingest.Daily.End.In(loc).Format("15:04 MST") + "\n" +
      ":link: *Daily* ➞ [Microsoft Teams](" + ingest.Daily.Location + ") \n" +
      ":airplane: " + ingest.TravellingPersons + "\n" +
      ":palm_tree: " + ingest.AbsentPersons,
    }

    logger.Info(
        fmt.Sprintf("Sent out daily digest with %d persons travelling " +
                    "and %d persons absent", strings.Count(ingest.TravellingPersons, ","),
                    strings.Count(ingest.AbsentPersons, ",")))

    webhook.Send(dailyMessage)

  }

}
