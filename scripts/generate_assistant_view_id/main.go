package main

import (
	"fmt"

	"github.com/purplior/edi-adam/lib/strgen"
)

func main() {
	suuid := strgen.ShortUniqueID()

	fmt.Println(suuid)
}
