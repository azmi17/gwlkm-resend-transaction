package helper

import (
	"math/rand"
	"strings"
	"time"
)

var (
	numbers    = "0123456789"
	allCharset = numbers
)

func GenerateSTAN() string {

	rand.Seed(time.Now().Unix())

	minSpecialChar := 0
	minNum := 0
	minUpperCase := 0
	textLength := 12

	randomId := generateNumStan(textLength, minSpecialChar, minNum, minUpperCase)

	var outPutsId = string(randomId)

	return outPutsId

}

func generateNumStan(textLength, minSpecialChar, minNum, minUpperCase int) string {

	var strBuilder strings.Builder

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numbers))
		strBuilder.WriteString(string(numbers[random]))
	}

	remainingLength := textLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharset))
		strBuilder.WriteString(string(allCharset[random]))
	}

	inRune := []rune(strBuilder.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}
