package unciv

import (
	"html/template"
	"log"
	"net/http"
)

type FrontServer struct {
	PageTitle   string
	ServerName  string
	Description string
	TOS         string
	URL         string
}

var temp = `<!DOCTYPE html>
<head>
<title>{{.PageTitle}}</title>
</head>
<html>
<body>
<h1>{{.PageTitle}}</h1>
<h2>{{.ServerName}}</h2>
<code>{{.URL}}</code>
<div>
	{{.Description}}
</div>
<h3>Terms of Service</h3>
<div>
	{{.TOS}}
</div>

</body>
</html>
`

func (f *FrontServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("foo").Parse(temp)
	if err != nil {
		log.Println(err)
		return
	}
	err = t.Execute(w, f)
	if err != nil {
		log.Panicln(err)
		return
	}
}
