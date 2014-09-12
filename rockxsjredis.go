package yaredis

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

//消息分隔符
const SEP = "\r\n"

type conn struct {
	conn         net.Conn
	br           *bufio.Reader
	bw           *bufio.Writer
	connTimeout  time.Duration
	writeTimeout time.Duration
	readTimeout  time.Duration
	mu           sync.Mutex
}

var redisConn *conn

func Conn(addr, port string) *conn {
	if redisConn != nil {
		return redisConn
	}
	newConn, err := net.Dial("tcp", addr+":"+port)

	if err != nil {
		log.Fatal(err)
		return nil
	}

	redisConn = &conn{
		conn: newConn,
		br:   bufio.NewReader(newConn),
		bw:   bufio.NewWriter(newConn),
	}
	return redisConn
}
func ConnTimeout(addr, port string, connTimeout, writeTimeout, readTimeout time.Duration) *conn {

	if redisConn != nil {
		return redisConn
	}
	var newConn net.Conn
	var err error
	if connTimeout == 0 {
		newConn, err = net.Dial("tcp", addr+":"+port)
	} else {
		newConn, err = net.DialTimeout("tcp", addr+":"+port, connTimeout)
	}

	if err != nil {
		log.Fatal(err)
		return nil
	}

	redisConn = &conn{
		conn:         newConn,
		br:           bufio.NewReader(newConn),
		bw:           bufio.NewWriter(newConn),
		connTimeout:  connTimeout,
		writeTimeout: writeTimeout,
		readTimeout:  readTimeout,
	}
	return redisConn
}

func (c *conn) send(cmd string, args ...interface{}) (interface{}, error) {
	cmdLen := 1 + len(args)
	sendCmd := "*" + strconv.Itoa(cmdLen) + SEP + "$" + strconv.Itoa(len(cmd)) + SEP + cmd + SEP
	for _, v := range args {
		switch v_val := v.(type) {
		case string:
			sendCmd = sendCmd + "$" + strconv.Itoa(len(v_val)) + SEP + v_val + SEP
		case int:
			sendCmd = sendCmd + "$" + strconv.Itoa(len(strconv.Itoa(v_val))) + SEP + strconv.Itoa(v_val) + SEP
		case int64:
			sendCmd = sendCmd + "$" + strconv.Itoa(len(strconv.FormatInt(v_val, 10))) + SEP + strconv.FormatInt(v_val, 10) + SEP
		}
	}
	if c.writeTimeout != 0 {
		c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	}

	if c.readTimeout != 0 {
		c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	}
	_, err := c.bw.Write([]byte(sendCmd))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	c.bw.Flush()
	return c.getReturn()
}

func (c *conn) getReturn() (interface{}, error) {
	back, _, err := c.br.ReadLine()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	data := string(back)
	firstChar := data[0]
	switch firstChar {
	case '+':
		return c.statusReply(data)
	case '-':
		return c.errorReply(data)
	case ':':
		return c.intReply(data)
	case '$':
		return c.bulkReply(data)
	default:
		return c.multiBulkReply(data)
	}
}

func (c *conn) statusReply(data string) (string, error) {
	return strings.TrimLeft(data, "+"), nil
}

func (c *conn) errorReply(data string) (interface{}, error) {
	slice := strings.SplitN(strings.TrimLeft(data, "-"), " ", 2)
	return nil, errors.New(slice[1])
}

func (c *conn) intReply(data string) (int64, error) {
	dataInt, err := strconv.ParseInt(strings.TrimLeft(data, ":"), 10, 0)
	if err != nil {
		return 0, err
	}
	return dataInt, nil
}

func (c *conn) bulkReply(data string) (interface{}, error) {
	if data[1] == '-' && data[2] == '1' {
		return nil, nil
	}
	dataString, _, err := c.br.ReadLine()
	if err != nil {
		return nil, err
	}
	return string(dataString), err
}

func (c *conn) multiBulkReply(data string) (interface{}, error) {
	n, err := strconv.ParseInt(string(data[1]), 10, 0)
	if err != nil {
		return nil, err
	}
	ret := make([]interface{}, n)
	for i := int64(0); i < n; i++ {
		ret[i], err = c.getReturn()
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (c *conn) Close() {
	c.mu.Lock()
	c.conn.Close()
	redisConn = nil
	c.mu.Unlock()
}

func (c *conn) Command(cmd string, args ...interface{}) (interface{}, error) {
	return c.send(strings.ToUpper(cmd), args...)
}
