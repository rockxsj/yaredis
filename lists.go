package yaredis

import (
//"fmt"
//"reflect"
)

func (c *conn) Lpush(list string, value interface{}) (interface{}, error) {
	//args := []interface{}{list, value}
	return c.send("Lpush", list, value)
}

func (c *conn) Lrange(list string, begin, end int64) (interface{}, error) {
	//args := []interface{}{list, begin, end}
	return c.send("Lrange", list, begin, end)
}

func (c *conn) Lpop(list string) (interface{}, error) {
	//args := []interface{}{list}
	return c.send("LPOP", list)
}

func (c *conn) Blpop(list string, timeout int64) (interface{}, error) {
	//args := []interface{}{list, timeout}
	return c.send("BLPOP", list, timeout)
}
