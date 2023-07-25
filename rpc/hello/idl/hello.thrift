namespace go api

struct Request {
	1: string message
}

struct Response {
	1: string message
}

service hello {
    Response hello(1: Request req)
}