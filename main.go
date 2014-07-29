package main

import (
	"fmt"
	"github.com/crackcomm/go-actions/action"
	"github.com/crackcomm/go-actions/core"
	"github.com/crackcomm/go-actions/encoding/yaml"
	"github.com/crackcomm/go-actions/local"
	_ "github.com/crackcomm/go-core"
	yml "gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	actions, err := fileToActions("cmds.yaml")
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
		core.Default.Registry.Add(name, a)
	}

	name := os.Args[1]

	fmt.Printf("Running %s\n", name)

	a := &action.Action{Name: name}

	res, err := local.Run(a)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	ybody, err := yml.Marshal(res)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("\n%s", indent(string(ybody)))
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
