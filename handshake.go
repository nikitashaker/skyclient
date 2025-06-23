package skyclient

import (
	"bytes"
	"strconv"
	"fmt"
)

func (c *Client) handshake() error {
	handshake := buildHandshake(c.Username, c.Password)
	_, err := c.Conn.Write(handshake)
	if err != nil {
		return fmt.Errorf("Error with sending handshake: %w", err)
	}

	resp := make([]byte, 4)
	_, err = c.Conn.Read(resp)
	if err != nil {
		return fmt.Errorf("Error with reading response: %w", err)
	}

	if bytes.Equal(resp, []byte{'H', 0x00, 0x00, 0x00}) {
		return nil
	} else if resp[0] == 'H' && resp[1] == 0x01 {
		return fmt.Errorf("Cannot auth: %d\n", resp[2])
	} else {
		return fmt.Errorf("Unexpectable response: %v\n", resp)
	}
}

func buildHandshake(username, password string) []byte {
	uLen := strconv.Itoa(len(username))
	pLen := strconv.Itoa(len(password))

	var buf bytes.Buffer
	buf.WriteByte(0x48)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteString(uLen)
	buf.WriteByte('\n')
	buf.WriteString(pLen)
	buf.WriteByte('\n')
	buf.WriteString(username)
	buf.WriteString(password)
	return buf.Bytes()
}