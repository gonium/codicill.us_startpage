/***
 * Copyright © 2014 Mathias Dalheimer. All Rights Reserved.
 * 
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 * 
 * 1. Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 * 
 * 2. Redistributions in binary form must reproduce the above copyright
 * notice, this list of conditions and the following disclaimer in the
 * documentation and/or other materials provided with the distribution.
 * 
 * 3. The name of the author may not be used to endorse or promote
 * products derived from this software without specific prior written
 * permission.
 * 
 * THIS SOFTWARE IS PROVIDED BY [LICENSOR] "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 **/

/***
 * Please note: This is my first program using the Go language. It may
 * contain bugs and/or explode. In any case, don't use it as an example
 * for good Go code.
 *
 * Build a binary using
 * $ go build src/mgmt/get_num_pads.go
 */

package main

import (
  "fmt"
	"flag"
	"io/ioutil"
	"log"
	"time"
	"os"
	"errors"
  "net/http"
  "encoding/json"
)

var host string
var apikey string
var verbose bool
var statsfile string

type PadListMessage struct {
  Code uint32 `json:"code"`
  Message string `json:"message"`
  Data struct {
    PadIDs []string `json:"padIDs"`
  } `json:"data"`
}

type PadDateCount struct {
    Date int64 `json:"date"`
    NumPads int `json:"numpads"`
}

type PadCountFileFormat struct {
  PadStats []PadDateCount `json:"padstats"`
}

func init() {
  const (
    default_host = "localhost"
    default_apikey = ""
    default_verbose = false
    default_statsfile = "padcount.json"
    host_usage = "the host where the api can be accessed"
    apikey_usage = "the API key"
    verbose_usage = "enable verbose output"
    statsfile_usage = "the data file to append current pad count to"
  )
  flag.StringVar(&host, "host", default_host, host_usage)
  flag.StringVar(&host, "h", default_host, host_usage+" (shorthand)")
  flag.StringVar(&apikey, "apikey", default_apikey, apikey_usage)
  flag.StringVar(&apikey, "k", default_apikey, apikey_usage+" (shorthand)")
  flag.BoolVar(&verbose, "verbose", default_verbose, verbose_usage)
  flag.BoolVar(&verbose, "v", default_verbose, verbose_usage+" (shorthand)")
  flag.StringVar(&statsfile, "statsfile", default_statsfile, statsfile_usage)
  flag.StringVar(&statsfile, "d", default_statsfile, statsfile_usage+" (shorthand)")
}

func retrievePadList () (PadListMessage, error) {
  var padList PadListMessage
  url := fmt.Sprintf("https://%s/api/1.2.7/listAllPads?apikey=%s", host, apikey);
  res, err := http.Get(url)
  if err != nil {
    return padList, errors.New("Failed during API GET request: " +
    err.Error())
  }
  json_padlist, err := ioutil.ReadAll(res.Body)
  res.Body.Close()
  if err != nil {
    return padList, errors.New("Failed during API GET request: " +
    err.Error())
  }
  err = json.Unmarshal(json_padlist, &padList)
  if err != nil {
    return padList, errors.New("Failed during API Response JSON decoding: " +
    err.Error())
  } 
  return padList, nil
}

func main() {
  flag.Parse()
  if (verbose) {
    fmt.Printf("Using host %s with api key %s\n", host, apikey)
    fmt.Printf("Storing stats in file %s\n", statsfile)
  }
  padList, err := retrievePadList();
  if err != nil {
    log.Fatal("Error retrieving pad list: " + err.Error());
  } else {
    //fmt.Printf("All: %+v\n", padList)
    if (padList.Code != 0) {
      fmt.Printf("API error: %s (code %d)\n", padList.Message, padList.Code)
    } else {
      if (verbose) {
        fmt.Printf("Found pad IDs: %s\n", padList.Data.PadIDs)
        fmt.Printf("Found %d pads.\n", len(padList.Data.PadIDs));
      }
      newentry := PadDateCount{time.Now().Unix(), len(padList.Data.PadIDs)}

      stats := PadCountFileFormat{}
      // Now: Check if the stats file already exists
      if _, err := os.Stat(statsfile); os.IsNotExist(err) {
        fmt.Printf("Creating new stats file %s \n", statsfile)
        stats.PadStats = make([]PadDateCount, 0)
      } else {
        // Read already stored content 
        content, err := ioutil.ReadFile(statsfile)
        if  err != nil {
          log.Fatal("Cannot read stats file: " + err.Error());
        } 
        if err = json.Unmarshal(content, &stats); err != nil {
          log.Fatal("Cannot parse stats file: " + err.Error());
        }
      }
      // Here, we either have an empty stats datastructure, or one 
      // populated by the contents of the stats file.
      // now: Append the new entry to the existing stats
      stats.PadStats = append(stats.PadStats, newentry)
      // Store the stats
      serializedFileContent, err := json.Marshal(stats);
      if err != nil {
        log.Fatal("Cannot serialize data: " + err.Error())
      } else {
        // Writing to file
        f, err := os.Create(statsfile)
        if err != nil {
          log.Fatal("Cannot create data file: " + err.Error())
        } else {
          defer f.Close();
          if _, err = f.Write(serializedFileContent); err != nil {
            log.Fatal("Failed to write to data file: " + err.Error())
          }
        }
      }
    }
  }
}

