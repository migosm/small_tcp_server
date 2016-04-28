package main

import (
    "fmt"
    "net"
    "io"
    "io/ioutil"
    "time"
    "flag"
    "strconv"
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

func handleConnection(conn net.Conn, storePath string) {
  msg := make([]byte, tmpBufSize)
  buff := []byte{}
  defer conn.Close()
  for {
    n, err := conn.Read(msg)
    if (err != nil) && (err != io.EOF)  {
      // Quit as we recieved io.EOF
      panic(err)
    }
    if n == 0 {
      break
    }
    buff = append(buff, msg[:n]...)
  }

// Write to file
  //fmt.Println(strconv.FormatInt(time.Now().Unix(), 10))
  filename := storePath + "/data" + strconv.FormatInt(time.Now().Unix(), 10)
  fmt.Println(filename)
  ioutil.WriteFile(filename, buff, 0777)
}

func listen(listenPort string, storePath string) {
  c, err := net.Listen("tcp", ":" + listenPort); checkErr(err)
  fmt.Println(c)
  for {
    conn, err := c.Accept(); checkErr(err)
    go handleConnection(conn, storePath)
  }
}

func main() {
  var storePath string
  var listenPort string
  flag.StringVar(&storePath, "path", "/home/forge/rfid_data", "")
  flag.StringVar(&listenPort, "port", "8085", "")
  flag.Parse()
  listen(listenPort, storePath)
}
