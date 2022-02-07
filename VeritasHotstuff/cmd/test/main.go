package main

import (
	"context"
	"fmt"

	pbv "hybrid/VeritasHotstuff/proto/veritas"
	"hybrid/VeritasHotstuff/storage"

	"google.golang.org/grpc"
)

func main() {
	addr := "127.0.0.1:40071"
	key := "user1"
	value := "66666666666666666666666666"
	redisaddr := "127.0.0.1:30071"
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
	}
	fmt.Println("Hotstuff Set done.")
	fmt.Println(res.String())

	val, err := rediskv.Get([]byte(key))
	fmt.Println("rediskv.Get done.")
	fmt.Println(string(val), err)

	res1, err := cli.Get(context.Background(), &pbv.GetRequest{Key: key})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Hotstuff Get done.")
	fmt.Println(res1.Value)

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
