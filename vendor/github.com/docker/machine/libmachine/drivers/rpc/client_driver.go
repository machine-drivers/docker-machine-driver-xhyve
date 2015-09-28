package rpcdriver

import (
	"net/rpc"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
)

type RpcClientDriver struct {
	Client *rpc.Client
}

func NewRpcClientDriver(rawDriverData []byte, serverEndpoint string) (*RpcClientDriver, error) {
	client, err := rpc.DialHTTP("tcp", serverEndpoint)
	if err != nil {
		return nil, err
	}

	c := &RpcClientDriver{
		Client: client,
	}

	if err := c.Client.Call("RpcServerDriver.SetConfigRaw", rawDriverData, nil); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *RpcClientDriver) MarshalJSON() ([]byte, error) {
	return c.GetConfigRaw()
}

func (c *RpcClientDriver) UnmarshalJSON(data []byte) error {
	return c.SetConfigRaw(data)
}

func (c *RpcClientDriver) Close() error {
	return c.Client.Call("RpcServerDriver.Close", struct{}{}, nil)
}

// Helper method to make requests which take no arguments and return simply a
// string, e.g. "GetIP".
func (c *RpcClientDriver) rpcStringCall(method string) (string, error) {
	var info string

	if err := c.Client.Call(method, struct{}{}, &info); err != nil {
		return "", err
	}

	return info, nil
}

func (c *RpcClientDriver) GetCreateFlags() []mcnflag.Flag {
	var flags []mcnflag.Flag

	if err := c.Client.Call("RpcServerDriver.GetCreateFlags", struct{}{}, &flags); err != nil {
		log.Warnf("Error attempting call to get create flags: %s", err)
	}

	return flags
}

func (c *RpcClientDriver) SetConfigRaw(data []byte) error {
	return c.Client.Call("RpcServerDriver.SetConfigRaw", data, nil)
}

func (c *RpcClientDriver) GetConfigRaw() ([]byte, error) {
	var data []byte

	if err := c.Client.Call("RpcServerDriver.GetConfigRaw", struct{}{}, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *RpcClientDriver) DriverName() string {
	driverName, err := c.rpcStringCall("RpcServerDriver.DriverName")
	if err != nil {
		log.Warnf("Error attempting call to get driver name: %s", err)
	}

	return driverName
}

func (c *RpcClientDriver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	return c.Client.Call("RpcServerDriver.SetConfigFromFlags", &flags, nil)
}

func (c *RpcClientDriver) GetURL() (string, error) {
	return c.rpcStringCall("RpcServerDriver.GetURL")
}

func (c *RpcClientDriver) GetMachineName() string {
	name, err := c.rpcStringCall("RpcServerDriver.GetMachineName")
	if err != nil {
		log.Warnf("Error attempting call to get machine name: %s", err)
	}

	return name
}

func (c *RpcClientDriver) GetIP() (string, error) {
	return c.rpcStringCall("RpcServerDriver.GetIP")
}

func (c *RpcClientDriver) GetSSHHostname() (string, error) {
	return c.rpcStringCall("RpcServerDriver.GetSSHHostname")
}

// TODO:  This method doesn't even make sense to have with RPC.
func (c *RpcClientDriver) GetSSHKeyPath() string {
	path, err := c.rpcStringCall("RpcServerDriver.GetSSHKeyPath")
	if err != nil {
		log.Warnf("Error attempting call to get SSH key path: %s", err)
	}

	return path
}

func (c *RpcClientDriver) GetSSHPort() (int, error) {
	var port int

	if err := c.Client.Call("RpcServerDriver.GetSSHPort", struct{}{}, &port); err != nil {
		return 0, err
	}

	return port, nil
}

func (c *RpcClientDriver) GetSSHUsername() string {
	username, err := c.rpcStringCall("RpcServerDriver.GetSSHUsername")
	if err != nil {
		log.Warnf("Error attempting call to get SSH username: %s", err)
	}

	return username
}

func (c *RpcClientDriver) GetState() (state.State, error) {
	var s state.State

	if err := c.Client.Call("RpcServerDriver.GetState", struct{}{}, &s); err != nil {
		return state.Error, err
	}

	return s, nil
}

func (c *RpcClientDriver) PreCreateCheck() error {
	return c.Client.Call("RpcServerDriver.PreCreateCheck", struct{}{}, nil)
}

func (c *RpcClientDriver) Create() error {
	return c.Client.Call("RpcServerDriver.Create", struct{}{}, nil)
}

func (c *RpcClientDriver) Remove() error {
	return c.Client.Call("RpcServerDriver.Remove", struct{}{}, nil)
}

func (c *RpcClientDriver) Start() error {
	return c.Client.Call("RpcServerDriver.Start", struct{}{}, nil)
}

func (c *RpcClientDriver) Stop() error {
	return c.Client.Call("RpcServerDriver.Stop", struct{}{}, nil)
}

func (c *RpcClientDriver) Restart() error {
	return c.Client.Call("RpcServerDriver.Restart", struct{}{}, nil)
}

func (c *RpcClientDriver) Kill() error {
	return c.Client.Call("RpcServerDriver.Kill", struct{}{}, nil)
}

func (c *RpcClientDriver) LocalArtifactPath(file string) string {
	var path string

	if err := c.Client.Call("RpcServerDriver.LocalArtifactPath", file, &path); err != nil {
		log.Warnf("Error attempting call to get LocalArtifactPath: %s", err)
	}

	return path
}

func (c *RpcClientDriver) GlobalArtifactPath() string {
	globalArtifactPath, err := c.rpcStringCall("RpcServerDriver.GlobalArtifactPath")
	if err != nil {
		log.Warnf("Error attempting call to get GlobalArtifactPath: %s", err)
	}

	return globalArtifactPath
}

func (c *RpcClientDriver) Upgrade() error {
	return c.Client.Call("RpcServerDriver.Upgrade", struct{}{}, nil)
}
