package cloudup

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strconv"

	"github.com/automattic/go/jaguar"
)

// TODO: use go http lib to make more generic

var baseURL = "https://api.cloudup.com"

type Client struct {
	BasicToken string
	OAuthToken string
}

type Item struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	DirectUrl   string `json:"direct_url"`
	Filename    string `json:"filename"`
	S3Key       string `json:"s3_key"`
	S3Policy    string `json:"s3_policy"`
	S3Signature string `json:"s3_signature"`
	S3Url       string `json:"s3_url"`
	S3AccessKey string `json:"s3_access_key"`
}

type Stream struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Items []string `json:"items"`
	Url   string   `json:"url"`
}

func apiURL(path string) string {
	return baseURL + path
}

func (client Client) newRequest() jaguar.Jaguar {
	j := jaguar.New()
	switch {
	case client.OAuthToken != "":
		j.Header.Add("Authorization", "Bearer "+client.OAuthToken)
	case client.BasicToken != "":
		j.Header.Add("Authorization", "Basic "+client.BasicToken)
	}
	return j
}

func (client Client) CreateItem(streamId, filename, title string) (ci Item, err error) {
	url := apiURL("/1/items")
	ext := filepath.Ext(filename)
	mimetype := mime.TypeByExtension(ext)

	j := client.newRequest()
	j.Params.Add("filename", filename)
	j.Params.Add("title", title)
	j.Params.Add("stream_id", streamId)
	j.Params.Add("mime", mimetype)

	resp, err := j.Post(url).Send()
	if err != nil {
		return ci, err
	}

	err = json.Unmarshal(resp.Bytes, &ci)
	if err != nil {
		log.Printf("Cloudup response: %v", resp.String())
	}

	ci.Filename = filename
	return ci, err
}

func (client Client) CompleteItem(ci Item) error {
	url := apiURL("/1/items/" + ci.Id)

	j := client.newRequest()
	j.JsonData = map[string]interface{}{
		"complete": true,
	}

	_, err := j.Patch(url).JsonRequest()
	if err != nil {
		return err
	}

	return err
}

func (client Client) CreateStream(title string) (cs Stream, err error) {
	url := apiURL("/1/streams")

	j := client.newRequest()
	j.Params.Add("title", title)

	resp, err := j.Post(url).Send()
	if err != nil {
		return cs, err
	}

	err = json.Unmarshal(resp.Bytes, &cs)
	if err != nil {
		fmt.Println("Error unmarshaling request: ", resp.String())
	}
	return cs, err
}

func (client Client) UploadToS3(ci Item) error {

	fileinfo, err := os.Stat(ci.Filename)
	if err != nil {
		return err
	}

	ext := filepath.Ext(ci.Filename)

	j := jaguar.New()
	j.Url(ci.S3Url)
	j.Params.Add("key", ci.S3Key)
	j.Params.Add("AWSAccessKeyId", ci.S3AccessKey)
	j.Params.Add("acl", "public-read")
	j.Params.Add("policy", ci.S3Policy)
	j.Params.Add("signature", ci.S3Signature)
	j.Params.Add("Content-Type", mime.TypeByExtension(ext))
	j.Params.Add("Content-Length", strconv.FormatInt(fileinfo.Size(), 10))
	j.Files["file"] = ci.Filename
	resp, err := j.Method("POST").Send()
	if err != nil {
		fmt.Println("Error uploading to S3:", err)
		return err
	}

	if resp.StatusCode != 204 {
		return errors.New(fmt.Sprintf("Invalid status code. %v \n", resp.String()))
	}

	return nil
}
