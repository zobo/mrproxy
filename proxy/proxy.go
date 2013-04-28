package proxy

import (
	"github.com/zobo/mrproxy/protocol"
)

type ProtocolProxy interface {
	Process(*protocol.McRequest) protocol.McResponse
}
