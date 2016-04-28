package main

import (
    "fmt"
    "net"
    "io"
    "io/ioutil"
    "time"
)

const(
  maxBufSize = 65535
  tmpBufSize = 256
)

func checkErr(err error) {
  if (err != nil) {
    fmt.Println(err)
  }
}

func handleConnection(conn net.Conn) {
  msg := make([]byte, tmpBufSize)
  buff := []byte{}
  defer conn.Close()
  for {
    n, err := conn.Read(msg)
    //fmt.Println(n, string(msg[:n]))
    if (err != nil) && (err != io.EOF)  {
      // Quit as we recieved io.EOF
      panic(err)
    }
    if n == 0 {
      break
    }
    buff = append(buff, msg[:n]...)
  }
//  fmt.Println(string(buff))
//  fmt.Println(len(buff))

// Write to file
  fmt.Println(time.Now())
  ioutil.WriteFile("test_file", buff, 0777)
}

func listen() {
  c, err := net.Listen("tcp", ":8085"); checkErr(err)
  fmt.Println(c)
  for {
    conn, err := c.Accept(); checkErr(err)
    go handleConnection(conn)
  }
}

func main() {
    listen()
}
