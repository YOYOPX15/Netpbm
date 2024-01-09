package main // Projet en cours

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// Ouverture du fichier
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Initialisation du scanner
	scanner := bufio.NewScanner(file)

	// Lecture nombre magique
	scanner.Scan()
	magicNumber := scanner.Text()

	// Lecture largeur et hauteur
	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	width, height := 0, 0
	fmt.Sscanf(dimensions[0], "%d", &width)
	fmt.Sscanf(dimensions[1], "%d", &height)

	// Lecture données de l'image
	var data [][]bool
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]bool, width)
		for i, char := range line {
			if char == '1' {
				row[i] = true
			}
		}
		data = append(data, row)
	}

	// Création et renvoi de la structure PBM
	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Size renvoie la largeur et la hauteur de l'image
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At renvoie la valeur du pixel en (x, y)
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set définit la valeur du pixel à (x, y)
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save enregistre l'image PBM dans un fichier et renvoie une erreur en cas de problème
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Écrit le nombre magique, la largeur et la hauteur
	file.WriteString(fmt.Sprintf("%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height))

	// Écrit les données de l'image
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pixel {
				file.WriteString("1")
			} else {
				file.WriteString("0")
			}
		}
		file.WriteString("\n")
	}

	return nil
}

// Invert inverse les couleurs de l'image PBM
func (pbm *PBM) Invert() {
	for y, row := range pbm.data {
		for x := range row {
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Flip retourne l'image PBM horizontalement
func (pbm *PBM) Flip() {
	for _, row := range pbm.data {
		for i, j := 0, pbm.width-1; i < j; i, j = i+1, j-1 {
			row[i], row[j] = row[j], row[i]
		}
	}
}

// Flop fait basculer l'image PBM verticalement
func (pbm *PBM) Flop() {
	for i, j := 0, pbm.height-1; i < j; i, j = i+1, j-1 {
		pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
	}
}

// SetMagicNumber définit le nombre magique de l'image PBM
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func main() {
	image, err := ReadPBM("example.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	width, height := image.Size()
	fmt.Printf("Image Size: %d x %d\n", width, height)

	// Operations
	image.Invert()
	image.Flip()
	image.Flop()

	// Sauvegarde image modifié
	err = image.Save("modified_example.pbm")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
