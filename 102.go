package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	a := "1c0111001f010100061a024b53535009181c"
	b := "686974207468652062756c6c277320657965"

	expectedRes := "746865206b696420646f6e277420706c6179"

	res := xorHexStrings(a, b)

	if string(res) != expectedRes {
		fmt.Printf("Expected %v, instead got %v\n", expectedRes, hex.EncodeToString(res))
	}
	fmt.Printf("Got encoded string %s\n", hex.EncodeToString(res))

}

func xorHexStrings(a string, b string) []byte {
	bytes1, _ := hex.DecodeString(a)
	bytes2, _ := hex.DecodeString(b)

	res := make([]byte, len(bytes1))

	for x := 0; x < len(bytes1); x++ {
		res[x] = bytes1[x] ^ bytes2[x]
	}
	return res
}
