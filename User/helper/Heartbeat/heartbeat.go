package heartbeat

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
	"user/userpb"

	"google.golang.org/grpc"
)

func Beat() {
	for {
		Heartbeat()
		time.Sleep(3 * time.Second)
	}
}

func Heartbeat() map[string]bool {
	var wg sync.WaitGroup
	var mt sync.Mutex
	m := make(map[string]bool)
	for i := 3001; i < 3006; i++ {
		s := strconv.Itoa(i)
		wg.Add(1)
		go func(s string, wg *sync.WaitGroup, mt *sync.Mutex) {
			defer wg.Done()
			status := Call_Server(s)
			mt.Lock()
			m[s] = status
			mt.Unlock()
		}(s, &wg, &mt)
	}
	wg.Wait()
	Display(m)
	fmt.Println("\n")
	return m
}

func Call_Server(s string) bool {
	host := fmt.Sprintf("localhost:%s", s)
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		log.Println(err)
		return false
	}
	defer conn.Close()
	client := userpb.NewFileUploadServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &userpb.Send_HeartBeat_Request{
		Msg: "checking server health",
	}
	rsp, err := client.Heartbeat(ctx, req)
	if err != nil {
		log.Println(err)
		return false
	}
	status := rsp.Status
	if status {
		return true
	} else {
		return false
	}
}

func Display(m map[string]bool) {
	for key, value := range m {
		if value {
			fmt.Println("The server %s is active", key)
		} else {
			fmt.Println("The server %s is inactive", key)
		}
	}
}
