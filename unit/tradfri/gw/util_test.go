package gateway

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////

func Test_Util_000(t *testing.T) {
	t.Log("Test_Util_000")
}

func Test_Util_001(t *testing.T) {
	if boolToUint(false) != 0 {
		t.Error("Unexpected boolToUint value")
	}
	if boolToUint(true) == 0 {
		t.Error("Unexpected boolToUint value")
	}
}

func Test_Util_002(t *testing.T) {
	if durationToTransition(0) != 0 {
		t.Error("Unexpected durationToTransition value")
	}
	if durationToTransition(time.Second) != 1000.0/100.0 {
		t.Error("Unexpected durationToTransition value")
	}
	if durationToTransition(time.Millisecond) != 1.0/100.0 {
		t.Error("Unexpected durationToTransition value")
	}
}

func Test_Util_003(t *testing.T) {
	if _, err := addrForService(gopi.RPCServiceRecord{}, 0); err == nil {
		t.Error("addrForService: expected error return")
	}
	if _, err := addrForService(gopi.RPCServiceRecord{
		Addrs: []net.IP{},
	}, 0); err == nil {
		t.Error("addrForService: expected error return")
	}
	if addr, err := addrForService(gopi.RPCServiceRecord{
		Addrs: []net.IP{
			net.ParseIP("127.0.0.1"),
		},
	}, 0); err != nil {
		t.Error("addrForService:", err)
	} else {
		t.Log(addr)
	}
	if addr, err := addrForService(gopi.RPCServiceRecord{
		Addrs: []net.IP{
			net.ParseIP("::1"),
			net.ParseIP("127.0.0.1"),
		},
	}, gopi.RPC_FLAG_INET_V6); err != nil {
		t.Error("addrForService:", err)
	} else if addr != fmt.Sprintf("[::1]:%v", DEFAULT_PORT) {
		t.Error("Expected", fmt.Sprintf("[::1]:%v", DEFAULT_PORT), "got", addr)
	} else {
		t.Log(addr)
	}
	if addr, err := addrForService(gopi.RPCServiceRecord{
		Addrs: []net.IP{
			net.ParseIP("::1"),
			net.ParseIP("127.0.0.1"),
		},
	}, gopi.RPC_FLAG_INET_V4); err != nil {
		t.Error("addrForService:", err)
	} else if addr != fmt.Sprintf("127.0.0.1:%v", DEFAULT_PORT) {
		t.Error("Expected", fmt.Sprintf("127.0.0.1:%v", DEFAULT_PORT), "got", addr)
	}
}
