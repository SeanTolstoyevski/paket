package main

import (
	"fmt"
	paket "github.com/SeanTolstoyevski/paket/pengine"
	"os"
)

func main() {
	encFiles, err := paket.New(paket.Option{Key: []byte("aalllauaidasjewnwoo210dkcnsjeiwjqijijvijiejiji3ji4ji3nini34"),
		Salt:      PaketSalt,
		PaketFile: "./../../testdata/GCM.data",
		Table:     PaketData,
		Mode:      paket.MODEGCM})
	if err != nil {
		fmt.Println("Error: creating new GCM paket:", err)
		return
	}

	fmt.Println("Paket: Mode: GCM - write all data in table to disk")
	fileCount := 0
	fileBytesTotal := 0
	for fileName := range PaketData {
		secureData, hash, err := encFiles.GetFile(fileName, true, true)
		if err != nil {
			fmt.Println("Error: getting data with Getfile(name, true, true):", err)
			os.Exit(1)
		}
		fmt.Println("Hash check:", hash,
			"\nFile name:", fileName,
			"\nFile lenght:", len(secureData),
		)
		if !hash {
			fmt.Println("Warning: in data hash check")
		}

		if writeFileObj, err := os.Create(fileName); err != nil {
			fmt.Println("Error: Create new file with os.Create")
			os.Exit(1)
		} else {
			writeFileObj.Write(secureData)
			fileCount++
			fileBytesTotal += len(secureData)
			writeFileObj.Close()
		}

	} // loop block

	fmt.Println("Paket: Mode: GCM - Write to disk finished")
	fmt.Println("File count:", fileCount,
		"\nTotal bytes:", fileBytesTotal)

	if length, err := encFiles.GetLen(); err == nil {
		if fileBytesTotal == length[0] {
			fmt.Println("Info: total bytes in table match the length  of bytes write  to disk:", length[0])
		} else {
			fmt.Println("Warning: The length  of bytes write to disk does not match the length of bytes in table.",
				"\nTable:", length[0],
				"\nBy this example:", fileBytesTotal)
		}
	} else {
		fmt.Println("Error: getting length of data from table")
	}

	// Although not true for all scenarios, Close should not be used with "defer".
	// When you  get all your data you should release the paket with close.
	if err := encFiles.Close(); err == nil {
		fmt.Println("Closed paket")
		return
	} else {
		fmt.Println("Error: while closing paket with Close.", err)
		os.Exit(1)
	}

}
