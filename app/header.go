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

	var headerTemplate = fmt.Sprintf(`
          %s
          |
 %s   %s    %s
  \ / \ / | \
   %s   %s  %s  %s

`, o, o, o, o, o, o, o, o)

	headerTemplate += fmt.Sprintf(" Running on port %s \n\n", lb("{{ .Port }}"))

	buf := new(bytes.Buffer)
	t := template.Must(template.New("header").Parse(headerTemplate))
	t.Execute(buf, &headerOpts{
		Port: port,
	})
	color.Println(buf.String())
}
