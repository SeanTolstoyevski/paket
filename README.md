# 🔒 Paket – A vault to packaging and encrypt/decrypt your files in golang! 🔑

[![Go Reference](https://pkg.go.dev/badge/github.com/SeanTolstoyevski/paket@master.svg)](https://pkg.go.dev/github.com/SeanTolstoyevski/paket@master) | [![Go Report Card](https://goreportcard.com/badge/github.com/SeanTolstoyevski/paket)](https://goreportcard.com/report/github.com/SeanTolstoyevski/paket)

## 🔊 Informations

Hey,  this is not for archiving your files like 7-zip.

We recommend that you take a look at the items below before using this module.  
The world of encryption is a **complex  topic**. It is important to know what you are doing and what this module actually does.

## What does it do?

Imagine you are producing a game. You will probably have carefully designed animations and sound effects. You do not want users to receive this data.  
If we think for this scenario; The package encrypts the files in the specified folder using AES with a key you specify. And it combines all encrypted files into a single file. Calculates the hash of the encrypted and unencrypted version of the file. Saves to a table. This is a little shield for people trying to deceive you.  
Then, you can easily retrieve the decoded or encrypted version of your file from the encrypted file.  
Normally you should create a system to securely encrypt and decrypt your files.  
This is a ready system 😎 .

## Installation

This module consists of two parts:
1. CMD tool – command-line tool for encrypting and packaging files.
2. "pengine" (paket engine) – subfolder that provides low-level APIs (reading encrypted datas, verifications etc...).

To create a paket you need to install  cmd tool:  
`go get -u github.com/SeanTolstoyevski/paket`

This install installs all dependencies. Golang compiles the package tool to your GOPATH/bin directory.

This command gets the latest tagged version (git tag).

To install the latest commit:  
`go get github.com/SeanTolstoyevski/paket@master`

You can type `paket -help` to be sure of the installation.  
If you don't see anything, the **paket could not be installed** Or **missing gopath/bin path**.

## First paket creation and CMD tool

We can read the **help text** to understand some things.  
The help text is simple and self explanatory.

```cmd
...>paket -help
Usage of paket:
  -a    anonymize file names. For example, the ''lion.zip'' file is written to the table with a name such as ''201bce5f''
  -f string
        Folder containing files to be encrypted. It is not recursive, Subfolders is not encrypted.
  -i uint
        Iteration count for pbkdf2. For less than 4096, 4096 will be selected.
        For modern CPUs values like 100000 may be appropriate. (default 4096)
  -k string
        Key for encrypting files. If this parameter is null, the tool generates one randomly byte  and prints value to the console.
  -m string
        The mode to be selected for encryption. Currently ''CFB'', ''CTR'', ''GCM'' and ''OFB'' are supported. (default "gcm")
  -o string
        The file to which your encrypted data will be written. If there is a file with the same name, you will be warned. (default "data.pack")
  -s    prints progress steps to the console. For example, which file is currently encrypting, etc. (default true)
  -t string
        The go file to be written for Paket to read. When compiling this file, you must import it into your program.
        It is created as "package main." (default "PaketTable.go")
```

***

You can review the list  below for more information and to know what's going on behind.

* `-a` – Anonymizing Filenames

Go compiler leaks many strings during compilation. You can view these strings in a simple hex editor or a code editor like Notepad++.  
When the names of your files are guessed, it's easier for those trying to tamper with the program.  
This was not designed to make the process impossible. Just an extra step.  
It can be enabled with `-a=1`. This will create a file named **anonymization-information.txt** in your working directory.  
Note: these names  are completely randomized. It is nothing like the hex encoding of the filename.  
Important note: giving up the readable names of your files can complicate the writing of the program.

* `-f` – Folder To Pack And Encrypt

The folder with the files we want to package.  
Subfolders will not be included.  
The name of these files is written to the table without the name of the folder.  
So when you think about it, the data1.eng file in the /datas folder is not written as data/data1.eng.  
If you suspect your filenames have been leaked and their purpose has been compromised, you can examine the "-a" flag.

* `-i` – Iteration for PBDFK2

In modern technologies, using a plaintex key is equivalent to suicide.  
PBDFK2 creates a much more complex brute forge scenario by repeatedly hashing the key in the specified number loop. For example, instead of plaintex key, a key that has been hashed 12000 times is used.  
The person trying to guess the key must know the iteration, find the salt, and guess the hash function correctly. All of this complicates the process.  
You can choose an iteration number by performing the appropriate tests according to the architecture you are targeting.  
For modern CPUs, hashing and loops appear to be simple functions. For this reason, values above 50000 can be considered good. However, relying only on PBDFK2 is not very accurate either.

* `-k` – AES Encryption Key

key to use for AES encryption.  
Paket currently only **uses AES256**.  
If it blank, the tool generates random bytes. So you can't just use it for packaging.  
There is no minimum character entry or maximum character entry limit.  
However; It is your responsibility to generate a complex, punctuated, mixed case key.  
Important note 1: When you forget this key, there is no way to access any data.  
Important note 2: Random or any key you specify will not be written to any file. The management of your keys belongs to you.

* `-m` – AES Encryption Mode

Allows you to choose one of the AES encryption modes.  
Each of these encryption standards has different advantages and different usage scenarios.  
As this topic is complex and lengthy enough, it is left to the user to make the right decision.  
However, **GCM is a good choice** as it supports embedded authendication and parallelism.

## Examples

You should visit the [examples folder](https://github.com/SeanTolstoyevski/paket/tree/master/examples) to see some use cases, how it works, and more.

## 😋 If you like this

* 📝🖊 Please consider creating a PR or emailing me for grammatical errors and other language issues in documents. English is **not my native language**.
	- And we have a few not fixed issues right now.  [🤔 Any ideas or code for these](https://github.com/SeanTolstoyevski/paket/blob/master/developing_and_contribute.md).
* 💰🤑 If you don't want to do any of them and want to give financial support (like a cup of french press), you can send an e-mail
* 🌟⭐ And send a star. This is the greatest favor. I'm looking for a job. Employers are not so optimistic in Turkey against disabled users. But any good project with a few stars can win over employers hearts.

## 🤔 FAQ

* **Q**: Is it really secure? How secure is it?

**A**: Frankly the person who wants to get the data can crack anything if tries.
Especially if the program you are distributing runs directly on the user's computer and all data is with the program.
However, what AES and Package do is complex enough.  
Don't Remember, **every executable file is sensitive to disassembly**.  
You can pass your files through other complex processes before encrypting them. However, this causes your program to load files into memory slowly at run time.

* **Q**: What encryption algorithm does it use?

**A**: AES CFB, CTR, GCM, OFB.  
If enough people write to add new algorithms, we will add new algorithms to the extent that golang supports it.

* **Q**: So this is the file format with a TOC?

**A**: No, there is no file format. The whole process is executed by the Golang runtime.  
Anyone trying to inspect the file will only see the binary data.  
The person who wants to get the data must disassemble  the Golang executable program.  
The only downside is that the more files are saved to the table, the greater the memory requirement. I guess nobody writes 10000 files.

