/*The Base64-encoded content in this file has been encrypted via AES-128 in ECB mode under the key

"YELLOW SUBMARINE".
(case-sensitive, without the quotes; exactly 16 characters; I like "YELLOW SUBMARINE" because it's exactly 16 bytes long, and now you do too).

Decrypt it. You know the key, after all.

Easiest way: use OpenSSL::Cipher and give it AES-128-ECB as the cipher.*/

package main

import(
  "crypto/aes"
  "bufio"
  "encoding/base64"
  "fmt"
  "os"
  "log"
)

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
  encoded := ReadFile("107.txt")
	bytes, _ := base64.StdEncoding.DecodeString(string(encoded))
  key := []byte("YELLOW SUBMARINE")
  aesCipher, err := aes.NewCipher(key)
  if err != nil {
    fmt.Println("Error: ", err)
  }
  for i := 0; i < len(bytes); i += 16 {
		aesCipher.Decrypt(bytes[i:i+16], bytes[i:i+16])
	}
	fmt.Println(string(bytes))
}
