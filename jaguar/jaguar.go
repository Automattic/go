// jaguar.go - http request helper
// https://github.com/mkaz/jaguar
//
// A library to make it a bit easier to do HTTP requests
//

package jaguar

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type Jaguar struct {
	RequestUrl    string
	RequestMethod string
	Params        url.Values
	Header        http.Header
	Files         map[string]string
	JsonData      map[string]interface{}
	VerifyCert    bool
}

type Response struct {
	StatusCode int
	Bytes      []byte
	Header     http.Header
}

// convenience function to get body result as string
func (r Response) String() string {
	return string(r.Bytes)
}

// New creates a Jaguar request instance
func New() Jaguar {
	j := Jaguar{}
	j.RequestMethod = "GET"
	j.Params = make(url.Values)
	j.Header = make(http.Header)
	j.Files = map[string]string{}
	j.VerifyCert = true
	return j
}

func (j *Jaguar) Url(url string) *Jaguar {
	j.RequestUrl = url
	return j
}

func (j *Jaguar) Method(method string) *Jaguar {
	j.RequestMethod = method
	return j
}

func (j *Jaguar) SkipVerify() *Jaguar {
	j.VerifyCert = false
	return j
}

func (j *Jaguar) Get(url string) *Jaguar {
	j.RequestUrl = url
	j.RequestMethod = "GET"
	return j
}

func (j *Jaguar) Post(url string) *Jaguar {
	j.RequestUrl = url
	j.RequestMethod = "POST"
	return j
}

func (j *Jaguar) Patch(url string) *Jaguar {
	j.RequestUrl = url
	j.RequestMethod = "PATCH"
	return j
}

func (j *Jaguar) Put(url string) *Jaguar {
	j.RequestUrl = url
	j.RequestMethod = "PUT"
	return j
}

func (j *Jaguar) Delete(url string) *Jaguar {
	j.RequestUrl = url
	j.RequestMethod = "DELETE"
	return j
}

// Send the request
func (j *Jaguar) Send() (resp Response, err error) {

	var requestBody io.Reader

	// check if multipart form, determined by j.FILES set
	if len(j.Files) > 0 {
		requestBody, err = j.createMultiPartBody()
		if err != nil {
			return
		}
	} else if j.RequestMethod == "GET" {
		j.RequestUrl = j.RequestUrl + "?" + j.Params.Encode()
		requestBody = nil
	} else if j.RequestMethod == "POST" || j.RequestMethod == "PATCH" || j.RequestMethod == "PUT" || j.RequestMethod == "DELETE" {
		j.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		formData := []byte(j.Params.Encode())
		requestBody = bytes.NewReader(formData)
	} else {
		err = errors.New("Unknown request method specified: " + j.RequestMethod)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !j.VerifyCert},
	}

	// build request object
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest(j.RequestMethod, j.RequestUrl, requestBody)
	if err != nil {
		return
	}

	request.Header = j.Header

	// execute request
	rs, err := client.Do(request)
	if err != nil {
		return
	}

	// process response
	defer rs.Body.Close()

	resp.StatusCode = rs.StatusCode
	resp.Header = rs.Header
	resp.Bytes, _ = ioutil.ReadAll(rs.Body)

	return resp, nil
}

// create body for post - includes files, params
func (j *Jaguar) createMultiPartBody() (body io.Reader, err error) {

	var b bytes.Buffer

	writer := multipart.NewWriter(&b)

	// add parameters first if there are parameters
	// Amazon doesn't like params after File
	// TODO: support multiple values for single parameter
	for k, _ := range j.Params {
		_ = writer.WriteField(k, j.Params.Get(k))
	}

	// add files if we are uploading a file
	for k, v := range j.Files {
		file, err := os.Open(v)
		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile(k, filepath.Base(v))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	body = &b

	// content type might be different due to file uploads
	j.Header.Set("Content-Type", writer.FormDataContentType())

	return body, nil
}

func (j Jaguar) JsonRequest() (resp Response, err error) {

	jsonStr, err := json.Marshal(j.JsonData)
	if err != nil {
		return resp, err
	}

	request, err := http.NewRequest(j.RequestMethod, j.RequestUrl, bytes.NewBuffer(jsonStr))

	// add additional content-type header for json
	j.Header.Add("Content-Type", "application/json")

	request.Header = j.Header

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !j.VerifyCert},
	}

	client := &http.Client{Transport: tr}
	rs, err := client.Do(request)
	if err != nil {
		return resp, err
	}
	defer rs.Body.Close()

	resp.StatusCode = rs.StatusCode
	resp.Header = rs.Header
	resp.Bytes, _ = ioutil.ReadAll(rs.Body)
	return resp, nil
}
