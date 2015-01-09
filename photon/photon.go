package photon

import (
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"net/url"
	"regexp"
)

var (
	cdnMaxSuffix        = 2
	photonScheme        = "https"
	hostnamePattern     = "i%d.wp.com"
	regexpStripProto, _ = regexp.Compile("^https?://")
	regexpPhotonURL, _  = regexp.Compile("^https?://i[0-9]+.wp.com")

	ErrEmptyHost        = errors.New("photon: URL is missing a hostname")
	ErrAlreadyPhotonURL = errors.New("photon: URL is already Photonized")
)

func IsPhotonURL(u string) bool {
	return regexpPhotonURL.MatchString(u)
}

type Options struct {
	Host   string
	Params url.Values
}

func GetPhotonURL(imageURL string, opts Options) (photonURL string, err error) {
	parsed, err := url.Parse(imageURL)
	if err != nil {
		return photonURL, err
	} else if parsed.Host == "" {
		return photonURL, ErrEmptyHost
	}

	// If already a photon URL, bail
	if IsPhotonURL(imageURL) {
		return photonURL, ErrAlreadyPhotonURL
	}

	// Strip any leading `http(s)://`
	path := regexpStripProto.ReplaceAllString(imageURL, "")

	// Allow host overrides
	var host string
	if opts.Host != "" {
		host = opts.Host
	} else {
		host = getHostname(path)
	}

	// TODO: allow hostname overrides
	u := url.URL{
		Scheme:   photonScheme,
		Host:     host,
		Path:     path,
		RawQuery: opts.Params.Encode(),
	}

	return u.String(), nil
}

func GetSupportedHostnames() []string {
	var hostnames []string
	for i := 0; i <= cdnMaxSuffix; i++ {
		host := fmt.Sprintf(hostnamePattern, i)
		hostnames = append(hostnames, host)
	}
	return hostnames
}

// Determine which Photon server to connect to: `i0`, `i1`, or `i2`.
// Statically hash the subdomain based on the URL, to optimize browser caches.
func getHostname(imageURL string) string {
	suffix := getSubdomainSuffix(imageURL)
	return fmt.Sprintf(hostnamePattern, suffix)
}

func getSubdomainSuffix(imageURL string) int {
	seed := crc32.ChecksumIEEE([]byte(imageURL))
	r := rand.New(rand.NewSource(int64(seed)))
	return r.Intn(cdnMaxSuffix)
}
