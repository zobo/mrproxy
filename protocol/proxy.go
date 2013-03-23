package protocol

type ProtocolProxy interface {
	Process(*McRequest) McResponse
}
