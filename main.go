package main

import (
	"flag"
	"github.com/vtils/gobundle/bindata"
)

func main() {
	var embed string
	var saveas string
	var folder bool
	flag.StringVar(&embed, "embed", "", "Embed file contents")
	flag.StringVar(&saveas, "saveas", "", "Save file as")
	flag.BoolVar(&folder, "folder", false, "Bundle given folder")
	flag.Parse()
	bd := &bindata.BinData{}
	if !folder {
		bd.ConvertAsGoBundle(embed, saveas, "agilebindata.go")
	} else {
		bd.ConvertFolderAsGoBundle(embed, "agilebindata.go")
	}

}
