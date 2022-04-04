package main

import (
	"log"
	"net"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var password = "amogus"

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:9013")
	if err != nil {
		log.Fatal("Failed to listen on port!")
		return
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection!")
			continue
		}
		connectionHandler(conn)
	}
}

func connectionHandler(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(10 * time.Second))

	input, err := ReadLine(conn)
	if err != nil {
		log.Fatal("Failed to read password")
		conn.Close()
		return
	}

	if (input == password) {
		conn.Write([]byte("Welcome!\n"))
		for {
			conn.SetDeadline(time.Now().Add(time.Second * 30))

			userData, err := user.Current()
			if err != nil {
				conn.Write([]byte("Failed to fetch current user data, continuing..."))
			}

			hostname, err := os.Hostname()
			if err != nil {
				conn.Write([]byte("Failed to fetch hostname, continuing..."))
			}

			pwd, err := os.Getwd()
			if err != nil {
				conn.Write([]byte("Failed to fetch current working directory, continuing..."))
			}

			promptSymbol := "$"
			if userData.Gid == "0" {
				promptSymbol = "#"
			}

			prompt := userData.Username + "@" + hostname + ":" + pwd + promptSymbol + " "

			conn.Write([]byte(prompt))
			input, err := ReadLine(conn)
			if err != nil {
				conn.Write([]byte("Failed to read command, try again!"))
				return
			}

			if (input == "exit") {
				conn.Write([]byte("Bye!\n"))
				return
			}

			args := strings.Split(input, " ")

			output, err := exec.Command(args[0], args[1:]...).Output()
			if err != nil {
				conn.Write([]byte(err.Error()))
				conn.Write([]byte("\n"))
			} else {
				conn.Write(output)
			}
		}
	} else {
		conn.Close()
	}
	
}

func ReadLine(conn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	bufPos := 0

	for {
		n, err := conn.Read(buf[bufPos:bufPos + 1])
		if err != nil || n != 1 {
			return "", err
		}
		if buf[bufPos] == '\r' || buf[bufPos] == '\t' {
			bufPos--
		} else if buf[bufPos] == '\n' || buf[bufPos] == '\x00' {
			return string(buf[:bufPos]), nil
		}
		bufPos++
	}
}