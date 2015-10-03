package gammit

import(
	"gopkg.in/yaml.v2"
	"fmt"
	"github.com/tdewolff/minify"
	"os"
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

type JavascriptAssets struct {
	asset []string
}

type Assets struct {
	Javascripts JavascriptAssets `yaml:"javascripts"`
}

type Gammit struct {
	Config map[interface{}]interface{}
	Minifier  *minify.Minify
	OutputPath string
}

func (g *Gammit) ReadYaml(data []byte) {
	yaml.Unmarshal(data, &g.Config)
}

func (g *Gammit) Process() {
	os.MkdirAll(g.OutputPath, 0755)

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

		outputFile, err := os.Create(g.OutputPath + "/" + group.(string) + fileType)
		g.check(err)
		defer outputFile.Close()

		for _, minifiedBytes := range minified {
			_, err := outputFile.Write(minifiedBytes)
			g.check(err)

			var b bytes.Buffer
			gw := gzip.NewWriter(&b)
			gw.Write(minifiedBytes)
			gw.Close()

			err = ioutil.WriteFile(g.OutputPath + "/" + group.(string) + fileType + ".gz", b.Bytes(), 0666)
			g.check(err)
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
		g.check(err)
		defer f.Close()

		buf := new(bytes.Buffer)
		g.Minifier.Minify(mediaType, buf, f)

		minified[i] = buf.Bytes()

		i += 1
	}
	return minified

}

func (g *Gammit) check(error error) {
	if error != nil {
		panic(error)
	}
}


