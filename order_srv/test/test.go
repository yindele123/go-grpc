package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"project/order_srv/proto"
)

func main()  {
	conn,err := grpc.Dial("192.168.111.1:8520",grpc.WithInsecure())
	if err!=nil{
		panic(err)
	}
	defer conn.Close()
	//c := proto.NewGoodsClient(conn)
	c := proto.NewOrderClient(conn)
	//添加用户
	r,err :=c.UpdateOrderStatus(context.Background(),&proto.OrderStatus{
		OrderSn: "a213123",
		Status: 2,
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r)
}