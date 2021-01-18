// go test -cover

// Copyright (C) 2021 SeanTolstoyevski -  mailto:s.tolstoyevski@protonmail.com
//
// The source code of this project is licensed under the MIT license.
// You can find the license on the repo's main folder.
// Provided without warranty of any kind.

// cmd tool for creating the package file.
package main

import (
	"crypto/sha256"
	"flag"
	"fmt"

	paket "github.com/SeanTolstoyevski/paket/pengine"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	randBytes, raerr = paket.CreateRandomBytes(16)
	sKeyDefault      = fmt.Sprintf("%x", randBytes)
)

func init() {
	//handle randBytes error
	if raerr != nil {
		panic(raerr)
	}
}

var (
	foldername      = flag.String("f", "", "Folder containing files to be encrypted.\n		Note: Your original files are not deleted. \nIt is not recursive, Subfolders is not encrypted.")
	outputfile      = flag.String("o", "data.pack", "The file to which your encrypted data will be written. \n If there is a file with the same name, you will be warned.")
	keyvalue        = flag.String("k", sKeyDefault, "Key for encrypting files. It must be 16, 24, or 32 lenght in bytes.\nIf you leave it null, the tool generates one randomly byte  and prints value to the console.")
	tablefile       = flag.String("t", "PaketTable.go", "The go file to be written for Paket to read. \n When compiling this file, you must import it into your program. \n It is created as 'package main.'")
	addshaval       = flag.Bool("h", true, "Writes hash of original and encrypted versions of the files to table.\nThis is required for security. \nIf left null, hash checks will not work.")
	showprogressval = flag.Bool("s", true, "prints progress steps to the console. For example, which file is currently encrypting, etc.")
)

func main() {
	flag.Parse()
	if *foldername == "" {
		fmt.Println("\"-fn\" parameter cannot be null.\nSee", os.Args[0], "-help")
		os.Exit(1)
	}
	skey := *keyvalue
	keyByte := []byte(skey)
	keyByteLen := len(keyByte)
	if !confirmatorLen(keyByteLen) {
		fmt.Println("Exiting. \nWrong key lenght. Lenght: (bytes)", keyByteLen)
		os.Exit(1)
	}

	if paket.Exists(*outputfile) {
		fmt.Printf("Exiting.\nThere is a file with this name (%s). You can rerun cmd tool  under a different name, rename the existing file, or delete it.", *outputfile)
		os.Exit(1)
	}
	if paket.Exists(*tablefile) {
		fmt.Println("The table file will be recreate.")
	}
	if *keyvalue == sKeyDefault {
		fmt.Println("❗❕ Warning! Your random key. Please note:", *keyvalue)
	} else {
		fmt.Println("❗❕ Warning! Your key is:", *keyvalue)
	}
	gotablefile, err := os.Create(*tablefile)
	defer gotablefile.Close()
	errHandler(err)

	packFile, err := os.OpenFile(*outputfile, os.O_RDWR|os.O_CREATE, 0666)
	defer packFile.Close()
	errHandler(err)

	var start, full, end int

	listFiles, err := ioutil.ReadDir(*foldername)
	errHandler(err)
	show := *showprogressval
	if show {
		fmt.Printf("%d files were found in %s folder.\n", len(listFiles), *foldername)
	}
	gotablefile.Write([]byte(fmt.Sprintf(toptemplate, *foldername)))
	for _, file := range listFiles {
		if !file.IsDir() {
			name := file.Name()
			if show {
				fmt.Printf("%s file is encrypting. Size: %0.004f MB\n", name, float32(file.Size())/100000.0)
			}
			content, err := ioutil.ReadFile(*foldername + "/" + name)
			errHandler(err)
			orgLen := len(content)
			encData, _ := paket.Encrypt(keyByte, content)
			encLen := len(encData)
			orgSha := fmt.Sprintf("%x", sha256.Sum256(content))
			encSha := fmt.Sprintf("%x", sha256.Sum256(encData))
			_, rerr := packFile.Write(encData)
			errHandler(rerr)
			start = full
			full += encLen
			end = full

			gotablefile.Write([]byte(fmt.Sprintf(goTemplate, name, strconv.Itoa(start), strconv.Itoa(end), strconv.Itoa(orgLen), strconv.Itoa(encLen), orgSha, encSha)))
		}
	}
	gotablefile.Write([]byte("}"))
}

func errHandler(err error) {
	if err != nil {
		panic(err)
	}
}

var toptemplate string = `//important: You can edit this file. However, you need to know what you are doing.
// *panic* may occur.

package main

import (
	paket "github.com/SeanTolstoyevski/paket/pengine"
)

//The map vault for datas. The init function writing the required data.
var Data = make(paket.Datas)

// The name of the folder from which the files were was taken. Information is writing by init.
var foldername string

func init() {
	foldername = "%s"
`

var goTemplate string = `	Data["%s"] = paket.Values{"%s", "%s", "%s", "%s", "%s", "%s"}
`

func confirmatorLen(l int) bool {
	if l == 16 || l == 24 || l == 32 {
		return true
	}
	return false
}
