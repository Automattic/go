package photon

import (
	"net/url"
	"testing"
)

func TestGetPhotonURL(t *testing.T) {
	if u, err := GetPhotonURL("invalid", nil); err == nil {
		t.Fatalf("invalid URL; expecting error; got: %v", err)
	} else if u != "" {
		t.Fatalf("invalid URL; expecting empty URL; got: %v", u)
	}

	if u, err := GetPhotonURL("/path/to/file.jpg", nil); err != ErrEmptyHost {
		t.Fatalf("relative URL; expecting error; got: %v", err)
	} else if u != "" {
		t.Fatalf("relative URL; expecting empty URL; got: %v", u)
	}

	if u, err := GetPhotonURL("https://i1.wp.com/example.com/file.jpg", nil); err != ErrAlreadyPhotonURL {
		t.Fatalf("already Photon URL; expecting error; got: %v", err)
	} else if u != "" {
		t.Fatalf("already Photon URL; expecting empty URL; got: %v", u)
	}

	tests := []struct {
		description string
		imageURL    string
		params      url.Values
		expectURL   string
	}{
		{
			description: "no params",
			imageURL:    "http://example.com/file.jpg",
			params:      nil,
			expectURL:   "https://i0.wp.com/example.com/file.jpg",
		},
		{
			description: "with params",
			imageURL:    "http://example.com/file.jpg",
			params:      url.Values{"w": {"123"}},
			expectURL:   "https://i0.wp.com/example.com/file.jpg?w=123",
		},
	}

	for _, test := range tests {
		testURL(t, test.description, test.imageURL, test.params, test.expectURL)
	}
}

func testURL(t *testing.T, description string, imageURL string, params url.Values, expectURL string) {
	gotURL, err := GetPhotonURL(imageURL, params)
	if err != nil {
		t.Errorf("Got err when generating URL; %v", err)
	}
	if gotURL != expectURL {
		t.Errorf("expected: %v; got: %v", expectURL, gotURL)
	}
}
