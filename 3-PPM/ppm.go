package Netpbm

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
)

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

	img := &PPM{}
	scanner := bufio.NewScanner(file)

	// Read magic number
	if scanner.Scan() {
		img.magicNumber = scanner.Text()
	} else {
		return nil, errors.New("failed to read magic number")
	}

	// Read width, height, and max value
	if scanner.Scan() {
		line := scanner.Text()
		fmt.Sscanf(line, "%d %d", &img.width, &img.height)
	} else {
		return nil, errors.New("failed to read width and height")
	}

	if scanner.Scan() {
		line := scanner.Text()
		fmt.Sscanf(line, "%d", &img.max)
	} else {
		return nil, errors.New("failed to read max value")
	}

	// Read pixel data
	img.data = make([][]Pixel, img.height)
	for y := 0; y < img.height; y++ {
		img.data[y] = make([]Pixel, img.width)
		for x := 0; x < img.width; x++ {
			if scanner.Scan() {
				line := scanner.Text()
				var r, g, b uint8
				fmt.Sscanf(line, "%d %d %d", &r, &g, &b)
				img.data[y][x] = Pixel{R: r, G: g, B: b}
			} else {
				return nil, errors.New("failed to read pixel data")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return img, nil
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

	// Write magic number
	_, err = writer.WriteString(ppm.magicNumber + "\n")
	if err != nil {
		return err
	}

	// Write width, height, and max value
	_, err = fmt.Fprintf(writer, "%d %d\n", ppm.width, ppm.height)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(writer, "%d\n", ppm.max)
	if err != nil {
		return err
	}

	// Write pixel data
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			_, err = fmt.Fprintf(writer, "%d %d %d\n", pixel.R, pixel.G, pixel.B)
			if err != nil {
				return err
			}
		}
	}

	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Invert inverts the colors of the PPM image.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			pixel.R = 255 - pixel.R
			pixel.G = 255 - pixel.G
			pixel.B = 255 - pixel.B
			ppm.data[y][x] = pixel
		}
	}
}

// Flip flips the PPM image horizontally.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width/2; x++ {
			leftPixel := ppm.data[y][x]
			rightPixel := ppm.data[y][ppm.width-1-x]
			ppm.data[y][x] = rightPixel
			ppm.data[y][ppm.width-1-x] = leftPixel
		}
	}
}

