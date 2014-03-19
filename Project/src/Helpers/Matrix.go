package Helpers

import()

type Matrix struct {
	rows int
	cols int
	mat [][]int
}

func New_Matrix(rows int, cols int) Matrix{
	m := new(Matrix)
	m.rows = rows
	m.cols = cols

	m.mat = make([][]int,rows)
	for i:=0;i<rows;i++{
		m.mat[i] = make([]int,cols)
	}
	
	return *m
}

func (m *Matrix) Set(row int, col int, value int){
	m.mat[row][col] = value
}

func (m *Matrix) Get(row int, col int) int{
	return m.mat[row][col]
}
