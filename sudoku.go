package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type List []uint8

func (l *List) RemoveElement(e uint8) {
	baseLength := len(*l)
	cutIndex := -1
	for i := 0; i < baseLength; i++ {
		if (*l)[i] == e {
			cutIndex = i
			break
		}
	}
	if cutIndex < 0 {
		return
	}
	removedVersion := append((*l)[:cutIndex], (*l)[cutIndex+1:]...)
	*l = removedVersion
}

func (l *List) ToString() string {
	baseLength := len(*l)
	ret := "["
	for i := 0; i < baseLength; i++ {
		ret += strconv.Itoa(int((*l)[i]))
	}
	return ret + "]"
}

func (l *List) contains(e uint8) bool {
	baseLength := len(*l)
	for i := 0; i < baseLength; i++ {
		if (*l)[i] == e {
			return true
		}
	}
	return false
}

type cell struct {
	value          uint8
	xIndex         uint8
	yIndex         uint8
	possibleValues *List
}

type column struct {
	index uint8
	cells [9]*cell
}

type line struct {
	index uint8
	cell  [9]*cell
}

type miniGrid struct {
	index uint8
	cells [9]*cell
}

func (c *cell) GetPosition() (uint8, uint8) {
	return c.xIndex, c.yIndex
}

type modification struct {
	index uint8
}

type hypothesis struct {
	valuesTested List  // list of values we've already tried for this cell
	cellIndex    uint8 // index of the cell we hypothesized the value of
}

// Grid is a representation of a sudoku grid
type Grid struct {
	cells      [81]cell
	history    []uint8      // History of modifications made to the grid
	hypotheses []hypothesis // Which steps did we start making hypotheses ?
	isPossible bool         // Is it possible to solve this grid ?
}

// ParseGrid parses a grid from an array of numbers
func ParseGrid(array [81]uint8) Grid {
	g := Grid{}
	g.isPossible = true
	for i := 0; i < 81; i++ {
		g.cells[i] = cell{
			xIndex:         uint8(i % 9),
			yIndex:         uint8(i / 9),
			value:          array[i],
			possibleValues: &List{1, 2, 3, 4, 5, 6, 7, 8, 9},
		}
	}
	return g
}

// Little helper with errors
func check(e error){
	if e != nil {
		panic(e)
	}
}
// LoadGrid loads a grid by parsing a file
func LoadGrid(fileName string) (Grid, error) {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil{
		return Grid{}, err
	}
	fmt.Print("Reading file " + fileName)
	strData := string(dat)
	splitStr := strings.Split(strData, " ")
	if len(splitStr) != 81{
		fmt.Println("File did not contain 81 numbers!")
	}
	numArray := [81]uint8{}
	for i := 0; i < 81; i++{
		num, e := strconv.Atoi(splitStr[i])
		if e != nil{
			return Grid{}, e
		}
		numArray[i] = uint8(num)
	}
	return ParseGrid(numArray), nil
}

// ToString conveniently converts the sudoku grid into a string for printing
func (g *Grid) ToString() string {
	ret := ""
	for i := 0; i < 81; i++ {
		if i%9 == 0 && i != 0 {
			ret += "\n"
			if i%27 == 0 {
				ret += "-----------\n"
			}

		}
		if i%9 == 3 || i%9 == 6 {
			ret += "|"
		}
		ret += strconv.FormatUint(uint64(g.cells[i].value), 10)
	}
	ret += "\n"
	// index := 1
	// ret += "Minigrid for cell " + strconv.Itoa(int(index)) + ": "
	// mini := g.getMiniGrid(uint8(index))
	// for i := 0; i < 9; i++ {
	// 	ret += " " + strconv.Itoa(int(mini[i].value))
	// }
	// ret += "\n"
	// ret += "Column for cell " + strconv.Itoa(int(index)) + ": "
	// mini = g.getColumn(uint8(index))
	// for i := 0; i < 9; i++ {
	// 	ret += " " + strconv.Itoa(int(mini[i].value))
	// }
	// ret += "\n"
	// ret += "Line for cell " + strconv.Itoa(int(index)) + ": "
	// mini = g.getLine(uint8(index))
	// for i := 0; i < 9; i++ {
	// 	ret += " " + strconv.Itoa(int(mini[i].value))
	// }
	// for i := 0; i < 81; i++ {
	// 	ret += strconv.Itoa(int(i)) + ": " + g.cells[i].possibleValues.ToString() + "\n"
	// }
	return ret
}

