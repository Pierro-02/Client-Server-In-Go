package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net/rpc"
	"sync"

	"client_server/proto"
)

// A structure for the coordinator
type Coordinator struct {
	mu      sync.Mutex
	workers map[string]int
}

// Adding a function to the coordinator struct (GPT This function definition if you don't understand)
func (c *Coordinator) RegisterWorker(workerAddr string, reply *string) error {
	c.mu.Lock()         // Simple Mutex Lock
	defer c.mu.Unlock() // "defer" means when function ends.

	// checking if worker with the incoming address exists.
	if _, exists := c.workers[workerAddr]; !exists {
		c.workers[workerAddr] = 0 // If it doesn't, create it with 0 tasks
		log.Println("Worker registered:", workerAddr)
	}

	*reply = "Worker registered successfully"
	return nil
}

func (c *Coordinator) getLeastBusyWorker() (string, error) {
	// Checking if any worker is available
	if len(c.workers) == 0 {
		return "", errors.New("no available workers")
	}

	var selectedWorker string
	minLoad := int(^uint(0) >> 1) //Max int value

	// Looping through workers to find the one with least tasks assigned
	for worker, load := range c.workers {
		if load < minLoad {
			selectedWorker = worker
			minLoad = load
		}
	}

	return selectedWorker, nil
}

// The main coordinator function that assigns the tasks to the workers
func (c *Coordinator) RequestComputation(req proto.MatrixRequest, res *proto.MatrixResponse) error {
	c.mu.Lock()
	workerAddr, err := c.getLeastBusyWorker() // Getting worker that is least busy / has least tasks assigned to it
	if err != nil {
		c.mu.Unlock()
		return err
	}

	// Incrementing its task value as going to assign it a task
	c.workers[workerAddr]++
	log.Println("Request Successfully forwarded to worker:", workerAddr)
	c.mu.Unlock()

	// establishing tcp connection with the worker
	client, err := rpc.Dial("tcp", workerAddr)
	if err != nil {
		c.handleWorkerFailure(workerAddr)
		return errors.New("worker unavailable, please try again")
	}
	defer client.Close() // close conn on function end

	// Calling the Perform Operation function in the worker, sending the response as pointer or as a reference
	err = client.Call("Worker.PerformOperation", req, res)

	c.mu.Lock()
	c.workers[workerAddr]-- // Decrementing task value from worker as task is done at this stage
	c.mu.Unlock()

	if err != nil {
		return errors.New("computation error: " + err.Error())
	}

	return nil
}

// Simple fault tolerance mechanism
func (c *Coordinator) handleWorkerFailure(workerAddr string) {
	// if worker fails, the coordinator does not handle it anymore and notifies that the worker has failed
	c.mu.Lock()
	delete(c.workers, workerAddr)
	log.Println("Removed failed worker:", workerAddr)
	c.mu.Unlock()
}

func main() {
	// Initialising the coordinator struct
	cert, err := tls.LoadX509KeyPair("tls/tls_cert.pem", "tls/tls_key.pem")
	if err != nil {
		log.Fatal("Failed to load TLS certificate:", err)
	}

	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	coordinator := &Coordinator{workers: make(map[string]int)}
	rpc.Register(coordinator) // registering it as rpc

	// Opening connection on port 5000
	listener, err := tls.Listen("tcp", "0.0.0.0:5000", tlsConfig)
	if err != nil {
		log.Fatal("Coordinator error:", err)
	}
	defer listener.Close() // close conn when function ends

	log.Print("Coordinator running with TLS on port 5000...")

	// Infinite function to accept accept connection requests and handle them in go routines
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
