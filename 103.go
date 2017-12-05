package main

import (
	"encoding/hex"
	"fmt"
	"math"
)

var idealFreqs = []float64{
	.0817, .0149, .0278, .0425, .1270, .0223, .0202, .0609, .0697, .0015, .0077, .0402, .0241,
	.0675, .0751, .0193, .0009, .0599, .0633, .0906, .0276, .0098, .0236, .0015, .0197, .0007}

func xorByte(a []byte, k byte) []byte {
	res := make([]byte, len(a))
	for i := range a {
		res[i] = a[i] ^ k
	}
	return res
}

func dotVec(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += a[i] * a[i]
	}
	return sum
}

func lenVec(a []float64) float64 {
	return math.Sqrt(dotVec(a, a))
}

func cosine(a, b []float64) float64 {
	return dotVec(a, b) / (lenVec(a) * lenVec(b))
}

func scoreText(a []byte) float64 {
	cts := make([]int, 26)
	for _, ch := range a {
		if 'A' <= ch && ch <= 'Z' {
			ch -= 32
		}
		if 'a' <= ch && ch <= 'z' {
			cts[int(ch)-'a']++
		}
	}
	amount := float64(len(a))
	score := 0.0
	freqs := make([]float64, 26)
	for i, c := range cts {
		freqs[i] = float64(c) / amount
		score += freqs[i]
	}
	return cosine(freqs, idealFreqs)
}

func break1Xor(a []byte) (byte, []byte) {
	var maxScore float64
	var maxKey byte
	var maxDecoded []byte
	for k := 0; k <= 255; k++ {
		decoded := xorByte(a, byte(k))
		score := scoreText(decoded)
		if score > maxScore {
			maxScore = score
			maxKey = byte(k)
			maxDecoded = decoded
			fmt.Println(k, score, string(decoded))
		}
	}
	return maxKey, maxDecoded
}

func main() {
	inp := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	b, _ := hex.DecodeString(inp)

	key, decoded := break1Xor(b)
	fmt.Println(key, string(decoded))
}
