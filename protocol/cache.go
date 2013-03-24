/*
Local cache if we cant connect to server
*/
package protocol

import (
	"time"
)

type McEntry struct {
	McValue
	Exptime time.Time
}

var cache map[string]*McEntry

func init() {
	cache = make(map[string]*McEntry)
}

func NewMcEntry(key, flags string, exptime int64, data []byte) *McEntry {
	var ex time.Time
	if exptime != 0 {
		ex = time.Unix(exptime, 0)
	}
	return &McEntry{McValue{key, flags, data}, ex}
}

func AddCache(e *McEntry) {
	cache[e.Key] = e

	// free mem?
	//if len(cache) > 10000 {
	//	
	//}
}

func GetCache(key string) *McEntry {
	e := cache[key]
	if e != nil && !e.Exptime.IsZero() && e.Exptime.Before(time.Now()) {
		//		delete(cache, key)
		//		return nil
	}
	return cache[key]
}
