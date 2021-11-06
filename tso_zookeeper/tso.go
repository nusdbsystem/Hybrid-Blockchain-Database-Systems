package tso

import (
	// "fmt"
	"fmt"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

const PathPrefix = "/veritasts"

type Client struct {
	conn  *zk.Conn
	lock  *sync.Mutex
	maxTS int64
	curTS int64
	path  string
}

var singleton *Client
var once sync.Once

func NewClient(address string) (*Client, error) {
	once.Do(func() {
		zk_servers := make([]string, 1)
		zk_servers[0] = address
		conn, _, err := zk.Connect(zk_servers, 10000000000)

		if err != nil {
			singleton = nil
			panic(err)
		}

		// create a new unique path
		t := time.Now()
		path := fmt.Sprintf("%s-%s", PathPrefix, t.Format("20211010090930"))
		conn.Create(path, []byte{1, 2, 3, 4}, 0, zk.WorldACL(zk.PermAll))

		singleton = &Client{
			conn:  conn,
			lock:  &sync.Mutex{},
			maxTS: 0,
			curTS: 0,
			path:  path,
		}

		// async function for look-ahead counter
		go func() {
			for true {
				_, stat, err := conn.Get(path)
				if err != nil {
					fmt.Printf("TSO ZK Get error %v\n", err)
					continue
				}
				stat, err = conn.Set(path, []byte{1, 2, 3, 4}, stat.Version)
				if err != nil {
					fmt.Printf("TSO ZK Set error %v\n", err)
					continue
				}
				singleton.maxTS += 100 * int64(stat.Version)
				// fmt.Printf("Max TS %d Cur TS %d\n", singleton.maxTS, singleton.curTS)
				time.Sleep(750 * time.Millisecond)
			}
		}()
	})

	return singleton, nil
}

// Close the client after all TS responses are returned
func (c *Client) Close() {
	c.conn.Close()
}

// naive implementation
func (c *Client) TSX() (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, stat, err := c.conn.Get(c.path)
	if err != nil {
		fmt.Printf("TSO ZK Get error %v\n", err)
		return -1, err
	}

	stat, err = c.conn.Set(c.path, []byte{1, 2, 3, 4}, stat.Version)
	if err != nil {
		fmt.Printf("TSO ZK Set error %v\n", err)
		return -1, err
	}

	return int64(stat.Version), nil
}

// faster implementation
func (c *Client) TS() (int64, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.curTS < c.maxTS {
		ret := c.curTS
		c.curTS += 1
		return ret, nil
	}
	for c.curTS == c.maxTS {
		time.Sleep(100 * time.Millisecond)
	}
	ret := c.curTS
	c.curTS += 1
	return ret, nil
}
