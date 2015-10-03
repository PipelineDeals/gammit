package main

import(
	"gopkg.in/yaml.v2"
  "log"
  "fmt"
)

var data = `
template_function: off
template_extension: .hbs

javascripts:

  ###########################################################
  # Global
  ###########################################################

  one:
    - test/test.js
    - test/test2.js

  two:
    - test/test3.js
    - test/test4.js
`

type JavascriptAssets struct {
	asset []string
}

type Assets struct {
	Javascripts JavascriptAssets `yaml:"javascripts"`
}

func main() {
    m := make(map[interface{}]interface{})
    err := yaml.Unmarshal([]byte(data), &m)
    if err != nil {
        log.Fatalf("error: %v", err)
    }

    fmt.Println("groups:")
    for group, files := range m["javascripts"].(map[interface{}]interface{}) {
        fmt.Println("Group " +group.(string))

        for _, file := range files.([]interface{}) {
            fmt.Println("   File: " +file.(string))
        }
    }
}

