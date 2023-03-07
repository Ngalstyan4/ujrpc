package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	pb "grpc_go/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcServer struct {
	pb.UnimplementedSumServiceServer
}

func (s *GrpcServer) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumResponse, error) {
	return &pb.SumResponse{Result: in.A + in.B}, nil
}

func Server() {
	s := grpc.NewServer()
	pb.RegisterSumServiceServer(s, &GrpcServer{})
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic("unable to create a listener")
	}
	s.Serve(l)
}

func conn() pb.SumServiceClient {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(fmt.Sprintf("unable to create a connection %v\n", err))
	}
	return pb.NewSumServiceClient(conn)
}

func bench(c pb.SumServiceClient, n int, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	for i := 0; i < n; i++ {
		_, err := c.Sum(context.Background(), &pb.SumRequest{A: 1, B: 2})
		if err != nil {
			panic("unable to send a request")
		}
	}

}

func Client(nreqs, nconns, nclients int) {
	var wg sync.WaitGroup
	t0 := time.Now()
	for nc := 0; nc < nconns; nc++ {
		c := conn()
		for i := 0; i < nclients; i++ {
			wg.Add(1)
			go bench(c, nreqs, &wg)
		}
	}
	wg.Wait()
	total := time.Since(t0)
	n := nclients * nconns * nreqs
	fmt.Printf("conns:%3d clients/conn:%3d clients/conn total reqs:%-7d total time:%-5.2fs reqs/s:%-9.2f",
		nconns, nclients, n, total.Seconds(), float64(n)/total.Seconds())
	if nconns == 1 && nclients == 1 {
		fmt.Printf(" lat:%-6.2f us\n", float64(total.Microseconds())/float64(n))
	} else {
		fmt.Println()
	}
}

func main() {
	server := flag.Bool("server", false, "run as server")
	nconns := flag.Int("conns", 32, "number of connections")
	nclients := flag.Int("clients", 10, "number of clients per connection")
	nreqs := flag.Int("reqs", 100_000, "number of requests per client per connection")
	flag.Parse()
	if *server {
		Server()
	} else {
		Client(*nreqs, *nconns, *nclients)
	}
}
