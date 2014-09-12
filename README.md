##yaredis的意思是yet another redis


###使用方法：

```golang
package main

import (
	"fmt"
	"yaredis"
    "time"
)

func main() {
    conn := yaredis.ConnTimeout("127.0.0.1", "6379", time.Millisecond, time.Millisecond, time.Millisecond)
    conn.Lrange("mylisttrue", 0, 3)
    conn.Get("abc")
    conn.Close()
    conn = yaredis.Conn("127.0.0.1", "6379")
    status, err := conn.Command("ping")
    status, err = conn.Get("haha")
    fmt.Println(status, err)
}
```

