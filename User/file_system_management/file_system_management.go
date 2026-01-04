package filesystemmanagement

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
	"user/userpb"

	"google.golang.org/grpc"
)

var chunkSize int

func Data_sender_init(shards [][]byte, filename string, alloted_servers []string, name string) bool {
	var wg sync.WaitGroup
	var mt sync.Mutex
	stat := true

	for i := 1; i < 6; i++ {
		wg.Add(1)
		s := alloted_servers[i-1]
		go func(s string, idx int, name string) {
			defer wg.Done()

			status := Send_file(s, shards[idx-1], filename, idx, name)

			mt.Lock()
			if status != nil {
				fmt.Println("Error sending to", s)
				stat = false
			}
			mt.Unlock()

		}(s, i, name)
	}

	wg.Wait()
	return stat
}

func Send_file(port string, shard []byte, filename string, part int, name string) error {
	chunkSize = 16 * 1024
	host := fmt.Sprintf("localhost:%s", port)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer conn.Close()
	client := userpb.NewFileUploadServiceClient(conn)
	stream, err := client.Upload(context.Background())
	if err != nil {
		log.Fatal(err)
		return err
	}
	for offset := 0; offset < len(shard); offset += chunkSize {
		end := offset + chunkSize
		if end > len(shard) {
			end = len(shard)
		}

		req := &userpb.UploadRequest{
			Filename:  filename,
			Chunk:     shard[offset:end],
			ChunkPart: int64(part) - 1,
			Username:  name,
		}

		if err := stream.Send(req); err != nil {
			log.Fatalf("Failed to send chunk: %v", err)
			return err
		}
	}
	resp, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		log.Fatalf("Failed to receive response: %v", err)
		return err
	}
	fmt.Println(resp.Message)
	fmt.Println(resp.TotalBytes)
	return nil
}

func Receive_file_init(filename string, name string, loc_detail map[string]string) [][]byte {
	var wg sync.WaitGroup
	var mt sync.Mutex
	shards := make([][]byte, 5)
	for key, value := range loc_detail {
		num, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		port := value
		wg.Add(1)
		go func(wg *sync.WaitGroup, mt *sync.Mutex, port string, filename string, name string, num int) {
			defer wg.Done()
			shard := Receive_file(port, num, filename, name)
			mt.Lock()
			shards[num] = shard
			mt.Unlock()
		}(&wg, &mt, port, filename, name, num)
	}
	wg.Wait()
	return shards
}

func Receive_file(port string, chunknum int, filename string, name string) []byte {
	var shards []byte
	host := fmt.Sprintf("localhost:%s", port)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer conn.Close()
	client := userpb.NewFileUploadServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	stream, err := client.Download(ctx, &userpb.DownloadRequest{
		Name:        name,
		Filename:    filename,
		ChunkNumber: int64(chunknum),
	})
	if err != nil {
		log.Println(err)
		return nil
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break // download complete
		}
		if err != nil {
			log.Println(err)
			return nil
		}
		shard := resp.Chunk
		shards = append(shards, shard...)
	}
	for i := 0; i < len(shards); i++ {
		fmt.Println(shards[i])
	}
	return shards
}
