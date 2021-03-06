package d11

import (
	"bufio"
	_ "embed"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/matthinz/aoc-golang"
)

//go:embed input
var defaultInput string

func New() aoc.Day {
	return aoc.NewDay(11, defaultInput, Puzzle1, Puzzle2)
}

func Puzzle1(r io.Reader, l *log.Logger) string {

	input := parseInput(r)
	totalFlashes := 0

	for stepIndex := 0; stepIndex < 100; stepIndex++ {

		nextInput, flashes := step(input)
		totalFlashes += flashes

		input = nextInput
	}

	return strconv.Itoa(totalFlashes)
}

func Puzzle2(r io.Reader, l *log.Logger) string {
	input := parseInput(r)

	height := len(input)
	width := len(input[0])
	area := width * height

	stepIndex := 0

	for {
		stepIndex++

		nextInput, flashes := step(input)

		if flashes == area {
			return strconv.Itoa(stepIndex)
		}

		input = nextInput
	}
}

func parseInput(r io.Reader) [][]int {
	var result [][]int
	s := bufio.NewScanner(r)

	width := 0

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if len(line) == 0 {
			continue
		}
		if width == 0 {
			width = len(line)
		} else {
			if len(line) != width {
				panic("line is the wrong width")
			}
		}
		row := make([]int, width)
		for i, r := range line {
			num, err := strconv.ParseInt(string(r), 10, 8)
			if err != nil {
				panic(err)
			}
			row[i] = int(num)
		}
		result = append(result, row)
	}

	return result
}

func step(input [][]int) ([][]int, int) {
	result := increment(input)

	totalFlashes := 0
	for {
		next, flashes := flash(result)
		result = next
		totalFlashes += flashes
		if flashes == 0 {
			break
		}
	}

	return reset(result), totalFlashes
}

func reset(input [][]int) [][]int {
	height := len(input)
	width := len(input[0])

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if input[y][x] > 9 || input[y][x] == -1 {
				input[y][x] = 0
			}
		}
	}

	return input
}

func increment(input [][]int) [][]int {
	for y := 0; y < len(input); y++ {
		for x := 0; x < len(input[y]); x++ {
			input[y][x]++
		}
	}
	return input
}

func flash(input [][]int) ([][]int, int) {

	flashes := 0

	height := len(input)
	width := len(input[0])

	// make a copy of the input for us to modify
	result := make([][]int, height)
	copy(result, input)

	// flash everything in the input that is > 9
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if input[y][x] <= 9 {
				continue
			}

			for deltaY := -1; deltaY <= 1; deltaY++ {
				for deltaX := -1; deltaX <= 1; deltaX++ {

					yValid := y+deltaY >= 0 && y+deltaY < height
					xValid := x+deltaX >= 0 && x+deltaX < width

					if yValid && xValid {

						hasAlreadyFlashed := result[y+deltaY][x+deltaX] == -1

						if hasAlreadyFlashed {
							// this octopus has already flashed on this tick
							continue
						}

						result[y+deltaY][x+deltaX]++
					}

				}
			}

			result[y][x] = -1 // sentinel value indicating "don't mess with this one again"

			flashes++
		}
	}

	return result, flashes
}

func printGrid(grid [][]int) {
	for y := 0; y < len(grid); y++ {
		fmt.Println()
		for x := 0; x < len(grid[y]); x++ {
			fmt.Print(grid[y][x])
		}
	}
	fmt.Println()
}
