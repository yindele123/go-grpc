package main

import (
	"context"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"project/common/register"
	"project/order_srv/global"
	"project/order_srv/handler"
	"project/order_srv/initialize"
	"project/order_srv/proto"
	"project/order_srv/utils"
	"syscall"
)

// HealthImpl 健康检查实现
type HealthImpl struct{}

// Check 实现健康检查接口，这里直接返回健康状态，这里也可以有更复杂的健康检查策略，比如根据服务器负载来返回
func (h *HealthImpl) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

//Watch 这个没用，只是为了让HealthImpl实现RegisterHealthServer内部的interface接口
func (h *HealthImpl) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return nil
}

func main() {
	initialize.InitConfig()
	//初始化数据库
	initialize.InitMysql()
	initialize.InitLogger()

	g := grpc.NewServer()
	orderServer := handler.OrderServer{}
	proto.RegisterOrderServer(g, &orderServer)
	grpc_health_v1.RegisterHealthServer(g, &HealthImpl{}) //比普通的grpc开启多了这一步

	//注册服务
	//register
	uuid, _ := uuid.NewV4()
	serviceId := fmt.Sprintf("%s", uuid)
	var nacosRegister register.Register = register.NacosRegister{
		Host:        global.NacosConfig.NacosInfo.Host,
		Port:        global.NacosConfig.NacosInfo.Port,
		NamespaceId: global.NacosConfig.NacosInfo.NamespaceId,
	}
	port, err := utils.GetFreePort()
	if err == nil {
		global.ServerConfig.Port = port
	}
	rerr := nacosRegister.Register(global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.ServiceName, map[string]string{"idc": "xindele", "name": "yindele123", "server": "order-srv"}, serviceId)
	if rerr != nil {
		zap.S().Panic("注册服务失败:", rerr.Error())
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", global.ServerConfig.Port))
	if err != nil {
		zap.S().Panic("启动失败:", err.Error())
	}
	go func() {
		_ = g.Serve(lis)
	}()

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := nacosRegister.Deregister(serviceId, global.ServerConfig.Host, global.ServerConfig.Port, global.ServerConfig.ServiceName); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
