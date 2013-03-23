/*
Local cache if we cant connect to server
*/
package protocol

import (
	"testing"
)

func TestCacheAdd(t *testing.T) {
	AddCache(NewMcEntry("test", "0", 1, nil))
}

func TestCacheGet(t *testing.T) {
	r := GetCache("test")
	t.Logf("Got %v", r)
}
