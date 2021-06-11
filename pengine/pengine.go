// Copyright (C) 2021 SeanTolstoyevski -  mailto:seantolstoyevski@protonmail.com
// The source code of this project is licensed under the MIT license.
// You can find the license on the repo's main folder.
// Provided without warranty of any kind.

//Package pengine low-level APIs for paket.
//
// Before using it, you need to create a file with the cmd tool. (If you are not creating a new tool or API).
//
// Users do not need functions and structures other than New and Paket methods.
//
//Other exported functions and variables are for the cmd tool.
// If you only want to read the package created with the cmd tool,
// you can create a new Paket method with New().
package pengine

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"os"
	"sync"

	"golang.org/x/crypto/pbkdf2"
)

var (
	// If there is no data in the map sent to New, the functions you use will return this error.
	ErrMinimumMapValue = errors.New("map cannot be less than 1 in length")

	// ErrInvalidMode returned is if the encryption mode is invalid or not currently supported.
	ErrInvalidMode = errors.New("invalid encrypt/decrypt mode")

	//
	ErrShortData = errors.New("data too short")

	//
	ErrNotFound = errors.New("paket not found")
)

// type declaration for map values.
type Values struct {
	// start position
	StartPos int

	// end position
	EndPos int

	// length of the original file.
	OriginalLenght int

	// length of the encrypted data.
	// It can also be calculated as the original length + aes.BlockSize.
	EncryptLenght int

	// Hash of the original file.
	//
	// For us to trust that the decrypted data is correct data.
	// It can be generated with a hash function such as sha256, sha512.
	HashOriginal []byte

	// Hash of encrypted data.
	// A guarantee that the encrypted data has not been changed.
	HashEncrypt []byte

	// for gcm mode
	// nil can be write  if GCM is not used.
	//
	// Usually nonce is added at the beginning of the first GCM block.
	// It may be added as an option in a future release.
	// It is currently being write to the table with the cmd tool.
	Nonce []byte
}

// type definition for the Paket.
//
// Paket reads the requested file through this map.
//
// string refers to the file name, values refers to information about the file (length, sha value etc.).
//
// minimum value is 1.
type Datas map[string]Values

