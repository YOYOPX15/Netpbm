package netpbm

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

// =============================================================================================
// ============================================ PBM ============================================
// =============================================================================================

// PBM represents a PBM image.
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	scanner.Scan()
	magicNumber := strings.TrimSpace(scanner.Text())

	// Read width and height
	var width, height int
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d %d", &width, &height)

	// Initialize PBM struct
	pbm := &PBM{
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		data:        make([][]bool, height),
	}

	// Read pixel data
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		for _, char := range line {
			pixelValue := char == '1'
			pbm.data[i] = append(pbm.data[i], pixelValue)
		}
	}

	return pbm, nil
}

// Size returns the width and height of the image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At returns the value of the pixel at (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save saves the PBM image to a file and returns an error if there was a problem.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write magic number, width, and height
	fmt.Fprintf(writer, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Write pixel data
	for _, row := range pbm.data {
		for _, pixel := range row {
			if pixel {
				fmt.Fprint(writer, "1")
			} else {
				fmt.Fprint(writer, "0")
			}
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

// Invert inverts the colors of the PBM image.
func (pbm *PBM) Invert() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width; x++ {
			pbm.data[y][x] = !pbm.data[y][x]
		}
	}
}

// Flip flips the PBM image horizontally.
func (pbm *PBM) Flip() {
	for y := 0; y < pbm.height; y++ {
		for x := 0; x < pbm.width/2; x++ {
			pbm.data[y][x], pbm.data[y][pbm.width-x-1] = pbm.data[y][pbm.width-x-1], pbm.data[y][x]
		}
	}
}

// Flop flops the PBM image vertically.
func (pbm *PBM) Flop() {
	for y := 0; y < pbm.height/2; y++ {
		pbm.data[y], pbm.data[pbm.height-y-1] = pbm.data[pbm.height-y-1], pbm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PBM image.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

// =============================================================================================
// ============================================ PGM ============================================
// =============================================================================================

// PGM represents a PGM image.
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	scanner.Scan()
	magicNumber := strings.TrimSpace(scanner.Text())

	// Read width, height, and max value
	scanner.Scan()
	var width, height int
	fmt.Sscanf(scanner.Text(), "%d %d", &width, &height)

	scanner.Scan()
	var maxValue int
	fmt.Sscanf(scanner.Text(), "%d", &maxValue)

	// Initialize PGM struct
	pgm := &PGM{
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         maxValue,
		data:        make([][]uint8, height),
	}

	// Read pixel data
	for i := 0; i < height; i++ {
		pgm.data[i] = make([]uint8, width)
		for j := 0; j < width; j++ {
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "%d", &pgm.data[i][j])
		}
	}

	return pgm, nil
}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write magic number, width, height, and max value
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", pgm.magicNumber, pgm.width, pgm.height, pgm.max)

	// Write pixel data
	for _, row := range pgm.data {
		for _, pixel := range row {
			fmt.Fprintf(writer, "%d ", pixel)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			pgm.data[y][x] = uint8(pgm.max - int(pgm.data[y][x]))
		}
	}
}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {
	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width/2; x++ {
			pgm.data[y][x], pgm.data[y][pgm.width-x-1] = pgm.data[y][pgm.width-x-1], pgm.data[y][x]
		}
	}
}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {
	for y := 0; y < pgm.height/2; y++ {
		pgm.data[y], pgm.data[pgm.height-y-1] = pgm.data[pgm.height-y-1], pgm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	rotatedData := make([][]uint8, pgm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	for y := 0; y < pgm.height; y++ {
		for x := 0; x < pgm.width; x++ {
			rotatedData[x][pgm.height-y-1] = pgm.data[y][x]
		}
	}

	pgm.width, pgm.height = pgm.height, pgm.width
	pgm.data = rotatedData
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
	pbm := &PBM{
		width:       pgm.width,
		height:      pgm.height,
		magicNumber: "P1", // Assuming binary PBM format
		data:        make([][]bool, pgm.height),
	}

	for y := 0; y < pgm.height; y++ {
		pbm.data[y] = make([]bool, pgm.width)
		for x := 0; x < pgm.width; x++ {
			pbm.data[y][x] = pgm.data[y][x] > uint8(pgm.max/2)
		}
	}

	return pbm
}

// =============================================================================================
// ============================================ PPM ============================================
// =============================================================================================

// PPM represents a PPM image.
type PPM struct {
	data          [][]Pixel
	width, height int
	magicNumber   string
	max           int
}

// Pixel represents a color pixel.
type Pixel struct {
	R, G, B uint8
}

// ReadPPM reads a PPM image from a file and returns a struct that represents the image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Read magic number
	scanner.Scan()
	magicNumber := strings.TrimSpace(scanner.Text())

	// Read width, height, and max value
	scanner.Scan()
	var width, height int
	fmt.Sscanf(scanner.Text(), "%d %d", &width, &height)

	scanner.Scan()
	var maxValue int
	fmt.Sscanf(scanner.Text(), "%d", &maxValue)

	// Initialize PPM struct
	ppm := &PPM{
		width:       width,
		height:      height,
		magicNumber: magicNumber,
		max:         maxValue,
		data:        make([][]Pixel, height),
	}

	// Read pixel data
	for i := 0; i < height; i++ {
		ppm.data[i] = make([]Pixel, width)
		for j := 0; j < width; j++ {
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "%d %d %d", &ppm.data[i][j].R, &ppm.data[i][j].G, &ppm.data[i][j].B)
		}
	}

	return ppm, nil
}

