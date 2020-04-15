package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	netListen, err := net.Listen("tcp", ":5000")
	CheckError(err)

	defer netListen.Close()

	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	allbuf := make([]byte, 0)
	buffer := make([]byte, 65535)
	for {
		readLen, err := conn.Read(buffer)
		//fmt.Println("readLen: ", readLen, len(allbuf))
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read error")
			return
		}

		if len(allbuf) != 0 {
			allbuf = append(allbuf, buffer...)
		} else {
			allbuf = buffer[:]
		}
		var readP int = 0
		for {
			//buffer长度小于6
			if readLen-readP < 6 {
				allbuf = buffer[readP:]
				break
			}
			// 包头 第 4,5位存的包长度
			msgLen := byte2int16(allbuf[readP+4 : readP+6])
			logLen := int(6 + msgLen)
			//fmt.Println(readP, readP+logLen)
			//buffer剩余长度>将处理的数据长度
			if len(allbuf[readP:]) >= logLen {
				//fmt.Println(string(allbuf[4:7]))
				// 第6位往后才是数据
				fmt.Println(string(allbuf[readP+6 : readP+logLen]))
				readP += logLen
				//fmt.Println(readP, readLen)
				if readP == readLen {
					allbuf = nil
					break
				}
			} else {
				allbuf = buffer[readP:]
				break
			}
		}
	}
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

//  [2]byte 转 uint16
func byte2int16(s []byte) uint16 {

	// []byte 转 数字
	var i uint16
	buf := bytes.NewBuffer(s)
	// 大端序列
	binary.Read(buf, binary.BigEndian, &i)
	return i
}
