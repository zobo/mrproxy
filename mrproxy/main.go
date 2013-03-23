package main

import (
	"bufio"
	"flag"
	"github.com/garyburd/redigo/redis"
	"github.com/zobo/mrproxy/protocol"
	"log"
	"net"
	"time"
)

const listenAddr = "0.0.0.0:11211"

// "10.13.37.106:6379"
var redis_server = flag.String("server", "127.0.0.1:6379", "Redis server to connect to")

func main() {

	flag.Parse()

	// move to global??
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			//c, err := redis.Dial("tcp", *redis_server)
			d, _ := time.ParseDuration("1s")
			c, err := redis.DialTimeout("tcp", *redis_server, d, d, d)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	l, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go processMc(c, pool)
	}
}

func processMc(c net.Conn, pool *redis.Pool) {
	defer log.Printf("%v end processMc", c)
	defer c.Close()

	// process
	br := bufio.NewReader(c)
	bw := bufio.NewWriter(c)

	// take it per need
	conn := pool.Get()
	defer conn.Close()

	redisProxy := protocol.NewRedisProxy(conn)

	for {
		req, err := protocol.ReadRequest(br)
		if perr, ok := err.(protocol.ProtocolError); ok {
			log.Printf("%v ReadRequest protocol err: %v", c, err)
			bw.WriteString("CLIENT_ERROR " + perr.Error() + "\r\n")
			bw.Flush()
			continue
		} else if err != nil {
			log.Printf("%v ReadRequest err: %v", c, err)
			return
		}
		log.Printf("%v Req: %+v\n", c, req)

		switch req.Command {
		case "quit":
			return
		default:
			res := redisProxy.Process(req)
			if !req.Noreply {
				log.Printf("%v Res: %+v\n", c, res)
				bw.WriteString(res.Protocol())
				bw.Flush()
			}
			/*
					case "get":
						for i:= range req.Keys {
							r, err := redis.Values(conn.Do("MGET", req.Keys[i], req.Keys[i]+"_mcflags"))
							if err != nil {
								// check for error type.
								e := protocol.GetCache(req.Keys[i])
			log.Println(e)
								if e != nil {
									bw.WriteString(fmt.Sprintf("VALUE %s %s %d\r\n", e.Key, e.Flags, len(e.Data)))
									bw.Write(e.Data)
									bw.WriteString("\r\n")
									continue
								}
								bw.WriteString("SERVER_ERROR "+err.Error()+"\r\n")
								bw.Flush()
								log.Println(err)
								continue;
							}
							if r[0] != nil {
								data, err := redis.Bytes(r[0], nil)
								flags, err := redis.String(r[1], err)
								if err != nil {
									bw.WriteString("SERVER_ERROR "+err.Error()+"\r\n")
									bw.Flush()
									log.Println(err)
									continue;
								}
								// todo, both can return error
								bw.WriteString(fmt.Sprintf("VALUE %s %s %d\r\n", req.Keys[i], flags, len(data)))
								bw.Write(data)
								bw.WriteString("\r\n")
							}
						}
						bw.WriteString("END\r\n");
						bw.Flush()
					case "set":
						protocol.AddCache(protocol.NewMcEntry(req.Key, req.Flags, req.Exptime, req.Data))

						r, err := redis.String(conn.Do("MSET", req.Key, req.Data, req.Key+"_mcflags", req.Flags))
						if err != nil || r != "OK" {
							bw.WriteString("SERVER_ERROR "+r+"\r\n")
							bw.Flush()
							log.Println(err)
							continue;
						}

						if req.Exptime != 0 {
							_, err = conn.Do("EXPIREAT", req.Key, req.Exptime)
							if err != nil {
								bw.WriteString("SERVER_ERROR storing expire time\r\n")
								bw.Flush()
								log.Println(err)
								continue;
							}
						}

						// err
						bw.WriteString("STORED\r\n")
						bw.Flush()
			*/
		}
	}
}
