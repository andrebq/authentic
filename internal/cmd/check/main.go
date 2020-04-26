package main

import (
	"crypto/tls"
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
	dir       = flag.String("workdir", ".", "Directory to process the tests")
	file      = flag.String("file", "main.lua", "Main file to execute")
	tlsConfig = &tls.Config{InsecureSkipVerify: false}

	exitCode int
)

func printFn(l *lua.LState) int {
	ac := l.GetTop()
	args := make([]string, 0, ac)
	for i := 1; i <= ac; i++ {
		val := l.Get(-i)
		args = append(args, fmt.Sprintf("%v", val))
	}
	l.Pop(ac)
	log.Print(strings.Join(args, " "))
	return 0
}

func failFn(l *lua.LState) int {
	printFn(l)
	exitCode = 1
	return 0
}

func fatalFn(l *lua.LState) int {
	printFn(l)
	log.Fatal("Abort!")
	return 0
}

func skipVerifyFn(l *lua.LState) int {
	tlsConfig.InsecureSkipVerify = true
	return 0
}

func main() {
	flag.Parse()
	st := lua.NewState()

	os.Chdir(*dir)

	st.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}).Loader)
	st.SetGlobal("println", st.NewFunction(lua.LGFunction(printFn)))
	st.SetGlobal("print", st.NewFunction(lua.LGFunction(printFn)))
	st.SetGlobal("fail", st.NewFunction(lua.LGFunction(failFn)))
	st.SetGlobal("fatal", st.NewFunction(lua.LGFunction(fatalFn)))
	st.SetGlobal("skipVerify", st.NewFunction(lua.LGFunction(skipVerifyFn)))

	code, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}
	err = st.DoString(string(code))
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}
