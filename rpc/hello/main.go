package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	etcd "github.com/kitex-contrib/registry-etcd"
	constants "github.com/shamesjen/orbital5/pkg/constants"
)

func main() {
	// Parse IDL file
	p, err := generic.NewThriftFileProvider("idl/hello.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	// Create etcd registry
	r, err := etcd.NewEtcdRegistry([]string{constants.EtcdAddress})
	if err != nil {
		log.Fatalf("Failed to create etcd registry: %v", err)
	}

	// Create and start servers
	servers := make([]server.Server, constants.NumServers)
	for i := 0; i < constants.NumServers; i++ {
		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("hellorpc:%d", 8888+i))
		if err != nil {
			log.Fatalf("Failed to resolve server address: %v", err)
		}

		impl := &GenericServiceImpl{ServerName: fmt.Sprintf("hello%d", i)} // Set the server name
		svr := genericserver.NewServer(
			impl,
			g,
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "hello"}),
			server.WithServiceAddr(addr),
			server.WithRegistry(r),
		)
		if err != nil {
			panic(err)
		}

		servers[i] = svr

		// Start server
		go func(svr server.Server) {
			err := svr.Run()
			if err != nil {
				log.Fatalf("Failed to start server: %v", err)
			}
		}(servers[i])
	}

	select {} // Prevent main from exiting
}

// GenericServiceImpl handles generic calls for the hello service.
type GenericServiceImpl struct {
	ServerName string
}

// GenericCall processes the hello request and constructs the response.
func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	log.Println("Request received on server:", g.ServerName)

	m := request.(string)
	var jsonRequest map[string]interface{}
	err = json.Unmarshal([]byte(m), &jsonRequest)
	if err != nil {
		log.Printf("Error unmarshalling JSON request: %v", err)
		return nil, fmt.Errorf("Invalid JSON request")
	}

	// Extract message field
	dataValue, ok := jsonRequest["message"].(string)
	if !ok {
		return nil, fmt.Errorf("Field 'message' is not a string")
	}

	// Construct response
	jsonRequest["message"] = "Hello!, " + dataValue
	jsonResponse, err := json.Marshal(jsonRequest)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling JSON response: %v", err)
	}

	return string(jsonResponse), nil
}
