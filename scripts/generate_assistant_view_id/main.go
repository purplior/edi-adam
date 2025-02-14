package main

import (
	"fmt"

	"github.com/purplior/sbec/lib/strgen"
)

func main() {
	suuid := strgen.ShortUniqueID()

	fmt.Println(suuid)
}
