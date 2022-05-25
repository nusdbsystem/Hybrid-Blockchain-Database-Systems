package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nusdbsystem/hybrid/dbconn"
	pbv "github.com/nusdbsystem/hybrid/veritas_hotstuff/proto/veritashs"
	"google.golang.org/grpc"
)

func main() {
	addr := "192.168.20.5:50001"
	key := "user6284781860667377211" //time.StampMilli       //"user1"
	redisaddr := "192.168.20.5:6379"
	value := time.Now().String() //"66666666666666666666666666"
	version, _ := strconv.ParseInt("3", 10, 64)
	//entry, err := veritas.Encode(value, version)

	rdb, err := dbconn.NewRedisConn(redisaddr, "", 1)
	if err != nil {
		fmt.Println("New server redis db fail: " + redisaddr)
	}
	// err1 := rdb.Set(context.Background(), key, value, 0)
	// fmt.Println(err1)
	res2, err2 := rdb.Get(context.Background(), key).Result()
	fmt.Println(res2)
	if err2 != nil {
		log.Println(err2)
	}
	// v, err := veritas.Decode(res2)
	// log.Println(v.Val)
	// log.Println(v.Version)
	// os.Exit(1)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli := pbv.NewVeritasNodeClient(conn)

	res, err := cli.Set(context.Background(), &pbv.SetRequest{Key: key, Value: value, Version: version})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.String())
	}
	fmt.Println("Hotstuff Set done.")
	time.Sleep(time.Duration(5) * time.Second)

	res1, err := cli.Get(context.Background(), &pbv.GetRequest{Key: key})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res1.Value)
		fmt.Println(res1.Version)
	}
	fmt.Println("Hotstuff Get done.")

	os.Exit(1)
	// rediskv, err := storage.NewRedisKV(redisaddr, "", 1)
	// if err != nil {
	// 	panic(err)
	// }
	// err = rediskv.Set([]byte(key), []byte(entry))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("rediskv.Set done.")
	// val, err := rediskv.Get([]byte(key))
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(string(val))
	// }
	// fmt.Println("rediskv.Get done.")

	// sets := []*pbv.SetRequest{}
	// for i := 0; i < 1000; i++ {
	// 	sets = append(sets, &pbv.SetRequest{Key: "key" + strconv.Itoa(i), Value: "value" + strconv.Itoa(i)})
	// }

	// res3, err3 := cli.BatchSet(context.Background(), &pbv.BatchSetRequest{Sets: sets})
	// if err != nil {
	// 	fmt.Println(res3, err3)
	// }
	// valnew, err := rediskv.Get([]byte("key0"))
	// fmt.Println("rediskv.Get done.")
	// fmt.Println(string(valnew), err)

	// res2, err := cli.Get(context.Background(), &pbv.GetRequest{Key: "user406358257191842439"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Hotstuff Get done.")
	// fmt.Println(res2.Value)

}
