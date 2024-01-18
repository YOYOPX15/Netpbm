package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// PGM représente une image PGM.
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// ReadPGM lit une image PGM à partir d'un fichier et renvoie une structure qui représente l'image.
func ReadPGM(filename string) (*PGM, error) {
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
	if magicNumber != "P2" && magicNumber != "P5" {
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
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: width and height must be positive")
	}

	// Read max value
	maxValue, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading max value: %v", err)
	}
	maxValue = strings.TrimSpace(maxValue)
	var max int
	_, err = fmt.Sscanf(maxValue, "%d", &max)
	if err != nil {
		return nil, fmt.Errorf("invalid max value: %v", err)
	}

	// Read image data
	data := make([][]uint8, height)
	expectedBytesPerPixel := 1

	if magicNumber == "P2" {
		// Read P2 format (ASCII)
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			rowData := make([]uint8, width)
			for x, field := range fields {
				if x >= width {
					return nil, fmt.Errorf("index out of range at row %d", y)
				}
				var pixelValue uint8
				_, err := fmt.Sscanf(field, "%d", &pixelValue)
				if err != nil {
					return nil, fmt.Errorf("error parsing pixel value at row %d, column %d: %v", y, x, err)
				}
				rowData[x] = pixelValue
			}
			data[y] = rowData
		}
	} else if magicNumber == "P5" {
		// Read P5 format (binary)
		for y := 0; y < height; y++ {
			row := make([]byte, width*expectedBytesPerPixel)
			n, err := reader.Read(row)
			if err != nil {
				if err == io.EOF {
					return nil, fmt.Errorf("unexpected end of file at row %d", y)
				}
				return nil, fmt.Errorf("error reading pixel data at row %d: %v", y, err)
			}
			if n < width*expectedBytesPerPixel {
				return nil, fmt.Errorf("unexpected end of file at row %d, expected %d bytes, got %d", y, width*expectedBytesPerPixel, n)
			}

			rowData := make([]uint8, width)
			for x := 0; x < width; x++ {
				pixelValue := uint8(row[x*expectedBytesPerPixel])
				rowData[x] = pixelValue
			}
			data[y] = rowData
		}
	}

	// Return the PGM struct
	return &PGM{data, width, height, magicNumber, max}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At renvoie la valeur du pixel en (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set définit la valeur du pixel à (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save enregistre l'image PGM dans un fichier et renvoie une erreur en cas de problème.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = fmt.Fprintln(writer, pgm.magicNumber)
	if err != nil {
		return fmt.Errorf("error writing magic number: %v", err)
	}

	// Write dimensions
	_, err = fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	if err != nil {
		return fmt.Errorf("error writing dimensions: %v", err)
	}

	// Write max value
	_, err = fmt.Fprintln(writer, pgm.max)
	if err != nil {
		return fmt.Errorf("error writing max value: %v", err)
	}
	for _, row := range pgm.data {
		if len(row) != pgm.width {
			return fmt.Errorf("inconsistent row length in data")
		}
	}

	// Write image data
	if pgm.magicNumber == "P2" {
		err = saveP2PGM(writer, pgm)
		if err != nil {
			return err
		}
	} else if pgm.magicNumber == "P5" {
		err = saveP5PGM(writer, pgm)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

// saveP2PGM saves the PGM image in P2 format (ASCII).
func saveP2PGM(file *bufio.Writer, pgm *PGM) error {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			// Write the pixel value
			_, err := fmt.Fprint(file, pgm.data[y][x])
			if err != nil {
				return fmt.Errorf("error writing pixel data at row %d, column %d: %v", y, x, err)
			}

			// Add a space after each pixel, except the last one in a row
			if x < pgm.width-1 {
				_, err = fmt.Fprint(file, " ")
				if err != nil {
					return fmt.Errorf("error writing space after pixel at row %d, column %d: %v", y, x, err)
				}
			}
		}
		// Add a newline after each row
		_, err := fmt.Fprintln(file)
		if err != nil {
			return fmt.Errorf("error writing newline after row %d: %v", y, err)
		}
	}
	return nil
}

// saveP5PGM saves the PGM image in P5 format (binary).
func saveP5PGM(file *bufio.Writer, pgm *PGM) error {
	for y := 0; y < pgm.height; y++ {
		row := make([]byte, pgm.width)
		for x := 0; x < pgm.width; x++ {
			row[x] = byte(pgm.data[y][x])
		}
		_, err := file.Write(row)
		if err != nil {
			return fmt.Errorf("error writing pixel data at row %d: %v", y, err)
		}
	}
	return nil
}

// Invert inverse les couleurs de l’image PGM.
func (pgm *PGM) Invert() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pixel := pgm.data[y][x]
			invertedPixel := uint8(pgm.max) - pixel
			pgm.data[y][x] = invertedPixel
		}
	}
}

// Flip retourne l'image PGM horizontalement.
func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width/2; x++ {
			// Swap pixels horizontally
			pgm.data[y][x], pgm.data[y][pgm.width-x-1] = pgm.data[y][pgm.width-x-1], pgm.data[y][x]
		}
	}
}

// Flop fait basculer l'image PGM verticalement.
func (pgm *PGM) Flop() {
	for y := 0; y < pgm.height/2; y++ {
		// Swap rows vertically
		pgm.data[y], pgm.data[pgm.height-y-1] = pgm.data[pgm.height-y-1], pgm.data[y]
	}
}

// SetMagicNumber définit le nombre magique de l'image PGM.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PGM.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			// Scale the pixel values based on the new max value
			scaledValue := float64(pgm.data[y][x]) * float64(maxValue) / float64(pgm.max)
			// Round to the nearest integer
			newValue := uint8(scaledValue)
			pgm.data[y][x] = newValue
		}
	}

	// Update the max value
	pgm.max = int(maxValue)
}

// Rotate90CW fait pivoter l'image PGM de 90° dans le sens des aiguilles d'une montre.
func (pgm *PGM) Rotate90CW() {
	if pgm.width <= 0 || pgm.height <= 0 {
		return
	}

	newData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		newData[i] = make([]uint8, pgm.height)
		for j := 0; j < pgm.height; j++ {
			newData[i][j] = pgm.data[pgm.height-j-1][i]
		}
	}
	pgm.data = newData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// ToPBM convertit l'image PGM en PBM.
func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{
		data:        make([][]bool, pgm.height),
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1",
	}

	for y := 0; y < pgm.height; y++ {
		pbm.data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			pbm.data[y][x] = pgm.data[y][x] < uint8(pgm.max/2)
		}
	}
	return pbm
}

func (pgm *PGM) PrintData() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			fmt.Printf("%d ", pgm.data[i][j])
		}
		fmt.Println()
	}
}
