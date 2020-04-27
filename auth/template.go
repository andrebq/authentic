package auth

import (
	"html/template"
	"io/ioutil"

	"github.com/markbates/pkger"
	pkg "github.com/markbates/pkger/pkging"

	// load builtin templates
	_ "github.com/andrebq/authentic/res"
)

// BuiltinTemplates returns the templates shipped with authentic
func BuiltinTemplates() *template.Template {
	t := template.New("root")
	_, err := t.New("login/new").Parse(readBuiltin(pkger.Open("/auth/templates/login/new.html")))
	if err != nil {
		panic(err.Error())
	}
	return t
}

func readBuiltin(f pkg.File, err error) string {
	if err != nil {
		panic(err)
	}
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
