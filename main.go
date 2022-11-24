package main

import (
	"os"

	"github.com/nce/ics2mattermost/icsparser"
	"github.com/nce/ics2mattermost/logger"
	"github.com/nce/ics2mattermost/mattermost"

	"strings"
)


func main() {
  logger.SetupLogging(strings.ToLower("debug"))

  cal := icsparser.Setup(
      os.Getenv("ICS_URL"),
      os.Getenv("ICS_USER"),
      os.Getenv("ICS_TOKEN"))

  cal.GetTodaysEvents()
  webhook := mattermost.Setup(os.Getenv("MATTERMOST_URL"))

  for _, foo := range cal.Events {

    logger.Info(foo.GetTravellingPersons())
  }
}
