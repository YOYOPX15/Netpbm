package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Read magic number
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Read dimensions
	dimensions, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading dimensions: %v", err)
	}
	var width, height int
	_, err = fmt.Sscanf(strings.TrimSpace(dimensions), "%d %d", &width, &height)
	if err != nil {
		return nil, fmt.Errorf("invalid dimensions: %v", err)
	}

	data := make([][]bool, height)

	for i := range data {
		data[i] = make([]bool, width)
	}

	if magicNumber == "P1" {
		// Read P1 format (ASCII)
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at row %d", y)
				}
				data[y][x] = field == "1"
			}
		}

	} else if magicNumber == "P4" {
		// Read P4 format (binary)
		expectedBytesPerRow := (width + 7) / 8
		for y := 0; y < height; y++ {
			row := make([]byte, expectedBytesPerRow)
			n, err := reader.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at row %d", y)
				}
				return nil, fmt.Errorf("error reading pixel data at row %d: %v", y, err)
			}
			if n < expectedBytesPerRow {
				return nil, fmt.Errorf("unexpected end of file at row %d, expected %d bytes, got %d", y, expectedBytesPerRow, n)
			}

			for x := 0; x < width; x++ {
				byteIndex := x / 8
				bitIndex := 7 - (x % 8)

				// Convert ASCII to decimal and extract the bit
				decimalValue := int(row[byteIndex])
				bitValue := (decimalValue >> bitIndex) & 1

				data[y][x] = bitValue != 0
			}
		}
	}

	return &PBM{data, width, height, magicNumber}, nil
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	if x < 0 || x >= pbm.width || y < 0 || y >= pbm.height {
		return false
	}
	return pbm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	if x < 0 || x >= pbm.width || y < 0 || y >= pbm.height {
		return
	}
	pbm.data[y][x] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write magic number
	_, err = file.WriteString(pbm.magicNumber + "\n")
	if err != nil {
		return err
	}

	// Write dimensions
	_, err = file.WriteString(strconv.Itoa(pbm.width) + " " + strconv.Itoa(pbm.height) + "\n")
	if err != nil {
		return err
	}

	// Write data
	if pbm.magicNumber == "P1" {
		// ASCII format
		for _, row := range pbm.data {
			for _, pixel := range row {
				if pixel {
					_, err = file.WriteString("1 ")
				} else {
					_, err = file.WriteString("0 ")
				}
				if err != nil {
					return err
				}
			}
			_, err = file.WriteString("\n")
			if err != nil {
				return err
			}
		}
	} else if pbm.magicNumber == "P4" {
		// Binary format
		for _, row := range pbm.data {
			// Reset the bytes slice for each row
			bytes := make([]byte, (pbm.width+7)/8)
			for x, pixel := range row {
				if pixel {
					byteIndex := x / 8
					bitIndex := uint(x % 8)
					bytes[byteIndex] |= 1 << (7 - bitIndex)
				}
			}
			_, err = file.Write(bytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
