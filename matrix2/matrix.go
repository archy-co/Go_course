package main

import (
	"errors"
	"fmt"
	"math"
)

type Matrix struct {
	rows, cols int
	data       [][]float64
}

func NewMatrix(rows, cols int) *Matrix {
	data := make([][]float64, rows)
	for i := range data {
		data[i] = make([]float64, cols)
	}
	return &Matrix{rows, cols, data}
}

func (m *Matrix) Add(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, errors.New("розміри матриць не співпадають")
	}
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] + other.data[i][j]
		}
	}
	return result, nil
}

func (m *Matrix) Subtract(other *Matrix) (*Matrix, error) {
	if m.rows != other.rows || m.cols != other.cols {
		return nil, errors.New("розміри матриць не співпадають")
	}
	result := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[i][j] = m.data[i][j] - other.data[i][j]
		}
	}
	return result, nil
}

func (m *Matrix) Multiply(other *Matrix) (*Matrix, error) {
	if m.cols != other.rows {
		return nil, errors.New("кількість стовпців першої матриці має дорівнювати кількості рядків другої")
	}
	result := NewMatrix(m.rows, other.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < other.cols; j++ {
			for k := 0; k < m.cols; k++ {
				result.data[i][j] += m.data[i][k] * other.data[k][j]
			}
		}
	}
	return result, nil
}

func (m *Matrix) Transpose() *Matrix {
	result := NewMatrix(m.cols, m.rows)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			result.data[j][i] = m.data[i][j]
		}
	}
	return result
}

func (m *Matrix) Determinant() (float64, error) {
	if m.rows != m.cols {
		return 0, errors.New("матриця не є квадратною")
	}
	return m.determinantRecursive(m.data), nil
}

func (m *Matrix) determinantRecursive(data [][]float64) float64 {
	if len(data) == 1 {
		return data[0][0]
	}
	det := 0.0
	sign := 1.0
	for j := 0; j < len(data); j++ {
		minor := m.getMinor(data, 0, j)
		det += sign * data[0][j] * m.determinantRecursive(minor)
		sign = -sign
	}
	return det
}

func (m *Matrix) getMinor(data [][]float64, row, col int) [][]float64 {
	minor := make([][]float64, len(data)-1)
	for i := range minor {
		minor[i] = make([]float64, len(data)-1)
	}
	for i := range data {
		if i == row {
			continue
		}
		for j := range data[i] {
			if j == col {
				continue
			}
			minorRow := i
			if i > row {
				minorRow--
			}
			minorCol := j
			if j > col {
				minorCol--
			}
			minor[minorRow][minorCol] = data[i][j]
		}
	}
	return minor
}

func (m *Matrix) Inverse() (*Matrix, error) {
	det, err := m.Determinant()
	if err != nil {
		return nil, err
	}
	if math.Abs(det) < 1e-10 {
		return nil, errors.New("матриця не є оборотньою (визначник = 0)")
	}
	cofactors := NewMatrix(m.rows, m.cols)
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			minor := m.getMinor(m.data, i, j)
			cofactors.data[i][j] = math.Pow(-1, float64(i+j)) * m.determinantRecursive(minor)
		}
	}
	cofactors = cofactors.Transpose()
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			cofactors.data[i][j] /= det
		}
	}
	return cofactors, nil
}

func (m *Matrix) Print() {
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			fmt.Printf("%10.2f ", m.data[i][j])
		}
		fmt.Println()
	}
}

func (m *Matrix) Input() {
	fmt.Println("Введіть елементи матриці:")
	for i := 0; i < m.rows; i++ {
		for j := 0; j < m.cols; j++ {
			fmt.Printf("Елемент [%d][%d]: ", i+1, j+1)
			fmt.Scan(&m.data[i][j])
		}
	}
}

type RowLexicographicSort struct {
	*Matrix
}

func (r *RowLexicographicSort) Sort() {
	for i := 0; i < r.rows-1; i++ {
		for j := i + 1; j < r.rows; j++ {
			if r.data[i][0] > r.data[j][0] {
				r.data[i], r.data[j] = r.data[j], r.data[i]
			}
		}
	}
}

type ColumnLexicographicSort struct {
	*Matrix
}

func (c *ColumnLexicographicSort) Sort() {
	for i := 0; i < c.cols-1; i++ {
		for j := i + 1; j < c.cols; j++ {
			for k := 0; k < c.rows; k++ {
				if c.data[k][i] > c.data[k][j] {
					c.data[k][i], c.data[k][j] = c.data[k][j], c.data[k][i]
				}
			}
		}
	}
}

