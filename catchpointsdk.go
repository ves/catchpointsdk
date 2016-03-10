package catchpointsdk

import (
	"os"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
  "github.com/kelseyhightower/envconfig"
)

type Config struct {
  ClientID string
  ClientSecret  string
  Endpoint  string
}

var c Config

func main() {
	log.SetHandler(json.New(os.Stderr))
  log.SetLevel(log.DebugLevel)

  err := envconfig.Process("catchpointsdk", &c)
  if err != nil { log.Fatal(err.Error()) }
}
