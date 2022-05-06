package main

import (
	"flag"
	"github.com/mft-labs/agilebindata/bindata"
)

func main() {
	var filename string
	var saveas string
	flag.StringVar(&filename, "embed", "", "Embed file contents")
	flag.StringVar(&saveas, "saveas", "", "Save file as")
	flag.Parse()
	bd := &bindata.BinData{}
	bd.ConvertToGoBundle(filename, saveas, "agilebindata.go")
}
