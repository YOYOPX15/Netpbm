package Netpbm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// PPM représente une image PPM.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

// PIxel représente un pixel de couleur.
type Pixel struct {
	R, G, B uint8
}

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
	x := radius
	y := 0
	err := 0

	for x >= y {
		ppm.Set(center.X+x, center.Y+y, color)
		ppm.Set(center.X+y, center.Y+x, color)
		ppm.Set(center.X-y, center.Y+x, color)
		ppm.Set(center.X-x, center.Y+y, color)
		ppm.Set(center.X-x, center.Y-y, color)
		ppm.Set(center.X-y, center.Y-x, color)
		ppm.Set(center.X+y, center.Y-x, color)
		ppm.Set(center.X+x, center.Y-y, color)

		if err <= 0 {
			y += 1
			err += 2*y + 1
		}

		if err > 0 {
			x -= 1
			err -= 2*x + 1
		}
	}
}

// DrawFilledCircle dessine un cercle rempli.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
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
	// Sort the points by y-coordinate
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y
	})

	// Create a list to store the x-coordinates of the intersections
	intersections := make([]int, ppm.height)

	// Initialize the intersections with the maximum x-coordinate
	for i := range intersections {
		intersections[i] = ppm.width
	}

	// Iterate over each pair of consecutive points
	for i := 0; i < len(points); i++ {
		// Get the current and next point
		current := points[i]
		next := points[(i+1)%len(points)]

		// Calculate the slope of the edge
		slope := float64(next.X-current.X) / float64(next.Y-current.Y)

		// Iterate over each scanline
		for y := current.Y; y <= next.Y; y++ {
			// Calculate the x-coordinate of the intersection
			x := int(float64(current.X) + slope*float64(y-current.Y))

			// Update the intersection if the x-coordinate is smaller
			if x < intersections[y] {
				intersections[y] = x
			}
		}
	}

	// Iterate over each scanline
	for y := 0; y < ppm.height; y++ {
		// Draw a horizontal line between the intersections
		ppm.DrawLine(Point{intersections[y], y}, Point{intersections[y+1], y}, color)
	}
}

/*
// DrawKochSnowflake dessine un flocon de neige Koch.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
	// N est le nombre d'itérations.
	// Le flocon de neige de Koch est une courbe de Koch 3 fois supérieure.
	// Start est le point culminant du flocon de neige.
	// Width est la largeur de toutes les lignes.
	// Color est la couleur des lignes.

	// Draw the initial Koch curve
	ppm.drawKochCurve(n, start, width, color)

	// Calculate the coordinates for the other two points of the equilateral triangle
	height := int(float64(width) * math.Sqrt(3) / 2)
	p2 := Point{start.X - width/2, start.Y + height}
	p3 := Point{start.X + width/2, start.Y + height}

	// Draw the other two Koch curves
	ppm.drawKochCurve(n, p2, width, color)
	ppm.drawKochCurve(n, p3, width, color)
}

// Additionnal Function
func (ppm *PPM) drawKochCurve(n int, start Point, length int, color Pixel) {
	if n == 0 {
		// Base case: draw a straight line
		end := Point{start.X + length, start.Y}
		ppm.DrawLine(start, end, color)
	} else {
		// Recursive case: divide the line into four segments and draw Koch curves on each segment
		segmentLength := length / 3

		// Calculate the coordinates for the four points of the segments
		p1 := start
		p2 := Point{start.X + segmentLength, start.Y}
		p3 := Point{start.X + 2*segmentLength, start.Y}
		p4 := Point{start.X + 3*segmentLength, start.Y}

		// Calculate the height of the equilateral triangle formed by the middle two segments
		height := int(float64(segmentLength) * math.Sqrt(3) / 2)

		// Calculate the coordinates for the middle point of the middle segment
		middle := Point{start.X + segmentLength*2, start.Y + height}

		// Recursively draw Koch curves on the four segments
		ppm.drawKochCurve(n-1, p1, segmentLength, color)
		ppm.drawKochCurve(n-1, p2, segmentLength, color)
		ppm.drawKochCurve(n-1, middle, segmentLength, color)
		ppm.drawKochCurve(n-1, p3, segmentLength, color)
		ppm.drawKochCurve(n-1, p4, segmentLength, color)
	}
}

// DrawSierpinskiTriangle draws a Sierpinski triangle.
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

// DrawPerlinNoise dessine le bruit Perlin.
// cette fonction Dessine un bruit perlin de toute l'image.
func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	// Color1 est la couleur de 0.
	// Color2 est la couleur de 1.

	// Iterate over each pixel in the image
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Calculate the Perlin noise value for the current pixel
			noise := perlin.Noise2D(float64(x)/float64(ppm.width), float64(y)/float64(ppm.height))

			// Interpolate between color1 and color2 based on the noise value
			color := interpolateColor(color1, color2, noise)

			// Set the pixel color in the PPM image
			ppm.SetPixel(x, y, color)
		}
	}
}

// Additionnal Function
// interpolateColor interpolates between two colors based on a value between 0 and 1
func interpolateColor(color1 Pixel, color2 Pixel, t float64) Pixel {
	r := uint8(float64(color1.R)*(1-t) + float64(color2.R)*t)
	g := uint8(float64(color1.G)*(1-t) + float64(color2.G)*t)
	b := uint8(float64(color1.B)*(1-t) + float64(color2.B)*t)
	return Pixel{R: r, G: g, B: b}
}

// KNearestNeighbors redimensionne l'image PPM à l'aide de l'algorithme des k-voisins les plus proches.
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int) {
	// Calculate the scaling factors for width and height
	widthScale := float64(ppm.width) / float64(newWidth)
	heightScale := float64(ppm.height) / float64(newHeight)

	// Create a new PPM image with the new dimensions
	newPPM := NewPPM(newWidth, newHeight)

	// Iterate over each pixel in the new image
	for newY := 0; newY < newHeight; newY++ {
		for newX := 0; newX < newWidth; newX++ {
			// Calculate the corresponding pixel coordinates in the original image
			oldX := int(float64(newX) * widthScale)
			oldY := int(float64(newY) * heightScale)

			// Get the color of the nearest neighbor pixel in the original image
			color := ppm.GetPixel(oldX, oldY)

			// Set the pixel color in the new image
			newPPM.SetPixel(newX, newY, color)
		}
	}

	// Replace the original image with the resized image
	*ppm = *newPPM
}
*/
