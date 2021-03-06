package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/crackcomm/go-actions/action"
	"github.com/crackcomm/go-actions/core"
	"github.com/crackcomm/go-actions/encoding/yaml"
	"github.com/crackcomm/go-actions/local"
	"github.com/golang/glog"

	_ "github.com/crackcomm/go-actions/source/file"
	_ "github.com/crackcomm/go-actions/source/http"
	_ "github.com/crackcomm/go-core"
)

var filename = "cmds.yaml" // cmds filename

func main() {
	defer glog.Flush()
	flag.Parse()
	actions, err := fileToActions(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) <= 1 {
		fmt.Println("cmds needs at least one argument\n")
		fmt.Println("available commands:\n")
		for name := range actions {
			fmt.Printf("  %s\n", name)
		}
		return
	}

	for name, a := range actions {
		core.Action(name, a)
	}

	res, err := local.Run(&action.Action{Name: os.Args[1]})
	if res != nil {
		res = mapBytes(res)
		printMap(res, 1)
		fmt.Printf("\n")
	}
	if err != nil {
		fmt.Printf("ERROR %v\n", err)
	} else {
		fmt.Println("OK")
	}
}

// return n * 2 spaces
func spaces(n int) string {
	return strings.Repeat(" ", n*2)
}

// print with spaces
func prints(n int, f string, a ...interface{}) {
	fmt.Printf("%s%s", spaces(n), fmt.Sprintf(f, a...))
}

func printValue(value interface{}, n int) {
	switch value.(type) {
	case action.Map:
		print("\n")
		printMap(value.(action.Map), n+1)
	case []interface{}:
		print("\n")
		for _, v := range value.([]interface{}) {
			prints(n+1, "-")
			printValue(v, n+1)
		}
	case []string:
		print("\n")
		for _, v := range value.([]string) {
			prints(n+1, "- %v\n", v)
		}
	case string:
		lines := strings.Split(value.(string), "\n")
		if len(lines) <= 1 {
			print(value.(string))
			print("\n")
		} else {
			print("\n")
			for _, ln := range lines {
				prints(n+1, "%s\n", ln)
			}
		}
	default:
		f := action.Format{value}
		if v, ok := f.String(); ok {
			prints(0, "%#v\n", v)
		} else {
			prints(0, "%#v\n", value)
		}
	}
}

func printKeyValue(key string, value interface{}, n int) {
	prints(n, "%s: ", key)
	printValue(value, n)
}

func printMap(m action.Map, n int) {
	for key, value := range m {
		printKeyValue(key, value, n)
	}
}

// mapBytes - Iterates through a result map and tranforms all byte arrays into strings.
func mapBytes(m action.Map) action.Map {
	// Iterate map
	for k, v := range m {
		// Pass if value is nil
		if v == nil {
			continue
		}

		// If it's a byte array or a reader we are gonna make it a string
		if rc, ok := v.(io.ReadCloser); ok {
			arr, _ := ioutil.ReadAll(rc)
			rc.Close()
			m[k] = string(arr)
		} else if arr, ok := v.([]byte); ok {
			m[k] = string(arr)
		}
	}
	return m
}

func indent(body string) string {
	res := ""
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		line = "  " + line + "\n"
		res = res + line
	}
	return res
}

func fileToActions(filename string) (map[string]*action.Action, error) {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return yaml.UnmarshalMany(body)
}

func init() {
	flag.StringVar(&filename, "file", filename, "Cmds file")
}
