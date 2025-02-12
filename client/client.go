package main

import (
	"client_server/proto"
	"client_server/shared"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

// Address of the coordinator / where coordinator is listeneing
const coordinatorAddr = "172.16.60.92:5000"

func main() {

	// client making a rpc connection with the coordinator
	client, err := rpc.Dial("tcp", coordinatorAddr)
	if err != nil {
		log.Fatal("Client couldn't establish a connection with the coordinator:", err)
	}
	// when the main function ends, close the connection
	defer client.Close()

	// Making 3 matrices
	mat1 := shared.Matrix{Rows: 2, Cols: 2, Data: [][]int{{2, 2}, {3, 4}}}
	mat2 := shared.Matrix{Rows: 2, Cols: 2, Data: [][]int{{4, 6}, {1, 2}}}
	mat3 := shared.Matrix{Rows: 2, Cols: 3, Data: [][]int{{2, 3, 5}, {1, 4, 6}}}

	// Performing Add Operation
	req1 := proto.MatrixRequest{
		Operation: "add",
		Mat1:      mat1,
		Mat2:      mat2,
	}

	// Performing Multiply Operation
	req2 := proto.MatrixRequest{
		Operation: "multiply",
		Mat1:      mat1,
		Mat2:      mat2,
	}

	// Performing Transpose Operation
	req3 := proto.MatrixRequest{
		Operation: "transpose",
		Mat1:      mat3,
	}

	// Variables for the response
	var res1, res2, res3 proto.MatrixResponse

	// Creating a re-usable request function, this is used to send requests to the coordinator
	request := func(req proto.MatrixRequest, res *proto.MatrixResponse) {
		err = client.Call("Coordinator.RequestComputation", req, res)
		if err != nil {
			log.Fatal("RPC error with request 1:", err)
		}

		fmt.Println("Computed Answer:")
		res.Result.Print()
	}

	// Sending 3 requests N times to the coordinator in the form of go routines
	// go rouines are basically light weight threads in go
	for i := 0; i < 10; i++ {
		go request(req1, &res1)
		go request(req2, &res2)
		go request(req3, &res3)
	}

	// waiting for all the go routines to finish before finishing execution of main
	time.Sleep(5 * time.Second)
}
