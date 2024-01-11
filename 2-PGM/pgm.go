package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	magicNumber := scanner.Text()

	scanner.Scan()
	size := strings.Split(scanner.Text(), " ")
	width, _ := strconv.Atoi(size[0])
	height, _ := strconv.Atoi(size[1])

	scanner.Scan()
	max, _ := strconv.Atoi(scanner.Text())

	data := make([][]uint8, height)
	for i := range data {
		data[i] = make([]uint8, width)
	}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			scanner.Scan()
			value, _ := strconv.Atoi(scanner.Text())
			data[i][j] = uint8(value)
		}
	}

	return &PGM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         max,
	}, nil
}

func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {

}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {

}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {

}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {

}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {

}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {

}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {

}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {

}
