package aes

import (
  "bytes"
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
)

// PKCS7Padding -
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// PKCS7UnPadding -
func PKCS7UnPadding (origData []byte) []byte {
  length := len(origData)
  unpadding := int(origData[length - 1])
  return origData[:(length - unpadding)]
}

// Encrypt -
func Encrypt (origData, key []byte) (enc []byte, err error) {
  block , err := aes.NewCipher(key)
  if err != nil {
    return
  }
  blockSize := block.BlockSize()
  origData = PKCS7Padding(origData, blockSize)
  ivByte := make([]byte, blockSize)
  _, err = rand.Read(ivByte)
  if err != nil {
    return
  }
  blockMode := cipher.NewCBCEncrypter(block, ivByte)
  enc = make([]byte, len(origData))
  blockMode.CryptBlocks(enc, origData)
  return
}

// Decrypt -
func Decrypt (encData, key []byte) (dec []byte, err error) {
  block, err := aes.NewCipher(key)
  if err != nil {
    return
  }
  blockSize := block.BlockSize()
  ivByte := encData[:blockSize]
  encData = encData[blockSize:]
  blockMode := cipher.NewCBCDecrypter(block, ivByte)
  dec = make([]byte, len(encData))
  blockMode.CryptBlocks(dec, encData)
  dec = PKCS7UnPadding(dec)
  return
}