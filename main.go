package main

import (
  "fmt"
  "os"
  "time"

  "github.com/nce/ics2mattermost/icsparser"
  "github.com/nce/ics2mattermost/logger"
  "github.com/nce/ics2mattermost/mattermost"
  "github.com/nce/ics2mattermost/confluence"

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

  var icsUrl, cUser, cToken, mattermostUrl string
  var cApi, cContentId, memberPage string
  var err error

  cApi = checkIfEmpty("CONFLUENCE_API")
  cUser = checkIfEmpty("CONFLUENCE_USER")
  cToken = checkIfEmpty("CONFLUENCE_TOKEN")

  icsUrl = checkIfEmpty("ICS_URL")
  mattermostUrl = checkIfEmpty("MATTERMOST_URL")
  cContentId = checkIfEmpty("CONFLUENCE_MEMBERPAGE_ID")

  webhook := mattermost.Setup(mattermostUrl)

  // get confluence memberPage
  // parse page to extract all emails
  //
  // convert emails to mmost handles
  // provide a function to get the next presenter

  memberPage, err = confluence.Init(cApi, cUser, cToken, cContentId)
  if err != nil {
    logger.Fatal(fmt.Sprintf("could not build confluence API: %s", err.Error()))
  }

  var mhandles []string
  for _, email := range confluence.ExtractAddresses(memberPage) {
    mhandle, err := confluence.Email2Mattermost(email)
    if err != nil {
      logger.Info(fmt.Sprintf("couldn't parse project members email %s", email))
    }

    mhandles = append(mhandles, mhandle)
    logger.Info(fmt.Sprintf("found %s", mhandle))
  }

  return

  s := gocron.NewScheduler(time.Local)

  //s.Every(1).Weeks().Monday().Tuesday().Wednesday().Thursday().At("8:30").Do(func() {
  s.Every(60).Seconds().Do(func() {

    cal := icsparser.Setup(
        icsUrl,
        cUser,
        cToken)

    cal.GetTodaysEvents()

    dailyMessage, err := cal.PrepareDailyIngest()
    if err == nil {
      webhook.Send(dailyMessage)
    } else {
      logger.Warn(
        fmt.Sprintf("could not prepare daily: %s", err.Error()))
    }

  })

  s.StartBlocking()
}