// Метод Гауса
func (m *Matrix) SolveSystem(b []float64) ([]float64, error) {
	if m.rows != m.cols {
		return nil, errors.New("матриця не є квадратною")
	}
	if len(b) != m.rows {
		return nil, errors.New("довжина вектора вільних членів не відповідає розміру матриці")
	}

	augmentedMatrix := make([][]float64, m.rows)
	for i := 0; i < m.rows; i++ {
		augmentedMatrix[i] = append(m.data[i], b[i])
	}

	for i := 0; i < m.rows; i++ {
		maxRow := i
		for k := i + 1; k < m.rows; k++ {
			if math.Abs(augmentedMatrix[k][i]) > math.Abs(augmentedMatrix[maxRow][i]) {
				maxRow = k
			}
		}

		augmentedMatrix[i], augmentedMatrix[maxRow] = augmentedMatrix[maxRow], augmentedMatrix[i]

		// Перевірка на нульовий елемент на діагоналі
		if math.Abs(augmentedMatrix[i][i]) < 1e-10 {
			return nil, errors.New("система не має єдиного розв'язку")
		}

		// Зведення елементів під діагоналлю до нуля
		for k := i + 1; k < m.rows; k++ {
			factor := augmentedMatrix[k][i] / augmentedMatrix[i][i]
			for j := i; j <= m.cols; j++ {
				augmentedMatrix[k][j] -= factor * augmentedMatrix[i][j]
			}
		}
	}

	// Зворотний хід методу Гауса
	solution := make([]float64, m.rows)
	for i := m.rows - 1; i >= 0; i-- {
		sum := augmentedMatrix[i][m.cols]
		for j := i + 1; j < m.cols; j++ {
			sum -= augmentedMatrix[i][j] * solution[j]
		}
		solution[i] = sum / augmentedMatrix[i][i]
	}

	return solution, nil
}

func main() {
	fmt.Println("Програма для роботи з матрицями")
	var rows, cols int
	fmt.Print("Введіть кількість рядків: ")
	fmt.Scan(&rows)
	fmt.Print("Введіть кількість стовпців: ")
	fmt.Scan(&cols)

	m := NewMatrix(rows, cols)
	m.Input()
	fmt.Println("Введена матриця:")
	m.Print()

	////////////////////////////////////////////
	fmt.Print("Введіть кількість рядків: ")
	fmt.Scan(&rows)
	fmt.Print("Введіть кількість стовпців: ")
	fmt.Scan(&cols)

	m2 := NewMatrix(rows, cols)
	m2.Input()
	fmt.Println("Введена матриця:")
	m2.Print()
	////////////////////////////////////////////

	fmt.Println("\nТранспонована матриця:")
	transposed := m.Transpose()
	transposed.Print()

	fmt.Println("\nДетермінант першої матриці:")
	det, err := m.Determinant()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%.2f\n", det)
	}

	fmt.Println("\nОбернута перша матриця:")
	inverse, err := m.Inverse()
	if err != nil {
		fmt.Println(err)
	} else {
		inverse.Print()
	}

	fmt.Println("\nДодавання:")
	result, err := m.Add(m2)
	if err != nil {
		fmt.Println(err)
	} else {
		result.Print()
	}

	fmt.Println("\nВіднімання:")
	result, err = m.Subtract(m2)
	if err != nil {
		fmt.Println(err)
	} else {
		result.Print()
	}

	fmt.Println("\nМноження:")
	result, err = m.Multiply(m2)
	if err != nil {
		fmt.Println(err)
	} else {
		result.Print()
	}

	fmt.Println("\nСортування по рядках:")
	r := &RowLexicographicSort{m}
	r.Sort()
	fmt.Println("Матриця після сортування:")
	m.Print()

	fmt.Println("\nСортування по стовпцях:")
	c := &ColumnLexicographicSort{m}
	c.Sort()
	fmt.Println("Матриця після сортування:")
	m.Print()

	fmt.Println("Розвʼязок СЛАР")
	b := make([]float64, m.rows)
	fmt.Print("Введіть вектор вільних членів: ")
	for i := 0; i < m.rows; i++ {
		fmt.Scan(&b[i])
	}
	solution, err := m.SolveSystem(b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Розвʼязок СЛАР:")
		for i, x := range solution {
			fmt.Printf("x%d = %.2f\n", i+1, x)
		}

		checkSolution := make([]float64, m.rows)
		for i := 0; i < m.rows; i++ {
			for j := 0; j < m.cols; j++ {
				checkSolution[i] += m.data[i][j] * solution[j]
			}
		}
		fmt.Println("Перевірка розв'язку:")
		for i, x := range checkSolution {
			fmt.Printf("x%d = %.2f\n", i+1, x)
		}
	}

}
