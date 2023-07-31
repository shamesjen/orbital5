package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/cloudwego/kitex/server/genericserver"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/opentracing/opentracing-go"
	constants "github.com/shamesjen/orbital5/pkg/constants"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
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
	servers := make([]server.Server, 3)
	for i := 0; i < 3; i++ {
		// Initialize tracer for this server instance
		serverName := fmt.Sprintf("hello%d", i)
		defer InitTracer(serverName).Close()

		addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("hellorpc:%d", 8888+i))
		if err != nil {
			log.Fatalf("Failed to resolve server address: %v", err)
		}

		impl := &GenericServiceImpl{ServerName: serverName} // Set the server name
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
	span, ctx := opentracing.StartSpanFromContext(ctx, "GenericCall")
	defer span.Finish()
	
	log.Println("Request received on server:", g.ServerName)

	m := request.(string)
	var jsonRequest map[string]interface{}
	err = json.Unmarshal([]byte(m), &jsonRequest)
	if err != nil {
		log.Printf("Error unmarshalling JSON request: %v", err)
		return nil, fmt.Errorf("invalid JSON request")
	}

	// Extract message field
	dataValue, ok := jsonRequest["message"].(string)
	if !ok {
		return nil, fmt.Errorf("field 'message' is not a string")
	}

	// Construct response
	jsonRequest["message"] = "Hello!, " + dataValue
	jsonResponse, err := json.Marshal(jsonRequest)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON response: %v", err)
	}

	return string(jsonResponse), nil
}

func InitTracer(serviceName string) io.Closer {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	closer, err := cfg.InitGlobalTracer(
		serviceName,
	)
	if err != nil {
		log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
	}
	return closer
}