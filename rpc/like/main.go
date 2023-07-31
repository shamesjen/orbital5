package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	etcd "github.com/kitex-contrib/registry-etcd"
	constants "github.com/shamesjen/orbital5/pkg/constants"
	"github.com/shamesjen/orbital5/pkg/tracer"
)

func main() {
	// Parse IDL file
	p, err := generic.NewThriftFileProvider("idl/like.thrift")
	if err != nil {
		panic(err)
	}
	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		panic(err)
	}

	// Create etcd registry
	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"})
	if err != nil {
		log.Fatalf("Failed to create etcd registry: %v", err)
	}

	// Create and start servers
	servers := make([]server.Server, constants.NumServers)
	for i := 0; i < constants.NumServers; i++ {
		// Initialize tracer for this server instance
		serverName := fmt.Sprintf("like%d", i)
		defer tracer.InitTracer(serverName).Close()

		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("likerpc:%d", 9000+i))
		if err != nil {
			log.Fatalf("Failed to resolve server address: %v", err)
		}

		impl := &GenericServiceImpl{ServerName: serverName} // Set the server name
		svr := genericserver.NewServer(
			impl,
			g,
			server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "like"}),
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

// GenericServiceImpl handles generic calls for the like service.
type GenericServiceImpl struct {
	ServerName string
}

// GenericCall processes the like request and constructs the response.
func (g *GenericServiceImpl) GenericCall(ctx context.Context, method string, request interface{}) (response interface{}, err error) {
	log.Println("Request received on server:", g.ServerName)

	m := request.(string)
	var jsonRequest map[string]interface{}
	err = json.Unmarshal([]byte(m), &jsonRequest)
	if err != nil {
		log.Printf("Error unmarshalling JSON request: %v", err)
		return nil, errors.New("invalid JSON request")
	}

	// Extract fields
	user, dataValue, err := extractFields(jsonRequest)
	if err != nil {
		return nil, err
	}

	// Construct response
	jsonRequest["message"] = fmt.Sprintf("%s has successfully liked VideoID: %s", user, dataValue)
	jsonResponse, err := json.Marshal(jsonRequest)
	if err != nil {
		return nil, errors.New("error marshalling JSON response")
	}

	return string(jsonResponse), nil
}

// extractFields extracts required fields from the JSON request.
func extractFields(jsonRequest map[string]interface{}) (user, dataValue string, err error) {
	user, ok := jsonRequest["message"].(string)
	if !ok {
		return "", "", errors.New("field 'message' is not a string")
	}

	dataValue, ok = jsonRequest["data"].(string)
	if !ok {
		return "", "", errors.New("field 'data' is not a string")
	}

	return user, dataValue, nil
}
