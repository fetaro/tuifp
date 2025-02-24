package main

import (
	"fmt"
	tui_file_picker "github.com/fetaro/gcal_forcerun_go"
)

func main() {
	fp := tui_file_picker.NewTuiFilePicker()
	s, err := fp.Pick()
	if err != nil {
		panic(err)
	}
	fmt.Printf("picked file path = %s\n", s)
}
