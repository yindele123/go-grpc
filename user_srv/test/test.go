package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"project/user_srv/proto"
)

func main()  {
	conn,err := grpc.Dial("192.168.111.1:57695",grpc.WithInsecure())
	if err!=nil{
		panic(err)
	}
	defer conn.Close()
	c := proto.NewUserClient(conn)
	//用户列表
	/*r,err := c.GetUserList(context.Background(),&proto.PageInfo{
		Pn:    1,
		PSize: 1,
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r.Data,r.Total)*/

	//通过手机号码获取用户信息
	r,err :=c.GetUserByMobile(context.Background(),&proto.MobileRequest{
		Mobile:"18602058150",
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r.Id,r.Mobile)
	//通过id获取用户信息
	/*r,err :=c.GetUserById(context.Background(),&proto.IdRequest{
		Id:1,
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r)*/
	//通过Id更新用户
	/*r,err :=c.UpdateUser(context.Background(),&proto.UpdateUserInfo{
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
	/*r,err :=c.CreateUser(context.Background(),&proto.CreateUserInfo{
		NickName :"好康人",
		PassWord :"12346789",
		Mobile :"18602058156",
	})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r)*/
}