// Size returns the width and height of the image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At returns the value of the pixel at (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.data[y][x]
}

// Set sets the value of the pixel at (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	ppm.data[y][x] = value
}

// Save saves the PPM image to a file and returns an error if there was a problem.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Write magic number, width, height, and max value
	fmt.Fprintf(writer, "%s\n%d %d\n%d\n", ppm.magicNumber, ppm.width, ppm.height, ppm.max)

	// Write pixel data
	for _, row := range ppm.data {
		for _, pixel := range row {
			fmt.Fprintf(writer, "%d %d %d ", pixel.R, pixel.G, pixel.B)
		}
		fmt.Fprintln(writer)
	}

	return writer.Flush()
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			ppm.data[y][x].R = uint8(ppm.max) - ppm.data[y][x].R
			ppm.data[y][x].G = uint8(ppm.max) - ppm.data[y][x].G
			ppm.data[y][x].B = uint8(ppm.max) - ppm.data[y][x].B
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			ppm.data[y][x], ppm.data[y][ppm.width-x-1] = ppm.data[y][ppm.width-x-1], ppm.data[y][x]
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		ppm.data[y], ppm.data[ppm.height-y-1] = ppm.data[ppm.height-y-1], ppm.data[y]
	}
}

// SetMagicNumber sets the magic number of the PPM image.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue sets the max value of the PPM image.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = int(maxValue)
}

// Rotate90CW rotates the PPM image 90° clockwise.
func (ppm *PPM) Rotate90CW() {
	rotatedData := make([][]Pixel, ppm.width)
	for i := range rotatedData {
		rotatedData[i] = make([]Pixel, ppm.height)
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			rotatedData[x][ppm.height-y-1] = ppm.data[y][x]
		}
	}

	ppm.width, ppm.height = ppm.height, ppm.width
	ppm.data = rotatedData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	pgm := &PGM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2", // Assuming ASCII PGM format
		max:         ppm.max,
		data:        make([][]uint8, ppm.height),
	}

	for y := 0; y < ppm.height; y++ {
		pgm.data[y] = make([]uint8, ppm.width)
		for x := 0; x < ppm.width; x++ {
			// Convert color to grayscale using Luminosity method
			pgm.data[y][x] = uint8(0.299*float64(ppm.data[y][x].R) + 0.587*float64(ppm.data[y][x].G) + 0.114*float64(ppm.data[y][x].B))
		}
	}

	return pgm
}

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	pbm := &PBM{
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1", // Assuming binary PBM format
		data:        make([][]bool, ppm.height),
	}

	for y := 0; y < ppm.height; y++ {
		pbm.data[y] = make([]bool, ppm.width)
		for x := 0; x < ppm.width; x++ {
			// Convert color to binary using average intensity
			averageIntensity := (int(ppm.data[y][x].R) + int(ppm.data[y][x].G) + int(ppm.data[y][x].B)) / 3
			pbm.data[y][x] = averageIntensity > ppm.max/2
		}
	}

	return pbm
}

// Point represents a point in the image.
type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
		if p1.X > p2.X {
			p1, p2 = p2, p1
		}
		for x := p1.X; x <= p2.X; x++ {
			y := p1.Y + (x-p1.X)*(p2.Y-p1.Y)/(p2.X-p1.X)
			ppm.Set(x, y, color)
		}
	} else {
		if p1.Y > p2.Y {
			p1, p2 = p2, p1
		}
		for y := p1.Y; y <= p2.Y; y++ {
			x := p1.X + (y-p1.Y)*(p2.X-p1.X)/(p2.Y-p1.Y)
			ppm.Set(x, y, color)
		}
	}
}

// DrawRectangle draws a rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle draws a filled rectangle.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for y := p1.Y; y < p1.Y+height; y++ {
		for x := p1.X; x < p1.X+width; x++ {
			ppm.Set(x, y, color)
		}
	}
}

// DrawCircle draws a circle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	for angle := 0; angle < 360; angle++ {
		x := int(float64(center.X) + float64(radius)*math.Cos(float64(angle)))
		y := int(float64(center.Y) + float64(radius)*math.Sin(float64(angle)))
		ppm.Set(x, y, color)
	}
}

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for y := center.Y - radius; y <= center.Y+radius; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			if (x-center.X)*(x-center.X)+(y-center.Y)*(y-center.Y) <= radius*radius {
				ppm.Set(x, y, color)
			}
		}
	}
}

