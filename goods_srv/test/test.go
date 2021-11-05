package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"project/goods_srv/proto"
)

func main() {
	conn, err := grpc.Dial("192.168.111.1:62102", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//c := proto.NewGoodsClient(conn)
	c := proto.NewGoodsClient(conn)
	//用户列表
	/*r,err := c.GetGoodsList(context.Background(),&proto.PageInfo{
		Pn:    1,
		PSize: 1,
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r.Data,r.Total)*/
	/*//通过手机号码获取用户信息
	r,err :=c.GetGoodsByMobile(context.Background(),&proto.MobileRequest{
		Mobile:"18602058150",
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r.Id,r.Mobile)*/
	//通过id获取用户信息
	/*r,err :=c.GetGoodsById(context.Background(),&proto.IdRequest{
		Id:1,
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r)*/
	//通过Id更新用户
	/*r,err :=c.UpdateGoods(context.Background(),&proto.UpdateGoodsInfo{
		Id:20,
		NickName: "试试",
		Gender: 1,
		BirthDay:1631332408,
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r)*/
	//添加用户
	s :=[] uint64 {1,2,3 }
	r, err := c.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: s,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(r)
}
