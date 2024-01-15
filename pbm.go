package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

/*
// Ouverture du fichier
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Initialisation du scanner
	scanner := bufio.NewScanner(file)

	// Lecture nombre magique et commentaires
	var magicNumber string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Ignore comments
			continue
		}
		magicNumber = line
		break
	}

	// Vérifie que le nombre magique a été trouvé
	if magicNumber == "" {
		return nil, errors.New("Nombre magique introuvable")
	}

	// Lecture largeur et hauteur
	var width, height int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Ignore comments
			continue
		}
		dimensions := strings.Fields(line)
		if len(dimensions) == 2 {
			width, err = strconv.Atoi(dimensions[0])
			if err != nil {
				return nil, err
			}
			height, err = strconv.Atoi(dimensions[1])
			if err != nil {
				return nil, err
			}
			break
		}
	}

	// Vérifie que les dimensions ont été trouvées
	if width == 0 || height == 0 {
		return nil, errors.New("Dimensions invalides")
	}

	// Lecture données de l'image
	var data [][]bool
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Ignore comments
			continue
		}
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

	// Création et renvoi de la structure PBM
	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}
*/

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

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	// Write magic number
	_, err = file.WriteString(pbm.magicNumber + "\n")
	if err != nil {
		return fmt.Errorf("error writing magic number: %v", err)
	}

	// Write dimensions
	_, err = file.WriteString(fmt.Sprintf("%d %d\n", pbm.width, pbm.height))
	if err != nil {
		return fmt.Errorf("error writing dimensions: %v", err)
	}

	// Write data
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pixel {
				_, err = file.WriteString("1")
			} else {
				_, err = file.WriteString("0")
			}
			if err != nil {
				return fmt.Errorf("error writing pixel data: %v", err)
			}
		}
		_, err = file.WriteString("\n")
		if err != nil {
			return fmt.Errorf("error writing newline: %v", err)
		}
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

/*
func main() {
	image, err := ReadPBM("duck.pbm")
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	width, height := image.Size()
	fmt.Println("Nombre Magique: ", image.magicNumber)
	fmt.Println("Taille Image: ", width, height)
	fmt.Println("Data: ")
	for _, row := range image.data {
		for _, pixel := range row {
			if pixel {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}

	fmt.Print("\n")
	fmt.Print("Image Modifié:\n")

	image.Invert()
	fmt.Println("Image Inversé Couleur")

	image.Flip()
	fmt.Println("Image Inversée Horizontalement")

	image.Flop()
	fmt.Println("Image Renversée Verticalement")

	fmt.Println("Data: ")
	for _, row := range image.data {
		for _, pixel := range row {
			if pixel {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}

	// Sauvegarde image modifié
	err = image.Save("modified_duck.pbm")
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}
}
*/
