package Netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var pgm PGM

	// Lit le Nombre Magique
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

	// Écrire le nombre magique
	fmt.Fprintf(writer, "%s\n", pgm.magicNumber)

	// Écrire la largeur et la hauteureight
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)

	// Écrire la valeur maximale
	fmt.Fprintf(writer, "%d\n", pgm.max)

	// Écrire les données de l'image
	for _, row := range pgm.data {
		for _, value := range row {
			fmt.Fprintf(writer, "%d ", value)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
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
	pgm.max = int(maxValue)
}

// Rotate90CW fait pivoter l'image PGM de 90° dans le sens des aiguilles d'une montre.
func (pgm *PGM) Rotate90CW() {
	// Créez une nouvelle tranche 2D pour stocker l'image pivotée
	rotatedData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	// Faire pivoter l'image
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			rotatedData[x][pgm.height-y-1] = pgm.data[y][x]
		}
	}

	// Mettre à jour la largeur et la hauteur de l'image
	pgm.width, pgm.height = pgm.height, pgm.width

	//Mettre à jour les données avec l'image pivotée
	pgm.data = rotatedData
}

// ToPBM convertit l'image PGM en PBM.
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

// PBM représente une image BitMap portable.
type PBM struct {
	data   [][]uint8
	width  int
	height int
}

// NewPBM crée une nouvelle instance PBM avec la largeur et la hauteur données.
func NewPBM(width, height int) *PBM {
	return &PBM{
		data:   make([][]uint8, height),
		width:  width,
		height: height,
	}
}

// Set définit la valeur du pixel à (x, y).
func (pbm *PBM) Set(x, y int, value uint8) {
	pbm.data[y][x] = value
}
