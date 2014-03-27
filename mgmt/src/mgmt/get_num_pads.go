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
var verbose bool

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
    default_verbose = false
    host_usage = "the host where the api can be accessed"
    apikey_usage = "the API key"
    verbose_usage = "enable verbose output"
  )
  flag.StringVar(&host, "host", default_host, host_usage)
  flag.StringVar(&host, "h", default_host, host_usage+" (shorthand)")
  flag.StringVar(&apikey, "apikey", default_apikey, apikey_usage)
  flag.StringVar(&apikey, "k", default_apikey, apikey_usage+" (shorthand)")
  flag.BoolVar(&verbose, "verbose", default_verbose, verbose_usage)
  flag.BoolVar(&verbose, "v", default_verbose, verbose_usage+" (shorthand)")
}

func main() {
  flag.Parse()
  if (verbose) {
    fmt.Printf("Using host %s with api key %s\n", host, apikey)
  }
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
  var padList PadListMessage
  err = json.Unmarshal(json_padlist, &padList)
  if err != nil {
    log.Fatal("Cannot decode response: " + err.Error())
  } else {
    //fmt.Printf("All: %+v\n", padList)
    if (padList.Code != 0) {
      fmt.Printf("API error: %s (code %d)\n", padList.Message, padList.Code)
    } else {
      if (verbose) {
        fmt.Printf("Found pad IDs: %s\n", padList.Data.PadIDs)
      }
      fmt.Printf("Found %d pads.\n", len(padList.Data.PadIDs));
    }
  }
}

