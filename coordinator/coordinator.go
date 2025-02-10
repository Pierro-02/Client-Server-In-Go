package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"sync"

	"client_server/proto"
)

type Coordinator struct {
	mu      sync.Mutex
	workers map[string]int
}

func (c *Coordinator) RegisterWorker(workerAddr string, reply *string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.workers[workerAddr]; !exists {
		c.workers[workerAddr] = 0
		log.Println("Worker registered:", workerAddr)
	}

	*reply = "Worker registered successfully"
	return nil
}

func (c *Coordinator) getLeastBusyWorker() (string, error) {
	if len(c.workers) == 0 {
		return "", errors.New("no available workers")
	}

	var selectedWorker string
	minLoad := int(^uint(0) >> 1) //Max int value

	for worker, load := range c.workers {
		if load < minLoad {
			selectedWorker = worker
			minLoad = load
		}
	}

	return selectedWorker, nil
}

func (c *Coordinator) RequestComputation(req proto.MatrixRequest, res *proto.MatrixResponse) error {
	c.mu.Lock()
	workerAddr, err := c.getLeastBusyWorker()
	if err != nil {
		c.mu.Unlock()
		return err
	}

	c.workers[workerAddr]++
	log.Println("Request Successfully forwarded to worker:", workerAddr)
	c.mu.Unlock()

	client, err := rpc.Dial("tcp", workerAddr)
	if err != nil {
		c.handleWorkerFailure(workerAddr)
		return errors.New("worker unavailable, please try again")
	}
	defer client.Close()

	err = client.Call("Worker.PerformOperation", req, res)

	c.mu.Lock()
	c.workers[workerAddr]--
	c.mu.Unlock()

	if err != nil {
		return errors.New("computation error: " + err.Error())
	}

	return nil
}

func (c *Coordinator) handleWorkerFailure(workerAddr string) {
	c.mu.Lock()
	delete(c.workers, workerAddr)
	log.Println("Removed failed worker:", workerAddr)
	c.mu.Unlock()
}

func main() {
	coordinator := &Coordinator{workers: make(map[string]int)}
	rpc.Register(coordinator)

	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatal("Coordinator error:", err)
	}
	defer listener.Close()

	log.Print("Coordinator running on port 5000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
