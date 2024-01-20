package Netpbm

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

// PPM représente une image PPM.
type PPM struct {
	data          [][]Pixel // Pixels de l'image PPM représentés par un tableau bidimensionnel de pixels.
	width, height int       // Largeur et hauteur de l'image
	magicNumber   string    // Nombre magique du format PBM ("P3" ou "P6")
	max           int       // Valeur maximale d'un pixel dans l'image.
}

// PIxel représente un pixel de couleur.
type Pixel struct {
	R, G, B uint8
}

// Point représente un point dans un plan 2D avec des coordonnées X et Y.
type Point struct {
	X, Y int
}

// ReadPPM lit une image PPM à partir d'un fichier et renvoie une structure qui représente l'image.
func ReadPPM(filename string) (*PPM, error) {
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
	if magicNumber != "P3" && magicNumber != "P6" {
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
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid dimensions: width and height must be positive")
	}

	// Lire la valeur maximale
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

	// Lire les données d'image
	data := make([][]Pixel, height)
	expectedBytesPerPixel := 3

	if magicNumber == "P3" {
		// Lire le format P3 (ASCII)
		for y := 0; y < height; y++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				return nil, fmt.Errorf("error reading data at row %d: %v", y, err)
			}
			fields := strings.Fields(line)
			rowData := make([]Pixel, width)
			for x := 0; x < width; x++ {
				if x*3+2 >= len(fields) {
					return nil, fmt.Errorf("index out of range at row %d, column %d", y, x)
				}
				var pixel Pixel
				_, err := fmt.Sscanf(fields[x*3], "%d", &pixel.R)
				if err != nil {
					return nil, fmt.Errorf("error parsing Red value at row %d, column %d: %v", y, x, err)
				}
				_, err = fmt.Sscanf(fields[x*3+1], "%d", &pixel.G)
				if err != nil {
					return nil, fmt.Errorf("error parsing Green value at row %d, column %d: %v", y, x, err)
				}
				_, err = fmt.Sscanf(fields[x*3+2], "%d", &pixel.B)
				if err != nil {
					return nil, fmt.Errorf("error parsing Blue value at row %d, column %d: %v", y, x, err)
				}
				rowData[x] = pixel
			}
			data[y] = rowData
		}
	} else if magicNumber == "P6" {
		// Lire le format P6 (binaire)
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

			rowData := make([]Pixel, width)
			for x := 0; x < width; x++ {
				pixel := Pixel{R: row[x*expectedBytesPerPixel], G: row[x*expectedBytesPerPixel+1], B: row[x*expectedBytesPerPixel+2]}
				rowData[x] = pixel
			}
			data[y] = rowData
		}
	}

	// Renvoie la structure PPM
	return &PPM{data, width, height, magicNumber, max}, nil
}

func (ppm *PPM) PrintPPM() {
	fmt.Printf("Magic Number: %s\n", ppm.magicNumber)
	fmt.Printf("Width: %d\n", ppm.width)
	fmt.Printf("Height: %d\n", ppm.height)
	fmt.Printf("Max Value: %d\n", ppm.max)

	fmt.Println("Pixel Data:")
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			fmt.Printf("(%d, %d, %d) ", pixel.R, pixel.G, pixel.B)
		}
		fmt.Println()
	}
}

// Size renvoie la largeur et la hauteur de l'image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At renvoie la valeur du pixel en (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set définit la valeur du pixel à (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save enregistre l'image PPM dans un fichier et renvoie une erreur en cas de problème.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if ppm.magicNumber == "P6" || ppm.magicNumber == "P3" {
		fmt.Fprintf(file, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)
	} else {
		err = fmt.Errorf("magic number error")
		return err
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			if ppm.magicNumber == "P6" {
				// Conversion inverse des pixels
				file.Write([]byte{pixel.R, pixel.G, pixel.B})
			} else if ppm.magicNumber == "P3" {
				// Conversion inverse des pixels
				fmt.Fprintf(file, "%d %d %d ", pixel.R, pixel.G, pixel.B)
			}
		}
		if ppm.magicNumber == "P3" {
			fmt.Fprint(file, "\n")
		}
	}

	return nil
}

// Invert inverse les couleurs de l’image PPM.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x].R = uint8(ppm.max) - ppm.data[y][x].R
			ppm.data[y][x].G = uint8(ppm.max) - ppm.data[y][x].G
			ppm.data[y][x].B = uint8(ppm.max) - ppm.data[y][x].B
		}
	}
}

// Flip retourne l'image PPM horizontalement.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

// Flop fait basculer l'image PPM verticalement.
func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x], ppm.data[ppm.height-y-1][x] = ppm.data[ppm.height-y-1][x], ppm.data[y][x]
		}
	}
}

