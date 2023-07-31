// Code generated by hertz generator.

package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	etcd "github.com/kitex-contrib/registry-etcd"
	//api "github.com/shamesjen/orbital5/biz/model/api"
)

// Hello .
// @router /hello [POST]
func Hello(ctx context.Context, c *app.RequestContext) {
	var IDLPATH string = "idl/hello.thrift"
	var jsonData map[string]interface{}
	var service = "hello"

	//return data in bytes
	response := c.GetRawData()

	err := json.Unmarshal(response, &jsonData)

	if err != nil {
		fmt.Println("Error", err)
		c.String(consts.StatusBadRequest, "bad post request")
		return
	}

	fmt.Println(jsonData)

	responseFromRPC, err := makeThriftCall(IDLPATH, service, jsonData, ctx)

	if err != nil {
		fmt.Println(err)
		c.String(consts.StatusBadRequest, "error in thrift call ")
		return
	}

	fmt.Println("Post request successful")

	c.JSON(consts.StatusOK, responseFromRPC)
}

func makeThriftCall(IDLPath string, service string, jsonData map[string]interface{}, ctx context.Context) (interface{}, error) {
	p, err := generic.NewThriftFileProvider(IDLPath)
	if err != nil {
		fmt.Println("error creating thrift file provider")
		return 0, err
	}

	g, err := generic.JSONThriftGeneric(p)
	if err != nil {
		return 0, errors.New(("error creating thrift generic"))
	}

	r, err := etcd.NewEtcdResolver([]string{"etcd:2379"})
	if err != nil {
		log.Fatal(err)
	}

	cli, err := genericclient.NewClient(service, g, client.WithResolver(r), client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()))
	// cli, err := genericclient.NewClient(service, g, client.WithResolver(r))

	if err != nil {
		return 0, errors.New(("invalid client name"))
	}

	jsonString, _ := json.Marshal(jsonData)

	resp, err := cli.GenericCall(ctx, service, string(jsonString))

	if err != nil {
		fmt.Println("error making generic call")
		return 0, err
	}

	respString, ok := resp.(string)
	if !ok {
		fmt.Println("resp is not a string. Actual value:", resp)
		return nil, errors.New("resp is not a string")
	}

	fmt.Println("generic call successful:", respString)

	var respData map[string]interface{}

	err = json.Unmarshal([]byte(respString), &respData)
	if err != nil {
		fmt.Println("error unmarshalling response", err)
		return nil, err
	}

	fmt.Println("response:", respData["message"])

	return respData, nil
}
