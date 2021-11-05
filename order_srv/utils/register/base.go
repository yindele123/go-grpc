package register

import "google.golang.org/grpc"

type Register interface {
	Register(address string, port int, name string, tags []string, id string) error
	Deregister(serviceId string) error
	GetServuce(serverName string) (*grpc.ClientConn, error)
}
