package xhyve

import (
	"testing"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
	"github.com/zchee/docker-machine-driver-xhyve/vmnet"
)

func TestDriverName(t *testing.T) {
	driverName := newTestDriver("default").DriverName()

	assert.Equal(t, "xhyve", driverName)
}

func TestDefaultSSHUsername(t *testing.T) {
	username := newTestDriver("default").GetSSHUsername()

	assert.Equal(t, "docker", username)
}

func TestPreCreateCheck(t *testing.T) {
	err := newTestDriver("default").PreCreateCheck()
	assert.NoError(t, err)
}

func TestGenerateUUID(t *testing.T) {
	uuid := uuidgen()
	t.Logf(uuid)
}

func TestConvertUUIDToMac(t *testing.T) {
	uuid := uuidgen()
	mac, _ := vmnet.GetMACAddressByUUID(uuid)
	t.Logf(uuid)
	t.Logf(mac)
}

func TestTrimMacAddress(t *testing.T) {
	// test MAC address 02:f0:0d:60:0f:30 and reverse
	testMacAddress := trimMacAddress("02:f0:0d:60:01:03")
	newMacAddress := "2:f0:d:60:1:03"

	if !assert.Equal(t, testMacAddress, newMacAddress) {
		t.Fatalf("expected different MacAddress \n  source %s\nreceived %s", testMacAddress, newMacAddress)
	}

	if !assert.Equal(t, reverse(testMacAddress), reverse(newMacAddress)) {
		t.Fatalf("expected different MacAddress \n  source %s\nreceived %s", testMacAddress, newMacAddress)
	}
}

func newTestDriver(name string) *Driver {
	return NewDriver(name, "")
}

func TestSetConfigFromFlags(t *testing.T) {
	driver := NewDriver("default", "path")

	checkFlags := &drivers.CheckDriverOptions{
		FlagsValues: map[string]interface{}{},
		CreateFlags: driver.GetCreateFlags(),
	}

	err := driver.SetConfigFromFlags(checkFlags)

	assert.NoError(t, err)
	assert.Empty(t, checkFlags.InvalidFlags)
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
