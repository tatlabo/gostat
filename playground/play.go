// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

func main() {

	password := []byte("myPassword")
	// English sorting
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("%v\n", string(hash))

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		fmt.Printf("%#v", err)
	} else {
		fmt.Println("Ok!")
	}

}

func sortWords() {

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
