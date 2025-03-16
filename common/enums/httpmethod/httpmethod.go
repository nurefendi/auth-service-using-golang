package enums

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
	POST   HttpMethod = "POST"
	PATCH  HttpMethod = "PATCH"
)

func (c HttpMethod) Name() string {
	switch c {
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case PATCH:
		return "PATCH"
	case GET:
	default:
		return "GET"
	}
	return "GET"
}

func GetValue(val string) HttpMethod {
	switch val {
	case "DELETE":
		return DELETE
	case "PATCH":
		return PATCH
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "GET":
	default:
		return GET
	}
	return GET
}