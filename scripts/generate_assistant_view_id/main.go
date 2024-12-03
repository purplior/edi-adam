package main

import (
	"fmt"

	"github.com/purplior/podoroot/lib/strgen"
)

func main() {
	suuid := strgen.ShortUniqueID()

	fmt.Println(suuid)
}
