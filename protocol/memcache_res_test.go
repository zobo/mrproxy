package protocol

import (
	"testing"
)

func TestRespEmpty(t *testing.T) {
	var res McResponse

	r := res.Protocol()

	if r != "\r\n" {
		t.Errorf("Empty response is not empty: %v", r)
	}
}

func TestResp1(t *testing.T) {
	res := McResponse{Response: "END"}
	r := res.Protocol()

	if r != "END\r\n" {
		t.Errorf("%v", r)
	}
}

func TestResp2(t *testing.T) {
	res := McResponse{
		"END",
		[]McValue{
			McValue{"k1", "f1", []byte("123")},
		},
	}
	r := res.Protocol()

	if r != "VALUE k1 f1 3\r\n123\r\nEND\r\n" {
		t.Errorf("%v", r)
	}
}

func TestResp3(t *testing.T) {
	res := McResponse{
		"END",
		[]McValue{
			McValue{"k1", "f1", []byte("123")},
			McValue{"k2", "f2", []byte("456")},
		},
	}
	r := res.Protocol()

	if r != "VALUE k1 f1 3\r\n123\r\nVALUE k2 f2 3\r\n456\r\nEND\r\n" {
		t.Errorf("%v", r)
	}
}
