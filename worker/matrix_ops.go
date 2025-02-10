package main

import (
	"client_server/shared"
	"errors"
)

func MatrixAdd(mat1, mat2 shared.Matrix) (shared.Matrix, error) {
	err := shared.ValidateMatrices(mat1, mat2, "Add")
	if err != nil {
		return shared.Matrix{}, err
	}

	sum := shared.NewMatrix(mat1.Rows, mat1.Cols)
	for i := 0; i < mat1.Rows; i++ {
		for j := 0; j < mat1.Cols; j++ {
			sum.Data[i][j] = mat1.Data[i][j] + mat2.Data[i][j]
		}
	}

	return sum, err
}

func MatrixMul(mat1, mat2 shared.Matrix) (shared.Matrix, error) {
	err := shared.ValidateMatrices(mat1, mat2, "Add")
	if err != nil {
		return shared.Matrix{}, err
	}

	result := shared.NewMatrix(mat1.Rows, mat2.Cols)
	for i := range result.Rows {
		for j := range result.Cols {
			for k := 0; k < mat1.Cols; k++ {
				result.Data[i][j] += mat1.Data[i][k] * mat2.Data[k][j]
			}
		}
	}

	return result, err
}

func MatrixTranspose(mat shared.Matrix) (shared.Matrix, error) {
	if mat.Rows == 0 && mat.Cols == 0 {
		return shared.Matrix{}, errors.New("matrix is empty")
	}

	transposed := shared.NewMatrix(mat.Cols, mat.Rows)
	for i := 0; i < mat.Cols; i++ {
		for j := 0; j < mat.Rows; j++ {
			transposed.Data[i][j] = mat.Data[j][i]
		}
	}
	return transposed, nil
}
