package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cjoudrey/gluahttp"
	lua "github.com/yuin/gopher-lua"
)

var (
	dir  = flag.String("workdir", ".", "Directory to process the tests")
	file = flag.String("file", "main.lua", "Main file to execute")
)

func printFn(l *lua.LState) int {
	ac := l.GetTop()
	args := make([]string, 0, ac)
	for i := 1; i <= ac; i++ {
		val := l.Get(-i)
		args = append(args, fmt.Sprintf("%v", val))
	}
	log.Print(strings.Join(args, " "))
	return 0
}

func main() {
	flag.Parse()
	st := lua.NewState()

	os.Chdir(*dir)

	st.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)
	st.SetGlobal("println", st.NewFunction(lua.LGFunction(printFn)))
	st.SetGlobal("print", st.NewFunction(lua.LGFunction(printFn)))

	code, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}
	err = st.DoString(string(code))
	if err != nil {
		log.Fatal(err)
	}
}
