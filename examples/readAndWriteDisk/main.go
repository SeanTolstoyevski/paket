// Copyright (C) 2021 SeanTolstoyevski -  mailto:s.tolstoyevski@protonmail.com
// The source code of this project is licensed under the MIT license.
// You can find the license on the repo's main folder.
// Provided without warranty of any kind.

package main

import (
	"fmt"
	paket "github.com/SeanTolstoyevski/paket/pengine"
	"os"
)

var (
	packfile = "./../data.pack"
	packkey  = []byte("bf40e1d71af5ca0be1e2b02bbcf42d3f")
)

func main() {
	encFiles, err := paket.New(packkey, packfile, Data)
	errHandler(err)
	defer encFiles.Close()

	for key := range Data {
		wf, err := os.Create(key)
		errHandler(err)
		file, _, err := encFiles.GetFile(key, true, true)
		errHandler(err)
		fmt.Printf("%s file is writing...\n", key)
		wf.Write(*file)
		cerr := wf.Close()
		errHandler(cerr)

	}
}

func errHandler(err error) {
	if err != nil {
		panic(err)
	}
}
