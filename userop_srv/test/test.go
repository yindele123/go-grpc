package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"project/userop_srv/proto"
)

func main() {
	conn, err := grpc.Dial("192.168.111.1:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//c := proto.NewGoodsClient(conn)
	c := proto.NewUserFavClient(conn)


	r, err := c.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: 1,

	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(r)
}
