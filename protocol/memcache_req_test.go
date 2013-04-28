package protocol

import (
	"bufio"
	"io"
	"reflect"
	"testing"
)

func testReq(in string, t *testing.T) (ret *McRequest, err error) {
	r, w := io.Pipe()
	defer r.Close()

	br := bufio.NewReader(r)

	go func() {
		w.Write([]byte(in))
		w.Close()
	}()

	return ReadRequest(br)
}

func TestSet(t *testing.T) {
	ret, err := testReq("set KEY 0 0 10\r\n1234567890\r\n", t)
	if err != nil {
		t.Fatalf("ReadRequest %+v", err)
	}

	if ret.Command != "set" {
		t.Errorf("Command %s", ret.Command)
	}
	if ret.Key != "KEY" {
		t.Errorf("Key %s", ret.Key)
	}
	if ret.Flags != "0" {
		t.Errorf("Flags %s", ret.Flags)
	}
	if ret.Exptime != 0 {
		t.Errorf("Exptime %s", ret.Exptime)
	}
	if string(ret.Data) != "1234567890" {
		t.Errorf("Data %s", ret.Data)
	}

	// Out &{Command:set Key:KEY Keys:[] Flags:0 Exptime:0 Data:[49 50 51 52 53 54 55 56 57 48] Noreply:false}Written 28
}

func TestGet(t *testing.T) {
	ret, err := testReq("get a bb c\r\n", t)
	if err != nil {
		t.Fatalf("ReadRequest %+v", err)
	}

	if ret.Command != "get" {
		t.Errorf("Command %s", ret.Command)
	}
	if !reflect.DeepEqual(ret.Keys, []string{"a", "bb", "c"}) {
		t.Errorf("Keys %v", ret.Keys)
	}
}

func TestCas(t *testing.T) {
	ret, err := testReq("cas KEY 0 0 10 UNIQ\r\n1234567890\r\n", t)
	if err != nil {
		t.Fatalf("ReadRequest %+v", err)
	}

	if ret.Command != "cas" {
		t.Errorf("Command %s", ret.Command)
	}
	if ret.Key != "KEY" {
		t.Errorf("Key %s", ret.Key)
	}
	if ret.Flags != "0" {
		t.Errorf("Flags %s", ret.Flags)
	}
	if ret.Exptime != 0 {
		t.Errorf("Exptime %s", ret.Exptime)
	}
	if ret.Cas != "UNIQ" {
		t.Errorf("Cas %s", ret.Exptime)
	}
	if string(ret.Data) != "1234567890" {
		t.Errorf("Data %s", ret.Data)
	}
}

func TestProtocolError(t *testing.T) {
	_, err := testReq("xxx KEY 0 0 10\r\n1234567890\r\n", t)
	if perr, ok := err.(ProtocolError); ok {
		t.Logf("Good error: %v", perr)
		return
	}
	t.Fatalf("ReadRequest did not return error")
}
