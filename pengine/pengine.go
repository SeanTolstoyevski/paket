// Copyright (C) 2021 SeanTolstoyevski -  mailto:s.tolstoyevski@protonmail.com
// The source code of this project is licensed under the MIT license.
// You can find the license on the repo's main folder.
// Provided without warranty of any kind.

//low level APIs for paket.
//
// Before using it, you need to create a file with the cmd tool. (If you are not creating a new tool or API).
//
// Users do not need functions and structures other than New and Paket methods.
//
//Other exported functions and variables are for the cmd tool.
// If you only want to read the package created with the cmd tool,
// you can create a new package method with New().
package pengine

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

// type declaration for map values.
//
// 0: start position
//
// 1: end position
//
// 2: length of the original file.
//
// 3: length of the encrypted data.
// ❕‼ Note: Gives [2]-[1] length of file. However, it may be more correct to write for security.
//
//4: Hash of the original file.
//
//5: Hash of encrypted data.
//
// 0, 1, 2 and 3 are required values. 4 and 5 should be written for security.
//
// If 4 and 5 index is null, sha controls will not work. This makes it difficult for you to know the security of your content.
type Values [6]string

// type definition for the Paket.
//
// Paket reads the requested file through this map.
//
// string refers to the file name, values refers to information about the file (lenght, sha value etc.).
type Datas map[string]Values

// generates random bytes of the given size.
//The maximum value should be 32 and the minimum value should be 16.
//
//Used to generate a random key if the user has not specified a key. (for cmd tool)
//
// Returns error for the wrong size or  creating bytes.
func CreateRandomBytes(l uint8) ([]byte, error) {
	if l < 16 || l > 32 {
		freeerr := errors.New("Minimum value for l is 16, maximum value for l is 32.")
		return nil, freeerr
	}
	res := make([]byte, l)
	_, err := rand.Read(res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Encrypt encrypts the data using the key.
//
// Uses the CFB mode.
//
// Key must be 16, 24 or 32 size.
// Otherwise, the cypher module returns an error.
//
// You can compare the data sended  to the function with the output data. It might be a good idea to make sure it's working properly.
//
//If everything is working correctly, it returns an encrypted bytes and nil error.
func Encrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	v := ciphertext[:aes.BlockSize]

	_, rerr := io.ReadFull(rand.Reader, v)
	if rerr  != nil {
		return nil, rerr 
    }

	s := cipher.NewCFBEncrypter(block, v)
s.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return ciphertext, nil
}

// It decrypts the encrypted data with the key.
//
// Uses the CFB mode.
//
// It doesn't matter whether you have the correct key or not. It decrypts data with the key given under any condition.
// So you should compare it with the original data with a suitable hash function (see sha256, sha512 module...).
// Otherwise, you can't be sure it is returning the correct data.
//
// If everything is working correctly, it returns  decrypted bytes and nil error.
func Decrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(data, data)
	return data, nil
}

// Struct that keeps the information of the file to be read.
// It should be created with New.
type Paket struct {
	// Key value for reading the file's data.
	// As a warning, you shouldn't just create a plaintext key.
	Key		[]byte
	// Map value that keep the information of files in Paket.
	// It must be at least 1 lenght.
// Otherwise, panic occurs at runtime.
	//
	// Usually created by the cmd tool.
	Table		Datas
//non-exported value created for access the file.
	// This value is opened by New with filename parameter.
	file		*os.File
}

// Creates a new Package method.
// This method should be used to read the files.
//
// key parameter refers to the encryption key. It must be 16, 24 or 32 lenght. Returns nil and error for keys of incorrect length.
//
// Panic occurs if the specified file does not exist.
//
// table parameter is defined in go file created by the cmd tool.
// There must be a minimum of 1 file in the table.
//
// After getting all the data you need, should be terminated with  Close.
func New(key []byte, file string, table Datas) (*Paket, error) {
	l := len(key)
	if l == 16 || l == 24 || l == 32 {
	if !Exists(file) {
		panic(string(file + " paket not found."))
	}

	f, err := os.Open(file)
		if err != nil {
			panic(err)
		}

		fInfo, ferr := f.Stat()
		if ferr != nil {
			panic(ferr)
		}

		if fInfo.Size() > 0 {
			return &Paket{file: f, Table: table, Key: key}, nil
		} else {
			perr := "There is no data in the file: " + f.Name()
			panic(perr)
		}
	}

	freeerr := errors.New("Key must be 16, 24 or 32 lenght.")
	return nil, freeerr
}

