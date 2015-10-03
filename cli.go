package main

import(
    "pipelinedeals.com/gammit/gammit"
    "github.com/tdewolff/minify"
    "github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
    "io/ioutil"
    "flag"
)

var (
    configFileLocation = flag.String("c", "assets.yml", "The config file location")
    outputLocation = flag.String("o", ".", "The output directory")
)

func main() {
    minifier := minify.New()
    minifier.AddFunc("text/css", css.Minify)
    minifier.AddFunc("text/javascript", js.Minify)

    data, err := ioutil.ReadFile(*configFileLocation)
    if err != nil {
        panic(err)
    }

    gammit := &gammit.Gammit{Minifier: minifier}
    gammit.ReadYaml(data)
    gammit.Process(*outputLocation)
}

