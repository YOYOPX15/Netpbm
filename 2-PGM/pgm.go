package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// PGM représente une image GrayMap portable
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// ReadPGM lit une image PGM à partir d'un fichier et renvoie une structure qui représente l'image
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

// Size renvoie la largeur et la hauteur de l'image
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At renvoie la valeur du pixel en (x, y)
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set définit la valeur du pixel à (x, y)
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Enregistre enregistre l'image PGM dans un fichier et renvoie une erreur en cas de problème
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Écrit le nombre magique, la largeur, la hauteur et la valeur maximale
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Écrit les données Pixel
	for _, row := range pgm.data {
		for _, value := range row {
			fmt.Fprintf(writer, "%d ", value)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

// Inverser inverse les couleurs de l’image PGM
func (pgm *PGM) Invert() {
	for i := range pgm.data {
		for j := range pgm.data[i] {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip retourne l'image PGM horizontalement
func (pgm *PGM) Flip() {
	for i := range pgm.data {
		for j, k := 0, pgm.width-1; j < k; j, k = j+1, k-1 {
			pgm.data[i][j], pgm.data[i][k] = pgm.data[i][k], pgm.data[i][j]
		}
	}
}

// Flop flops the PGM image vertically
func (pgm *PGM) Flop() {
	for i, j := 0, pgm.height-1; i < j; i, j = i+1, j-1 {
		pgm.data[i], pgm.data[j] = pgm.data[j], pgm.data[i]
	}
}

// Flop fait basculer l'image PGM verticalement
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PGM
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise
func (pgm *PGM) Rotate90CW() {
	// Transpose the image
	for i := 0; i < pgm.height; i++ {
		for j := i + 1; j < pgm.width; j++ {
			pgm.data[i][j], pgm.data[j][i] = pgm.data[j][i], pgm.data[i][j]
		}
	}

	// Rotate90CW fait pivoter l'image PGM de 90° dans le sens des aiguilles d'une montre
	for i := range pgm.data {
		for j, k := 0, pgm.width-1; j < k; j, k = j+1, k-1 {
			pgm.data[i][j], pgm.data[i][k] = pgm.data[i][k], pgm.data[i][j]
		}
	}

	// Intervertir la largeur et la hauteur
	pgm.width, pgm.height = pgm.height, pgm.width
}

// ToPBM convertit l'image PGM en PBM
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

func main() {
	filename := "duck.pgm"

	// Lit le fichier PGM
	pgmImage, err := ReadPGM(filename)
	if err != nil {
		fmt.Println("Error reading PGM file:", err)
		return
	}

	// Affiche les données de l'image
	width, height := pgmImage.Size()
	fmt.Printf("Original Image Size: %dx%d\n", width, height)
	fmt.Println("Nombre Magique: ", pgmImage.magicNumber)
	fmt.Println("Taille Image: ", width, height)
	fmt.Println("Max: ", pgmImage.max)
	fmt.Println("Data:")
	for _, row := range pgmImage.data {
		for _, pixel := range row {
			if pixel == 0 {
				fmt.Print("0")
			} else if pixel == 10 {
				fmt.Print("10")
			} else {
				fmt.Print("7")
			}
		}
		fmt.Println()
	}

	// Invert colors
	pgmImage.Invert()
	fmt.Println("Image Inverted")

	// Flip horizontally
	pgmImage.Flip()
	fmt.Println("Image Flipped Horizontally")

	// Flop vertically
	pgmImage.Flop()
	fmt.Println("Image Flopped Vertically")

	// Rotate 90 degrees clockwise
	pgmImage.Rotate90CW()
	fmt.Println("Image Rotated 90° Clockwise")

	// Display modified image size
	width, height = pgmImage.Size()
	fmt.Printf("Modified Image Size: %dx%d\n", width, height)

	// Save sauvegarde l'image modifié
	err = pgmImage.Save("modified_duck.pgm")
	if err != nil {
		fmt.Println("Error saving modified image:", err)
		return
	}

	// Convert PGM to PBM
	pbmImage := pgmImage.ToPBM()

	// Display PBM image size
	pbmWidth, pbmHeight := pbmImage.width, pbmImage.height
	fmt.Printf("PBM Image Size: %dx%d\n", pbmWidth, pbmHeight)
}
