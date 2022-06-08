package main

import (
	"flag"
	"github.com/vtils/gobundle/bindata"
)

var (
	extract bool
)

func main() {
	var embed string
	var saveas string
	var folder bool
	flag.StringVar(&embed, "embed", "", "Embed file contents")
	flag.StringVar(&saveas, "saveas", "", "Save file as")
	flag.BoolVar(&folder, "folder", false, "Bundle given folder")
	flag.BoolVar(&extract,"extract",false,"Extract bundle")
	flag.Parse()
	var bd interface{}
	bd = &bindata.BinData{}

	if embed != "" && saveas != "" {
		if !folder {
			bd.(*bindata.BinData).ConvertAsGoBundle(embed, saveas, "bindata/agilebindata.go")
		} else {
			bd.(*bindata.BinData).ConvertFolderAsGoBundle(embed, "bindata/agilebindata.go")
		}
	} else  {
		if extract {
			if obj, ok := bd.(interface{ExtractAssets()}); ok { 
				obj.ExtractAssets() 
			} 
		}
		
	}
}
