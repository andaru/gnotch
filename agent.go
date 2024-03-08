package gnotch

import (
	"context"

	"github.com/andaru/gnotch/device"
	gnotch "github.com/andaru/gnotch/proto"
	"github.com/andaru/gnotch/session"
)

// Agent is the Gnotch agent server implementation
type Agent struct {
	sm session.Manager

	gnotch.UnimplementedGnotchServer
}

type AgentConfig struct {
	AgentEndpoint        string
	SessionManageOptions []session.ManagerOption
}

// NewAgent returns a new gnotch agent with a given config
func NewAgent(config AgentConfig) *Agent {
	agent := &Agent{
		sm: *session.NewManager(config.SessionManageOptions...),
	}
	return agent
}

func (a *Agent) Command(ctx context.Context, req *gnotch.CommandRequest) (resp *gnotch.CommandResponse, err error) {
	var session *session.Session
	session, err = a.sm.Session(req.Device)
	if err != nil {
		panic("TODO: confirm error details")
		// return
	}
	commander, ok := session.Device.(device.Commander)
	if !ok {
		panic("TODO: not a commander")
	}

	resp = &gnotch.CommandResponse{}
	resp.Repsonse, err = commander.Command(string(req.Command))
	if err != nil {
		panic("TODO: confirm error details")
	}
	return
}

// static interface checks
var (
	_ gnotch.GnotchServer = &Agent{}
)
