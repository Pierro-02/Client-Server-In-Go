# A Basic client-server application.

## The client

- Programs request the coordinator for computation on their data.
- The server acts as a coordinator for n worker processes (n ≥ 3).
- The worker processes are responsible for performing matrix operations, limited to:
  - Addition
  - Multiplication
  - Transpose
- The coordinator is responsible for assigning tasks to the workers and sending the results back to the client.

## Client’s Role

- The client initiates computation requests from the coordinator process.
- The client requests services via RPC (Remote Procedure Calls), ensuring that the client and server programs run on different physical devices

## Coordinator’s Responsibilities

- The coordinator (server) schedules the tasks on a First-Come, First-Served (FCFS) basis.
- The coordinator assigns tasks to the least busy workers first (Load Balancing).
- In case of a worker’s failure, the server assigns the tasks to the next available worker ensuring basic level Fault Tolerance.
- The coordinator gathers the results from workers and sends them back to the client.

## Worker’s Responsibilities

- The worker is responsible for performing matrix operations on the data received from the coordinator process.

## How to run

- First run the coordinator
  - `go run coordinator/coordinator.go`
- Then run as many worker nodes as you want
  - `go run worker/worker.go worker/matrix.go` (To run a single worker)
- Finally run the client
  - `go run client/client.go`
