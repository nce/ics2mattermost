package confluence

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmail2Mattermost(t *testing.T) {
  mail := "foo.bar@example.org"
  var test string
  var err error

  test, err = email2Mattermost(mail)
  assert.Equal(t, "@fbar", test)
  assert.Nil(t, err)

  wrongEmail := "ssss@example.org"
  test, err = email2Mattermost(wrongEmail)
  assert.Equal(t, "", test)
  assert.NotNil(t, err)

  wrongDomain := "ssss@a"
  test, err = email2Mattermost(wrongDomain)
  assert.Equal(t, "", test)
  assert.NotNil(t, err)


  subdomainEmail := "foo.bar@exp.example.org"
  test, err = email2Mattermost(subdomainEmail)
  assert.Equal(t, "@fbar", test)
  assert.Nil(t, err)

}
