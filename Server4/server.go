package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	file_management "server4/helper"
	"server4/userpb"

	"google.golang.org/grpc"
)

type FileUploadServer struct {
	userpb.UnimplementedFileUploadServiceServer
}

func (s *FileUploadServer) Upload(stream userpb.FileUploadService_UploadServer) error {
	var shards []byte
	var file_name string
	var chunknumer int
	var Username string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		file_name = req.Filename
		shards = append(shards, req.Chunk...)
		chunknumer = int(req.ChunkPart)
		Username = req.Username
	}
	file_management.Save_Chunk(shards, file_name, chunknumer, Username)
	return stream.SendAndClose(&userpb.UploadResponse{
		Message:    "Combined successfully",
		TotalBytes: int64(len(shards)),
	})
}
func (s *FileUploadServer) Heartbeat(context.Context, *userpb.Send_HeartBeat_Request) (*userpb.Reply_HeartBeat, error) {
	return &userpb.Reply_HeartBeat{
		Status: true,
	}, nil

}
func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:3004")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterFileUploadServiceServer(grpcServer, &FileUploadServer{})

	fmt.Println("gRPC server running on :3004")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
		return
	}
	return
}
