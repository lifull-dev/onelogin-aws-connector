package cmd

import (
	"fmt"
	"os"
)

func errorExit(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(-1)
}
