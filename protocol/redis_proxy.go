package protocol

import (
	"github.com/garyburd/redigo/redis"
)

type RedisProxy struct {
	conn redis.Conn
}

func NewRedisProxy(conn redis.Conn) *RedisProxy {
	r := new(RedisProxy)
	r.conn = conn
	return r
}

func serverError(err error) McResponse {
	return McResponse{Response: "SERVER_ERROR " + err.Error()}
}

// process a request and generate a resonse
func (p *RedisProxy) Process(req *McRequest) McResponse {

	switch req.Command {
	case "get":
		res := McResponse{}
		for i := range req.Keys {

			r, err := redis.Values(p.conn.Do("MGET", req.Keys[i], req.Keys[i]+"_mcflags"))
			if err != nil {
				// hmm, barf errors, or just ignore?
				return serverError(err)
			}
			if r[0] != nil {
				data, err := redis.Bytes(r[0], nil)
				flags, err := redis.String(r[1], err)
				if err != nil {
					return serverError(err)
				}
				// todo, both can return error
				res.Values = append(res.Values, McValue{req.Keys[i], flags, data})
			}
		}
		res.Response = "END"
		return res

	case "set":
		r, err := redis.String(p.conn.Do("MSET", req.Key, req.Data, req.Key+"_mcflags", req.Flags))
		if err != nil || r != "OK" {
			return serverError(err)
		}

		if req.Exptime != 0 {
			_, err = p.conn.Do("EXPIREAT", req.Key, req.Exptime)
			if err != nil {
				return serverError(err)
			}
		}

		return McResponse{Response: "STORED"}

	case "delete":
		r, err := redis.Int(p.conn.Do("DEL", toInterface(req.Keys)...))
		if err != nil {
			return serverError(err)
		}
		if r>0 {
			return McResponse{Response: "DELETED"}
		}
		return McResponse{Response: "NOT_FOUND"}

	// todo "touch"...
	}

	return McResponse{Response: "ERROR"}

}

func toInterface(s []string) []interface{} {

ret := make([]interface{}, len(s))
for i,v:= range s {
 ret[i] = interface{}(v)
}
return ret
}
