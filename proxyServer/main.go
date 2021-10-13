package main

import (
	"encoding/binary"
	"github.com/LXY1226/codeupAgent/AliAgent"
	"io"
	"log"
	"net"
)

var lAddr = &net.TCPAddr{Port: 8966}
var rAddr = &net.TCPAddr{IP: net.IP{182,92,29,39}, Port: 8000}
func main() {
	log.SetFlags(log.Ltime)
	listener, err := net.ListenTCP("tcp4", lAddr)
	log.Println("监听于", lAddr)
	if err != nil {
		panic(err)
	}
	for {
		child, err := listener.AcceptTCP()
		if err != nil {
			panic(err)
		}
		master, err := net.DialTCP("tcp", nil, rAddr)
		if err != nil {
			log.Panicln("无法连接Remote", err)
		}
		log.Println(child.RemoteAddr(),  "->", rAddr)
		go decryptAndSend("↑", child, master)
		go decryptAndSend("↓", master, child)
	}
}

func decryptAndSend(TAG string, src io.ReadCloser, dst io.WriteCloser) {
	defer src.Close()
	defer dst.Close()
	for {
		header, body, err := ProxyPacket(src, dst)
		if err != nil {
			log.Println(TAG, err)
			return
		}
		log.Println(TAG, "H", header)
		log.Println(TAG, "B", body)
	}
}


func ProxyPacket(src io.ReadCloser, dst io.WriteCloser) (header, body string, err error) {
	b := make([]byte, 4)
	_, err = io.ReadFull(src, b)
	if err != nil { return "", "", err }
	_, err = dst.Write(b)
	if err != nil { return "", "", err }
	lenOfPacket := binary.BigEndian.Uint32(b)
	encryptPacket := make([]byte, lenOfPacket)
	_, err = io.ReadFull(src, encryptPacket)
	if err != nil { return "", "", err }
	_, err = dst.Write(encryptPacket)
	if err != nil { return "", "", err }
	lenOfHeader := binary.BigEndian.Uint32(encryptPacket)
	header = AliAgent.Decrypt(encryptPacket[4:4+lenOfHeader])
	body = AliAgent.Decrypt(encryptPacket[4+lenOfHeader:])
	return

}