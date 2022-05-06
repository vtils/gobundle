package bindata

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type BinData struct {
}

func (bd *BinData) ConvertAsGoBundle(filename, saveas, target string) error {
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

var filesList = make([]string, 0)

func (bd *BinData) ConvertFolderAsGoBundle(folder, target string) error {
	err := filepath.Walk(folder, walkFn)
	if err != nil {
		return err
	}
	output := "package main\n\n" +
		"import \"flag\"\n" +
		"import \"log\"\n" +
		"import \"strings\"\n" +
		"import \"bytes\"\n" +
		"import \"compress/gzip\"\n" +
		"import \"encoding/base64\"\n" +
		"import \"os\"\n" +
		"import \"runtime\"\n" +
		"import \"io/ioutil\"\n\n" +
		"var extract bool\n\n"
	output += "var bindata_contents = make(map[string]string)\n"
	output += "var bindata_files = make(map[string]string)\n"
	output += "\n"
	initFunction := "\n" +
		"func init(){\n" +
		"	flag.BoolVar(&extract, \"extract\", false, \"Extract file contents\")\n"

	for _, fpath := range filesList {
		log.Printf("%v", fpath)
		fpath2 := strings.Replace(fpath, "\\", "/", -1)
		contents, err := ioutil.ReadFile(fpath)
		if err == nil {
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
			b64Text := base64.StdEncoding.EncodeToString(b.Bytes())
			initFunction += fmt.Sprintf(`	bindata_contents["%v"]="%v"`, fpath2, b64Text)
			initFunction += "\n\n"
			fpath3 := strings.Replace(fpath, "\\", "\\\\", -1)
			initFunction += fmt.Sprintf(`	bindata_files["%v"]="%v"`, fpath2, fpath3)
			initFunction += "\n\n"
		}

	}
	initFunction += "\n}\n\n"
	initFunction += "\n\n"
	output += initFunction
	output += "\n\n"
	output += "func ExtractAssets(){\n" +
		//"	arr := strings.Split(fileslist,\",\")\n" +
		" 	for key, _ := range bindata_files{\n" +
		" 		data, _ := base64.StdEncoding.DecodeString(bindata_contents[key])\n" +
		"		rdata := bytes.NewReader(data)\n" +
		"		r,_ := gzip.NewReader(rdata)\n" +
		"		sout, _ := ioutil.ReadAll(r)\n" +
		"		saveas := bindata_files[key]" +
		"		//ioutil.WriteFile(filename,[]byte(data),0700)\n" +
		"		fpath := strings.Replace(saveas,\"\\\\\\\\\",\"\\\\\",-1)\n" +
		"		if runtime.GOOS == \"windows\" {\n" +
		"			arr := strings.Split(fpath,\"\\\\\")\n" +
		"			levels := len(arr)-1\n" +
		"			folderpath := \"\"\n" +
		"			count := 0\n" +
		"			for _, elem := range arr {\n" +
		"				count ++\n" +
		"				folderpath += elem \n" +
		"				if count == levels {\n" +
		"					break\n" +
		"				}\n" +
		"				folderpath +=  \"\\\\\"\n" +
		"			}\n" +
		"			if folderpath[0] == '/' {\n" +
		"				folderpath = folderpath[1:]\n" +
		"			}\n" +
		"			if folderpath[0] == '\\\\' {\n" +
		"				folderpath = folderpath[1:]\n" +
		"			}\n" +
		"			log.Printf(\"%v\",folderpath)\n" +
		"			_, err := os.Stat(folderpath)\n" +
		"			if os.IsNotExist(err) {\n" +
		"				errDir := os.MkdirAll(folderpath, 0755)\n" +
		"				if errDir != nil {\n" +
		"					continue\n" +
		"				}\n" +
		"			}\n" +
		"		}\n" +
		"		if runtime.GOOS == \"linux\" {\n" +
		"			fpath = strings.Replace(saveas,\"\\\\\\\\\",\"/\",-1)\n" +
		"			fpath = strings.Replace(saveas,\"\\\\\",\"/\",-1)\n" +
		"			arr := strings.Split(fpath,\"/\")\n" +
		"			levels := len(arr) - 1\n" +
		"			folderpath := \"\"\n" +
		"			count := 0\n" +
		"			for _, elem := range arr {\n" +
		"				count ++\n" +
		"				folderpath += elem \n" +
		"				if count == levels {\n" +
		"					break\n" +
		"				}\n" +
		"				folderpath += \"/\"\n" +
		"			}\n" +
		"			if folderpath[0] == '/' {\n" +
		"				folderpath = folderpath[1:]\n" +
		"			}\n" +
		"			log.Printf(\"%v\",folderpath)\n" +
		"			_, err := os.Stat(folderpath)\n" +
		"			if os.IsNotExist(err) {\n" +
		"				errDir := os.MkdirAll(folderpath, 0755)\n" +
		"				if errDir != nil {\n" +
		"					continue\n" +
		"				}\n" +
		"			}\n" +
		"		}\n" +
		"		if fpath[0] == '/' {\n" +
		"			fpath = fpath[1:]\n" +
		"		}\n" +
		"		if fpath[0] == '\\\\' {\n" +
		"			fpath = fpath[1:]\n" +
		"		}\n" +
		"		log.Printf(\"Extracting %v\",fpath)\n" +
		"		ioutil.WriteFile(fpath,sout,0700)\n" +
		"	}\n" +
		"}\n"
	ioutil.WriteFile(target, []byte(output), 0700)
	return err
}

func walkFn(path string, fi os.FileInfo, err error) (e error) {
	if !fi.IsDir() {
		filesList = append(filesList, path)
	}
	return nil
}
