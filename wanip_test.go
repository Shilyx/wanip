package wanip

import (
	"testing"
)

func TestGetViaDNS(t *testing.T) {
	addr8 := getViaDNS("8.8.8.8")
	addr114 := getViaDNS("114.114.114.114")
	addr220 := getViaDNS("208.67.220.220")
	addr222 := getViaDNS("208.67.222.222")

	if len(addr8) == 0 && len(addr114) == 0 && len(addr220) > 0 && addr220 == addr222 {
		return
	}

	t.Error(addr220)
}

func TestGet(t *testing.T) {
	var addr string
	for i := 0; i < 10; i++ {
		s := Get()

		if addr == "" {
			addr = s
			continue
		}

		if addr != s {
			t.Error(s)
		}
	}
}
