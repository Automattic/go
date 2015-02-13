## go-gravatar

Simple go library to interact with Gravatar.

## Usage

If you already have the hash, build a new Gravatar object and generate its url like so:

    g := NewGravatar()
    g.Hash = "hash-goes-here"
    url := g.GetURL()

If you have an email address:

    g := NewGravatarFromEmail("foo@example.com")
    url := g.GetURL()

## Contributors

You? Pull requests welcome.


## License

[GNU General Public License, version 2](http://www.gnu.org/licenses/gpl-2.0.html)