// Returns the content of the requested file.
//
// If the file cannot be found in the map and the length cannot be read, a panic occurs.
//
// All errors except these errors return with error.
//
// If decrypt is true, it is decrypted. If not, encrypted bytes are returned.
//
// If value of shaControl is true, the hash of the decrypted data is compared with hash of the original file.
//
// If decrypt is false and shaControl is true, the hash of the encrypted file in the table is compared with the encrypted hash of the read file.
//
// If the hash comparison is true, the second value is set to true.
//
// If hashControl is false, checks are skipped. Returns False.
//
// Both values do not have to be true. However, it may be good to generate a control mechanism like hash with your own work.
// The decrypt (bool) value has been added for convenience. As a recommendation,
// it is better to pass both values to true to this function.
func (p *Paket) GetFile(filename string, decrypt, shaControl bool) (*[]byte, bool, error) {
	file, found := p.Table[filename]
	var content []byte
	if !found {
		panic(string("dosya map'te bulunamadı: " + filename))
	}
	lenght, err := strconv.Atoi(file[3])
	if err != nil {
		panic(err)
	}
		start, err := strconv.Atoi(file[0])
	if err != nil {
		panic(err)
	}
	content = make([]byte, lenght)
	_, serr := p.file.Seek(int64(start), 0)
	if serr != nil {
		return nil, false, serr
	}
	_, rerr := p.file.Read(content)
	if rerr != nil {
		return nil, false, rerr
	}
	switch shaControl {
	case true:
		decdata, err := Decrypt(p.Key, content)
		if err != nil {
			return nil, false, err
		}
		if shaControl {
			decSha := []byte(fmt.Sprintf("%x", sha256.Sum256(decdata)))
			encSha := []byte(file[4])
			return  &decdata, bytes.Equal(decSha, encSha), nil
		}
		return &decdata, false, nil
	case false:
		if shaControl {
			forgSha := []byte(file[5])
			corgSha := []byte(fmt.Sprintf("%x", sha256.Sum256(content)))
			return &content, bytes.Equal(corgSha, forgSha), nil
		}
		return &content, false, nil
	default:
		return &content, false, nil
	}
}

// Returns the original and encrypted lengths of all files contained in Paket.
// First variable refers to the original, second variable to the encrypted data.
// In the meantime, no control is made. The same will return as the values are written into the table.
//
// Normally values should be in bytes.
//
// returns an error if lenght is less than 1. In this case, other variables are 0.
func (p *Paket) GetLen() (int, int, error) {
	var orgval int
	var encval int
if len(p.Table) > 0 {
		for _, value := range p.Table {
			// oi = original integer
			// ei = encrypted integer
			oi, _ := strconv.Atoi(value[2])
			ei, _ := strconv.Atoi(value[3])
			orgval += oi
			encval += ei
		}
		return orgval, encval, nil
	}
	freeerr := errors.New("Map cannot be less than 1 in length.")
	return 0, 0, freeerr
}

// Close Closes the opened file (see Paket.file (non-exported)).
//
// Use this function when all your transactions are done (so you shouldn't use it with defer or something like that).
// Otherwise, you must create a new Paket method.
//
// When you call Close, you cannot access the Package again.
//
// Returns error for unsuccessful events.
func (p *Paket) Close() error {
	err := p.file.Close()
	if err != nil {
		return nil
	}
	return err
}



// Source: https://stackoverflow.com/a/12527546/13431469
//
//a guarantee about the existence of file.
func Exists(name string) bool {
    if _, err := os.Stat(name); err != nil {
        if os.IsNotExist(err) {
            return false
        }
    }
    return true
}
