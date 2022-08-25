package main

//import "github.con/eyedeekay/soap/lib"
import (
	"net/http"

	unciv "./lib"

	onramp "github.com/eyedeekay/onramp"
)

func main() {
	ucs := &unciv.UncivServer{}
	//listener
	garlic, err := onramp.NewGarlic("unciv", "127.0.0.1:7656", nil)
	if err != nil {
		panic(err)
	}
	//ln, err := garlic.ListenTLS()
	ln, err := garlic.Listen()
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	if err := http.Serve(ln, ucs); err != nil {
		panic(err)
	}
}
