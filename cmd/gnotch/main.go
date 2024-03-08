package main

import (
	"log"
	"net"
	"os"

	"github.com/andaru/gnotch"
	pb "github.com/andaru/gnotch/proto"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// runAgentAction returns a urfave/cli Action function for starting the agent
func runAgentAction(cfg *gnotch.AgentConfig) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		log.Printf("starting gRPC agent\n")
		log.Printf("endpoint address %s\n", cfg.AgentEndpoint)
		lis, err := net.Listen("tcp", cfg.AgentEndpoint)
		if err != nil {
			return err
		}
		agent := gnotch.NewAgent(*cfg)
		server := grpc.NewServer()
		pb.RegisterGnotchServer(server, agent)
		server.Serve(lis)
		return nil
	}
}

func main() {
	cfg := &gnotch.AgentConfig{}
	app := &cli.App{
		Name:  "gnotch",
		Usage: "The Generic Network Operator's Command Helper gRPC agent",
		Description: `Gnotch provides an interface to network equipment admin interfaces.

This gRPC service is intended to be called by CLI aggregators or
gRPC proxies (for low-latency access across very large networks).`,

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				Usage:       "gRPC server endpoint in host:port format",
				Value:       ":42069",
				Destination: &cfg.AgentEndpoint,
			},
		},
		Action: runAgentAction(cfg),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
