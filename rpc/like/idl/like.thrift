namespace go api

struct Request {
	1: string message
}

struct Response {
	1: string message
}

service like {
    Response like(1: Request req)
}
