package main
import (
  "net"
  "context"
  "bytes"
  "log"
  "net/http"
  "time"
  "sync"
  "strconv"
  "strings"
  "io"
  "io/ioutil"
  "flag"
  "encoding/json"
)


// Generic Steem API call structure.
type Call struct {
  Jsonrpc string        `json:"jsonrpc"`
  Id      string        `json:"id"`
  Method  string        `json:"method"`
  Params  []interface{} `json:"params"`
}


// Call get block on the upstream.
func get_block(client *http.Client, url string, httpmethod string, i int, verbose bool) int {
  mark := time.Now()

  // Build request.
  params := []interface{}{"condenser_api", "get_block", []interface{}{strconv.Itoa(i)}}
  requestJson, err := json.Marshal(&Call{ Jsonrpc:"2.0", Id : "0", Method: "call", Params: params})

  // Submit request.
  req, err := http.NewRequest(httpmethod, url, bytes.NewBuffer(requestJson))
  req.Header.Set("Content-Type", "application/json")
  resp, err := client.Do(req)
  if err != nil {
    if verbose {  log.Println(err)  }
    return 0
  }

  // Decode response as a json.
  var iface interface{}
  rdecoder := json.NewDecoder(resp.Body)
  err = rdecoder.Decode(&iface)
  if err != nil {
    if verbose {  log.Println(err)  }
    return 0
  }
  resultJson := iface.(map[string]interface{})

  // Check for upstream json error.
  if _, found := resultJson["error"]; found {
    log.Println(resultJson)
    return 0
  }

  // Clean-up and stats.
  elapsed := time.Since(mark);
  if rdecoder.More() {  io.Copy(ioutil.Discard, resp.Body)  }
  if verbose {  log.Printf("Block %d: %s ", i, elapsed)  }
  return 1
}


func main() {
  url     := flag.String("u", "http://127.0.0.1:8090", "url")
  blocks  :=    flag.Int("b", 4096,                    "number of blocks to get")
  wgsize  :=    flag.Int("c", 64,                      "concurrency")
  method  := flag.String("m", "POST",                  "http method")
  verbose :=   flag.Bool("v", false,                   "verbose blocks")
  flag.Parse()

  var client *http.Client
  var httpurl = "http://unix"

  // Build http client for unix socket or tcp socket.
  if strings.HasPrefix(*url, "unix:") {
    client = &http.Client {
      Transport: &http.Transport{
        DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
          return net.Dial("unix", strings.TrimPrefix(*url, "unix:"))
        },
      },
    }
  } else {
    client = &http.Client {
      Transport: &http.Transport{
            MaxIdleConns:       100000,
            IdleConnTimeout:    30 * time.Second,
        MaxIdleConnsPerHost: 100000,
      },
    }
    httpurl = *url
  }

  // Set up concurrency pool.
  test_loop := (*blocks)/(*wgsize)
  succ_ch := make(chan int, *blocks)
  var wg sync.WaitGroup
  wg.Add(*wgsize)

  mark := time.Now()
  for i := 0; i < *wgsize; i++ {
    go func(i int) {
      defer wg.Done()
      succ := 0
      // Split blocks by thread id and job index.
      for j := 0; j < test_loop; j++ {
        succ += get_block(client, httpurl, *method, 10000000 + i*test_loop + j, *verbose)
      }
      log.Printf("%d / %d", succ, test_loop)
      succ_ch <- succ
    } (i)
	}
  wg.Wait()
  elapsed := time.Since(mark);

  // Sum successes and print stats.
  close(succ_ch)
  total := 0;  for val := range succ_ch {  total += val  }
  log.Printf("Wall clock time: %s. Successes: %d/%d", elapsed, total, *blocks)
  log.Printf("Blocks per second: %d", (time.Duration(total) * time.Second) /(elapsed))
}
