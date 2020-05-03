package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/cjoudrey/gluahttp"
	"github.com/cjoudrey/gluaurl"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/net/publicsuffix"
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
		val := l.Get(i)
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

func noRedirect(_ *http.Request, _ []*http.Request) error {
	return http.ErrUseLastResponse
}

type (
	testJar struct {
		http.CookieJar
	}
)

func (e *testJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	e.CookieJar.SetCookies(u, cookies)
}

func (e *testJar) Cookies(u *url.URL) []*http.Cookie {
	res := e.CookieJar.Cookies(u)
	return res
}

func (e *testJar) Clear() {
	j, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		panic(err)
	}
	e.CookieJar = j
}

func main() {
	flag.Parse()
	st := lua.NewState()

	os.Chdir(*dir)

	jar := &testJar{}
	jar.Clear()

	st.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		//CheckRedirect: noRedirect,
		Jar: jar,
	}).Loader)
	st.PreloadModule("url", gluaurl.Loader)
	st.SetGlobal("println", st.NewFunction(lua.LGFunction(printFn)))
	st.SetGlobal("print", st.NewFunction(lua.LGFunction(printFn)))
	st.SetGlobal("fail", st.NewFunction(lua.LGFunction(failFn)))
	st.SetGlobal("fatal", st.NewFunction(lua.LGFunction(fatalFn)))
	st.SetGlobal("skipVerify", st.NewFunction(lua.LGFunction(skipVerifyFn)))
	st.SetGlobal("clear_cookies", st.NewFunction(lua.LGFunction(func(_ *lua.LState) int {
		jar.Clear()
		return 0
	})))

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
