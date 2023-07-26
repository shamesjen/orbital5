namespace go api

struct Request {
	1: string message
	2: string data
}

struct Response {
	1: string message
}

service like {
    Response like(1: Request req)
}
