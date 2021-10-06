package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	. "github.com/mjwong/sudoku1/linkedlist"
)

const (
	SQ = 3
	N  = 9
)

type (
	intmat [N][N]int
	pmat   [N][N][]int // matrix where cells contain list of possible values
)

var (
	iterCnt   int    // no. of iterations in iterMat
	mat, mat3 intmat // mat is the original matrix, mat3 hold the current result
)

func main() {
	var s string

	fmt.Println("Enter sudoku string (. represents empty square)")
	fmt.Scanf("%s\n", &s)

	if len(s) != N*N {
		log.Fatalf("Expected %d but got %d.\n", N*N, len(s))
	}

	mat = toGrid(s)
	printSudoku(mat)
	posmat, cnt := getPossibilityMat(mat) // empty cells contain a list of possibile values
	fmt.Printf("Starting empty cells = %d\n", cnt)
	printPossibleMat(posmat)

	emptyList := fillEmptyList(posmat)
	fmt.Printf("Nodes in empty list = %d\n", emptyList.CountNodes())

	mat3 = mat
	iterMat(emptyList.Head)
}

func iterMat(node *Cell) {
	startTime := time.Now()
	if node != nil {
		for _, num := range node.Vals {
			iterCnt++
			if isSafe(mat3, node.Row, node.Col, num) {
				mat3[node.Row][node.Col] = num
				if node.Next != nil {
					iterMat(node.Next)
					mat3[node.Row][node.Col] = 0
				} else {
					fmt.Printf("Iterations: %d. Time elapsed: %v\n", iterCnt, time.Since(startTime))
					color.LightRed.Println("******* Finished iterMat *******")
					checkSums(mat3)
					printSudoku(mat3)
				}
			}
		}
	}
}

func isSafe(m intmat, row, col, num int) bool {
	if !inRow(m, row, num) && !inCol(m, col, num) && !inSqu(m, row, col, num) {
		return true
	}
	return false
}

func toGrid(s string) intmat {
	var (
		index  int
		err    error
		substr string
		m      intmat
	)
	m = intmat{}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			substr = s[index : index+1]
			if strings.Compare(substr, ".") == 0 {
				m[i][j] = 0
			} else {
				m[i][j], err = strconv.Atoi(substr)
			}

			if err != nil {
				log.Fatalf("Error in reading char at %d,%d.\n", i, j)
			}
			index++
		}
	}
	fmt.Println()
	return m
}

func getPossibilityMat(m intmat) (pmat, int) {
	var (
		cnt int
		pm  pmat
	)

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if m[i][j] == 0 {
				for num := 1; num <= N; num++ {
					if !inRow(m, i, num) && !inCol(m, j, num) && !inSqu(m, i, j, num) {
						pm[i][j] = append(pm[i][j], num)
					}
				}
				cnt++
			}
		}
	}
	return pm, cnt
}

func fillEmptyList(pm pmat) *LinkedList {
	list := &LinkedList{}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if pm[i][j] != nil {
				list.AddCell(i, j, pm[i][j])
			}
		}
	}
	return list
}

func contains(a, b []int) bool {
	for _, v := range b {
		for _, w := range a {
			if v == w {
				return true
			}
		}
	}
	return false
}

func inRow(m intmat, row, num int) bool {
	for c := 0; c < N; c++ {
		if m[row][c] == num {
			return true
		}
	}
	return false
}

func inCol(m intmat, col, num int) bool {
	for r := 0; r < N; r++ {
		if m[r][col] == num {
			return true
		}
	}
	return false
}

func inSqu(m intmat, row, col, num int) bool {
	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if m[x][y] == num {
				return true
			}
		}
	}
	return false
}

func getSqList(m intmat, x, y int) []int {
	var (
		startRow, startCol int
		list               []int
	)
	startRow = x * SQ
	startCol = y * SQ

	for i := startRow; i < startRow+SQ; i++ {
		for j := startCol; j < startCol+SQ; j++ {
			if m[i][j] != 0 {
				list = append(list, m[i][j])
			}
		}
	}
	return list
}

// check the row, col and square sums
func checkSums(m intmat) bool {
	const (
		totalSum = 405
	)
	var (
		val     int
		rowSums []int
		colSums []int
		sqSums  []int
		success bool
	)
	success = false

	for i := 0; i < N; i++ {
		rSum := 0
		cSum := 0
		for j := 0; j < N; j++ {
			val = m[i][j]
			if val > 0 {
				rSum += val
			}

			val = m[j][i]
			if val > 0 {
				cSum += val
			}
		}
		rowSums = append(rowSums, rSum)
		colSums = append(colSums, cSum)

	}

	for i := 0; i < SQ; i++ {
		for j := 0; j < SQ; j++ {
			sSum := 0
			startrow := i * SQ
			startcol := j * SQ
			for x := startrow; x < startrow+SQ; x++ {
				for y := startcol; y < startcol+SQ; y++ {
					val = m[x][y]
					if val > 0 {
						sSum += val
					}
				}
			}
			sqSums = append(sqSums, sSum)
		}
	}

	fmt.Printf("Total sums: %2v\n", totalSum)
	fmt.Printf("Row sums: %2v\n", rowSums)
	fmt.Printf("Col sums: %2v\n", colSums)
	fmt.Printf("Squ sums: %2v\n", sqSums)

	if sumArr(rowSums) == totalSum && sumArr(colSums) == totalSum && sumArr(sqSums) == totalSum {
		color.New(color.FgLightBlue, color.OpBold).Println("Finished!")
		success = true
	}
	return success
}

func sumArr(arr []int) int {
	result := 0
	for _, v := range arr {
		result += v
	}
	return result
}

func printSudoku(m intmat) {
	var sqi, sqj int

	for i := 0; i < N; i++ {
		sqi = (i / SQ) % 2
		for j := 0; j < N; j++ {
			sqj = (j / SQ) % 2
			if (sqi == 0 && sqj == 1) || (sqi == 1 && sqj == 0) {
				color.LightBlue.Printf("%d ", m[i][j])
			} else {
				color.LightGreen.Printf("%d ", m[i][j])
			}
		}
		fmt.Println()
	}
}

func printPossibleMat(m pmat) {
	const (
		scols = 3
		ncols = 9
	)

	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------------------------------------------")

	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
			fmt.Printf("|%-16v ", arr2String(m[i][j], ","))
		}
		fmt.Println("|")
		if i != 0 && i%scols == 2 {
			fmt.Println("-------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		}
	}
}

func arr2String(a []int, delim string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

func intArrayEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
