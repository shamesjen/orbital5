namespace go api

struct Request {
	1: string name (api.query="name");
}

struct commentRequest {
    1: string name (api.query="name")
    2: string comment (api.query="comment")
}

struct Response {
	1: string message
}

service thriftCall {
    Response hello(1: Request req) (api.post="/hello")
    Response like(1: Request req) (api.post="/like")
    Response unlike(1: Request req) (api.delete="/unlike")
}

service thriftCallComments {
    Response hello(1: Request req) (api.post="/hello")
    Response comment(1: commentRequest req) (api.post="/comment")
    Response edit(1: commentRequest req) (api.put="/edit")
}
