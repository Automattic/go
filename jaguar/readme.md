
# jaguar.go

Marcus Kazmierczak, [mkaz.com][1]

Jaguar is a library to make HTTP requests a little easier in Go. I use it to get
and send requests to API. This replaces my previous [Fetcher][2] library but is
not API compatible hence the new name.

In general, requests aren't too bad using the `net/http` library. I built
this library to help with testing a REST API which required setting headers,
uploading images and creating more complex requests, the basic `net/http`
package becomes a bit challenging and verbose.

The previous library did not set headers correctly for multipart and basic form
requests, Jaguar is an attempt to clean that up. I've also improved the
interface, but will likely tweak a little as I go. 

If you don't need file upload, check out https://github.com/parnurzeal/gorequest
which otherwise has a bit more features and a nice chainable interface.


## Install

```
$ go get github.com/mkaz/jaguar
```


## Usage

### GET Example

Here's a basic example using fetcher and GET request

```go
import "github.com/mkaz/jaguar"

j := jaguar.New()
resp, err := j.Get("https://google.com/").Send()
if err != nil {
    fmt.Println("Error fetching:", err)
}

fmt.Println("Status Code: ", resp.StatusCode)
```

### POST Example

Example using fetcher to POST params to a form

```go
j := jaguar.New()
j.Params.Add("q", "golang")
j.Url("https://google.com/")
resp, err := j.Method("POST").Send()
if err != nil {
    fmt.Println("Error Fetching:", err)
}
fmt.Println(resp.String())
```

### File Upload Example

Example using fetcher to upload files, set parameters and header variable

```go
j := jaguar.New()
j.Url("/upload-file")
j.Header.Add("X-Auth", "my-secret-token")
j.Params.Add("foo", "bar")
j.Params.Add("baz", "foz")
j.Files["filedata"] = "/home/mkaz/tmp/upload.jpg"
resp, err := j.Method("POST").Send()
if err != nil {
    fmt.Println("Error Fetching:", err)
}
fmt.Println(resp.String())
```

## License

This software is licensed under the MIT License.


[1]: https://mkaz.com/
[2]: https://github.com/mkaz/fetcher

