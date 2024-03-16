package gnotch

import (
	"context"
	"log"

	"github.com/andaru/gnotch/device"
	gnotch "github.com/andaru/gnotch/proto"
	"github.com/andaru/gnotch/session"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Agent is the Gnotch agent server implementation
type Agent struct {
	sm session.Manager

	gnotch.UnimplementedGnotchServer
}

type AgentConfig struct {
	AgentEndpoint        string
	SessionManageOptions []session.ManagerOption
	Verbose              bool
}

// NewAgent returns a new gnotch agent with a given config
func NewAgent(config AgentConfig) *Agent {
	agent := &Agent{
		sm: *session.NewManager(config.SessionManageOptions...),
	}
	return agent
}

func (a *Agent) Command(ctx context.Context, req *gnotch.CommandRequest) (resp *gnotch.CommandResponse, err error) {
	defer func() {
		if err != nil {
			log.Printf("debug: %v\n", err)
		}
	}()
	var session *session.Session
	switch session, err = a.sm.Session(req.Device); {
	case err != nil:
	case session == nil:
		err = status.Errorf(codes.NotFound, "session could not be obtained for request %s", req)
	case session.Device == nil:
		err = status.Errorf(codes.NotFound, "device %s not found", req.Device)
	}
	if err != nil {
		log.Printf("err %v\n", err)
		return
	}

	switch commander, ok := session.Device.(device.Commander); ok {
	case false:
		err = status.Errorf(codes.Internal, "bad dev: device model does not implement device.Commander")
	default:
		resp = &gnotch.CommandResponse{}
		if resp.Repsonse, err = commander.Command(string(req.Command)); err != nil {
			log.Printf("err %v\n", err)
		}
	}
	return
}

// static interface checks
var (
	_ gnotch.GnotchServer = &Agent{}
)
