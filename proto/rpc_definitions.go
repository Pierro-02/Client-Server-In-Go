package proto

import "client_server/shared"

type MatrixRequest struct {
	Operation string
	Mat1      shared.Matrix
	Mat2      shared.Matrix
}

type MatrixResponse struct {
	Result shared.Matrix
	Error  string
}

type WorkerService interface {
	PerformOperation(req MatrixRequest, res *MatrixResponse) error
}
