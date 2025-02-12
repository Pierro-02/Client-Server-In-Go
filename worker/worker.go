package main

import (
	"client_server/proto"
	"client_server/shared"
	"fmt"
	"log"
	"net"
	"net/rpc"
)

const coordinatorAddr = "172.16.60.92:5000"

type Worker struct{}

// Adding a function tot the Worker struct. Simply calls the neccessary operation function
func (w *Worker) PerformOperation(req proto.MatrixRequest, res *proto.MatrixResponse) error {
	var result shared.Matrix
	var err error

	switch req.Operation {
	case "add":
		result, err = MatrixAdd(req.Mat1, req.Mat2)
	case "multiply":
		result, err = MatrixMul(req.Mat1, req.Mat2)
	case "transpose":
		result, err = MatrixTranspose(req.Mat1)
	default:
		err = fmt.Errorf("unknown operation %s", req.Operation)
	}

	if err != nil {
		res.Error = err.Error()
		return err
	}

	res.Result = result
	return nil
}

func registerWithCoordinator(workerAddr string) {
	client, err := rpc.Dial("tcp", coordinatorAddr) // establish connection with coordinator
	if err != nil {
		log.Fatal("Failed to establish a connection with the coordinator:", err)
	}
	defer client.Close()

	var reply string
	err = client.Call("Coordinator.RegisterWorker", workerAddr, &reply) // When conn established, call the coordinators function to register it
	if err != nil {
		log.Fatal("Error registering with the coordinator:", err)
	}

	log.Println("Worker registered response:", reply)
}

func main() {
	worker := new(Worker)
	rpc.Register(worker)

	listener, err := net.Listen("tcp", ":0") // listens on a random free port

	if err != nil {
		log.Fatal("Error starting up a worker server:", err)
	}
	defer listener.Close()

	workerAddr := listener.Addr().String() // getting the worker address and port
	registerWithCoordinator(workerAddr)    // register it with the coordinator

	log.Printf("Worker server started on port %s...\n", workerAddr)

	// open connection in an infinite loop to handle all possible requests from the client
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Connection error:", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
