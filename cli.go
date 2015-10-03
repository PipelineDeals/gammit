package main

import(
	"gopkg.in/yaml.v2"
    "log"
    "fmt"
    "github.com/web-assets/go-jsmin"
    "os"
    "bytes"
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

        fileList := files.([]interface{})

        minified := minifyFilesInGroup(fileList)
        fmt.Println(minified)

    }
}

func minifyFilesInGroup(fileList []interface{}) ([][]byte) {
    i := 0
    minified := make([][]byte, len(fileList))

    for _, file := range fileList {
        fmt.Println("   File: " +file.(string))
        os.Open(file.(string))

        f, err := os.Open(file.(string))
        if err != nil {
            panic(err)
        }

        buf := new(bytes.Buffer)
        jsmin.Min(f, buf)
        f.Close()

        // read the output into our minified slice
        minified[i] = buf.Bytes()

        i += 1
    }
    return minified
}

