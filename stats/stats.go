package stats

import (
	"bytes"
	"github.com/zobo/mrproxy/protocol"
	"github.com/zobo/mrproxy/proxy"
	"os"
	"strconv"
	"time"
)

// the available operations
type statsOps struct {
	data       chan chan statsData
	closing    chan bool
	connect    chan bool
	disconnect chan bool
	cmd_get    chan bool
	cmd_set    chan bool
}

// the operations
var ops statsOps = statsOps{
	data:       make(chan chan statsData),
	closing:    make(chan bool),
	connect:    make(chan bool, 100),
	disconnect: make(chan bool, 100),
	cmd_get:    make(chan bool, 100),
	cmd_set:    make(chan bool, 100),
}

var stats = statsData{}
var startTime = time.Now()

type statsData struct {
	cmd_get           int
	cmd_set           int
	curr_connections  int
	total_connections int
	get_hits          int
	get_misses        int
}

// for select loop
func loop() {
	for {
		select {
		case <-ops.closing:
			// we should care about errors, but we don't
			close(ops.closing)
			return
		case <-ops.connect:
			stats.curr_connections++
			stats.total_connections++
		case <-ops.disconnect:
			stats.curr_connections--
		case <-ops.cmd_get:
			stats.cmd_get++
		case <-ops.cmd_set:
			stats.cmd_set++
		case datac := <-ops.data:
			datac <- stats
		}
	}
}

// init things, create the stats process
func init() {
	// start the stats loop
	go loop()
}

// type to process "stats" request and do command counting
type StatsProxy struct {
	next proxy.ProtocolProxy
}

func NewStatsProxy(next proxy.ProtocolProxy) *StatsProxy {
	return &StatsProxy{next: next}
}

func Connect() {
	ops.connect <- true
}

func Disconnect() {
	ops.disconnect <- true
}

// Process the memcache request
func (proxy *StatsProxy) Process(req *protocol.McRequest) protocol.McResponse {

	// TODO count bytes in, bytes out
	datac := make(chan statsData)
	ops.data <- datac
	stats := <-datac

	switch req.Command {
	case "get":
		ops.cmd_get <- true
		ret := proxy.next.Process(req)
		stats.get_hits += len(ret.Values)
		stats.get_misses += len(req.Keys) - len(ret.Values)
		return ret
	case "set":
		ops.cmd_set <- true
		return proxy.next.Process(req)
	case "stats":
		var b bytes.Buffer
		b.WriteString("STAT pid ")
		b.WriteString(strconv.Itoa(os.Getpid()))
		b.WriteString("\r\n")

		b.WriteString("STAT uptime ")
		b.WriteString(strconv.Itoa(int(time.Now().Sub(startTime).Seconds())))
		b.WriteString("\r\n")

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
