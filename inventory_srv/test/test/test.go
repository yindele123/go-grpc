package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"project/inventory_srv/proto"
)

func main()  {
	conn,err := grpc.Dial("192.168.111.1:8888",grpc.WithInsecure())
	if err!=nil{
		panic(err)
	}
	defer conn.Close()
	//c := proto.NewGoodsClient(conn)
	c := proto.NewInventoryClient(conn)

	r,err :=c.Sell(context.Background(),&proto.SellInfo{
		GoodsInfo:[]*proto.GoodsInvInfo{{GoodsId: 1,Num: 3},{GoodsId: 2,Num: 5},{GoodsId: 33,Num: 10}},
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r)
}