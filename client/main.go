package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {

	server := "127.0.0.1:5000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	for i := 0; i < 50; i++ {
		//msg := strconv.Itoa(i)
		msg := RandString(i)
		//msgLen := fmt.Sprintf("%03s", strconv.Itoa(len(msg)))
		//fmt.Println(msg, msgLen)
		//words := "aaaa" + msgLen + msg
		msgLen := int2byte(int64(len(msg)))
		//words := append([]byte("aaaa"), []byte(msgLen), []byte(msg))
		fmt.Println(msgLen)
		data := append([]byte(`aaaa`), msgLen...)
		fmt.Println(string(append(data, []byte(msg)...)))
		//  发送的前 6 byte 是包头， 包头后2位为长度，需要转为 uint16, 最大 65535-6， 后面才是数据
		conn.Write(append(data, []byte(msg)...))
	}
}

/**
*生成随机大写字符
**/
func RandString(length int) string {
	rand.Seed(time.Now().UnixNano())
	rs := make([]byte, length)
	for start := 0; start < length; start++ {
		rs[start] = byte(rand.Int31n(26) + 65)
	}
	return string(rs)
}

// int64 转 []byte (0-65535)
func int2byte(l int64) []byte {

	s1 := make([]byte, 2)
	buf := bytes.NewBuffer([]byte{})

	// 数字转 []byte, 网络字节序为大端字节序
	err := binary.Write(buf, binary.BigEndian, &l)
	if err != nil {
		panic(err)

	}
	//buf.Reset()
	s1[0] = buf.Bytes()[6]
	s1[1] = buf.Bytes()[7]
	return s1
}
