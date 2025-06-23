package skyclient

import(
	"net"
)

type Client struct {
	Conn     net.Conn
	Username string
	Password string
	Addr     string
}

type Request struct {
	Query  string
	Params []interface{}
}

type Response struct {
	Type         byte            // Response type from server
	DataSingle   []interface{}   // Data for simple response
	DataPipeline [][]interface{} // Data for pipeline
	ErrMsg       string          // Error message
}

type ValueResponse struct {
	Value []interface{}
}