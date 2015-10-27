## Cloudup Library for Golang

A library to interact with Cloudup API using Go

A work in progress - Initially only supports create and upload 


## Example

A simple example ignoring errors to make concise, see `cloudup_test.go` for a
more complete example.

```
filename := "/path/to/local"
authToken := "--" // base64("username:password")
title := "My Title"

client := cloudup.Client{}
client.BasicToken = authToken

stream, _ := client.CreateStream(title)
item, _ := client.CreateItem(stream.Id, filename, title)

client.UploadToS3(item)

client.CompleteItem(item)
```

