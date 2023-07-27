// Code generated by hertz generator.

package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	//api "github.com/shamesjen/orbital5/biz/model/api"
	"github.com/shamesjen/orbital5/pkg/constants"
)

// Like .
// @router /like [POST]
func Like(ctx context.Context, c *app.RequestContext) {
	var IDLPATH string = "idl/like.thrift"
	var jsonData map[string]interface{}
	var service = "like"

	//return data in bytes
	response := c.GetRawData()

	err := json.Unmarshal(response, &jsonData)

	if err != nil {
		fmt.Println("Error", err)
		c.String(consts.StatusBadRequest, "bad post request")
		return
	}

	jsonData["data"] = constants.CurrentVideoID

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

// Unlike .
// @router /unlike [POST]
func Unlike(ctx context.Context, c *app.RequestContext) {
	var IDLPATH string = "idl/unlike.thrift"
	var jsonData map[string]interface{}
	var service = "unlike"

	//return data in bytes
	response := c.GetRawData()

	err := json.Unmarshal(response, &jsonData)

	if err != nil {
		fmt.Println("Error", err)
		c.String(consts.StatusBadRequest, "bad post request")
		return
	}

	jsonData["data"] = constants.CurrentVideoID

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
