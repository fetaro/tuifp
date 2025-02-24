# golang TUI (Terminal User Interface) File Picker 

Simple file picker through Terminal User Interface.

![](doc/tuifp.gif)

### Install

```
go get github.com/fetaro/tuifp
```

### Sample

```golang
package main

import (
	"fmt"
	"github.com/fetaro/tuifp"
)

func main() {
	fp := tuifp.NewTuiFilePicker()
	s, err := fp.Pick()
	if err != nil {
		panic(err)
	}
	fmt.Printf("picked file path = %s\n", s)
}
```
