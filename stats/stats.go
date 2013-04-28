package stats

// TODO move out of package

import (
	"bytes"
	"github.com/zobo/mrproxy/protocol"
	"github.com/zobo/mrproxy/proxy"
	"strconv"
)

// type to process "stats" request and do command counting
type StatsProxy struct {
	next proxy.ProtocolProxy
}

func NewStatsProxy(next proxy.ProtocolProxy) *StatsProxy {
	return &StatsProxy{next: next}
}

var stats = statsData{}

type statsData struct {
	cmd_get           int
	cmd_set           int
	curr_connections  int
	total_connections int
	get_hits          int
	get_misses        int
}

func Connect() {
	stats.curr_connections++
	stats.total_connections++
}

func Disconnect() {
	stats.curr_connections--
}

// Process the memcache request
func (proxy *StatsProxy) Process(req *protocol.McRequest) protocol.McResponse {

	// TODO count bytes in, bytes out
	switch req.Command {
	case "get":
		stats.cmd_get++
		ret := proxy.next.Process(req)
		stats.get_hits += len(ret.Values)
		stats.get_misses += len(req.Keys) - len(ret.Values)
		return ret
	case "set":
		stats.cmd_set++
		return proxy.next.Process(req)
	case "stats":
		var b bytes.Buffer
		b.WriteString("STAT cmd_get ")
		b.WriteString(strconv.Itoa(stats.cmd_get))
		b.WriteString("\r\n")

		b.WriteString("STAT cmd_set ")
		b.WriteString(strconv.Itoa(stats.cmd_set))
		b.WriteString("\r\n")

		b.WriteString("STAT curr_connections ")
		b.WriteString(strconv.Itoa(stats.curr_connections))
		b.WriteString("\r\n")

		b.WriteString("STAT total_connections ")
		b.WriteString(strconv.Itoa(stats.total_connections))
		b.WriteString("\r\n")

		b.WriteString("STAT get_hits ")
		b.WriteString(strconv.Itoa(stats.get_hits))
		b.WriteString("\r\n")

		b.WriteString("STAT get_misses ")
		b.WriteString(strconv.Itoa(stats.get_misses))
		b.WriteString("\r\n")

		b.WriteString("END")

		return protocol.McResponse{Response: b.String()}

	default:
		return proxy.next.Process(req)
	}
}
