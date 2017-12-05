package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func main() {
	s := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"

	expectedRes := "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"

	b, _ := hex.DecodeString(s)

	res := base64.StdEncoding.EncodeToString(b)

	fmt.Println(string(b))

	if res != expectedRes {
		fmt.Printf("Expected %s, instead got %s\n", expectedRes, res)
	}
	fmt.Printf("Got encoded string %s\n", res)
}
