package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/andaru/gnotch"
	pb "github.com/andaru/gnotch/proto"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// loadTLS
func loadTLS() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}
	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}

// runAgentAction returns a urfave/cli Action function for starting the agent
func runAgentAction(cfg *gnotch.AgentConfig) func(*cli.Context) error {
	return func(ctx *cli.Context) error {
		log.Printf("starting gRPC agent\n")
		log.Printf("endpoint address %s\n", cfg.AgentEndpoint)
		lis, err := net.Listen("tcp", cfg.AgentEndpoint)
		if err != nil {
			return err
		}

		// creds, err := loadTLS()
		// if err != nil {
		// 	return err
		// }
		agent := gnotch.NewAgent(*cfg)
		server := grpc.NewServer(
		// grpc.Creds(creds),
		)
		pb.RegisterGnotchServer(server, agent)
		reflection.Register(server)
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
				Name:        "grpc-endpoint",
				Usage:       "gRPC server endpoint in host:port format",
				Value:       ":42069",
				Destination: &cfg.AgentEndpoint,
				EnvVars:     []string{"GRPC_ENDPOINT"},
			},
			&cli.BoolFlag{
				Name:        "verbose",
				Aliases:     []string{"v"},
				Usage:       "emit agent logs to stdout",
				Value:       false,
				Destination: &cfg.Verbose,
				EnvVars:     []string{"VERBOSE"},
			},
		},
		Action: runAgentAction(cfg),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
