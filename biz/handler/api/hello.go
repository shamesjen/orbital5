package api

import (
	"context"
	"encoding/json"
	//"errors"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/opentracing/opentracing-go"
	tracer "github.com/shamesjen/orbital5/pkg/tracer"
)

// Hello handles a POST request to the /hello endpoint. It reads JSON from the request,
// makes a Thrift call using the provided service, and returns the response as JSON.
// @router /hello [POST]
func Hello(ctx context.Context, c *app.RequestContext) {
	const IDLPATH = "idl/hello.thrift"
	var jsonData map[string]interface{}
	var service = "hello"

	// Setup Jaegar tracing
	defer tracer.InitTracer("client").Close()

	// Retrieve raw response data
	response := c.GetRawData()

	err := json.Unmarshal(response, &jsonData)
	if err != nil {
		log.Printf("Error unmarshalling request data: %v", err)
		c.String(consts.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	fmt.Println(jsonData)

	responseFromRPC, err := makeThriftCall(IDLPATH, service, jsonData, ctx)
	if err != nil {
		log.Printf("Error in Thrift call: %v", err)
		c.String(consts.StatusBadRequest, "Internal server error during Thrift call")
		return
	}

	fmt.Println("Post request successful")
	c.JSON(consts.StatusOK, responseFromRPC)
}

// makeThriftCall performs a Thrift call to the specified service using the provided IDL file and JSON data.
func makeThriftCall(IDLPath string, service string, jsonData map[string]interface{}, ctx context.Context) (interface{}, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "makeThriftCall")
	defer span.Finish()	
	
	p, err := generic.NewThriftFileProvider(IDLPath)

	if err != nil {
		return nil, fmt.Errorf("error creating Thrift file provider: %w", err)
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		return nil, fmt.Errorf("error creating Thrift generic: %w", err)
	}

	r, err := etcd.NewEtcdResolver([]string{"etcd:2379"})
	if err != nil {
		return nil, fmt.Errorf("error creating Etcd resolver: %w", err)
	}

	cli, err := genericclient.NewClient(service, g, client.WithResolver(r), client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()))
	// cli, err := genericclient.NewClient(service, g, client.WithResolver(r))
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	jsonString, _ := json.Marshal(jsonData)
	resp, err := cli.GenericCall(ctx, service, string(jsonString))
	if err != nil {
		return nil, fmt.Errorf("error making generic call: %w", err)
	}

	respString, ok := resp.(string)
	if !ok {
		return nil, fmt.Errorf("response is not a string. Actual value: %v", resp)
	}

	fmt.Println("Generic call successful:", respString)

	var respData map[string]interface{}
	err = json.Unmarshal([]byte(respString), &respData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	fmt.Println("Response:", respData["message"])
	return respData, nil
}
