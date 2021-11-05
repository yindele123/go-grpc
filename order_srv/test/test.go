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
	r,err :=c.CreateOrder(context.Background(),&proto.OrderRequest{
		UserId: 1,
	})
	/*r,err :=c.CreateCartItem(context.Background(),&proto.CartItemRequest{
		UserId: 1,
		GoodsId:5,
		Nums: 1,
		Checked: true,
	})*/
	if err!=nil{
		panic(err)
	}
	fmt.Println(r)
}