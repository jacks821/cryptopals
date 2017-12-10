package main

import (
	"encoding/base64"
	"fmt"
	"math"
  "os"
  "log"
  "bufio"
  "bytes"
)

/*Decrypt it.

Here's how:

Let KEYSIZE be the guessed length of the key; try values from 2 to (say) 40.
Write a function to compute the edit distance/Hamming distance between two strings. The Hamming distance is just the number of differing bits. The distance between:
this is a test
and
wokka wokka!!!
is 37. Make sure your code agrees before you proceed.


For each KEYSIZE, take the first KEYSIZE worth of bytes, and the second KEYSIZE worth of bytes,
 and find the edit distance between them. Normalize this result by dividing by KEYSIZE.

The KEYSIZE with the smallest normalized edit distance is probably the key.
You could proceed perhaps with the smallest 2-3 KEYSIZE values.
Or take 4 KEYSIZE blocks instead of 2 and average the distances.

Now that you probably know the KEYSIZE: break the ciphertext into blocks of KEYSIZE length.
Now transpose the blocks: make a block that is the first byte of every block,
and a block that is the second byte of every block, and so on.

Solve each block as if it was single-character XOR. You already have code to do this.

For each block, the single-byte XOR key that produces the best looking histogram is the repeating-key
XOR key byte for that block. Put them together and you have the key.


This code is going to turn out to be surprisingly useful later on. Breaking repeating-key XOR ("Vigenere")
statistically is obviously an academic exercise,
a "Crypto 101" thing. But more people "know how" to break it than can actually break it,
and a similar technique breaks something much more important.*/

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
		sum += a[i] * b[i]
	}
	return sum
}

func lenVec(a []float64) float64 {
	return math.Sqrt(dotVec(a, a))
}

func cosine(a, b []float64) float64 {
	return dotVec(a, b) / (lenVec(a) * lenVec(b))
}

// scoreText returns integer representing how likely seq a to be a regular english text
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
	freqs := make([]float64, 26)
	for i, c := range cts {
		freqs[i] = float64(c) / amount
	}
	// fmt.Println(freqs)
	return cosine(freqs, idealFreqs)
}

// return most likely key for the sequence of bytes XORed with 1 byte
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

func findKey(data []byte, keySize int) ([]byte) {
	var buffer bytes.Buffer
	blocks := make([][]byte, keySize)
	for x := 0; x < keySize; x++ {
			blocks[x] = make([]byte, (len(data)/keySize) + 1)
	}
	for i, b := range data {
		fmt.Println(i%keySize, i/keySize)
		blocks[i%keySize][i/keySize] = b
	}
	for _, block := range blocks {
		key, _ := break1Xor(block)
		buffer.WriteString(string(key))
	}
	return buffer.Bytes()
}

func decipherText(data []byte, key []byte) []byte {
  var buffer bytes.Buffer
  for x := range data {
    buffer.WriteByte(data[x] ^ key[x%len(key)])
  }
  return buffer.Bytes()
}

func main() {
	encoded := ReadFile("106.txt")
	bytes, _ := base64.StdEncoding.DecodeString(string(encoded))

	fmt.Println("bytes:", len(bytes))
	fmt.Println("test: ", calcHamming([]byte("this is a test"), []byte("wokka wokka!!!")))
	size, _ := checkKey(bytes)
	key := findKey(bytes, size)
	fmt.Printf("key:\n%s\n", string(key))
	decoded := decipherText(bytes, key)
	fmt.Printf("decoded:\n%s\n", string(decoded))
}
