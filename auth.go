package catchpointsdk

import (
  "net/http"
  "github.com/kelseyhightower/envconfig"
  "encoding/json"
  "fmt"
  "log"
  "bytes"
  "encoding/base64"
  "github.com/syndtr/goleveldb/leveldb"
  "time"
  "strings"
  "strconv"
)

const catchpointTokenExpireTime int = 1800

type Authentication struct {
  ClientID string
  ClientSecret  string
  Endpoint  string
}

type AuthResponse struct {
  AccessToken string `json:"access_token"`
  TokenType string `json:"token_type"`
  ExpiresIn int `json:"expires_in"`
}


/*
Returns the base64 encoded Catchpoint auth token
*/
func Authenticate() string {
  return checkToken()
}

/*
HTTP request to get a new Catchpoint token; base64 encode the result
*/
func authToCatchpoint() (bearerToken string, accessToken int) {
  fmt.Printf("calling catchpoint")
  var a Authentication
  err := envconfig.Process("catchpointsdk", &a)
  if err != nil { log.Fatal(err.Error()) }
  b := []byte(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", a.ClientID, a.ClientSecret))
  req, err := http.NewRequest("POST", a.Endpoint, bytes.NewBuffer(b))
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("Accept", "*/*")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil { panic(err) }
  defer resp.Body.Close()

  if resp.Status == "200 OK" {
    decoder := json.NewDecoder(resp.Body)
    var ar AuthResponse
    _ = decoder.Decode(&ar)
    bearerToken := base64.StdEncoding.EncodeToString([]byte(ar.AccessToken))
    //log.Debug("Received Bearer token: %s - [base64 encoded] %s", ar.AccessToken, bearerToken)
    return bearerToken, int(time.Now().Unix())
  } else {
    log.Fatal("Invalid response to authentication request received from Catchpoint")
    panic(err)
  }
}

/*
Check to see if a valid Bearer token exists within the LevelDB store; if it
does, return it, else authenticate to Catchpoint and store the token
*/
func checkToken() (token string) {
  db, err := leveldb.OpenFile("catchpoint.state", nil)
  if err != nil { panic(err) }
  defer db.Close()
  data, err := db.Get([]byte("token"), nil)
  // no stored token found in LevelDB
  if err != nil {
    token, timeint := authToCatchpoint()
    err = db.Put([]byte("token"), []byte(fmt.Sprintf("%v::%s", timeint, token)), nil)
    if err != nil { panic(err) }
  // there is a stored token; let's see if it's less than 1800 seconds old
  } else {
    s := strings.Split(string(data), "::")
    timestring, token := s[0], s[1]
    timeint, _ := strconv.Atoi(timestring)
    diff := int(time.Now().Unix()) - timeint
    // ask for a new token if we're within 10 seconds of token expiry
    if diff > (catchpointTokenExpireTime - 10) {
      token, timeint = authToCatchpoint()
      err = db.Put([]byte("token"), []byte(fmt.Sprintf("%v::%s", timeint, token)), nil)
    }
  }
  return token
}
