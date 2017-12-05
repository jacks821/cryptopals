package main

import (
	"encoding/hex"
	"fmt"
)

func xorRepeatKey(inp string, key string) []byte {
	n := len(inp)
	res := make([]byte, len(inp))

	for x := 0; x < n; x++ {
		res[x] = inp[x] ^ key[x%len(key)]
	}
	return res
}

func main() {
	inp := "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	key := "ICE"

	expectedRes := "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"

	res := xorRepeatKey(inp, key)

	if hex.EncodeToString(res) != expectedRes {
		fmt.Printf("Expected %v, instead got %v\n", expectedRes, hex.EncodeToString(res))
	}
	fmt.Printf("Got encoded string %s\n", hex.EncodeToString(res))

}
