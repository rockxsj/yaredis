package yaredis

import (
//"fmt"
//"log"
)

func (c *conn) Get(key string) (interface{}, error) {
	//args := []interface{}{key}
	args := key
	value, error := c.send("GET", args)
	switch value.(type) {
	case string:
		return value.(string), error
	case nil:
		return nil, error
	default:
		return nil, error
	}
}

func (c *conn) Set(key string, value interface{}) (interface{}, error) {
	//args := []interface{}{key, value}
	status, err := c.send("SET", key, value)
	return status.(string), err
}
