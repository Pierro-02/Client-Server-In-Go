package shared

import (
	"errors"
	"fmt"
)

type Matrix struct {
	Rows int
	Cols int
	Data [][]int
}

func NewMatrix(rows, cols int) Matrix {
	data := make([][]int, rows)
	for i := range data {
		data[i] = make([]int, cols)
	}
	return Matrix{Rows: rows, Cols: cols, Data: data}
}

func (m *Matrix) Print() {
	for _, row := range m.Data {
		fmt.Println(row)
	}
}

func ValidateMatrices(mat1 Matrix, mat2 Matrix, operation string) error {
	if operation == "add" {
		if mat1.Rows != mat2.Rows || mat1.Cols != mat2.Cols {
			return errors.New("Matrix Addition not possible. Mismatch between dimensions")
		}
	}
	if operation == "multiply" {
		if mat1.Cols != mat2.Rows {
			return errors.New("Matrix Multiplication not possible. Mismatch between dimensions")
		}
	}

	return nil
}
