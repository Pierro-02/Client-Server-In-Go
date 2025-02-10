package main

import (
	"client_server/proto"
	"client_server/shared"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

const coordinatorAddr = "localhost:5000"

func main() {
	client, err := rpc.Dial("tcp", coordinatorAddr)
	if err != nil {
		log.Fatal("Client couldn't establish a connection with the coordinator:", err)
	}
	defer client.Close()

	mat1 := shared.Matrix{Rows: 2, Cols: 2, Data: [][]int{{2, 2}, {3, 4}}}
	mat2 := shared.Matrix{Rows: 2, Cols: 2, Data: [][]int{{4, 6}, {1, 2}}}
	mat3 := shared.Matrix{Rows: 2, Cols: 3, Data: [][]int{{2, 3, 5}, {1, 4, 6}}}

	req1 := proto.MatrixRequest{
		Operation: "add",
		Mat1:      mat1,
		Mat2:      mat2,
	}

	req2 := proto.MatrixRequest{
		Operation: "multiply",
		Mat1:      mat1,
		Mat2:      mat2,
	}

	req3 := proto.MatrixRequest{
		Operation: "transpose",
		Mat1:      mat3,
	}

	var res1, res2, res3 proto.MatrixResponse

	request := func(req proto.MatrixRequest, res *proto.MatrixResponse) {
		err = client.Call("Coordinator.RequestComputation", req, res)
		if err != nil {
			log.Fatal("RPC error with request 1:", err)
		}

		fmt.Println("Computed Answer:")
		res.Result.Print()
	}

	for i := 0; i < 10; i++ {
		go request(req1, &res1)
		go request(req2, &res2)
		go request(req3, &res3)
	}

	time.Sleep(5 * time.Second)
}