// SetMagicNumber définit le nombre magique de l'image PPM.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PPM.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Mettre à l'échelle les valeurs RVB en fonction de la nouvelle valeur maximale
			ppm.data[y][x].R = uint8(float64(ppm.data[y][x].R) * float64(maxValue) / float64(ppm.max))
			ppm.data[y][x].G = uint8(float64(ppm.data[y][x].G) * float64(maxValue) / float64(ppm.max))
			ppm.data[y][x].B = uint8(float64(ppm.data[y][x].B) * float64(maxValue) / float64(ppm.max))
		}
	}

	// Mettre à jour la valeur maximale
	ppm.max = int(maxValue)
}

// Rotate90CW fait pivoter l'image PPM de 90° dans le sens des aiguilles d'une montre.
func (ppm *PPM) Rotate90CW() {
	rotated := make([][]Pixel, ppm.width)
	for i := range rotated {
		rotated[i] = make([]Pixel, ppm.height)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			rotated[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotated
}

// ToPGM convertit l'image PPM en PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Créer une nouvelle image PGM avec la même largeur et hauteur que l'image PPM
	pgm := &PGM{
		data:        make([][]uint8, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         255,
	}

	// Convertir chaque pixel de RVB en niveaux de gris et attribuez-le à l'image PGM
	for y := 0; y < ppm.height; y++ {
		pgm.data[y] = make([]uint8, ppm.width)
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			gray := uint8((float32(pixel.R) + float32(pixel.G) + float32(pixel.B)) / 3)
			pgm.data[y][x] = gray
		}
	}

	return pgm
}

// ToPBM convertit l'image PPM en PBM.
func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	pbm.data = make([][]bool, ppm.height)
	for i := range pbm.data {
		pbm.data[i] = make([]bool, ppm.width)
	}

	//Définir un seuil pour la conversion binaire
	threshold := uint8(ppm.max / 2)

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Calculer l'intensité moyenne des valeurs RVB
			average := (uint16(ppm.data[y][x].R) + uint16(ppm.data[y][x].G) + uint16(ppm.data[y][x].B)) / 3
			// Définir la valeur binaire en fonction du seuil
			pbm.data[y][x] = average < uint16(threshold)
		}
	}
	return pbm
}

// SetPixel définit la couleur d'un pixel en un point donné.
func (ppm *PPM) SetPixel(p Point, color Pixel) {
	// Vérifier si le point se trouve dans les dimensions PPM
	if p.X >= 0 && p.X < ppm.width && p.Y >= 0 && p.Y < ppm.height {
		ppm.data[p.Y][p.X] = color
	}
}

// DrawLine trace une ligne entre deux points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	// Algorithme tracé de Bresenham
	x1, y1 := p1.X, p1.Y
	x2, y2 := p2.X, p2.Y

	dx := abs(x2 - x1)
	dy := abs(y2 - y1)

	var sx, sy int

	if x1 < x2 {
		sx = 1
	} else {
		sx = -1
	}

	if y1 < y2 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		ppm.SetPixel(Point{x1, y1}, color)

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x1 += sx
		}

		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DrawRectangle dessine un rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	// Dessiner les quatre côtés du rectangle à l'aide de DrawLine
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle dessine un rectangle rempli.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// Vérifier les dimensions valides
	if width <= 0 || height <= 0 {
		return
	}

	// Définir les coins du rectangle
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	// Dessiner les quatre côtés du rectangle à l'aide de DrawLine
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)

	// Déterminer la zone de remplissage
	minX := min(p1.X, min(p2.X, min(p3.X, p4.X)))
	maxX := max(p1.X, max(p2.X, max(p3.X, p4.X)))
	minY := min(p1.Y, min(p2.Y, min(p3.Y, p4.Y)))
	maxY := max(p1.Y, max(p2.Y, max(p3.Y, p4.Y)))

	// Remplir la zone de manière horizontale
	for y := minY + 1; y < maxY; y++ {
		for x := minX + 1; x < maxX; x++ {
			ppm.SetPixel(Point{x, y}, color)
		}
	}
}

// DrawCircle dessine un cercle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {

	for x := 0; x < ppm.height; x++ {
		for y := 0; y < ppm.width; y++ {
			dx := float64(x) - float64(center.X)
			dy := float64(y) - float64(center.Y)
			distance := math.Sqrt(dx*dx + dy*dy)

			if math.Abs(distance-float64(radius)) < 1.0 && distance < float64(radius) {
				ppm.Set(x, y, color)
			}
		}
	}
	ppm.Set(center.X-(radius-1), center.Y, color)
	ppm.Set(center.X+(radius-1), center.Y, color)
	ppm.Set(center.X, center.Y+(radius-1), color)
	ppm.Set(center.X, center.Y-(radius-1), color)
}

// DrawFilledCircle dessine un cercle rempli.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	ppm.DrawCircle(center, radius, color)

	for i := 0; i < ppm.height; i++ {
		var positions []int
		var numberPoints int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				numberPoints += 1
				positions = append(positions, j)
			}
		}
		if numberPoints > 1 {
			for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
				ppm.data[i][k] = color
			}
		}
	}
}

// DrawTriangle dessine un triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle dessine un triangle rempli.
func interpolate(p1, p2 Point, y int) float64 {
	return float64(p1.X) + (float64(y-p1.Y)/float64(p2.Y-p1.Y))*(float64(p2.X-p1.X))
}

