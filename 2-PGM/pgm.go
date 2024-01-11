package main

type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// ReadPGM reads a PGM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {

}

// Size returns the width and height of the image.
func (pgm *PGM) Size() (int, int) {

}

// At returns the value of the pixel at (x, y).
func (pgm *PGM) At(x, y int) uint8 {

}

// Set sets the value of the pixel at (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {

}

// Save saves the PGM image to a file and returns an error if there was a problem.
func (pgm *PGM) Save(filename string) error {

}

// Invert inverts the colors of the PGM image.
func (pgm *PGM) Invert() {

}

// Flip flips the PGM image horizontally.
func (pgm *PGM) Flip() {

}

// Flop flops the PGM image vertically.
func (pgm *PGM) Flop() {

}

// SetMagicNumber sets the magic number of the PGM image.
func (pgm *PGM) SetMagicNumber(magicNumber string) {

}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {

}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {

}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {

}
