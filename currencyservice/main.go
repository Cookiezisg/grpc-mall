package main

import (
	handler "currencyservice/handler"
	pb "currencyservice/proto"
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
)

const PORT = 50012
const ADDRESS = "127.0.0.1"
const serviceID = "currencyservice-50012"

func main() {
	ipport := ADDRESS + ":" + strconv.Itoa(PORT)

	consulConfig := api.DefaultConfig()

	consulClient, err_consul := api.NewClient(consulConfig)
	if err_consul != nil {
		fmt.Println("Consul创建对象报错:", err_consul)
		return
	}

	reg := api.AgentServiceRegistration{
		Tags:    []string{"currencyservice"},
		Name:    "currencyservice",
		Address: ADDRESS,
		Port:    PORT,
		ID:      serviceID,
		Check: &api.AgentServiceCheck{
			TCP:                            ipport,
			Interval:                       "5s",
			Timeout:                        "3s",
			DeregisterCriticalServiceAfter: "10s",
		},
	}

	err_agent := consulClient.Agent().ServiceRegister(&reg)
	if err_agent != nil {
		fmt.Println("Consul注册服务报错:", err_agent)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCurrencyServiceServer(grpcServer, new(handler.CurrencyService))

	listen, err := net.Listen("tcp", ipport)
	if err != nil {
		fmt.Println("监听端口报错:", err)
		return
	}
	defer listen.Close()

	fmt.Printf("货币服务已启动，监听端口: %d\n", PORT)

	err_grpc := grpcServer.Serve(listen)
	if err_grpc != nil {
		fmt.Println("启动服务报错:", err_grpc)
		return
	}
}
