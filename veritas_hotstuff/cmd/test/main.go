package main

import (
	"context"
	"fmt"
	"time"

	pbv "github.com/nusdbsystem/hybridveritas_hotstuff/proto/veritashs"
	"github.com/nusdbsystem/hybridveritas_hotstuff/storage"

	"google.golang.org/grpc"
)

func main() {
	addr := "192.168.20.2:50001"
	key := time.StampMilli       //"user1"
	value := time.Now().String() //"66666666666666666666666666"
	redisaddr := "192.168.20.2:6379"
	rediskv, err := storage.NewRedisKV(redisaddr, "", 1)
	if err != nil {
		panic(err)
	}
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cli := pbv.NewVeritasNodeClient(conn)

	res, err := cli.Set(context.Background(), &pbv.SetRequest{Key: key, Value: value})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.String())
	}
	fmt.Println("Hotstuff Set done.")
	time.Sleep(time.Duration(5) * time.Second)
	// err = rediskv.Set([]byte(key), []byte(value))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("rediskv.Set done.")
	val, err := rediskv.Get([]byte(key))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(val))
	}
	fmt.Println("rediskv.Get done.")

	res1, err := cli.Get(context.Background(), &pbv.GetRequest{Key: key})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res1.Value)
	}
	fmt.Println("Hotstuff Get done.")

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
