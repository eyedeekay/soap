package main

//import "github.con/eyedeekay/soap/lib"
import (
	"log"
	"net"
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
	log.Printf("Multiplayer Unciv Server listening on:\n\t%s", ln.Addr().String())
	defer ln.Close()
	fs := &unciv.FrontServer{
		PageTitle:   "Unciv Server",
		ServerName:  "Default Server",
		Description: "TODO",
		TOS:         "TODO",
		URL:         ln.Addr().String(),
	}
	go func() {
		garlic2, err := onramp.NewGarlic("unciv-display", "127.0.0.1:7656", nil)
		if err != nil {
			panic(err)
		}
		fsln, err := garlic2.Listen()
		if err != nil {
			panic(err)
		}
		log.Printf("Unciv Server sharing server:\n\t%s", fsln.Addr().String())
		defer fsln.Close()
		if err := http.Serve(fsln, fs); err != nil {
			panic(err)
		}
	}()

	go func() {
		fsln, err := net.Listen("tcp", "127.0.0.1:7699")
		if err != nil {
			panic(err)
		}
		defer fsln.Close()
		if err := http.Serve(fsln, fs); err != nil {
			panic(err)
		}
	}()

	if err := http.Serve(ln, ucs); err != nil {
		panic(err)
	}
}
