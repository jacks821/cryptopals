package main

/*
One of the 60-character strings in 4.txt has been encrypted by single-character XOR.
Find it.
*/

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

var idealFreqs = []float32{.0817, .0149, .0278, .0425, .1270, .0223, .0202, .0609, .0697, .0015, .0077, .0402, .0241, .0675, .0751, .0193, .0009, .0599, .0633, .0906, .0276, .0098, .0236, .0015, .0197, .0007}

func xorByte(a []byte, k byte) []byte {
	res := make([]byte, len(a))
	for i := range a {
		res[i] = a[i] ^ k
	}
	return res
}

// scoreText returns integer representing how likely seq a to be a regular english text
func scoreText(a []byte) float32 {
	cts := make([]int, 26)
	for _, ch := range a {
		if 'A' <= ch && ch <= 'Z' {
			ch -= 32
		}
		if 'a' <= ch && ch <= 'z' {
			cts[int(ch)-'a']++
		}
	}
	amount := float32(len(a))
	var score float32
	freqs := make([]float32, 26)
	for i, num := range cts {
		freqs[i] = float32(num) / amount
		score += freqs[i]
	}
	return score
}

// return most likely key for the sequence of bytes XORed with 1 byte
func crack1Xor(a []byte) (float32, byte, []byte) {
	var maxScore float32
	var maxKey byte
	var maxDecoded []byte
	for k := 0; k <= 255; k++ {
		decoded := xorByte(a, byte(k))
		score := scoreText(decoded)
		if score > maxScore {
			maxScore = score
			maxKey = byte(k)
			maxDecoded = decoded
			// fmt.Println(k, score, string(decoded))
		}
	}
	return maxScore, maxKey, maxDecoded
}

func GrabLines(args string) string {
	var text string
	file, err := os.Open(args)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text = text + scanner.Text() + " "
	}
	return text
}

func main() {
	file, err := os.Open("104.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var maxScore float32
	var maxS string
	var encodedMax string

	for scanner.Scan() {
		inp := scanner.Text()
		// fmt.Println(inp)
		b, _ := hex.DecodeString(inp)
		score, _, decoded := crack1Xor(b)
		// s := hex.EncodeToString(c)
		if score > maxScore {
			maxS = string(decoded)
			maxScore = score
			encodedMax = inp
			fmt.Println(score, string(decoded), inp)
		}
	}
	fmt.Println(maxScore, maxS, encodedMax)
}
