package tso

import (
	// "fmt"
	"github.com/go-zookeeper/zk"
)

const path = "/veritasts"

type Client struct {
	conn     *zk.Conn
}

func NewClient(address string) (*Client, error) {
	zk_servers := make([]string, 1)
	zk_servers[0] = address
	conn, _, err := zk.Connect(zk_servers, 10000000000)

	if err != nil {
		return nil, err
	}

	conn.Create(path, []byte{1, 2, 3, 4}, 0, zk.WorldACL(zk.PermAll))
	
	cl := &Client{
		conn:     conn,
	}
	return cl, nil
}

// Close the client after all TS responses are returned
func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) TS() (int64, error) {
	_, stat, err := c.conn.Get(path);
	if err != nil {
		return -1, err
	}

	stat, err = c.conn.Set(path, []byte{1, 2, 3, 4}, stat.Version)
	if err != nil {
		return -1, err
	}

	return int64(stat.Version), nil
}
