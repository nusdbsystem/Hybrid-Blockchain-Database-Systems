package tso

import (
	"bufio"
	"errors"
	"log"
	"net"
)

type Client struct {
	shutdown bool
	req      chan chan int64
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
}

func NewClient(address string) (*Client, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	cl := &Client{
		shutdown: false,
		conn:     conn,
		writer:   bufio.NewWriter(conn),
		reader:   bufio.NewReader(conn),
		req:      make(chan chan int64, 100000),
	}
	go cl.start()

	return cl, nil
}

// Close the client after all TS responses are returned
func (c *Client) Close() {
	c.shutdown = true
}

func (c *Client) TS() (int64, error) {
	if c.shutdown {
		return -1, errors.New("close")
	}
	ch := make(chan int64)
	c.req <- ch
	if ts := <-ch; ts >= 0 {
		return ts, nil
	} else {
		return -1, errors.New("invalid ts")
	}
}

func (c *Client) GetTS(num int32) (int64, error) {
	if c.shutdown {
		return -1, errors.New("already close")
	}

	getTs := &GetTS{num}

	c.writer.WriteByte(byte(GET))
	getTs.Marshal(c.writer)
	c.writer.Flush()

	msgType, err := c.reader.ReadByte()
	if err != nil {
		return -1, err
	}
	switch uint8(msgType) {
	case REPLY:
		replyTs := new(ReplyTS)
		if err := replyTs.Unmarshal(c.reader); err != nil {
			return -1, err
		}
		return replyTs.Timestamp, nil
	default:
		return -1, errors.New("unknown msg type")
	}
}

func (c *Client) start() {
	for !c.shutdown {
		ch := <-c.req
		l := len(c.req)
		// batch count
		// log.Println(l)
		ts, err := c.GetTS(int32(l + 1))
		if err != nil {
			log.Println("get ts error", err)
			c.shutdown = true
			c.conn.Close()
			break
		}
		ch <- ts
		for i := 1; i <= l; i++ {
			ch = <-c.req
			ch <- ts - int64(i)
		}
	}

	close(c.req)
}