// Flop flops the PPM image vertically.
func (ppm *PPM) Flop() {
	for y := 0; y < ppm.height/2; y++ {
		for x := 0; x < ppm.width; x++ {
			topPixel := ppm.data[y][x]
			bottomPixel := ppm.data[ppm.height-1-y][x]
			ppm.data[y][x] = bottomPixel
			ppm.data[ppm.height-1-y][x] = topPixel
		}
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
	// Create a new 2D slice to store the rotated image
	rotatedData := make([][]Pixel, ppm.width)
	for i := 0; i < ppm.width; i++ {
		rotatedData[i] = make([]Pixel, ppm.height)
	}

	// Rotate the image
	for y := 0; y < ppm.height; y++ {
		for x := 0; x < ppm.width; x++ {
			rotatedData[x][ppm.height-1-y] = ppm.data[y][x]
		}
	}

	// Update the width and height of the image
	ppm.width, ppm.height = ppm.height, ppm.width

	// Update the data with the rotated image
	ppm.data = rotatedData
}

// ToPGM converts the PPM image to PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Create a new PGM image with the same width and height as the PPM image
	pgm := &PGM{
		data:        make([][]uint8, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P2",
		max:         255,
	}

	// Convert each pixel from RGB to grayscale and assign it to the PGM image
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

// ToPBM converts the PPM image to PBM.
func (ppm *PPM) ToPBM() *PBM {
	// Create a new PBM image with the same width and height as the PPM image
	pbm := &PBM{
		data:        make([][]bool, ppm.height),
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: "P1",
	}

	// Convert each pixel to binary and assign it to the PBM image
	for y := 0; y < ppm.height; y++ {
		pbm.data[y] = make([]bool, ppm.width)
		for x := 0; x < ppm.width; x++ {
			pixel := ppm.data[y][x]
			// Check if the pixel is black or white
			isBlack := (pixel.R == 0 && pixel.G == 0 && pixel.B == 0)
			pbm.data[y][x] = isBlack
		}
	}

	return pbm
}

// Point represent a point in the image
type Point struct {
	X, Y int
}

// DrawLine draws a line between two points.
func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	// Determine the direction of the line
	var xInc, yInc int
	if dx > 0 {
		xInc = 1
	} else if dx < 0 {
		xInc = -1
	} else {
		xInc = 0
	}

	if dy > 0 {
		yInc = 1
	} else if dy < 0 {
		yInc = -1
	} else {
		yInc = 0
	}

	// Calculate the absolute differences
	dx = int(math.Abs(float64(dx)))
	dy = int(math.Abs(float64(dy)))

	// Determine the number of steps required
	steps := max(dx, dy)

	// Calculate the increments for each coordinate
	xInc *= dx / steps
	yInc *= dy / steps

	// Draw the line
	x := p1.X
	y := p1.Y
	for i := 0; i <= steps; i++ {
		ppm.Set(x, y, color)
		x += xInc
		y += yInc
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

// DrawFilledCircle draws a filled circle.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
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
	// Sort the points by y-coordinate
	points := []Point{p1, p2, p3}
	sort.Slice(points, func(i, j int) bool {
		return points[i].Y < points[j].Y
	})

	// Calculate the slopes of the edges
	slope1 := float64(points[1].X-points[0].X) / float64(points[1].Y-points[0].Y)
	slope2 := float64(points[2].X-points[0].X) / float64(points[2].Y-points[0].Y)

	// Initialize the x-coordinates of the edges
	x1 := float64(points[0].X)
	x2 := float64(points[0].X)

	// Iterate over each scanline
	for y := points[0].Y; y <= points[2].Y; y++ {
		// Draw a horizontal line between the x-coordinates
		ppm.DrawLine(Point{int(x1), y}, Point{int(x2), y}, color)

		// Update the x-coordinates based on the slopes
		x1 += slope1
		x2 += slope2
	}
}

// DrawPolygon draws a polygon.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// Iterate over each pair of consecutive points
	for i := 0; i < len(points)-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	// Connect the last point with the first point to close the polygon
	ppm.DrawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon draws a filled polygon.
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

// DrawKochSnowflake draws a Koch snowflake.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
	// N is the number of iterations.
	// Koch snowflake is a 3 times a Koch curve.
	// Start is the top point of the snowflake.
	// Width is the width of all the lines.
	// Color is the color of the lines.

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
	// N is the number of iterations.
	// Start is the top point of the triangle.
	// Width is the width all the lines.
	// Color is the color of the lines.
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

// DrawPerlinNoise draws perlin noise.
// this function Draw a perlin noise of all the image.
func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	// Color1 is the color of 0.
	// Color2 is the color of 1.

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

// interpolateColor interpolates between two colors based on a value between 0 and 1
func interpolateColor(color1 Pixel, color2 Pixel, t float64) Pixel {
	r := uint8(float64(color1.R)*(1-t) + float64(color2.R)*t)
	g := uint8(float64(color1.G)*(1-t) + float64(color2.G)*t)
	b := uint8(float64(color1.B)*(1-t) + float64(color2.B)*t)
	return Pixel{R: r, G: g, B: b}
}

// KNearestNeighbors resizes the PPM image using the k-nearest neighbors algorithm.
func (ppm *PPM) KNearestNeighbors(newWidth, newHeight int) {
	// Calculate the scaling factors
	scaleX := float64(ppm.width) / float64(newWidth)
	scaleY := float64(ppm.height) / float64(newHeight)

	// Create a new PPM image with the new dimensions
	resizedPPM := NewPPM(newWidth, newHeight)

	// Iterate over each pixel in the resized image
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			// Calculate the corresponding pixel coordinates in the original image
			srcX := int(float64(x) * scaleX)
			srcY := int(float64(y) * scaleY)

			// Get the nearest neighbor pixel from the original image
			color := ppm.GetPixel(srcX, srcY)

			// Set the pixel color in the resized image
			resizedPPM.SetPixel(x, y, color)
		}
	}

	// Replace the original image with the resized image
	*ppm = *resizedPPM
}
