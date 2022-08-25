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
	garlic1, err := onramp.NewGarlic("unciv", "127.0.0.1:7656", nil)
	if err != nil {
		panic(err)
	}
	//ln, err := garlic.ListenTLS()
	ln, err := garlic1.Listen()
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	go func() {
		garlic2, err := onramp.NewGarlic("unciv-display", "127.0.0.1:7656", nil)
		if err != nil {
			panic(err)
		}
		fsln, err := garlic2.Listen()
		if err != nil {
			panic(err)
		}
		defer fsln.Close()
		fs := &unciv.FrontServer{
			PageTitle:   "Unciv Server",
			ServerName:  "Default Server",
			Description: "TODO",
			TOS:         "TODO",
			URL:         ln.Addr().String(),
		}
		if err := http.Serve(fsln, fs); err != nil {
			panic(err)
		}
	}()

	if err := http.Serve(ln, ucs); err != nil {
		panic(err)
	}
}
