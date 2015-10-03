package main

import(
	"gopkg.in/yaml.v2"
    "fmt"
    "github.com/tdewolff/minify"
    "github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
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

type Gammit struct {
    Config map[interface{}]interface{}
    Minifier  *minify.Minify
}

func (g *Gammit) ReadYaml() {
    yaml.Unmarshal([]byte(data), &g.Config)
}

func (g *Gammit) Process() {
    g.processJavascripts();
}

func (g *Gammit) processJavascripts() {
    for group, files := range g.Config["javascripts"].(map[interface{}]interface{}) {
        fmt.Println("Group " + group.(string))

        fileList := files.([]interface{})
        minified := g.minifyFilesInGroup("text/javascript", fileList)

        outputFile, err := os.Create(group.(string) + ".js")
        check(err)
        defer outputFile.Close()

        for _, minifiedBytes := range minified {
            _, err := outputFile.Write(minifiedBytes)
            check(err)
        }
        outputFile.Close()
    }
}

func (g *Gammit) minifyFilesInGroup(mediaType string, fileList []interface{}) ([][]byte) {
    i := 0
    minified := make([][]byte, len(fileList))

    for _, file := range fileList {
        fmt.Println("   File: " +file.(string))
        os.Open(file.(string))

        f, err := os.Open(file.(string))
        check(err)
        defer f.Close()

        buf := new(bytes.Buffer)
        g.Minifier.Minify(mediaType, buf, f)

        minified[i] = buf.Bytes()

        i += 1
    }
    return minified

}


func main() {
    minifier := minify.New()
    minifier.AddFunc("text/css", css.Minify)
    minifier.AddFunc("text/javascript", js.Minify)

    gammit := &Gammit{Minifier: minifier}
    gammit.ReadYaml()
    gammit.Process()
}

func check(error error) {
    if error != nil {
        panic(error)
    }
}

