// Copyright (C) 2021 SeanTolstoyevski - mailto:seantolstoyevski@protonmail.com
//
// The source code of this project is licensed under the MIT license.
// You can find the license on the repo's main folder.
// Provided without warranty of any kind.

// cmd tool for creating the paket file.
//
// A typical use case looks like this.:
// 	paket -f=a_folder_path -k=my_secret_key -m=cfb -i=24000
//
// This command encrypts all the files in the 'a_folder_path' folder with 'my_secret_key' using AES 256, then write the hash information for each file in a table.
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt" // for random salt
	"golang.org/x/crypto/pbkdf2"

	paket "github.com/SeanTolstoyevski/paket/pengine"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

var (
	randBytes, raerr = paket.CreateRandomBytes(32)
	keyDefault       = fmt.Sprintf("%x", sha256.Sum256(randBytes))

	foldername      = flag.String("f", "", "Folder containing files to be encrypted. It is not recursive, Subfolders is not encrypted.")
	outputfile      = flag.String("o", "data.pack", "The file to which your encrypted data will be written. If there is a file with the same name, you will be warned.")
	keyvalue        = flag.String("k", "", "Key for encrypting files. If this parameter is null, the tool generates one randomly bytes and prints value to the console.")
	anonFileName    = flag.Bool("a", false, "anonymize file names. For example, the ''lion.zip'' file is written to the table with a name such as ''201bce5f''\nThis writes the names as ''original   	   random'' in a txt for you to remember later.")
	eMode           = flag.String("m", "gcm", "The mode to be selected for encryption. Currently ''CFB'', ''CTR'', ''GCM'' and ''OFB'' are supported.")
	pbkdf2Iter      = flag.Uint("i", 4096, "Iteration count for pbkdf2. For less than 4096, 4096 will be selected.\nFor modern CPUs values like 100000 may be appropriate.")
	tablefile       = flag.String("t", "PaketTable.go", "The go file to be written for Paket to read. When compiling this file, you must import it into your program.\nIt is created as \"package main.\"")
	showprogressval = flag.Bool("s", true, "prints progress steps to the console. For example, which file is currently encrypting, etc.")
)

