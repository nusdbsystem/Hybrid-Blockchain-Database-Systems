package main

import (
	"fmt"
	"hybrid/dbconn"
	"strconv"

	"github.com/mediocregopher/radix/v3"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	addr    = kingpin.Flag("addr", "Address of redisql").Default("localhost:6379").String()
	connNum = kingpin.Flag("pool-size", "Size of connection pool").Default("10").Int()
)

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	kingpin.Parse()
	r, err := dbconn.NewRedisqlConn(*addr, *connNum)
	if err != nil {
		panic(err)
	}
	if err := r.Do(radix.Cmd(nil, "REDISQL.CREATE_DB", "BENCH")); err != nil {
		panic(err)
	}
	if err := r.Do(radix.Cmd(
		nil,
		"REDISQL.EXEC",
		"BENCH",
		"CREATE TABLE IF NOT EXISTS test(key TEXT, value TEXT);",
	)); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		fmt.Println("set", i)
		if err := r.Do(radix.Cmd(
			nil,
			"REDISQL.EXEC",
			"BENCH",
			fmt.Sprintf("INSERT INTO test VALUES(%s, %s);", strconv.Itoa(i), strconv.Itoa(i)),
		)); err != nil {
			panic(err)
		}
	}

	for i := 0; i < 10; i++ {
		fmt.Println("get", i)
		var items [][]string
		if err := r.Do(radix.Cmd(
			&items,
			"REDISQL.EXEC",
			"BENCH",
			fmt.Sprintf("SELECT * FROM test WHERE test.key=\"%s\";", strconv.Itoa(i)),
		)); err != nil {
			panic(err)
		}
		fmt.Println(items)
	}
}
