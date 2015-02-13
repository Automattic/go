package gravatar

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strconv"
)

const (
	defaultScheme   = "https"
	defaultHostname = "www.gravatar.com"
)

func NewGravatarFromEmail(email string) Gravatar {
	hasher := md5.Sum([]byte(email))
	hash := hex.EncodeToString(hasher[:])

	g := NewGravatar()
	g.Hash = hash
	return g
}

func NewGravatar() Gravatar {
	return Gravatar{
		Scheme: defaultScheme,
		Host:   defaultHostname,
	}
}

type Gravatar struct {
	Scheme  string
	Host    string
	Hash    string
	Default string
	Rating  string
	Size    int
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
