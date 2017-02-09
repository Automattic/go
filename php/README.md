A package for unserializing values serialized with PHPs serialize() function

[Documentation Link](https://godoc.org/github.com/Automattic/go/php)

Example use

```go
import (
        "fmt"

        "github.com/Automattic/go/php"
)

func main() {
        v, _ := php.Unmarshal([]byte(`O:8:"stdClass":1:{s:7:"message";s:11:"hello world";}`))
        member, _ := v.GetKey("message")
        message, _ := member.String()
        fmt.Println(v.KindString(), "message:", message)
}
```

Example output

```
object message: hello world
```
