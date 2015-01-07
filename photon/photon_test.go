package photon

import (
	"net/url"
	"testing"
)

func TestGetPhotonURL(t *testing.T) {
	var emptyOpts Options
	if u, err := GetPhotonURL("invalid", emptyOpts); err == nil {
		t.Fatalf("invalid URL; expecting error; got: %v", err)
	} else if u != "" {
		t.Fatalf("invalid URL; expecting empty URL; got: %v", u)
	}

	if u, err := GetPhotonURL("/path/to/file.jpg", emptyOpts); err != ErrEmptyHost {
		t.Fatalf("relative URL; expecting error; got: %v", err)
	} else if u != "" {
		t.Fatalf("relative URL; expecting empty URL; got: %v", u)
	}

	if u, err := GetPhotonURL("https://i1.wp.com/example.com/file.jpg", emptyOpts); err != ErrAlreadyPhotonURL {
		t.Fatalf("already Photon URL; expecting error; got: %v", err)
	} else if u != "" {
		t.Fatalf("already Photon URL; expecting empty URL; got: %v", u)
	}

	tests := []struct {
		description string
		imageURL    string
		opts        Options
		expectURL   string
	}{
		{
			description: "no params",
			imageURL:    "http://example.com/file.jpg",
			opts:        emptyOpts,
			expectURL:   "https://i0.wp.com/example.com/file.jpg",
		},
		{
			description: "with params",
			imageURL:    "http://example.com/file.jpg",
			opts: Options{
				Params: url.Values{"w": {"123"}},
			},
			expectURL: "https://i0.wp.com/example.com/file.jpg?w=123",
		},
		{
			description: "host override",
			imageURL:    "http://example.com/file.jpg",
			opts: Options{
				Host: "i6.wp.com",
			},
			expectURL: "https://i6.wp.com/example.com/file.jpg",
		},
	}

	for _, test := range tests {
		testURL(t, test.description, test.imageURL, test.opts, test.expectURL)
	}
}

func testURL(t *testing.T, description string, imageURL string, opts Options, expectURL string) {
	gotURL, err := GetPhotonURL(imageURL, opts)
	if err != nil {
		t.Errorf("Got err when generating URL; %v", err)
	}
	if gotURL != expectURL {
		t.Errorf("expected: %v; got: %v", expectURL, gotURL)
	}
}

func TestGetSupportedHostnames(t *testing.T) {
	expected := []string{
		"i0.wp.com",
		"i1.wp.com",
		"i2.wp.com",
	}

	got := GetSupportedHostnames()

	if len(got) != len(expected) {
		t.Fatalf("length mismatch; expected: %+v; got: %+v", len(expected), len(got))
	}

	for i := 0; i < len(expected); i++ {
		if got[i] != expected[i] {
			t.Fatalf("value mismatch at index %d; expected: %+v; got: %+v", i, expected[i], got[i])
		}
	}
}