func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	vertices := []Point{p1, p2, p3}
	sort.Slice(vertices, func(i, j int) bool {
		return vertices[i].Y < vertices[j].Y
	})

	for y := vertices[0].Y; y <= vertices[2].Y; y++ {
		x1 := interpolate(vertices[0], vertices[2], y)
		x2 := interpolate(vertices[1], vertices[2], y)

		ppm.DrawLine(Point{X: int(x1), Y: y}, Point{X: int(x2), Y: y}, color)
	}
}

// DrawPolygon dessine un polygone.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}

	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon dessine un polygone rempli.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	ppm.DrawPolygon(points, color)
	for i := 0; i < ppm.height; i++ {
		var positions []int
		var numberPoints int
		for j := 0; j < ppm.width; j++ {
			if ppm.data[i][j] == color {
				numberPoints += 1
				positions = append(positions, j)
			}
		}
		if numberPoints > 1 {
			for k := positions[0] + 1; k < positions[len(positions)-1]; k++ {
				ppm.data[i][k] = color

			}
		}
	}
}

// DrawKochSnowflake dessine un flocon de neige Koch.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, end Point, width int, color Pixel) {
	// N est le nombre d'itérations.
	// Le flocon de neige de Koch est une courbe de Koch 3 fois supérieure.
	// Start est le point culminant du flocon de neige.
	// Width est la largeur de toutes les lignes.
	// Color est la couleur des lignes.
	if n == 0 {
		ppm.DrawLine(start, end, color)
	} else {
		dx := (end.X - start.X) / 3.0
		dy := (end.Y - start.Y) / 3.0

		// Calculate the points for the segments
		p1 := Point{start.X + dx, start.Y + dy}
		p3 := Point{start.X + 2*dx, start.Y + 2*dy}

		angle := math.Pi / 3.0 // 60 degrees
		p2 := Point{
			X: int(float64((p1.X+p3.X)/2) - math.Sin(angle)*float64(p1.Y-p3.Y)/2),
			Y: int(float64(p1.Y+p3.Y)/2 + math.Sin(angle)*float64(p1.X-p3.X)/2),
		}

		// Recursively draw the four segments
		ppm.DrawKochSnowflake(n-1, start, p1, width, color)
		ppm.DrawKochSnowflake(n-1, p1, p2, width, color)
		ppm.DrawKochSnowflake(n-1, p2, p3, width, color)
		ppm.DrawKochSnowflake(n-1, p3, end, width, color)
	}
}

// DrawSierpinskiTriangle dessine un triangle de Sierpinski.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	// N est le nombre d'itérations.
	// Start est le point culminant du triangle.
	// Width est la largeur de toutes les lignes.
	// Color est la couleur des lignes.
	if n == 0 {
		// Base case: draw a triangle
		p1 := start
		p2 := Point{start.X + width, start.Y}
		p3 := Point{start.X + width/2, start.Y + int(float64(width)*math.Sqrt(3)/2)}

		ppm.DrawLine(p1, p2, color)
		ppm.DrawLine(p2, p3, color)
		ppm.DrawLine(p3, p1, color)
	} else {
		// Recursive case: divide the triangle into three smaller triangles and draw Sierpinski triangles on each
		halfWidth := width / 2
		halfHeight := int(float64(halfWidth) * math.Sqrt(3) / 2)

		// Calculate the coordinates for the three top points of the smaller triangles
		top1 := start
		top2 := Point{start.X + halfWidth, start.Y}
		top3 := Point{start.X + halfWidth/2, start.Y + halfHeight}

		// Recursively draw Sierpinski triangles on the three smaller triangles
		ppm.DrawSierpinskiTriangle(n-1, top1, halfWidth, color)
		ppm.DrawSierpinskiTriangle(n-1, top2, halfWidth, color)
		ppm.DrawSierpinskiTriangle(n-1, top3, halfWidth, color)
	}
}

func perlinNoise(x, y float64) float64 {
	return (x + y) / 2
}

// DrawPerlinNoise dessine le bruit Perlin.
// Cette fonction dessine un bruit perlin de toute l'image.
func DrawPerlinNoise(img *image.RGBA, color1 color.Color, color2 color.Color) {
	// Color1 is the color of 0.
	// Color2 is the color of 1.
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Generate Perlin noise
			noise := perlinNoise(float64(x), float64(y))

			// Calculate the interpolation coefficient
			t := 0.5 + 0.5*math.Cos(math.Pi*noise)

			// Linearly interpolate between the two colors
			r1, g1, b1, _ := color1.RGBA()
			r2, g2, b2, _ := color2.RGBA()
			r := uint8(float64(r1)*(1-t) + float64(r2)*t)
			g := uint8(float64(g1)*(1-t) + float64(g2)*t)
			b := uint8(float64(b1)*(1-t) + float64(b2)*t)

			// Set the pixel color
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
}

// KNearestNeighbors redimensionne l'image PPM à l'aide de l'algorithme des k-voisins les plus proches.
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int) {
	// ...
}
