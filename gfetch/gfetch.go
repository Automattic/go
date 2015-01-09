// gfetch.go
// https://github.com/automattic/go/gfetch

// A library to make it a bit easier to do HTTP fetches using Google AppEngine
// supports adding headers, posting forms, parameters and uploading files

package gfetch

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"appengine"
	"appengine/urlfetch"
)

type Fetcher struct {
	Context       appengine.Context
	Params        url.Values
	Header, Files map[string]string
	Data          map[string]interface{}
}

// NewFetcher creates a fetcher request instance
func NewFetcher(context appengine.Context) (f Fetcher) {
	f.Context = context
	f.Params = make(url.Values)
	f.Header = map[string]string{}
	f.Files = map[string]string{}
	return f
}

type Response struct {
	StatusCode int
	BodyText   []byte
	Header     http.Header
}

// default Fetch returned results as a string
func (f Fetcher) Fetch(url, method string) (result string, err error) {
	bytes, err := f.FetchBytes(url, method)
	return string(bytes), err
}

func (f Fetcher) JsonRequest(url, method string) (resp Response, err error) {

	jsonStr, err := json.Marshal(f.Data)
	if err != nil {
		return resp, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))

	// add additional user header values
	for k, v := range f.Header {
		request.Header.Add(k, v)
	}

	// add in application/json header
	request.Header.Set("Content-Type", "application/json")

	client := urlfetch.Client(f.Context)
	rs, err := client.Do(request)
	if err != nil {
		return resp, err
	}
	defer rs.Body.Close()

	resp.StatusCode = rs.StatusCode
	resp.Header = rs.Header
	resp.BodyText, _ = ioutil.ReadAll(rs.Body)
	return resp, nil
}

// Return results as byte array
// Useful for unmarshaling json, so don't need to cast back to a byte array
func (f Fetcher) FetchBytes(url, method string) (result []byte, err error) {

	var reqBody io.Reader
	var contentType string

	// check if post and add post params
	if method == "POST" || method == "PATCH" {
		reqBody, contentType, err = f.createPostBody()
		if err != nil {
			return
		}
	} else {
		method = "GET"
		url = url + "?" + f.Params.Encode()

	}

	// build request object
	client := urlfetch.Client(f.Context)
	request, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return
	}

	// need to add header to request for content-type
	// this sets boundaries and builds proper header type
	if method == "POST" {
		request.Header.Add("Content-Type", contentType)
	}

	// add additional user header values
	for k, v := range f.Header {
		request.Header.Add(k, v)
	}

	// execute request
	res, err := client.Do(request)
	if err != nil {
		return
	}

	// process response
	defer res.Body.Close()
	result, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	return
}

// create body for post - includes files, params
func (f Fetcher) createPostBody() (body io.Reader, contentType string, err error) {

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// add parameters first if there are parameters
	// Amazon doesn't like params after File
	// TODO: support multiple values for single parameter
	for k, _ := range f.Params {
		_ = writer.WriteField(k, f.Params.Get(k))
	}

	// add files if we are uploading a file
	for k, v := range f.Files {
		file, err := os.Open(v)
		if err != nil {
			return nil, "", err
		}

		part, err := writer.CreateFormFile(k, filepath.Base(v))
		if err != nil {
			return nil, "", err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, "", err
		}
	}

	err = writer.Close()
	if err != nil {
		return
	}

	// content type might be different due to file uploads
	contentType = writer.FormDataContentType()
	body = &b
	return body, contentType, nil

}
