// memcached protocol parser
package protocol

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// structure to hold parsed memcache packet
// Exptime will always be 0 or epoch
type McRequest struct {
	Command string
	Key     string
	Keys    []string
	Flags   string
	Exptime int64
	Data    []byte
	Value   int64
	Cas     string
	Noreply bool
}

type ProtocolError struct {
	Description string
}

func (e ProtocolError) Error() string {
	return fmt.Sprintf("Protocol error: %s", e.Description)
}

func NewProtocolError(description string) ProtocolError {
	return ProtocolError{description}
}

func ReadRequest(r *bufio.Reader) (req *McRequest, err error) {

	// todo use a panic error handling pattern

	lineBytes, _, err := r.ReadLine() // todo pref
	if err != nil {
		return nil, err
	}
	line := string(lineBytes)
	arr := strings.Fields(line)
	if len(arr) < 1 {
		return nil, NewProtocolError("empty line")
	}
	// arr[0] = strings.ToLower(arr[0])
	switch arr[0] {
	case "set", "add", "replace", "append", "prepend":
		// <command name> <key> <flags> <exptime> <bytes> [noreply]\r\n
		// <data block>\r\n
		if len(arr) < 5 {
			return nil, NewProtocolError(fmt.Sprintf("too few params to command %q", arr[0]))
		}
		req := &McRequest{}
		req.Command = arr[0]
		req.Key = arr[1]
		req.Flags = arr[2]
		req.Exptime, err = strconv.ParseInt(arr[3], 10, 64)
		if err != nil {
			return nil, NewProtocolError("cannot read exptime " + err.Error())
		}
		if req.Exptime > 0 {
			if req.Exptime < time.Now().Unix() {
				req.Exptime = time.Now().Unix() + req.Exptime
			}
		}
		bytes, err := strconv.Atoi(arr[4])
		if err != nil {
			return nil, NewProtocolError("cannot read bytes " + err.Error())
		}
		if len(arr) > 5 && arr[5] == "noreply" {
			req.Noreply = true
		}
		req.Data = make([]byte, bytes)
		n, err := r.Read(req.Data)
		if err != nil {
			return nil, err
		}
		if n != bytes {
			return nil, NewProtocolError(fmt.Sprintf("Read only %d bytes of %d bytes of expected data", n, bytes))
		}
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if c != '\r' {
			return nil, NewProtocolError("expected \\r")
		}
		c, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		if c != '\n' {
			return nil, NewProtocolError("expected \\n")
		}
		return req, nil
	case "cas":
		// cas <key> <flags> <exptime> <bytes> <cas unique> [noreply]\r\n
		// <data block>\r\n
		if len(arr) < 6 {
			return nil, NewProtocolError(fmt.Sprintf("too few params to command %q", arr[0]))
		}
		req := &McRequest{}
		req.Command = arr[0]
		req.Key = arr[1]
		req.Flags = arr[2]
		req.Exptime, err = strconv.ParseInt(arr[3], 10, 64)
		if err != nil {
			return nil, NewProtocolError("cannot read exptime " + err.Error())
		}
		bytes, err := strconv.Atoi(arr[4])
		if err != nil {
			return nil, err
		}
		req.Cas = arr[5]
		if len(arr) > 6 && arr[6] == "noreply" {
			req.Noreply = true
		}
		req.Data = make([]byte, bytes)
		n, err := r.Read(req.Data)
		if err != nil {
			return nil, err
		}
		if n != bytes {
			return nil, NewProtocolError(fmt.Sprintf("Read only %d bytes of %d bytes of expected data", n, bytes))
		}
		c, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if c != '\r' {
			return nil, NewProtocolError("expected \\r")
		}
		c, err = r.ReadByte()
		if err != nil {
			return nil, err
		}
		if c != '\n' {
			return nil, NewProtocolError("expected \\n")
		}
		return req, nil
	case "delete":
		// delete <key> [noreply]\r\n
		fallthrough
	case "get":
		// get <key>*\r\n
		fallthrough
	case "gets":
		// gets <key>*\r\n
		if len(arr) < 2 {
			return nil, NewProtocolError(fmt.Sprintf("too few params to command %q", arr[0]))
		}
		req := &McRequest{}
		req.Command = arr[0]
		req.Keys = arr[1:]
		return req, nil
	case "incr", "decr":
		// incr <key> <value> [noreply]\r\n
		// decr <key> <value> [noreply]\r\n
		if len(arr) < 3 {
			return nil, NewProtocolError(fmt.Sprintf("too few params to command %q", arr[0]))
		}
		req := &McRequest{}
		req.Command = arr[0]
		req.Key = arr[1]

		req.Value, err = strconv.ParseInt(arr[2], 10, 64)
		if err != nil {
			return nil, NewProtocolError("cannot read value " + err.Error())
		}
		return req, nil
	case "touch":
		// touch <key> <exptime> [noreply]\r\n
	case "version":
		// version\r\n
		return &McRequest{Command: arr[0]}, nil
	case "quit":
		// quit\r\n
		return &McRequest{Command: arr[0]}, nil
	case "stats":
		// stats\r\n
		// TODO stats <args>\r\n
		return &McRequest{Command: arr[0]}, nil
	}
	return nil, NewProtocolError(fmt.Sprintf("unknown command %q", arr[0]))
}