// PrintPossibleValues shows a list of possible values
// for the remaining cells of the grid
func (g *Grid) PrintPossibleValues() string {
	ret := ""
	for i := 0; i < 81; i++ {
		ret += strconv.Itoa(int(i)) + ": " + g.cells[i].possibleValues.ToString() + "\n"
	}
	return ret
}

func (g *Grid) getMiniGrid(index uint8) [9]*cell {
	bigLine := index / 27        // Getting which line of big cells the cell is in.
	bigColumn := (index % 9) / 3 // getting the column of big cells
	startIndex := bigLine*27 + bigColumn*3
	return [9]*cell{
		&g.cells[startIndex], &g.cells[startIndex+1], &g.cells[startIndex+2],
		&g.cells[startIndex+9], &g.cells[startIndex+1+9], &g.cells[startIndex+2+9],
		&g.cells[startIndex+18], &g.cells[startIndex+1+18], &g.cells[startIndex+2+18],
	}
}

func (g *Grid) getColumn(index uint8) [9]*cell {
	columnNumber := index % 9
	var toRet [9]*cell
	for i := 0; i < 9; i++ {
		toRet[i] = &g.cells[columnNumber+9*uint8(i)]
	}
	return toRet
}

func (g *Grid) getLine(index uint8) [9]*cell {
	lineNumber := index / 9
	var toRet [9]*cell
	for i := 0; i < 9; i++ {
		toRet[i] = &g.cells[9*lineNumber+uint8(i)]
	}
	return toRet
}

func (g *Grid) setCellValue(x uint8, y uint8, nValue uint8) {
	index := x*9 + y
	c := g.cells[index]
	if c.possibleValues.contains(nValue) {
		c.value = nValue
		c.possibleValues = &List{nValue}
	}
}

// For all the cells in the grid, update their possible values
func (g *Grid) updatePossibleValues() {
	var i uint8
	for i = 0; i < 81; i++ {
		cell := &(g.cells[i])
		cell.possibleValues = &List{1, 2, 3, 4, 5, 6, 7, 8, 9}
		if cell.value != 0 {
			for j := 1; j < 10; j++ {
				cell.possibleValues.RemoveElement(uint8(j))
			}
		}
		{
			cellLine := g.getLine(i)
			cellColumn := g.getColumn(i)
			cellMiniGrid := g.getMiniGrid(i)
			var j uint8
			for j = 0; j < 9; j++ {
				cell.possibleValues.RemoveElement(cellLine[j].value)
				cell.possibleValues.RemoveElement(cellColumn[j].value)
				cell.possibleValues.RemoveElement(cellMiniGrid[j].value)
			}
		}

	}
}

// Fill in the cells depending on their possible values.
func (g *Grid) fillValues() {
	var i uint8
	for i = 0; i < 81; i++ {
		cell := &(g.cells[i])
		if len(*cell.possibleValues) == 1 && cell.value == 0 {
			//fmt.Println("Filled index " + strconv.Itoa(int(i)))
			cell.value = (*cell.possibleValues)[0]
			g.history = append(g.history, uint8(i))
		}
		if cell.value == 0 && len(*cell.possibleValues) == 0 {
			c_x, c_y := cell.GetPosition()
			fmt.Println("No solution to this sudoku: the cell (" + strconv.Itoa(int(c_x)) + "," + strconv.Itoa(int(c_y)) + ")")
			g.isPossible = false
		}
		// For the minigrid, get all the possibilities.
		cellMiniGrid := g.getMiniGrid(i)
		possibilities := [9]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0}
		for j := 0; j < 9; j++ {
			poss := *(cellMiniGrid[j].possibleValues)
			for k := 0; k < 9; k++ {
				if poss.contains(uint8(k)) {
					possibilities[k]++
				}
			}
		}
		for k := 0; k < 9; k++ {
			if possibilities[k] == 1 && cell.possibleValues.contains(uint8(k)) {
				cell.value = uint8(k)
				g.history = append(g.history, uint8(i))
			}
		}
		// For the column, get all the possibilities.
		cellColumn := g.getColumn(i)
		possibilities = [9]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0}
		for j := 0; j < 9; j++ {
			poss := *(cellColumn[j].possibleValues)
			for k := 0; k < 9; k++ {
				if poss.contains(uint8(k)) {
					possibilities[k]++
				}
			}
		}
		for k := 0; k < 9; k++ {
			if possibilities[k] == 1 && cell.possibleValues.contains(uint8(k)) {
				cell.value = uint8(k)
				g.history = append(g.history, uint8(i))
			}
		}
		// For the line, get all the possibilities.
		cellLine := g.getLine(i)
		possibilities = [9]uint8{0, 0, 0, 0, 0, 0, 0, 0, 0}
		for j := 0; j < 9; j++ {
			poss := *(cellLine[j].possibleValues)
			for k := 0; k < 9; k++ {
				if poss.contains(uint8(k)) {
					possibilities[k]++
				}
			}
		}
		for k := 0; k < 9; k++ {
			if possibilities[k] == 1 && cell.possibleValues.contains(uint8(k)) {
				cell.value = uint8(k)
				g.history = append(g.history, uint8(i))
			}
		}
	}
}