// CreateRandomBytes generates random bytes of the given size.
//The maximum value should be 32 and the minimum value should be 16.
//
//Used to generate a random key if the user has not specified a key. (for cmd tool)
//
// Returns error for the wrong size or  creating bytes.
func CreateRandomBytes(l uint8) ([]byte, error) {
	if l < 16 || l > 32 {
		return nil, errors.New("minimum value for l is 16, maximum value for l is 32")
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
// Key must be 16, 24 or 32 length.
// Otherwise, the cypher module returns an error.
//
// If the data is encrypted with GCM mode selected, you should pass  the nonce.
// For other modes it can be nonce nil.
// Paket does not add the nonce to the beginning of the block.
// Nonce is written to the table by the cmd tool.
// If the package was created using the cmd tool with  selecting GCM,
// nonce will be saved in   table.
//
// No authendication is provided in any mode except GCM mode.
// You have to provide these implementations yourself.
// paket  relies on hash generators for authendication.
//
// You can compare the data sended  to the function with the output data.
// It might be a good idea to make sure it's working properly.
//
//If everything is working correctly, it returns an encrypted bytes and nil error.
func Encrypt(key, nonce, data []byte, mode MODE) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	ciphertext := []byte{}
	v := []byte{}
	if mode != MODEGCM {
		ciphertext = make([]byte, aes.BlockSize+len(data))
		v = ciphertext[:aes.BlockSize]
	} else {
		ciphertext = nil
		v = nil
	}

	if _, err := io.ReadFull(rand.Reader, v); err != nil {
		return nil, err
	}

	switch mode {
	case MODECBC:
		cbcMode := cipher.NewCBCEncrypter(block, v)
		cbcMode.CryptBlocks(ciphertext[aes.BlockSize:], data)
		return ciphertext, nil

	case MODECFB:
		s := cipher.NewCFBEncrypter(block, v)
		s.XORKeyStream(ciphertext[aes.BlockSize:], data)
		return ciphertext, nil

	case MODECTR:
		s := cipher.NewCTR(block, v)
		s.XORKeyStream(ciphertext[aes.BlockSize:], data)
		return ciphertext, nil

	case MODEOFB:
		s := cipher.NewOFB(block, v)
		s.XORKeyStream(ciphertext[aes.BlockSize:], data)
		return ciphertext, nil

	case MODEGCM:
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		return aesGCM.Seal(nil, nonce, data, nil), nil

	default:
		return nil, ErrInvalidMode
	}

}

// Decrypt decrypts the encrypted data with the key.
//
// It doesn't matter whether you have the correct key or not. It decrypts data with the key given under any condition.
// So you should compare it with the original data with a suitable hash function (see sha256, sha512 module...).
// Otherwise, you can't be sure it is returning the correct data.
//
// If everything is working correctly, it returns  decrypted bytes and nil error.
func Decrypt(key, nonce, data []byte, mode MODE) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, ErrShortData
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	v := data[:aes.BlockSize]

	var raw []byte
	if mode == MODEGCM {
		raw = nil // make([]byte, len(data)-16)
	} else {
		raw = make([]byte, len(data)-aes.BlockSize)
	}

	switch mode {
	case MODECBC:
		if len(data)%aes.BlockSize != 0 {
			return nil, errors.New("cbc: data is not a multiple of the block size")
		}

		cbcMode := cipher.NewCBCDecrypter(block, v)
		cbcMode.CryptBlocks(raw, data[aes.BlockSize:])
		return raw, nil

	case MODECFB:
		modeCFB := cipher.NewCFBDecrypter(block, v)
		modeCFB.XORKeyStream(raw, data[aes.BlockSize:])
		data = nil
		return raw, nil

	case MODECTR:
		modeCTR := cipher.NewCTR(block, v)
		modeCTR.XORKeyStream(raw, data[aes.BlockSize:])
		data = nil
		return raw, nil

	case MODEOFB: // OFB
		modeOFB := cipher.NewOFB(block, v)
		modeOFB.XORKeyStream(raw, data[aes.BlockSize:])
		data = nil
		return raw, nil

	case MODEGCM:
		aesGCM, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		ret, err := aesGCM.Open(nil, nonce, data, nil)
		if err != nil {
			data = nil
			return nil, err
		}
		data = nil
		return ret, nil

	default:
		return nil, ErrInvalidMode
	}

}

// Paket that keeps the information of the file to be read.
// It should be created with New.
type Paket struct {
	//
	key []byte

	//
	paketFileName string

	//
	mode MODE

	//
	table Datas

	// created for access the file.
	// This value is opened by New with filename parameter.
	// file released with the Close function.
	file *os.File

	// Used to prevent conflicts in GetFile. For files requested at the same time.
	globMut sync.Mutex
}

type Option struct {
	// Key value for reading the file's data
	Key []byte

	// PBDFK2 iteration
	Iteration uint

	Salt string

	// paket file path
	PaketFile string

	// encrypt/decrypt mode
	Mode MODE

	// Map value that keep the information of files in Paket.
	// It must be at least 1 length.
	// Otherwise, panic occurs at runtime.
	//
	// Usually created by the cmd tool.
	Table Datas
}

// New Creates a new Paket.
// This method should be used to read the files.
//
// key parameter refers to the encryption key.
//
// After getting all the data you need, should be terminated with  Close.
func New(o Option) (*Paket, error) {
	if !Exists(o.PaketFile) {
		return nil, ErrNotFound
	}

	f, err := os.Open(o.PaketFile)
	if err != nil {
		return nil, err
	}

	fInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if fInfo.Size() < 13+16 {
		return nil, errors.New("very short file")
	}

	p := new(Paket)
	p.file = f
	p.table = o.Table
	if o.Iteration < 4096 {
		o.Iteration = 4096
	}
	p.key = pbkdf2.Key(o.Key, []byte(o.Salt), int(o.Iteration), 32, sha256.New)
	p.paketFileName = o.PaketFile
	p.mode = o.Mode
	return p, nil
}

