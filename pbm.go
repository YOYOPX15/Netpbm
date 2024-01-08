package netpbm // Projet en cours

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Lire nombre magique
	scanner.Scan()
	magicNumber := scanner.Text()

	// Lire largeur et hauteur
	scanner.Scan()
	dimensions := strings.Fields(scanner.Text())
	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, err
	}
	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, err
	}

	// Lire données de l'image
	var data [][]bool
	for scanner.Scan() {
		line := scanner.Text()
		var row []bool
		for _, char := range line {
			if char == '1' {
				row = append(row, true)
			} else if char == '0' {
				row = append(row, false)
			}
		}
		data = append(data, row)
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

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