func (g *Grid) makeHypothesis() {
	for i := 0; i < 81; i++ {
		cell := &(g.cells[i])
		alreadyConsidered := -1
		possibilities := cell.possibleValues
		if len(*(cell.possibleValues)) == 2 {
			for k := 0; k < len(g.hypotheses); k++ {
				if g.hypotheses[k].cellIndex == uint8(i) {
					alreadyConsidered = k
					// Removing previously tested hypotheses
					for j := 1; j < 10; j++ {
						if g.hypotheses[k].valuesTested.contains(uint8(j)) {
							possibilities.RemoveElement(uint8(j))
						}
					}
				}
			}
			if len(*possibilities) > 0 {
				newValue := (*possibilities)[0]
				cell.value = newValue
				g.history = append(g.history, uint8(i))
				if alreadyConsidered >= 0 {
					g.hypotheses[alreadyConsidered].valuesTested = append(g.hypotheses[alreadyConsidered].valuesTested, newValue)
				}
				{
					g.hypotheses = append(g.hypotheses, hypothesis{cellIndex: uint8(i), valuesTested: List{newValue}})
				}
				g.isPossible = true
				fmt.Println("Making a hypothesis : " + strconv.Itoa(int(cell.value)) + " for cell " + strconv.Itoa(i))
				return
			}

		}
	}
}

func (g *Grid) revertLastHypothesis() {
	// Getting a list of all the previously tested hypotheses.
	fmt.Println("Reverting")
	listModifiedHypotheses := List{}
	for i := 0; i < 0; i++ {
		listModifiedHypotheses = append(listModifiedHypotheses, g.hypotheses[i].cellIndex)
	}
	for i := len(g.history) - 1; i >= 0; i-- {
		fmt.Println(strconv.Itoa(int(g.history[i])) + " : " + strconv.Itoa(int(g.cells[g.history[i]].value)))
		g.cells[g.history[i]].value = 0
		cellIndex := g.history[i]
		g.history = g.history[:i]
		if listModifiedHypotheses.contains(cellIndex) {
			break
		}
	}
}

func (g *Grid) checkCompletion() bool {
	test := true
	for i := 0; i < 81; i++ {
		if g.cells[i].value == 0 {
			test = false
		}
	}
	return test
}

// Solve is a solving method for the sudoku grid
func (g *Grid) Solve() {
	solved := false
	counter := 0
	for !solved {
		currentMoves := len(g.history) // number of moves already made
		g.updatePossibleValues()
		g.fillValues()
		newMoveNumber := len(g.history) // new number of moves made
		if !g.isPossible {
			// If something is wrong, revert till the previous hypothesis
			g.revertLastHypothesis()
			g.updatePossibleValues()
			newMoveNumber = currentMoves
		}
		// solved = g.checkCompletion()
		if newMoveNumber-currentMoves == 0 {
			// If no progress has been made this round, we look for a hypothesis to make
			g.makeHypothesis()
			if !g.isPossible {
				// If no hypothesis was found
				fmt.Println("No solution found!")
				solved = true
			}
		}
		counter++
		//fmt.Println(strconv.Itoa(counter) + " loop(s) done!")
		if counter == 100 {
			solved = true
		}
	}
}


func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter the full path to the .grid file.")
	fileName, _ := reader.ReadString('\n')
	g, err := LoadGrid(fileName[:len(fileName)-1])
	if err != nil{
		fmt.Println("The file could not be loaded, please re-try it!")
		fmt.Println(err)
		return
	}
	fmt.Println(g.ToString())
	g.Solve()
	fmt.Println(g.ToString())
}
