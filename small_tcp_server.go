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
    "encoding/xml"
    "text/template"
    "bytes"
    "errors"
)

type xmlReq struct {
  XMLName xml.Name  `xml:"rfidVisibilityRequest"`
  Host string `xml:"host,attr"`
  Num int `xml:"num,attr"`
  EventTag string `xml:"eventTag,attr"`
}

var responseTemplate = "<{{.Tag}} num={{.Num}} />\n"

type xmlResp struct {
  Tag string `xml:"rfidVisibilityResponse"`
  Num int `xml:"num,attr"`
}

func checkErr(err error) {
  if (err != nil) {
    fmt.Println(err)
  }
}

func parseRequest(msg []byte) xmlReq {
  var xmlData xmlReq
  xml.Unmarshal(msg, &xmlData)
  return xmlData
}

func handleXML(inmsg []byte) ([]byte, error) {
  parsedXML := parseRequest(inmsg)
  if parsedXML.XMLName.Local == "" {
    return nil, errors.New("XMLName is not defined")
  }
  resp := xmlResp{Tag: "rfidVisibilityResponse", Num: parsedXML.Num}
  tmpl, err := template.New("xmlResponse").Parse(responseTemplate); checkErr(err)
  var buff bytes.Buffer
  err = tmpl.Execute(&buff, resp); checkErr(err)
  return buff.Bytes(), nil
}

func handleConnection(conn net.Conn, storePath string) {
  defer conn.Close()
  for {
    msg, err := bufio.NewReader(conn).ReadBytes('\n')
    if (err != nil) && (err != io.EOF)  {
      panic(err)
    }
    fmt.Println(string(msg))
    filename := storePath + "/data" + strconv.FormatInt(time.Now().Unix(), 10)
    fmt.Println(filename)
    ioutil.WriteFile(filename, msg, 0777)
    resp, err := handleXML(msg)
    fmt.Println(string(resp))
    _, err = conn.Write(resp); checkErr(err)
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
