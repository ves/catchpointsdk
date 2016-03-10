package catchpointsdk

import (
  "net/http"
  "encoding/json"
  "fmt"
  log "github.com/apex/log"
  "bytes"
  "encoding/base64"
  "github.com/syndtr/goleveldb/leveldb"
  "time"
  "strings"
  "strconv"
  "github.com/kelseyhightower/envconfig"
)

const catchpointTokenExpireTime int = 1800

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
  log.Debugf("Calling catchpoint for a new authtoken")
  envconfig.Process("catchpointsdk", &c)
  payload := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", c.ClientID, c.ClientSecret)
  uri := fmt.Sprintf("%s/ui/api/token", c.Endpoint)
  fmt.Println(c.Endpoint)
  b := []byte(payload)
  log.Debugf("Sending Catchpoint the following payload: %s", payload)
  req, err := http.NewRequest("POST", uri, bytes.NewBuffer(b))
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Add("Accept", "*/*")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil { panic(err) }
  defer resp.Body.Close()
  log.Debugf("Catchpoint returned a status code: %s", resp.Status)

  if resp.Status == "200 OK" {
    decoder := json.NewDecoder(resp.Body)
    var ar AuthResponse
    _ = decoder.Decode(&ar)
    bearerToken := base64.StdEncoding.EncodeToString([]byte(ar.AccessToken))
    log.Debugf("Received Bearer token: %s - [base64 encoded] %s", ar.AccessToken, bearerToken)
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
    return token
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
    return token
  }
}