func main() {
	if *foldername == "" {
		fmt.Println("\"-f (folder)\" parameter cannot be null.\nSee", os.Args[0], "-help")
		return
	}

	// mode check
	var mode paket.MODE = 0
	switch strings.ToLower(*eMode) {
	case "cbc":
		mode = paket.MODECBC
	case "cfb":
		mode = paket.MODECFB
	case "ctr":
		mode = paket.MODECTR
	case "ofb":
		mode = paket.MODEOFB
	case "gcm":
		mode = paket.MODEGCM
	default:
		fmt.Printf("%s is invalid encryption mode", *eMode)
		return
	}

	if *pbkdf2Iter < 4096 {
		*pbkdf2Iter = 4096
	}

	if *showprogressval {
		fmt.Println("--- INFO ---")
		fmt.Println("Mode:", *eMode)
		fmt.Println("PBDFK2 iteration:", *pbkdf2Iter)
		fmt.Println("Anonymizing file names:", *anonFileName)
	}

	var useKey []byte

	randSalt, err := bcrypt.GenerateFromPassword(randBytes, 10)
	errHandler(err)

	if *keyvalue == "" {
		useKey = []byte(keyDefault)
		fmt.Printf("Your random key: %s\n", keyDefault)
	} else {
		useKey = []byte(*keyvalue)
		fmt.Printf("Your key is: %s\n", *keyvalue)
	}
	useKey = pbkdf2.Key(useKey, randSalt, int(*pbkdf2Iter), 32, sha256.New)

	if paket.Exists(*outputfile) {
		fmt.Printf("There is a file with this name (%s). You can rerun cmd tool  under a different name, rename the existing file, or delete it.", *outputfile)
		return
	}

	if paket.Exists(*tablefile) {
		fmt.Println("The table file will be recreate.")
	}

	gotablefile, err := os.Create(*tablefile)
	errHandler(err)
	defer gotablefile.Close()

	packFile, err := os.Create(*outputfile)
	errHandler(err)
	defer packFile.Close()

	var anonInfos *os.File
	if *anonFileName {
		var err error
		anonInfos, err = os.Create("anonymization-information.txt")
		errHandler(err)
		defer anonInfos.Close()
		anonInfos.Write([]byte("original   \t   anonymous\r\n\r\n"))
	}

	allFolderFiles, err := ioutil.ReadDir(*foldername)
	errHandler(err)

	fileList := []os.FileInfo{}
	for _, file := range allFolderFiles {
		if !file.IsDir() {
			fileList = append(fileList, file)
		}
	}

	if *showprogressval {
		fmt.Printf("%d files were found in %s folder.\n", len(fileList), *foldername)
	}

	gotablefile.Write([]byte(fmt.Sprintf(toptemplate, string(randSalt))))

	var start, full, end int = 0, 0, 0

	for _, file := range fileList {
		gcmNonce := make([]byte, 12)
		if mode == 5 {
			if _, err := io.ReadFull(rand.Reader, gcmNonce); err != nil {
				errHandler(err)
				return
			}
		}

		name := file.Name()

		if *showprogressval {
			fmt.Printf("%s is  encrypting - size: %0.03f MB\n", name, float64(file.Size())/1024.0/1024.0)
		}

		content, err := ioutil.ReadFile(*foldername + "/" + name)
		errHandler(err)
		orgLen := len(content)
		encData, err := paket.Encrypt(useKey, gcmNonce, content, mode)
		errHandler(err)
		encLen := len(encData)
		originalHash := sha256.Sum256(content)
		EncryptedHash := sha256.Sum256(encData)
		orgStringTemplate := "[]byte{"
		encStringTemplate := "[]byte{"

		for _, oHI := range originalHash {
			orgStringTemplate += fmt.Sprint(oHI, ", ")
		}
		orgStringTemplate = orgStringTemplate[:len(orgStringTemplate)-2] + "}"

		for _, eHI := range EncryptedHash {
			encStringTemplate += fmt.Sprint(eHI, ", ")
		}
		encStringTemplate = encStringTemplate[:len(encStringTemplate)-2] + "}"

		nonceTableString := "[]byte{"
		if mode == paket.MODEGCM {
			for _, nonNum := range gcmNonce {
				nonceTableString += fmt.Sprint(nonNum, ", ")
			}
			nonceTableString = nonceTableString[:len(nonceTableString)-2] + "}"
		}

		if _, err := packFile.Write(encData); err != nil {
			errHandler(err)
			return
		}

		start = full
		full += encLen
		end = full

		if *anonFileName {
			randNames16, _ := paket.CreateRandomBytes(16)
			randNames16 = randNames16[:7]
			rname := fmt.Sprintf("%x", randNames16)
			rname = rname[:7]
			anonInfos.Write([]byte(name + "   \t   " + rname + "\r\n"))
			name = rname
		}

		if mode == paket.MODEGCM {
			gotablefile.Write([]byte(fmt.Sprintf(goTemplate, name, strconv.Itoa(start), strconv.Itoa(end), strconv.Itoa(orgLen), strconv.Itoa(encLen), orgStringTemplate, encStringTemplate, nonceTableString)))
		} else {
			gotablefile.Write([]byte(fmt.Sprintf(goTemplate, name, strconv.Itoa(start), strconv.Itoa(end), strconv.Itoa(orgLen), strconv.Itoa(encLen), orgStringTemplate, encStringTemplate, "nil")))
		}

	}

	gotablefile.Write([]byte("}"))
}

func errHandler(err error) {
	if err != nil {
		panic(err)
	}
}

var toptemplate string = `// **DO NOT EDIT this file**. It is generated automatically and contains sensitive data.

// Copyright (C) 2021 SeanTolstoyevski - mailto:seantolstoyevski@protonmail.com

package main

import (
	paket "github.com/SeanTolstoyevski/paket/pengine"
)

// salt
const PaketSalt string = "%s"

// The map vault for datas.
var PaketData = map[string]paket.Values{
`

var goTemplate string = `	"%s" : {StartPos : %s, EndPos : %s, OriginalLenght : %s, EncryptLenght : %s, HashOriginal : %s, HashEncrypt : %s, Nonce: %s},
`

func init() {
	flag.Parse()
	//handle randBytes error
	if raerr != nil {
		panic(raerr)
	}
}
