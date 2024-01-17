package Netpbm

// Pour le test personnel
/*
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

	// Convertit PGM to PBM
	pbmImage := pgmImage.ToPBM()

	// Affiche la taille image PBM
	pbmWidth, pbmHeight := pbmImage.width, pbmImage.height

	fmt.Printf("Taille Image PBM: %dx%d\n", pbmWidth, pbmHeight)
	fmt.Println("Data :")
	for _, row := range pgmImage.ToPBM().data {
		for _, pixel := range row {
			fmt.Printf("%d ", pixel)
		}
		fmt.Println()
	}

	fmt.Print("\n")
	fmt.Print("Image Modifié:\n")

	pgmImage.Invert()
	fmt.Println("Image Inversé Couleur")

	pgmImage.Flip()
	fmt.Println("Image Inversée Horizontalement")

	pgmImage.Flop()
	fmt.Println("Image Renversée Verticalement")

	pgmImage.Rotate90CW()
	fmt.Println("Image pivotée de 90° dans le sens des aiguilles d'une montre")

	// Save sauvegarde l'image modifié
	err = pgmImage.Save("modified_duck.pgm")
	if err != nil {
		fmt.Println("Erreur sauvegarde image modifié:", err)
		return
	}

	// Affiche l'image modifié
	fmt.Println("Data:")
	for _, row := range pgmImage.data {
		for _, pixel := range row {
			fmt.Printf("%d ", pixel)
		}
		fmt.Println()
	}
}
*/
