// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func main() {

	// English sorting
	words := []string{"apple", "banana", "cherry"}
	cl := collate.New(language.English)
	cl.SortStrings(words)
	// Result: [apple banana cherry]

	// Polish sorting with special characters
	polishWords := []string{"Bug", "żółw", "zebra", "żabka", "zając", "ząb", "bąk", "bak", "buk", "bóg"}
	clPL := collate.New(language.Polish)
	clPL.SortStrings(polishWords)
	// Result respects Polish collation rules

	fmt.Printf("%v", polishWords)
}