// DrawTriangle draws a triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle draws a filled triangle.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Barycentric coordinates
	for x := min(p1.X, min(p2.X, p3.X)); x <= max(p1.X, max(p2.X, p3.X)); x++ {
		for y := min(p1.Y, min(p2.Y, p3.Y)); y <= max(p1.Y, max(p2.Y, p3.Y)); y++ {
			alpha := float64((p2.Y-p3.Y)*(x-p3.X)+(p3.X-p2.X)*(y-p3.Y)) / float64((p2.Y-p3.Y)*(p1.X-p3.X)+(p3.X-p2.X)*(p1.Y-p3.Y))
			beta := float64((p3.Y-p1.Y)*(x-p3.X)+(p1.X-p3.X)*(y-p3.Y)) / float64((p2.Y-p3.Y)*(p1.X-p3.X)+(p3.X-p2.X)*(p1.Y-p3.Y))
			gamma := 1 - alpha - beta

			if alpha >= 0 && beta >= 0 && gamma >= 0 {
				ppm.Set(x, y, color)
			}
		}
	}
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon draws a filled polygon.
func minYmaxY(points []Point) (int, int) {
	minY := points[0].Y
	maxY := points[0].Y

	for _, point := range points {
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	return minY, maxY
}

// DrawFilledPolygon draws a filled polygon.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Filling polygon using scanline algorithm
	yMin, yMax := minYmaxY(points)
	activeEdges := make(map[int]geometry.Edge)

	for y := yMin; y <= yMax; y++ {
		activeEdges = updateActiveEdges(activeEdges, points, y)
		fillScanline(activeEdges, y, color)
	}
}

// DrawKochSnowflake draws a Koch snowflake.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
	// Drawing a Koch snowflake using recursive Koch curve function
	length := float64(width)
	angle := 60.0

	end1 := Point{start.X + int(length*cos(0)), start.Y - int(length*sin(0))}
	end2 := Point{start.X + int(length*cos((angle*math.Pi)/180.0)), start.Y - int(length*sin((angle*math.Pi)/180.0))}
	end3 := Point{start.X + int(length*cos((2*angle*math.Pi)/180.0)), start.Y - int(length*sin((2*angle*math.Pi)/180.0))}

	ppm.DrawKochCurve(n, start, end1, color)
	ppm.DrawKochCurve(n, end1, end2, color)
	ppm.DrawKochCurve(n, end2, end3, color)
	ppm.DrawKochCurve(n, end3, start, color)
}

// DrawSierpinskiTriangle draws a Sierpinski triangle.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	// Drawing a Sierpinski triangle using recursive function
	length := float64(width)
	angle := 60.0

	p2 := Point{start.X + int(length*cos((angle*math.Pi)/180.0)), start.Y - int(length*sin((angle*math.Pi)/180.0))}
	p3 := Point{start.X + int(length*cos((2*angle*math.Pi)/180.0)), start.Y - int(length*sin((2*angle*math.Pi)/180.0))}

	ppm.drawSierpinskiTriangleRecursive(n, start, p2, p3, color)
}

func (ppm *PPM) drawSierpinskiTriangleRecursive(n int, p1, p2, p3 Point, color Pixel) {
	if n == 0 {
		ppm.DrawTriangle(p1, p2, p3, color)
	} else {
		// Calculate midpoints of edges
		mid1 := midpoint(p1, p2)
		mid2 := midpoint(p2, p3)
		mid3 := midpoint(p3, p1)

		// Recursive calls
		ppm.drawSierpinskiTriangleRecursive(n-1, p1, mid1, mid3, color)
		ppm.drawSierpinskiTriangleRecursive(n-1, mid1, p2, mid2, color)
		ppm.drawSierpinskiTriangleRecursive(n-1, mid3, mid2, p3, color)
	}
}

// DrawPerlinNoise draws Perlin noise on the entire image.
func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	// Generating Perlin noise
	noise := generatePerlinNoise(ppm.width, ppm.height)

	// Scaling the noise values to fit the color range
	scale := func(value float64, fromMin float64, fromMax float64, toMin float64, toMax float64) float64 {
		return (value-fromMin)/(fromMax-fromMin)*(toMax-toMin) + toMin
	}

	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			// Map noise values to colors
			value := noise[x][y]
			col := Pixel{
				R: uint8(scale(value, -1, 1, float64(color1.R), float64(color2.R))),
				G: uint8(scale(value, -1, 1, float64(color1.G), float64(color2.G))),
				B: uint8(scale(value, -1, 1, float64(color1.B), float64(color2.B))),
			}
			ppm.Set(x, y, col)
		}
	}
}

// KNearestNeighbors resizes the PPM image using the k-nearest neighbors algorithm.
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int) {
	// Calculate scaling factors
	scaleX := float64(ppm.width) / float64(newWidth)
	scaleY := float64(ppm.height) / float64(newHeight)

	// Create a new image with the desired dimensions
	newImage := NewPPM(newWidth, newHeight)

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Find the nearest neighbor in the original image
			nearestX := int(float64(x) * scaleX)
			nearestY := int(float64(y) * scaleY)

			// Set the color of the pixel in the new image
			newImage.Set(x, y, ppm.At(nearestX, nearestY))
		}
	}

	// Update the original image with the resized image
	ppm.data = newImage.data
	ppm.width = newWidth
	ppm.height = newHeight
}
