package register

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"project/user_srv/global"
)

type ConsulRegister struct {
}

func (c ConsulRegister) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", address, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name
	registration.ID = id
	registration.Port = port
	registration.Tags = tags
	registration.Address = address
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	}
	return nil
}
func (c ConsulRegister) Deregister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}

func (c ConsulRegister) GetAllAervice() (str map[string]*api.AgentService) {
	return
}

func (c ConsulRegister) FilterService(filter string) (str map[string]*api.AgentService){
	return
}