// GetFile returns the content of the requested file.
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
func (p *Paket) GetFile(filename string, decrypt, shaControl bool) ([]byte, bool, error) {
	file, found := p.table[filename]
	if !found {
		return nil, false, errors.New("File not found on map: " + filename)
	}

	p.globMut.Lock()
	defer p.globMut.Unlock()

	// We need the length of the encrypted data to be able to load to memory the file
	length := file.EncryptLenght
	// The position where our new file starts. Should be calculated based on the encrypted file length rather than the original file
	start := file.StartPos

	content := make([]byte, length)

	// We go to the position of file
	if _, err := p.file.Seek(int64(start), 0); err != nil {
		return nil, false, err
	}

	// We read it to the position we want. So in this case, up to the position  where the encrypted data ends. We Alocated the *content* variable
	if _, err := p.file.Read(content); err != nil {
		return nil, false, err
	}

	switch decrypt {
	case true:
		decryptedData, err := Decrypt(p.key, file.Nonce, content, p.mode)
		if err != nil {
			return nil, false, err
		}

		if shaControl {
			getOriginalHash := sha256.Sum256(decryptedData)
			tableOriginalHash := file.HashOriginal
			return decryptedData, bytes.Equal(getOriginalHash[:], tableOriginalHash), nil
		}
		return decryptedData, false, nil
	case false:
		if shaControl {
			tableHash := file.HashEncrypt
			getEncryptedHash := sha256.Sum256(content)
			return content, bytes.Equal(getEncryptedHash[:], tableHash), nil
		}
		return content, false, nil
	default:
		return content, false, nil
	}
}

// GetGoroutineSafe created to securely retrieve data when using with multiple goroutines.
// In any case, it only returns decrypted data.
//
// It does not do any hash checking.
func (p *Paket) GetGoroutineSafe(name string) ([]byte, error) {
	file, found := p.table[name]
	if !found {
		return nil, errors.New("File not found on map: " + name)
	}
	length := file.EncryptLenght
	encryptedLenght, _ := p.GetLen()
	if length > encryptedLenght[1] {
		return nil, errors.New("more length than file size")
	}
	start := file.StartPos

	f, err := os.Open(p.paketFileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Seek(int64(start), 0); err != nil {
		return nil, err
	}
	content := make([]byte, length)
	if _, err := f.Read(content); err != nil {
		return nil, err
	}
	decryptedData, err := Decrypt(p.key, file.Nonce, content, p.mode)
	if err != nil {
		content = nil // I don't understand what the gc of Go does sometimes. A guarantee
		return nil, err
	}

	content = nil // I don't understand what the gc of Go does sometimes. A guarantee
	return decryptedData, nil
}

// GetLen Returns the original and encrypted lengths of all files contained in Paket.
// 0 index refers to the original, 1  index to the encrypted data.
// In the meantime, no control is made. The same will return as the values are written into the table.
//
// Normally values should be in bytes.
//
// returns an error if length is less than 1(see ErrMinimumMapValue). This case, other  things are 0.
func (p *Paket) GetLen() ([2]int, error) {
	values := [2]int{}
	if len(p.table) < 1 {
		return values, ErrMinimumMapValue
	}
	for _, value := range p.table {
		values[0] += value.OriginalLenght
		values[1] += value.EncryptLenght
	}
	return values, nil
}

// Close Closes the opened Paket.
//
// Use this function when all your transactions are done (so you shouldn't use it with defer or something like that).
// Otherwise, you must create a new Paket method.
//
// When you call Close, you cannot access the Package again.
//
// Returns error for unsuccessful events.
func (p *Paket) Close() error {
	err := p.file.Close()
	p.key = nil
	p.table = nil
	p.file = nil
	p = nil
	return err
}

// Exists a guarantee about the existence of file.
//
// Source: https://stackoverflow.com/a/12527546/13431469
// Thanks to SO user.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
