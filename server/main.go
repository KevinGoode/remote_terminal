package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"os/exec"
	"strings"

	"../remote_terminal_api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement RemoteTerminalServer
type server struct{}

func (s *server) ExecuteCommand(ctx context.Context, in *remote_terminal_api.CommandRequest) (*remote_terminal_api.CommandResponse, error) {
	command := strings.Split(in.Cmd, " ")
	cmd := exec.Command(command[0], command[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	response := ""
	err := cmd.Run()
	if err != nil {
		response = "rt Error- Failed to execute command on server: " + err.Error()
	} else {
		response = out.String()
	}
	return &remote_terminal_api.CommandResponse{Response: response}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	remote_terminal_api.RegisterRemoteTerminalServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
