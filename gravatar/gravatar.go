package gravatar

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	defaultScheme   = "https"
	defaultHostname = "www.gravatar.com"
)

type Gravatar struct {
	Scheme  string
	Host    string
	Hash    string
	Default string
	Rating  string
	Size    int
}

type Account struct {
	Domain    string `json:"domain"`
	Display   string `json:"display"`
	Url       string `json:"url"`
	Username  string `json:"username"`
	Verified  bool   `json:"verified"`
	Shortname string `json:"shortname"`
}

type GravatarProfile struct {
	Hash        string    `json:"hash"`
	ProfileUrl  string    `json:"profileUrl"`
	DisplayName string    `json:"displayName"`
	AboutMe     string    `json:"aboutMe"`
	Accounts    []Account `json:"accounts"`
}

type GravatarResponse struct {
	Entry []GravatarProfile `json:"entry"`
}

func NewGravatarFromEmail(email string) Gravatar {
	hasher := md5.Sum([]byte(email))
	hash := hex.EncodeToString(hasher[:])

	g := NewGravatar()
	g.Hash = hash
	return g
}

// username is based off wordpress.com username
func NewGravatarFromUsername(username string) Gravatar {
	g := NewGravatar()

	// fetch gravatar.com/username.json
	gp, err := FetchGravatarProfileByUsername(username)
	if err != nil {
		return g
	}

	g.Hash = gp.Hash
	return g
}

func NewGravatar() Gravatar {
	return Gravatar{
		Scheme: defaultScheme,
		Host:   defaultHostname,
	}
}

func (g Gravatar) GetURL() string {
	path := "/avatar/" + g.Hash

	v := url.Values{}
	if g.Size > 0 {
		v.Add("s", strconv.Itoa(g.Size))
	}

	if g.Rating != "" {
		v.Add("r", g.Rating)
	}

	if g.Default != "" {
		v.Add("d", g.Default)
	}

	url := url.URL{
		Scheme:   g.Scheme,
		Host:     g.Host,
		Path:     path,
		RawQuery: v.Encode(),
	}

	return url.String()
}

func FetchGravatarProfileByUsername(username string) (gp GravatarProfile, err error) {
	url := fmt.Sprintf("http://en.gravatar.com/%s.json", username)
	response, err := http.Get(url)
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return gp, errors.New("Error fetching Gravatar data")
	}

	// slurp the entire body of the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return gp, err
	}

	// decode body to GravatarProfile
	var gr GravatarResponse
	if json.Unmarshal(body, &gr); err != nil {
		return gp, err
	}

	if len(gr.Entry) != 1 {
		return gp, errors.New("Invalid Gravatar Response")
	}

	gp = gr.Entry[0]

	return gp, err
}
