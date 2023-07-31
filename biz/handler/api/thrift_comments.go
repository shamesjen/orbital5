package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	api "github.com/shamesjen/orbital5/biz/model/api"
	constants "github.com/shamesjen/orbital5/pkg/constants"
	tracer "github.com/shamesjen/orbital5/pkg/tracer"
)

// Comment handles a POST request to add a comment to a video, identified by the current VideoID.
// The request body must contain JSON data.
// @router /comment [POST]
func Comment(ctx context.Context, c *app.RequestContext) {
	const IDLPATH = "idl/comment.thrift"
	var jsonData map[string]interface{}
	var service = "comment"
	
	// Setup Jaegar tracing
	defer tracer.InitTracer(service).Close()

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

// Edit handles a PUT request to edit a comment. The request is expected to contain
// a CommentRequest object, and the response includes an API Response object.
// @router /edit [PUT]
func Edit(ctx context.Context, c *app.RequestContext) {
	// Setup Jaegar tracing
	defer tracer.InitTracer("edit").Close()

	var req api.CommentRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		log.Printf("Error binding and validating request: %v", err)
		c.String(consts.StatusBadRequest, "Invalid request format")
		return
	}

	// TODO: Implement the actual edit logic here
	resp := new(api.Response)

	c.JSON(consts.StatusOK, resp)
}
