namespace go api

struct Request {
	1: string message
	2: string data
    3: string comment
}

struct Response {
	1: string message
}

service comment {
    Response comment(1: Request req)
}
