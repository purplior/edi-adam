package main

import (
	"fmt"

	"github.com/podossaem/podoroot/lib/strgen"
)

func main() {
	suuid := strgen.ShortUniqueID()

	fmt.Println(suuid)
}
