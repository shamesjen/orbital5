namespace go api

struct Request {
	1: string message (api.query="msg");
}

struct Response {
	1: string message
}

service hello {
    Response hello(1: Request req) (api.post="/hello")
}