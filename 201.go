package main

/*
Implement PKCS#7 padding
A block cipher transforms a fixed-sized block (usually 8 or 16 bytes) of plaintext into ciphertext.
But we almost never want to transform a single block; we encrypt irregularly-sized messages.

One way we account for irregularly-sized messages is by padding, creating a plaintext that is an even
multiple of the blocksize. The most popular padding scheme is called PKCS#7.

So: pad any block to a specific block length, by appending the number
of bytes of padding to the end of the block. For instance,

"YELLOW SUBMARINE"
... padded to 20 bytes would be:

"YELLOW SUBMARINE\x04\x04\x04\x04"*/

import(
  "bytes"
  "fmt"
)

func padMessage(msg string, length int) []byte {
  var buffer bytes.Buffer
  buffer.WriteString(msg)
  size := length - len(msg)
  padding := make([]byte, size)
  for x := 0; x < size; x++ {
    padding[x] = byte(size)
  }
  buffer.WriteString(string(padding))
  return buffer.Bytes()
}

func main() {
  msg := "YELLOW SUBMARINE"
  fmt.Println(padMessage(msg, 20))
  fmt.Println(string(padMessage(msg, 20)))
}
