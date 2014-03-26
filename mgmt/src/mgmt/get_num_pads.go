package main

import (
  "fmt"
	"flag"
	"io/ioutil"
	"log"
  "net/http"
  "encoding/json"
)

var host string
var apikey string

type PadListMessage struct {
  Code uint32 `json:"code"`
  Message string `json:"message"`
  Data struct {
    PadIDs []string `json:"padIDs"`
  } `json:"data"`
}

func init() {
  const (
    default_host = "localhost"
    default_apikey = ""
    host_usage = "the host where the api can be accessed"
    apikey_usage = "the API key"
  )
  flag.StringVar(&host, "host", default_host, host_usage)
  flag.StringVar(&host, "h", default_host, host_usage+" (shorthand)")
  flag.StringVar(&apikey, "apikey", default_apikey, apikey_usage)
  flag.StringVar(&apikey, "k", default_apikey, apikey_usage+" (shorthand)")
}

func main() {
  flag.Parse()
  fmt.Printf("Using host %s with api key %s\n", host, apikey)
  url := fmt.Sprintf("https://%s/api/1.2.7/listAllPads?apikey=%s", host, apikey);
  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  json_padlist, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("%s\n", json_padlist)
  var padList PadListMessage
  err = json.Unmarshal(json_padlist, &padList)
  if err != nil {
    log.Fatal("Cannot decode response: " + err)
  } else {
    fmt.Printf("Code: %d\n", padList.Code)
    fmt.Printf("Message: %s\n", padList.Message)
    fmt.Printf("Data: %s\n", padList.Data)
    fmt.Printf("IDs: %s\n", padList.Data.PadIDs)
    fmt.Printf("Found %d pads.\n", len(padList.Data.PadIDs));
    //fmt.Printf("All: %+v\n", padList)
  }
}

