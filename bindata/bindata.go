package bindata

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

type BinData struct {
}

func (bd *BinData) ConvertToGoBundle(filename, saveas, target string) error {
	contents, err := ioutil.ReadFile(filename)
	if err == nil {
		//log.Printf("%v\n", string(contents))
		var b bytes.Buffer
		gz := gzip.NewWriter(&b)
		if _, err := gz.Write(contents); err != nil {
			return err
		}
		if err := gz.Flush(); err != nil {
			return err
		}
		if err := gz.Close(); err != nil {
			return err
		}
		//b64Text := base64.StdEncoding.EncodeToString(contents)
		b64Text := base64.StdEncoding.EncodeToString(b.Bytes())
		//log.Printf("%v\n", b64Text)
		output := "package main\n\n" +
			"import \"flag\"\n" +
			"import \"strings\"\n" +
			"import \"bytes\"\n" +
			"import \"compress/gzip\"\n" +
			"import \"encoding/base64\"\n" +
			"import \"io/ioutil\"\n\n" +
			"var extract bool\n\n"

		output += fmt.Sprintf("var fileslist=\"%v\"\n", filename)
		output += "var bindata_contents = make(map[string]string)\n"
		output += "var bindata_files = make(map[string]string)\n"
		output += "\n" +
			"func init(){\n" +
			"	flag.BoolVar(&extract, \"extract\", false, \"Extract file contents\")\n"
		output += fmt.Sprintf(`	bindata_contents["%v"]="%v"`, filename, b64Text)
		output += "\n\n"
		output += fmt.Sprintf(`	bindata_files["%v"]="%v"`, filename, saveas)
		output += "\n}\n\n"
		output += "\n\n"
		output += "func ExtractAssets(){\n" +
			"	arr := strings.Split(fileslist,\",\")\n" +
			" 	for _, filename := range arr{\n" +
			" 		data, _ := base64.StdEncoding.DecodeString(bindata_contents[filename])\n" +
			"		rdata := bytes.NewReader(data)\n" +
			"		r,_ := gzip.NewReader(rdata)\n" +
			"		sout, _ := ioutil.ReadAll(r)\n" +
			"		saveas := bindata_files[filename]" +
			"		//ioutil.WriteFile(filename,[]byte(data),0700)\n" +
			"		ioutil.WriteFile(saveas,sout,0700)\n" +
			"	}\n" +
			"}\n"
		ioutil.WriteFile(target, []byte(output), 0700)
	}
	return err
}
