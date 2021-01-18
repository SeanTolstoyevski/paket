# 🔑 Paket - A vault to packaging and encrypt your files in golang!

 **Godoc**: [CMD Tool](https://godoc.org/github.com/SeanTolstoyevski/paket) | [pengine (paket engine)](https://godoc.org/github.com/SeanTolstoyevski/paket/pengine/)

**Main Repository**: [github.com/SeanTolstoyevski/paket](https://github.com/SeanTolstoyevski/paket/)

* * * * * *

Welcome to  Paket!

We recommend that you take a look at the items below before using this module.  
The world of encryption and encryption is a **sensitive topic**. It is important to know what you are doing and what this module actually does.

* What does it do, what can I use it for?

Imagine you are producing a game. You will probably have carefully designed animations and sound effects. You do not want users to receive this data. If we think for this scenario; The package encrypts the files in the specified folder using AES with a key you specify. And it combines all encrypted files into a single file. Calculates the hash of the encrypted and unencrypted version of the file. Saves to a table. This is a little shield for people trying to deceive you.  
Then, you can easily retrieve the decoded or encrypted version of your file from the encrypted file.  
Normally, you should create a system to securely encrypt and decrypt your files.  
This is a ready system :) .

* Is it really secure? How secure is it?

To be honest, the person who wants to get the data can crack anything if he or she  tries. Especially if the program you are distributing runs directly on the user's computer and all data is with the program. However, what AES and Package do is complex enough. Don't Remember, every executable file is sensitive to disassembly.  
You can pass your files through other complex processes before encrypting them. However, this causes your program to load files into memory slowly at run time.

* What encryption algorithm does it use?

AES CFB.  
If enough people write to add new algorithms, we will add new algorithms to the extent that golang supports it.

## Installation:

This module consists of two parts:
1. CMD tool - command- ine tool for encrypting and packaging files.
2. "pengine" (paket engine) - subfolder that provides low-level APIs (reading encrypted datas, verifications etc...).

To use Paket, you need to create a package file with the cmd tool.

You can install it like a normal golang module:  
`go get github.com/SeanTolstoyevski/paket`

### Usage

```
cmd>paket -help
Usage of paket.exe:
  -f string
        Folder containing files to be encrypted.
                        Note: Your original files are not deleted.
        It is not recursive, Subfolders is not encrypted.
  -h    Writes hash of original and encrypted versions of the files to table.
        This is required for security.
        If left null, hash checks will not work. (default true)
  -k string
        Key for encrypting files. It must be 16, 24, or 32 lenght in bytes.
        If this parameter is null, the tool generates one randomly byte  and prints value to the console. (default "8882af0e97faefe11953660eec14ea94")
  -o string
        The file to which your encrypted data will be written.
         If there is a file with the same name, you will be warned. (default "data.pack")
  -s    prints progress steps to the console. For example, which file is currently encrypting, etc. (default true)
  -t string
        The go file to be written for Paket to read.
         When compiling this file, you must import it into your program.
         It is created as 'package main.' (default "PaketTable.go")
```

**Warning**: If you leave the key blank, the system randomly generates a key.

Make sure the gopath / bin directory is appended to the path variables.

`paket -f=mydatas -o=data.dat`

Example output:

```
❗❕ Warning! Your random key. Please note: 092f8e0b25b0eeea32037e716dfcf2bc
3 files were found in mydatas folder.
Comedy of Errors (complete text) - Shakespeare.txt file is encrypting. Size: 0.9117 MB
George Orwell - Animal Farm.pdf file is encrypting. Size: 5.5276 MB
openal_soft_readme.md file is encrypting. Size: 0.0288 MB
```

If you don't want the details, you can pass  `-s=0` parameter.  
so:  
`paket -f=mydatas -o=data.dat -s=0`

Next, a go file like this is created.  
This is the table that keeps the information of your files.

Default name: PaketTable.go (see `paket -help`)

```
//important: You can edit this file. However, you need to know what you are doing.
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
	foldername = "datas"
	Data["Comedy of Errors (complete text) - Shakespeare.txt"] = paket.Values{"0", "91189", "91173", "91189", "2aa62dd2d930ed5d8e1c3a33fba4d8525e16448b12d567f5808452b94cacf693", "063b28be3d49e30710546c06b845e87ef9af811f01f7ef716be1f4516657d2d3"}
	Data["George Orwell - Animal Farm.pdf"] = paket.Values{"91189", "643961", "552756", "552772", "2d8d5810046a78daea56adcf73497b6f331023a0a2cb700db4bb029ca1425573", "86a5e5508ce4f8912f6f62b7c06c51134beb86722fa6ba670751ce727c3e081f"}
	Data["openal_soft_readme.md"] = paket.Values{"643961", "646858", "2881", "2897", "4034ec4242e7a700e2586f6520941599230e7bc8509ca60950e570df213c49ae", "0f16e1d5e7bbc82b1cb067190db8abc6aa8f00507395710095cc5cd45deb4d2a"}
}
```

**Great**, we created our first package. We're going to write some code now.  
If you want you can examine the codes in the [examples folder](https://github.com/SeanTolstoyevski/paket/examples).

## ⛑ If you like this:

* 📝🖊 Please consider creating a PR or emailing me for grammatical errors and other language issues in documents. English is **not my native language**. I'm advancing.

* 🤗 If you can test for Linux and Darwin, that would be a great  for me. I am a blind software developer. I cannot set up an environment in Linux that can develop and test these projects. Linux's accessibility is not as good as Windows. No mac. I'll try to test as much as possible with the built-in `go test` though.

* 💰🤑 If you don't want to do any of them and want to give financial support (like a cup of french press), you can send an e-mail

* 🌟⭐ And send a star. This is the greatest favor. I'm looking for a job. Employers are not so optimistic in Turkey against disabled users. But any good project with a few stars can win over employers' hearts.


## FAQ:

* Why name **Paket**?

> Package in English, "paket" in Turkish.  
I was looking for a module where I could package my data. This name came first to mind.

* Which Go versions are compatible?

> Tested with Go 15.3 64 bit on windows 10 64 bit.  
Probably compatible up to go1.12-13.
