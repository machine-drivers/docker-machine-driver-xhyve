// Originally Author is @ailispaw
// import by https://github.com/ailispaw/xhyvectl/tree/master/vmnet
package vmnet

// #cgo CFLAGS: -framework vmnet
// #cgo LDFLAGS: -framework vmnet
// #include "vmnet.c"
import "C"

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

const (
	CONFIG_PLIST = "/Library/Preferences/SystemConfiguration/com.apple.vmnet"
	NET_ADDR_KEY = "Shared_Net_Address"
	NET_MASK_KEY = "Shared_Net_Mask"
)

func GetMACAddressByUUID(uuid string) (string, error) {
	mac := C.vmnet_get_mac_address_from_uuid(C.CString(uuid))
	if mac == nil {
		return "", fmt.Errorf("Could not get a MAC address for %s", uuid)
	}
	hw, err := net.ParseMAC(C.GoString(mac))
	if err != nil {
		return "", err
	}
	return hw.String(), nil
}

func getNetAddr() (net.IP, error) {
	out, err := exec.Command("defaults", "read", CONFIG_PLIST, NET_ADDR_KEY).Output()
	if err != nil {
		return nil, err
	}
	ip := net.ParseIP(strings.TrimSpace(string(out)))
	if ip == nil {
		return nil, fmt.Errorf("Could not get the network address for vmnet")
	}
	return ip, nil
}

func getNetMask() (net.IPMask, error) {
	out, err := exec.Command("defaults", "read", CONFIG_PLIST, NET_MASK_KEY).Output()
	if err != nil {
		return nil, err
	}
	mask := net.ParseIP(strings.TrimSpace(string(out)))
	if mask == nil {
		return nil, fmt.Errorf("Could not get the network mask for vmnet")
	}
	return net.IPMask(mask.To4()), nil
}

func GetIPNet() (*net.IPNet, error) {
	ip, err := getNetAddr()
	if err != nil {
		return nil, err
	}

	mask, err := getNetMask()
	if err != nil {
		return nil, err
	}

	return &net.IPNet{
		IP:   ip.Mask(mask),
		Mask: mask,
	}, nil
}
