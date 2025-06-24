package skyclient

import (
	"bytes"
	"strconv"
	"fmt"
)

func codeByType(par interface{}) ([]byte, error) {
	var encoded []byte
	switch v := par.(type) {
	case int:
		encoded = codeToInt(int64(v))
	case string:
		encoded = codeToString(v)
	case float64:
		encoded = codeToFloat(v)
	case bool:
		encoded = codeToBool(v)
	case uint8:
		encoded = codeToUnsInt(uint64(v))
	case uint16:
		encoded = codeToUnsInt(uint64(v))
	case uint32:
		encoded = codeToUnsInt(uint64(v))
	case uint64:
		encoded = codeToUnsInt(uint64(v))
	case []byte:
		encoded = codeToBinary(v)
	case nil:
		encoded = codeToNull()
	default:
		return nil, fmt.Errorf("Unknown type: %T", v)
	}
	return encoded, nil
}

func codeToNull() []byte {
	return []byte{0x00}
}

func codeToString(str string) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x06)
	buf.WriteString(strconv.Itoa(len([]byte(str))))
	buf.WriteByte('\n')
	buf.WriteString(str)
	return buf.Bytes()
}

func codeToBool(par bool) []byte {
	var bytepar byte
	if par == true {
		bytepar = 0x01
	} else {
		bytepar = 0x00
	}
	var buf bytes.Buffer
	buf.WriteByte(0x01)
	buf.WriteByte(bytepar)
	buf.WriteByte('\n')
	return buf.Bytes()
}

func codeToUnsInt(num uint64) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x02)
	buf.WriteString(strconv.FormatUint(num, 10))
	buf.WriteByte('\n')
	return buf.Bytes()
}

func codeToInt(num int64) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x03)
	buf.WriteString(strconv.FormatInt(num, 10))
	buf.WriteByte('\n')
	return buf.Bytes()
}

func codeToFloat(num float64) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x04)
	buf.WriteString(strconv.FormatFloat(num, 'g', 17, 64))
	buf.WriteByte('\n')
	return buf.Bytes()
}

func codeToBinary(data []byte) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x05)
	buf.WriteString(strconv.Itoa(len(data)))
	buf.WriteByte('\n')
	buf.Write(data)
	return buf.Bytes()
}