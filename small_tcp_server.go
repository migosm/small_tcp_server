package main

import (
    "fmt"
    "net"
    "io"
    "io/ioutil"
    "bufio"
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
  defer conn.Close()
  for {
    msg, err := bufio.NewReader(conn).ReadBytes('\n')
    if (err != nil) && (err != io.EOF)  {
      // Quit as we recieved io.EOF
      panic(err)
    }
    filename := storePath + "/data" + strconv.FormatInt(time.Now().Unix(), 10)
    fmt.Println(filename)
    ioutil.WriteFile(filename, msg, 0777)
  }
}

func run(listenHost string, listenPort string, storePath string) {
  c, err := net.Listen("tcp", listenHost + ":" + listenPort); checkErr(err)
  fmt.Println(c)
  for {
    conn, err := c.Accept(); checkErr(err)
    go handleConnection(conn, storePath)
  }
}

func main() {
  var storePath, listenHost, listenPort string
  //var listenPort string
  flag.StringVar(&storePath, "path", "/home/forge/rfid_data", "")
  flag.StringVar(&listenHost, "host", "localhost", "")
  flag.StringVar(&listenPort, "port", "8085", "")
  flag.Parse()
  run(listenHost, listenPort, storePath)
}
