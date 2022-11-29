package icsparser

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/apognu/gocal"
	"github.com/stretchr/testify/assert"
)

const event = `BEGIN:VCALENDAR
BEGIN:VEVENT
DTSTART:%s
DTEND:%s
DTSTAMP:20141203T130000Z
UID:0002@google.com
END:VEVENT`

var oneEventYesterday = fmt.Sprintf(
  event,
  time.Now().Add(-23*time.Hour).Format("20060102T150407Z"),
  time.Now().Add(-24*time.Hour).Format("20060102T150407Z"))


var oneEventToday = fmt.Sprintf(
  event,
  time.Now().Format("20060102T150407Z"),
  time.Now().Add(time.Hour).Format("20060102T150407Z"))

var oneEventTomorrow = fmt.Sprintf(
  event,
  time.Now().Add(24*time.Hour).Format("20060102T150407Z"),
  time.Now().Add(25*time.Hour).Format("20060102T150407Z"))


func TestGetTodaysEvents(t *testing.T) {
  var cal *gocal.Gocal 
  var c Calendar

  cal = gocal.NewParser(strings.NewReader(oneEventYesterday))
  c = Calendar{cal, nil}
  c.GetTodaysEvents()

  assert.Equal(t, 0, len(c.Events))

  cal = gocal.NewParser(strings.NewReader(oneEventToday))
  c = Calendar{cal, nil}
  c.GetTodaysEvents()

  assert.Equal(t, 1, len(c.Events))

  cal = gocal.NewParser(strings.NewReader(oneEventTomorrow))
  c = Calendar{cal, nil}
  c.GetTodaysEvents()

  assert.Equal(t, 0, len(c.Events))
}

func TestGetPersonsByCategory(t *testing.T) {
  var e Event
  var test string
  var err error

  e.Categories = []string{"travel"}
  e.Attendees = []gocal.Attendee{{Cn: "asdf"}, {Cn: "foo"}}
  test, err = e.GetPersonsByCategory("travel")

  assert.Equal(t, "asdf, foo", test)
  assert.Nil(t, err)

  // wrong category
  e.Categories = []string{"other"}
  e.Attendees = []gocal.Attendee{{Cn: "asdf"}, {Cn: "foo"}}
  test, err = e.GetPersonsByCategory("travel")

  assert.Equal(t, "", test)
  assert.NotNil(t, err)

  // no attendee
  e.Categories = []string{"travel"}
  e.Attendees = []gocal.Attendee{{}, {}}
  test, err = e.GetPersonsByCategory("travel")

  assert.Equal(t, "", test)
  assert.NotNil(t, err)

  // no attendee
  e.Categories = []string{"travel"}
  e.Attendees = []gocal.Attendee{}
  test, err = e.GetPersonsByCategory("travel")

  assert.Equal(t, "", test)
  assert.NotNil(t, err)

  // one attendee on different category
  e.Categories = []string{"leaves"}
  e.Attendees = []gocal.Attendee{{Cn: "Bar Foo"}}
  test, err = e.GetPersonsByCategory("leaves")

  assert.Equal(t, "Bar Foo", test)
  assert.Nil(t, err)

}
