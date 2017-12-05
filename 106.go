package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"os"
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
			if score > 0.5 {
				//fmt.Println(k, score, string(decoded))
			}
		}
	}
	return maxKey, maxDecoded
}

func calcHamming(a, b []byte) int {
	al := len(a)
	bl := len(b)

	if al != bl {

		fmt.Println("strings are not equal (len(a)=%d, len(b)=%d)", al, bl)
		return -1
	}

	var hamming = 0

	for i := range a {
		r := a[i] ^ b[i]
		for r > 0 {
			if r&1 == 1 {
				hamming++
			}
			r = r >> 1
		}
	}
	return hamming
}

func checkKey(text []byte) (int, float64) {
	likelyKey := 40
	lowHam := 10000.00
	numBlocks := len(text) / 40
	for key := 1; key < 40; key++ {
		ham := 0.0
		for i := 0; i < numBlocks; i++ {
			a := i * key
			b := (i + 1) * key
			c := (i + 2) * key
			ham += float64(calcHamming(text[a:b], text[b:c])) / float64(key)
		}
		ham /= float64(numBlocks)
		if ham < lowHam {
			lowHam = ham
			likelyKey = key
			fmt.Printf("Bits: %v Size: %v\n", lowHam, likelyKey)
		}
	}
	fmt.Printf("Bits: %v Size: %v\n", lowHam, likelyKey)
	return likelyKey, lowHam
}

func ReadFile(filename string) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	scan := bufio.NewScanner(file)
	bytes := make([]byte, 0, 1024)

	for scan.Scan() {
		input := []byte(scan.Text())
		bytes = append(bytes, input...)
	}
	return bytes
}

func main() {
	encoded := ReadFile("106.txt")
	bytes, _ := base64.StdEncoding.DecodeString(string(encoded))

	a := "this is a test"
	b := "wokka wokka!!!"

	fmt.Println("Test: ", calcHamming([]byte(a), []byte(b)))

	keySize, ham := checkKey(bytes)
	fmt.Println(keySize, ham)

}
