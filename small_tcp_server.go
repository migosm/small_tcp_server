package main

import (
    _ "fmt"
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
    "os"
    "log"
    "github.com/nutrun/lentil"
)

// Log handler
var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
    responseTemplate = "<{{.Tag}} num={{.Num}} />\n"
)

type xmlReq struct {
  XMLName xml.Name  `xml:"rfidVisibilityRequest"`
  Host string `xml:"host,attr"`
  Num int `xml:"num,attr"`
  EventTag string `xml:"eventTag,attr"`
}

type xmlResp struct {
  Tag string `xml:"rfidVisibilityResponse"`
  Num int `xml:"num,attr"`
}

//
func initLog(
    traceHandle io.Writer,
    infoHandle io.Writer,
    warningHandle io.Writer,
    errorHandle io.Writer) {

    Trace = log.New(traceHandle,
        "TRACE: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Info = log.New(infoHandle,
        "INFO: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Warning = log.New(warningHandle,
        "WARNING: ",
        log.Ldate|log.Ltime|log.Lshortfile)

    Error = log.New(errorHandle,
        "ERROR: ",
        log.Ldate|log.Ltime|log.Lshortfile)
}

func checkErr(err error) {
  if (err != nil) {
    Error.Println(err)
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

func handleConnection(conn net.Conn, storePath string, beanstalkdChan chan string) error {
  connectionReader := bufio.NewReader(conn)
  defer conn.Close()
  msg, err := connectionReader.ReadBytes('\n')
  if (err != nil) && (err != io.EOF)  {
    Error.Println(err)
    return err
  }
  Info.Println("Msg length: ", len(msg))
  if len(msg) == 0 {
    return errors.New("Don't do anything with zero length messages")
  }
  filename := storePath + "/data" + strconv.FormatInt(time.Now().Unix(), 10)
  beanstalkdChan <-filename
  ioutil.WriteFile(filename, msg, 0777)
  resp, err := handleXML(msg)
  Info.Println(string(resp))
  _, err = conn.Write(resp); checkErr(err)
  return nil
}

func handleBeanstalkd(msgChan chan string) {
  conn, err := lentil.Dial("0.0.0.0:11300")
  if err != nil {
    Error.Fatal(err)
  }
  conn.Use("rfid_data")
  Info.Println(conn)
  for {
    msg := <-msgChan
    Info.Println("Beanstalkd: ", msg)
    conn.Put(0, 0, 0, []byte(msg))
  }

}

func run(listenHost string, listenPort string, storePath string) {
  // Prepare Beanstalkd
  beanstalkdChan := make(chan string, 256)
  go handleBeanstalkd(beanstalkdChan)
  sock, err := net.Listen("tcp", listenHost + ":" + listenPort); checkErr(err)
  Info.Println(sock)
  for {
    conn, err := sock.Accept(); checkErr(err)
    go handleConnection(conn, storePath, beanstalkdChan)
  }
}

func main() {
  initLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
  var storePath, listenHost, listenPort string
  //var listenPort string
  flag.StringVar(&storePath, "path", "/home/forge/rfid_data", "")
  flag.StringVar(&listenHost, "host", "", "")
  flag.StringVar(&listenPort, "port", "8085", "")
  flag.Parse()
  run(listenHost, listenPort, storePath)
}
