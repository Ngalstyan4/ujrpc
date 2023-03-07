package main

import (
	"fmt"
	"os/exec"
	"testing"
	"time"
)

func TestSumPerformance(t *testing.T) {
	fmt.Println("****** TESTS ON 8-THREADED SERVER ******")
	srv := exec.Command("./grpc_go", "-server")
	srv.Env = append(srv.Env, "GOMAXPROCS=8")
	go srv.Run()
	NREQ := 100_000
	time.Sleep(300 * time.Millisecond)
	Client(NREQ, 1, 1)
	Client(NREQ, 32, 1)
	Client(NREQ/10, 32, 10)
	srv.Process.Kill()
	time.Sleep(900 * time.Millisecond)

	fmt.Println()
	fmt.Println("****** TESTS ON SINGLE-THREADED SERVER ******")
	srv = exec.Command("./grpc_go", "-server")
	srv.Env = append(srv.Env, "GOMAXPROCS=1")
	go srv.Run()
	time.Sleep(300 * time.Millisecond)
	Client(NREQ, 1, 1)
	Client(NREQ, 32, 1)
	Client(NREQ/10, 32, 10)
	srv.Process.Kill()
}
