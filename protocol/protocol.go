package protocol

import (
	"bufio"
	"crypto/sha256"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/MayurUbarhande0/Simple-cloud-storage/IO"
	"github.com/MayurUbarhande0/Simple-cloud-storage/gateway/helper"
)

type Server struct {
	listenAddr string
	ln         net.Listener
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	return s.Acceptloop()

}

func (s *Server) Acceptloop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
		go s.HandleTCPHandshake(conn)
	}
}
func (s *Server) HandleTCPHandshake(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn) //Using this to directly pass the incoming bytes to compress func to avoid Ram and Disk heating
	EN_KEY := os.Getenv("ENCRYPTION_KEY")
	AUTH_KEY := os.Getenv("AUTH_KEY")
	hash_key := sha256.Sum256([]byte(AUTH_KEY))
	//Goingt to pass the token as dynamic size
	sizeBuffer := make([]byte, 1)
	if _, err := io.ReadFull(reader, sizeBuffer); err != nil {
		return
	}
	dynamicTokenSize := int(sizeBuffer[0])

	token_buff := make([]byte, dynamicTokenSize)
	if _, err := io.ReadFull(reader, token_buff); err != nil {
		return
	}
	//Decrypting AUTH token to do hash comparing
	DE_Token, err := helper.Decrypt(token_buff, []byte(EN_KEY))
	if err != nil {
		return
	}
	if sha256.Sum256(DE_Token) != hash_key {
		return
	}
	//as the size of the handshake is variable , see handshake string at Server/sp.go
	//we are going to read till \n
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSuffix(metadata, "\n")
	parts := strings.Split(metadata, "|")
	if len(parts) < 3 {
		return
	}
	command := parts[0]
	storagePath := parts[1]
	fileSize, _ := strconv.ParseInt(parts[2], 10, 64)
	if command == "0x01" {
		IO.Compress(reader, storagePath, fileSize)
	}
}
