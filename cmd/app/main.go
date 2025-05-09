package main

import (
	"fmt"

	"github.com/finahdinner/tidal/internal/preferences"
)

func main() {
	// gui.Gui.App.Run()
	fmt.Print(preferences.GetPreferences())
	// pref, err := preferences.GetPreferences()
	// if err != nil {
	// 	fmt.Println("error!")
	// 	log.Fatal(err)
	// }
	// fmt.Println(pref)
}
