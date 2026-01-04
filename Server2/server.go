package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	file_management "server2/helper"
	"server2/userpb"

	"google.golang.org/grpc"
)

type FileUploadServer struct {
	userpb.UnimplementedFileUploadServiceServer
}

func (s *FileUploadServer) Upload(stream userpb.FileUploadService_UploadServer) error {
	var shards []byte
	var file_name string
	var chunknumber int
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
		chunknumber = int(req.ChunkPart)
		Username = req.Username
	}
	file_management.Save_Chunk(shards, file_name, chunknumber, Username)
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
func (s *FileUploadServer) Download(req *userpb.DownloadRequest, stream userpb.FileUploadService_DownloadServer) error {
	fmt.Println("Hi")

	chunkSize := 16 * 1024
	name := req.Name
	file_name := req.Filename
	chunknum := req.ChunkNumber
	shard := file_management.Get_Chunk(name, file_name, int(chunknum))
	for offset := 0; offset < len(shard); offset += chunkSize {
		end := offset + chunkSize
		if end > len(shard) {
			end = len(shard)
		}
		reply := &userpb.DownloadResponse{
			Filemane: file_name,
			Chunk:    shard[offset:end],
		}
		err := stream.Send(reply)
		if err != nil {
			return err
		}
	}
	return nil

}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:3002")
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	userpb.RegisterFileUploadServiceServer(grpcServer, &FileUploadServer{})

	fmt.Println("gRPC server running on :3002")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
		return
	}
	return
}
