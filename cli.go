package main

import(
	"gopkg.in/yaml.v2"
    "fmt"
    "github.com/tdewolff/minify"
    "github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
    "os"
    "bytes"
    "compress/gzip"
    "io/ioutil"
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
    g.processGroup("javascripts", "text/javascript", ".js");
    g.processGroup("stylesheets", "text/css", ".css");
}

func (g *Gammit) processGroup(section string, mediaType string, fileType string) {
    if g.Config[section] == nil {
        panic("Could not find section '"+section+"' in the config file!")
    }
    for group, files := range g.Config[section].(map[interface{}]interface{}) {
        fmt.Println("Group " + group.(string))

        fileList := files.([]interface{})
        minified := g.minifyFilesInGroup(mediaType, fileList)

        outputFile, err := os.Create(group.(string) + fileType)
        check(err)
        defer outputFile.Close()

        for _, minifiedBytes := range minified {
            _, err := outputFile.Write(minifiedBytes)
            check(err)

            var b bytes.Buffer
            gw := gzip.NewWriter(&b)
            gw.Write(minifiedBytes)
            gw.Close()

            err = ioutil.WriteFile(group.(string) + fileType + ".gz", b.Bytes(), 0666)
            check(err)
        }
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

