package eos

import (
	"encoding/json"

	"github.com/ybbus/jsonrpc"
)

// Device is an Arista EOS device.
//
// This type communicates with the EOS device via Arista's JSON/RPC interface
type Device struct {
	url    string
	client jsonrpc.RPCClient
}

// New returns a new EOS device object whose EAPI (JSON-RPC) interface is available at url
func New(url string) *Device { return &Device{url: url} }

// Command executes a command on the device using EAPI, returning the JSON results or an error
func (d *Device) Command(command string) ([]byte, error) {
	return runCmds(d, []string{command}, "json", false)
}

// Connected returns true if a JSON-RPC connection is open to the EAPI service
func (d *Device) Connected() bool { return d.client != nil }

// Close closes the EAPI connection if one is underway, implementing io.Closer
func (d *Device) Close() error {
	close(d)
	return nil
}

func runCmds(d *Device, commands []string, format string, ts bool) ([]byte, error) {
	connect(d)
	rpcResponse, err := d.client.Call("runCmds", map[string]interface{}{
		"version":    "1",
		"cmds":       commands,
		"format":     format,
		"timestamps": ts,
	})
	if err != nil {
		close(d)
		return nil, err
	}
	return json.Marshal(rpcResponse.Result)
}

func connect(d *Device) {
	if d.client == nil {
		d.client = jsonrpc.NewClient(d.url)
	}
}

func close(d *Device) { d.client = nil }
