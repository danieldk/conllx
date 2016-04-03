package conllx

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

func ExampleNewToken() {
	// Construct a token using the builder-pattern
	token := NewToken().SetForm("apples").SetLemma("apple").SetPosTag("noun")

	// Append the token to a sentence.
	var sent Sentence
	sent = append(sent, *token)

	fmt.Println(sent)
	// Output: 1	apples	apple	_	noun	_	_	_	_	_
}

func ExampleToken_Features() {
	input := `1	test	_	_	_	f1:v1|f2:v2`
	strReader := strings.NewReader(input)
	reader := NewReader(bufio.NewReader(strReader))

	sent, err := reader.ReadSentence()
	if err != nil {
		log.Fatal("Error reading sentence")
	}

	features, ok := sent[0].Features()
	if !ok {
		log.Fatal("Token should have features")
	}

	fmt.Println(features.FeaturesMap()["f1"])
	fmt.Println(features.FeaturesMap()["f2"])

	// Output:
	// v1
	// v2
}

func ExampleToken_SetFeatures() {

	token := NewToken().SetFeatures(map[string]string{
		"num":   "sg",
		"tense": "past",
	})

	features, _ := token.Features()

	fmt.Println(features.FeaturesMap()["num"])
	fmt.Println(features.FeaturesMap()["tense"])

	// Output:
	// sg
	// past
}
