package skyclient

import (
	"fmt"
	"net"
	"reflect"
	"errors"
)

func NewConnection(addr, username, password string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Connection error %w", err)
	}

	client := &Client{
		Conn: conn,
		Username: username,
		Password: password,
		Addr: addr,
	}

	if err := client.handshake(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("Error with auth: %w", err)
	}

	return client, nil
}

func (c *Client) Query(query string, params ...interface{}) ([]interface{}, error) {
	resp, err := sendQuery(c.Conn, query, params...)
	if err != nil {
		return nil, fmt.Errorf(resp.ErrMsg)
	}
	if resp.ErrMsg != "" {
		return nil, fmt.Errorf(resp.ErrMsg)
	}
	return resp.DataSingle, nil
}

// func (c *Client) Query(query string, params ...interface{}) (*Response, error) {
// 	return sendQuery(c.Conn, query, params...)
// }

func (c *Client) QueryParse(query string, out interface{}, params ...interface{}) error {
    resp, err := sendQuery(c.Conn, query, params...)
    if err != nil {
        return fmt.Errorf(resp.ErrMsg)
    }

    if len(resp.DataSingle) == 0 {
        return errors.New("empty response")
    }

    v := reflect.ValueOf(out)
    if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
        return errors.New("out must be a pointer to a structure")
    }

    v = v.Elem()
    t := v.Type()

    for i := 0; i < v.NumField() && i < len(resp.DataSingle); i++ {
        field := v.Field(i)
        if !field.CanSet() {
            continue
        }

        val := reflect.ValueOf(resp.DataSingle[i])

        if val.Type().ConvertibleTo(field.Type()) {
            field.Set(val.Convert(field.Type()))
        } else {
            return fmt.Errorf("Cannot cast field %s (%s) a type %s or the structure does not match the received values" , t.Field(i).Name, val.Type(), field.Type())
        }
    }

    return nil
}

func (c *Client) Pipeline(req ...Request) ([][]interface{}, error) {
	resp, err := sendPipeline(c.Conn, req...)
	if err != nil {
		return nil, fmt.Errorf(resp.ErrMsg)
	}
	if resp.ErrMsg != "" {
		return nil, fmt.Errorf(resp.ErrMsg)
	}
	return resp.DataPipeline, nil
}

func (c *Client) Close() error {
	return c.Conn.Close()
}