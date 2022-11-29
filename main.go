package main

import (
  "fmt"
  "os"
  "time"

  "github.com/nce/ics2mattermost/icsparser"
  "github.com/nce/ics2mattermost/logger"
  "github.com/nce/ics2mattermost/mattermost"

  "github.com/go-co-op/gocron"

  "strings"
  _ "embed"
)

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

  var icsUrl, icsUser, icsToken, mattermostUrl string

  icsUrl = checkIfEmpty("ICS_URL")
  icsUser = checkIfEmpty("ICS_USER")
  icsToken = checkIfEmpty("ICS_TOKEN")
  mattermostUrl = checkIfEmpty("MATTERMOST_URL")

  webhook := mattermost.Setup(mattermostUrl)

  s := gocron.NewScheduler(time.Local)

  //s.Every(1).Weeks().Monday().Tuesday().Wednesday().Thursday().At("8:30").Do(func() {
  s.Every(20).Seconds().Do(func() {

    cal := icsparser.Setup(
        icsUrl,
        icsUser,
        icsToken)

    cal.GetTodaysEvents()

    dailyMessage, err := cal.PrepareDailyIngest()
    if err == nil {
      webhook.Send(dailyMessage)
    } else {
      logger.Error(
        fmt.Sprintf("could not prepare daily: %s", err.Error()))
    }

  })

  s.StartBlocking()
}
