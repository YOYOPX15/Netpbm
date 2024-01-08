package netpbm // Projet en cours

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

func (pbm *PBM) Save(filename string) error {

}

func (pbm *PBM) Invert() {

}

func (pbm *PBM) Flip() {

}

func (pbm *PBM) Flop() {

}

func (pbm *PBM) SetMagicNumber(magicNumber string) {

}

func main() {

}
