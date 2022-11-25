package helper

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
	// lowerCharSet = "abcdefghij"
	upperCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberSet    = "0123456789"
	allCharSet   = upperCharSet + numberSet
)

func HashSha1PassByInput(pwd string) (encrypted string) {

	var sha = sha1.New()
	sha.Write([]byte(pwd))
	var hashing = sha.Sum(nil)
	var hashingPass = fmt.Sprintf("%x", hashing)
	encrypted = hashingPass

	return encrypted
}

func HashSha1Pass() (string, string) {

	rand.Seed(time.Now().Unix())

	minSpecialChar := 0
	minNum := 0
	minUpperCase := 1
	passwordLength := 8

	pwd := generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase)

	var sha = sha1.New()
	sha.Write([]byte(pwd))
	var hashing = sha.Sum(nil)

	var plainTextPass = string(pwd)
	var hashingPass = fmt.Sprintf("%x", hashing)

	// plainTextPass := string(pwd)
	// hashingPass := string(hashing)

	return plainTextPass, hashingPass
}

func generatePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {

	var password strings.Builder

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}

	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})

	return string(inRune)
}
