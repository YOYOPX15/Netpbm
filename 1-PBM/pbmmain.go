package Netpbm

// Pour le test personnel
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
