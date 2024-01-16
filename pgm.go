package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var pgm PGM

	// Lit Nombre Magique
	scanner.Scan()
	pgm.magicNumber = scanner.Text()

	// Lire la largeur et la hauteur
	scanner.Scan()
	pgm.width, _ = strconv.Atoi(scanner.Text())
	scanner.Scan()
	pgm.height, _ = strconv.Atoi(scanner.Text())

	// Lit la valaur Max
	scanner.Scan()
	pgm.max, _ = strconv.Atoi(scanner.Text())

	// Lit les données Pixel
	pgm.data = make([][]uint8, pgm.height)
	for i := range pgm.data {
		pgm.data[i] = make([]uint8, pgm.width)
		for j := range pgm.data[i] {
			scanner.Scan()
			value, _ := strconv.Atoi(scanner.Text())
			pgm.data[i][j] = uint8(value)
		}
	}

	return &pgm, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write the magic number
	fmt.Fprintf(writer, "%s\n", pgm.magicNumber)

	// Write the width and height
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)

	// Write the maximum value
	fmt.Fprintf(writer, "%d\n", pgm.max)

	// Write the image data
	for _, row := range pgm.data {
		for _, value := range row {
			fmt.Fprintf(writer, "%d ", value)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pixel := pgm.data[y][x]
			invertedPixel := uint8(pgm.max) - pixel
			pgm.data[y][x] = invertedPixel
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width/2; x++ {
			// Swap pixels horizontally
			pgm.data[y][x], pgm.data[y][pgm.width-x-1] = pgm.data[y][pgm.width-x-1], pgm.data[y][x]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for y := 0; y < pgm.height/2; y++ {
		// Swap rows vertically
		pgm.data[y], pgm.data[pgm.height-y-1] = pgm.data[pgm.height-y-1], pgm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	// Create a new 2D slice to store the rotated image
	rotatedData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	// Rotate the image
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			rotatedData[x][pgm.height-y-1] = pgm.data[y][x]
		}
	}

	// Update the width and height of the image
	pgm.width, pgm.height = pgm.height, pgm.width

	// Update the data with the rotated image
	pgm.data = rotatedData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	threshold := uint8(pgm.max / 2)
	pbm := NewPBM(pgm.width, pgm.height)
	for i := range pgm.data {
		pbm.data[i] = make([]uint8, pgm.width)
		for j := range pgm.data[i] {
			if pgm.data[i][j] > threshold {
				pbm.Set(j, i, 1) // Blanc
			} else {
				pbm.Set(j, i, 0) // Noir
			}
		}
	}
	return pbm
}

// PBM représente une image BitMap portable
type PBM struct {
	data   [][]uint8
	width  int
	height int
}

// NewPBM crée une nouvelle instance PBM avec la largeur et la hauteur données
func NewPBM(width, height int) *PBM {
	return &PBM{
		data:   make([][]uint8, height),
		width:  width,
		height: height,
	}
}

// Set définit la valeur du pixel à (x, y)
func (pbm *PBM) Set(x, y int, value uint8) {
	pbm.data[y][x] = value
}
