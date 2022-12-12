package confluence

import (
  "fmt"
  "strings"
  "errors"

  "regexp"
  "golang.org/x/net/html"

  "github.com/nce/ics2mattermost/logger"
  goconfluence "github.com/virtomize/confluence-go-api"
)

type Confluence struct {
  confluenceAPI string
  authEmail string
  authToken string
  contentID string
}

func Init(cApi string, cUser string, cToken string, cContentId string) (string, error) {

  api, err := goconfluence.NewAPI(cApi, cUser, cToken)
  if err != nil {
    return "", errors.New(
        fmt.Sprintf("error building confluence api: %s", err.Error()))
  }

  // get content by content id
  c, err := api.GetContentByID(cContentId, goconfluence.ContentQuery{
    SpaceKey: "VIP",
    Expand:   []string{"body.storage", "version"},
    })

  if err != nil {
    return "", errors.New(
        fmt.Sprintf("error getting confluence page: %s", err.Error()))
  }

  return c.Body.Storage.Value, nil
}

func ExtractAddresses(htmlpage string) []string {
  doc, err := html.Parse(strings.NewReader(htmlpage))
  if err != nil {
    logger.Error(fmt.Sprintf("error in parsing confluence: %s", err.Error()))
  }

  // traverse DOM tree and extract email addresses
  var emailAddresses []string
  var traverse func(*html.Node)

  traverse = func(n *html.Node) {
    // check if node is an <a> element with valid href attribute
    if n.Type == html.ElementNode && n.Data == "a" && n.FirstChild != nil && n.FirstChild.Type == html.TextNode {
      for _, a := range n.Attr {
        if a.Key == "href" && strings.HasPrefix(a.Val, "mailto:") {
          // remove the mailto:
          emailAddresses = append(emailAddresses, a.Val[7:])
        }
      }
    }

    // traverse child nodes
    for c := n.FirstChild; c != nil; c = c.NextSibling {
      traverse(c)
    }
  }

  traverse(doc)

  return emailAddresses
}

// transform emailadresses to (our) mattermost handles
func Email2Mattermost(email string) (string, error) {
  var mhandle string
  var name []string
  pattern := `\.[a-zA-Z0-9+^_{|}~-]*@`
  r, err := regexp.Compile(pattern)
  if err != nil {
    return "", err
  }

  if ! r.MatchString(email) {
    return "", errors.New(fmt.Sprintf("no valid email address: %s - skipping", email))
  }

  // clear domain part
  tmp := strings.Split(email, "@")[0]

  // split in first and surname
  name = strings.Split(tmp, ".")

  // truncate firstname to 1 letter
  name[0] = name[0][:1]

  // add @ in front
  mhandle = "@" + name[0] + name[1]

  // to lower
  mhandle = strings.ToLower(mhandle)

  return mhandle, nil
}
