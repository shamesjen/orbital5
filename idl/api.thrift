namespace go api

struct Request {
	1: string message (api.query="msg");
}

struct commentRequest {
    1: string message (api.query="msg")
    2: string comment (api.query="comment")
}

struct Response {
	1: string message
}

service hello {
    Response hello(1: Request req) (api.post="/hello")
}

service thriftCall {
    Response like(1: Request req) (api.post="/like")
    Response unlike(1: Request req) (api.delete="/unlike")
}

service thriftComments {
    Response comment(1: commentRequest req) (api.post="/comment")
    Response edit(1: commentRequest req) (api.put="/edit")
}
