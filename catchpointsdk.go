package catchpointsdk

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/kelseyhightower/envconfig"
	"os"
)

type Config struct {
	ClientID     string
	ClientSecret string
	Endpoint     string
}

var c Config
var (
	division_id = os.Getenv("CATCHPOINTSDK_DIVISION_ID")
)



func main() {
	log.SetHandler(json.New(os.Stderr))
	log.SetLevel(log.DebugLevel)

	err := envconfig.Process("catchpointsdk", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
}
