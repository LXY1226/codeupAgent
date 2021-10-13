package AliAgent

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"io"
)


const (
	keyMaterial = "\x74\x27\x04\x1C\x07\x30\x56\x0C\x00\x72\x16\x0C\x36\x32\x5A\x61"
)


func ReadPacket(reader io.Reader) (header, body string, err error) {
	b := make([]byte, 4)
	_, err = io.ReadFull(reader, b)
	if err != nil { return "", "", err }
	lenOfPacket := binary.BigEndian.Uint32(b)
	encryptPacket := make([]byte, lenOfPacket)
	_, err = io.ReadFull(reader, encryptPacket)
	if err != nil { return "", "", err }
	lenOfHeader := binary.BigEndian.Uint32(encryptPacket)
	header = Decrypt(encryptPacket[4:4+lenOfHeader])
	body = Decrypt(encryptPacket[4+lenOfHeader:])
	return
}

var aliCrypto = func() cipher.Block {
	block, _ := aes.NewCipher([]byte(keyMaterial))
	return block
}()

func Decrypt(data []byte) string {
	if len(data) == 0 {
		return ""
	}
	decrypted := make([]byte, len(data))

	for i := 0; i < len(data); i += aliCrypto.BlockSize() {
		aliCrypto.Decrypt(decrypted[i:], data[i:])
	}
	l := len(decrypted) - int(decrypted[len(decrypted)-1])
	decrypted = decrypted[:l]
	return string(decrypted)
}