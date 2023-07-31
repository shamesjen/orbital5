package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/shamesjen/orbital5/pkg/constants"
)

// Like handles a POST request to like a video, identified by the current VideoID.
// The request body is expected to contain JSON data.
// @router /like [POST]
func Like(ctx context.Context, c *app.RequestContext) {
	const IDLPATH = "idl/like.thrift"
	var jsonData map[string]interface{}
	var service = "like"
	var VideoID = constants.CurrentVideoID

	// Retrieve and unmarshal the request data
	response := c.GetRawData()
	err := json.Unmarshal(response, &jsonData)
	if err != nil {
		log.Printf("Error unmarshalling request data: %v", err)
		c.String(consts.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	// Add VideoID to the JSON data
	jsonData["data"] = VideoID
	fmt.Println(jsonData)

	// Make a Thrift call and handle the response
	responseFromRPC, err := makeThriftCall(IDLPATH, service, jsonData, ctx)
	if err != nil {
		log.Printf("Error in Thrift call: %v", err)
		c.String(consts.StatusBadRequest, "Internal server error during Thrift call")
		return
	}

	fmt.Println("Post request successful")
	c.JSON(consts.StatusOK, responseFromRPC)
}

// Unlike handles a POST request to unlike a video, identified by the current VideoID.
// The request body is expected to contain JSON data.
// @router /unlike [POST]
func Unlike(ctx context.Context, c *app.RequestContext) {
	const IDLPATH = "idl/unlike.thrift"
	var jsonData map[string]interface{}
	var service = "unlike"

	// Retrieve and unmarshal the request data
	response := c.GetRawData()
	err := json.Unmarshal(response, &jsonData)
	if err != nil {
		log.Printf("Error unmarshalling request data: %v", err)
		c.String(consts.StatusBadRequest, "Invalid JSON in request body")
		return
	}

	// Add VideoID to the JSON data
	jsonData["data"] = constants.CurrentVideoID
	fmt.Println(jsonData)

	// Make a Thrift call and handle the response
	responseFromRPC, err := makeThriftCall(IDLPATH, service, jsonData, ctx)
	if err != nil {
		log.Printf("Error in Thrift call: %v", err)
		c.String(consts.StatusBadRequest, "Internal server error during Thrift call")
		return
	}

	fmt.Println("Post request successful")
	c.JSON(consts.StatusOK, responseFromRPC)
}
