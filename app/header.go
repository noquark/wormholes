package app

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/gookit/color"
)

type headerOpts struct {
	Port int
}

func ShowHeader(port int) {
	b := color.Blue.Render
	lb := color.New(color.Gray).Render
	o := b("o")

	var headerTmpl = fmt.Sprintf(`
          %s
          |
 %s   %s    %s
  \ / \ / | \
   %s   %s  %s  %s

`, o, o, o, o, o, o, o, o)

	headerTmpl += fmt.Sprintf(" Running on port %s \n\n", lb("{{ .Port }}"))
	headerTmpl += fmt.Sprintf(" %s     %s       - to create links\n", b("PUT"), lb("/api/v1/links"))
	headerTmpl += fmt.Sprintf(" %s    %s  - to update link\n", b("POST"), lb("/api/v1/links/<id>"))
	headerTmpl += fmt.Sprintf(" %s     %s  - to get link data\n", b("GET"), lb("/api/v1/links/<id>"))
	headerTmpl += fmt.Sprintf(" %s  %s  - to delete link\n", b("DELETE"), lb("/api/v1/links/<id>"))

	buf := new(bytes.Buffer)
	t := template.Must(template.New("header").Parse(headerTmpl))
	t.Execute(buf, &headerOpts{
		Port: port,
	})
	color.Println(buf.String())
}
