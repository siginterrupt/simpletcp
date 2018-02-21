package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Error listening: %s\n", err.Error())
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %s\n", err.Error())
		}
		go HandleRequest(conn)
	}
}

func HandleRequest(conn net.Conn) {
	defer conn.Close()
	log.Printf("Accepted connection from %s\n", conn.RemoteAddr())
	page := ReadRequest(conn)
	WriteResponse(conn, page)
}

func ReadRequest(conn net.Conn) string {
	line := 0
	page := ""
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if line == 0 {
			req := strings.Split(ln, " ")
			page = req[1]
		}
		if ln == "" {
			break
		}
		line++
	}
	if page == "/" {
		page = "/index.html"
	}
	return page
}

func WriteResponse(conn net.Conn, page string) {
	destPage := "public" + page
	content, err := ioutil.ReadFile(destPage)
	if err == nil {
		log.Println("Serving " + destPage)
		fmt.Fprintf(conn, "HTTP/1.1 200 OK\r\nContent-Length: %d\r\nContent-Type: text/html\r\n\r\n%s", len(content), content)
	} else {
		log.Println("Page not found")
		fmt.Fprint(conn, "HTTP/1.1 404 Page not found\r\nContent-Type: text/html\r\n\r\n")
	}
}
