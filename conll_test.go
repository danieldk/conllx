package conllx

import "fmt"

func ExampleNewToken() {
	// Construct a token using the builder-pattern
	token := NewToken().SetForm("apples").SetLemma("apple").SetPosTag("noun")

	// Append the token to a sentence.
	var sent Sentence
	sent = append(sent, *token)

	fmt.Println(sent)
	// Output: 1	apples	apple	_	noun	_	_	_	_	_

}
