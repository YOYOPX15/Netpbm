package Netpbm // ✨ PBM

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// PBM représente une image PBM.
type PBM struct {
	data          [][]bool // Matrice de données représentant les pixels de l'image (true pour blanc, false pour noir)
	width, height int      // Largeur et hauteur de l'image
	magicNumber   string   // Nombre magique du format PBM ("P1" ou "P4")
}

// ReadPBM lit une image PBM à partir d'un fichier et renvoie une structure qui représente l'image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// Lire le nombre magique
	magicNumber, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("error reading magic number: %v", err)
	}
	magicNumber = strings.TrimSpace(magicNumber)
	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("invalid magic number: %s", magicNumber)
	}

	// Lire les dimensions
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
		// Lire le format P1 (ASCII)
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
		// Lire le format P4 (binaire)
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

				// Convertir ASCII en décimal et extraire le bit
				decimalValue := int(row[byteIndex])
				bitValue := (decimalValue >> bitIndex) & 1

				data[y][x] = bitValue != 0
			}
		}
	}

	return &PBM{data, width, height, magicNumber}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At renvoie la valeur du pixel en (x, y).
func (pbm *PBM) At(x, y int) bool {
	if x < 0 || x >= pbm.width || y < 0 || y >= pbm.height {
		return false
	}
	return pbm.data[y][x]
}

// Set définit la valeur du pixel à (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	if x < 0 || x >= pbm.width || y < 0 || y >= pbm.height {
		return
	}
	pbm.data[y][x] = value
}

// Save enregistre l'image PBM dans un fichier et renvoie une erreur en cas de problème.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrire un nombre magique
	_, err = file.WriteString(pbm.magicNumber + "\n")
	if err != nil {
		return err
	}

	// Écrire les dimensions
	_, err = file.WriteString(strconv.Itoa(pbm.width) + " " + strconv.Itoa(pbm.height) + "\n")
	if err != nil {
		return err
	}

	// Écrire des données
	if pbm.magicNumber == "P1" {
		// Format ASCII
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
		// Format binaire
		for _, row := range pbm.data {
			// Réinitialiser la tranche d'octets pour chaque ligne
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

// Invert inverse les couleurs de l'image PBM.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip retourne l'image PBM horizontalement.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width/2; j++ {
			pbm.data[i][j], pbm.data[i][pbm.width-j-1] = pbm.data[i][pbm.width-j-1], pbm.data[i][j]
		}
	}
}

// Flop fait basculer l'image PBM verticalement.
func (pbm *PBM) Flop() {
	for i := 0; i < pbm.height/2; i++ {
		pbm.data[i], pbm.data[pbm.height-i-1] = pbm.data[pbm.height-i-1], pbm.data[i]
	}
}

// SetMagicNumber définit le nombre magique de l'image PBM.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}
