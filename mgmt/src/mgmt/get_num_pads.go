package main

import (
  "fmt"
	"flag"
	//"io/ioutil"
	//"log"
  //"net/http"
)

var host string
var apikey string

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
//	res, err := http.Get("http://www.google.com/robots.txt")
//	if err != nil {
//		log.Fatal(err)
//	}
//	robots, err := ioutil.ReadAll(res.Body)
//	res.Body.Close()
//	if err != nil {
//		log.Fatal(err)
//	}
//  //fmt.Printf("%s", robots)
  flag.Parse()
  fmt.Printf("Using host %s with api key %s", host, apikey)

}

