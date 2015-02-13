package gravatar

import (
	"testing"
)

func TestNewGravatarFromEmail(t *testing.T) {
	email := "foo@example.com"
	emailHashed := "b48def645758b95537d4424c84d1a9ff"
	g := NewGravatarFromEmail(email)
	if g.Hash != emailHashed {
		t.Errorf("got hash: %q; expected: %q", g.Hash, emailHashed)
	}
}

func TestGravatarGetURL(t *testing.T) {
	tests := []struct {
		hash     string
		def      string
		rating   string
		size     int
		expected string
	}{
		{"b48def645758b95537d4424c84d1a9ff", "", "", 0, "https://www.gravatar.com/avatar/b48def645758b95537d4424c84d1a9ff"},
		{"b48def645758b95537d4424c84d1a9ff", "", "", 100, "https://www.gravatar.com/avatar/b48def645758b95537d4424c84d1a9ff?s=100"},
		{"b48def645758b95537d4424c84d1a9ff", "404", "", 0, "https://www.gravatar.com/avatar/b48def645758b95537d4424c84d1a9ff?d=404"},
		{"b48def645758b95537d4424c84d1a9ff", "http://wp.com/wp-includes/images/blank.gif", "", 0, "https://www.gravatar.com/avatar/b48def645758b95537d4424c84d1a9ff?d=http%3A%2F%2Fwp.com%2Fwp-includes%2Fimages%2Fblank.gif"},
		{"b48def645758b95537d4424c84d1a9ff", "", "pg", 0, "https://www.gravatar.com/avatar/b48def645758b95537d4424c84d1a9ff?r=pg"},
		{"b48def645758b95537d4424c84d1a9ff", "mm", "r", 200, "https://www.gravatar.com/avatar/b48def645758b95537d4424c84d1a9ff?d=mm&r=r&s=200"},
	}

	for _, test := range tests {
		g := NewGravatar()
		g.Hash = test.hash
		g.Default = test.def
		g.Rating = test.rating
		g.Size = test.size

		if url := g.GetURL(); url != test.expected {
			t.Errorf("got url: %q; expected: %q", url, test.expected)
		}
	}
}
