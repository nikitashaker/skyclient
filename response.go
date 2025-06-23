package skyclient

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
)

func getResponse(conn net.Conn) (*Response, error) {
	header := make([]byte, 1)
	_, err := conn.Read(header)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response header: %w", err)
	}
	// fmt.Println(header) 
	// /*
	switch header[0] {
	case 0x10:
		errcode := make([]byte, 1)
		_, err := conn.Read(errcode)
		if err != nil {
			return nil, fmt.Errorf("Failed to read error code: %w", err)
		}
		return &Response{Type: 0x10, ErrMsg: fmt.Sprintf("Error. Code: %d", errcode[0])}, nil
	case 0x11:
		result, err := getRows(conn)
		if err != nil {
			return nil, err
		}
		return &Response{Type: 0x11, DataSingle: result.Value}, nil
	case 0x12:
		return &Response{Type: 0x12, DataSingle: nil}, nil
	case 0x13:
		result, err := getRows(conn)
		if err != nil {
			return nil, err
		}
		return &Response{Type: 0x13, DataSingle: result.Value}, nil
	case 0xFF:
		return &Response{Type: 0xFF, ErrMsg: "Pipeline error"}, nil
	default:
		return nil, fmt.Errorf("Unknown response type: 0x%X", header[0])
	}
	// */
	// return nil, nil
}

func getRows(conn net.Conn) (*ValueResponse, error) {
	colCountByte, err := getLenghtOrValue(conn)
	if err != nil {
		return nil, err
	}
	colCount, _ := strconv.Atoi(string(colCountByte))

	var values []interface{}

	for i := 0; i < colCount; i++ {
		dataType := make([]byte, 1)
		_, err := io.ReadFull(conn, dataType)
		if err != nil {
			return nil, err
		}
		switch dataType[0] {
		case 0x00: // null
			values = append(values, nil)
		case 0x01: // boolean
			respByte := make([]byte, 1)
			_, err := conn.Read(respByte)
			if err != nil {
				return nil, err
			}
			if respByte[0] == 0x01 {
       			values = append(values, true)
			} else if respByte[0] == 0x00 {
				values = append(values, false)
			} else {
				return nil, fmt.Errorf("Invalid boolean value: expected '0' or '1', got '%c'", respByte[0])
			}
		case 0x02: // unsigned
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 8)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x03:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 16)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x04:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 32)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x05:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 64)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x06: // signed
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 8)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x07:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 16)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x08:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 32)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x09:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 64)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x0C: // binary
			lenBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			len, _ := strconv.Atoi(string(lenBytes))

			Buf := make([]byte, len)
			_, err = conn.Read(Buf)
			if err != nil {
				return nil, err
			}

			values = append(values, Buf)
		case 0x0D: // string
			lenBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			length, _ := strconv.Atoi(string(lenBytes))

			strBuf := make([]byte, length)
			_, err = conn.Read(strBuf)
			if err != nil {
				return nil, err
			}
			values = append(values, string(strBuf))
		case 0x0B: // float64
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseFloat(string(numBytes), 64)
			if err != nil {
				return nil, err
			}
			values = append(values, num)
		case 0x0E: // list
			lenBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			var listValue []interface{}
			len, _ := strconv.Atoi(string(lenBytes))
			for i := 0; i < len; i++ {
				result, err := readValueForList(conn)
				if err != nil {
					return nil, err
				}
				listValue = append(listValue, result)
			}
			values = append(values, listValue)
			default:
				return nil, fmt.Errorf("Undefined type")
		}
	}
	return &ValueResponse{Value: values}, nil
}

func getLenghtOrValue(conn net.Conn) ([]byte, error) {
	var buf bytes.Buffer
	tmp := make([]byte, 1)
	for {
		_, err := io.ReadFull(conn, tmp)
		if err != nil {
			return nil, err
		}
		if tmp[0] == '\n' {
			break
		}
		buf.WriteByte(tmp[0])
	}

	return buf.Bytes(), nil
}

func readValueForList(conn net.Conn) (interface{}, error) {
	typeByte := make([]byte, 1)
	_, err := io.ReadFull(conn, typeByte)
	if err != nil {
		return nil, err
	}

	switch typeByte[0] {
		case 0x00: // null
			return nil, nil
		case 0x01: // boolean
			respByte, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			if respByte[0] == '1' {
       			return true, nil
			} else if respByte[0] == '0' {
				return false, nil
			} else {
				return nil, fmt.Errorf("Invalid boolean value: expected '0' or '1', got '%c'", respByte[0])
			}
		case 0x02: // unsigned
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 8)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x03:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 16)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x04:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 32)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x05:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseUint(string(numBytes), 10, 64)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x06: // signed
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 8)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x07:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 16)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x08:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 32)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x09:
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseInt(string(numBytes), 10, 64)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x0C: // binary
			lenBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			len, _ := strconv.Atoi(string(lenBytes))

			Buf := make([]byte, len)
			_, err = conn.Read(Buf)
			if err != nil {
				return nil, err
			}

			return Buf, nil
		case 0x0D: // string
			lenBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			length, _ := strconv.Atoi(string(lenBytes))

			strBuf := make([]byte, length)
			_, err = conn.Read(strBuf)
			if err != nil {
				return nil, err
			}
			return string(strBuf), nil
		case 0x0B: // float64
			numBytes, err := getLenghtOrValue(conn)
			if err != nil {
				return nil, err
			}
			num, err := strconv.ParseFloat(string(numBytes), 64)
			if err != nil {
				return nil, err
			}
			return num, nil
		case 0x0E: // list
			return readList(conn)
		default:
				return nil, fmt.Errorf("Undefined type")
	}
}

func readList(conn net.Conn) ([]interface{}, error) {
	lenBytes, err := getLenghtOrValue(conn)
	if err != nil {
		return nil, err
	}
	l, _ := strconv.Atoi(string(lenBytes))
	var result []interface{}
	for i := 0; i < l; i++ {
		val, err := readValueForList(conn)
		if err != nil {
			return nil, err
		}
		result = append(result, val)
	}
	return result, nil
}