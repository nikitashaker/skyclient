package skyclient

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
)

// Simple Response (S)
func sendQuery(conn net.Conn, query string, params ...interface{}) (*Response, error) {
	var serParams [][]byte

	for _, p := range params {
		encoded, err := codeByType(p)
		if err != nil {
			return nil, err
		}
		serParams = append(serParams, encoded)
    }

	var temp bytes.Buffer
	queryBytes := []byte(query)
	queryLen := len(queryBytes)
	temp.WriteString(strconv.Itoa(queryLen))
	temp.WriteByte('\n')
	temp.Write(queryBytes)

	if len(serParams) > 0 {
		for _, p := range serParams {
			temp.Write(p)
		}
	}

	payload := temp.Bytes()

	var buf bytes.Buffer
	buf.WriteByte('S')
	buf.WriteString(strconv.Itoa(len(payload)))
	buf.WriteByte('\n')
	buf.Write(payload)

	// fmt.Printf("FULL RAW REQUEST BYTES:\n%q\n", buf.Bytes())

	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Error while sending request: %w", err)
	}

	// fmt.Println("Запрос отправлен:", query)

	resp, err := getResponse(conn)
	if err != nil {
		return nil, fmt.Errorf("Error while sending response: %w", err)
	}

	return resp, nil
}

// Pipeline (P)
func sendPipeline(conn net.Conn, reqmas ...Request) (*Response, error) {
	var payload bytes.Buffer

	for _, qp := range reqmas {
		var temp bytes.Buffer
		for _, p := range qp.Params {
			serp, err := codeByType(p)
			if err != nil {
				return nil, err
			}
			temp.Write(serp)
		}
		paramPayload := temp.Bytes()

		queryBytes := []byte(qp.Query)
		queryLen := len(queryBytes)

		payload.WriteString(strconv.Itoa(queryLen))
		payload.WriteByte('\n')
		payload.WriteString(strconv.Itoa(len(paramPayload)))
		payload.WriteByte('\n')
		payload.Write(queryBytes)
		payload.Write(paramPayload)
	}

	finalPayload := payload.Bytes()

	var buf bytes.Buffer
	buf.WriteByte('P')
	buf.WriteString(strconv.Itoa(len(finalPayload)))
	buf.WriteByte('\n')
	buf.Write(finalPayload)

	// fmt.Printf("FULL RAW REQUEST BYTES:\n%q\n", buf.Bytes())

	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Error while sending request: %w", err)
	}

	// fmt.Println("Запрос отправлен")

	var finalResponse [][]interface{}
	for i := 0; i < len(reqmas); i++ {
		resp, err := getResponse(conn)
		if err != nil {
			return nil, fmt.Errorf("Error while sending response: %w", err)
		}
		finalResponse = append(finalResponse, resp.DataSingle)
	}
	
	return &Response{DataPipeline: finalResponse}, nil
}