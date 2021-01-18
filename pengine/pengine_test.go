package pengine_test

import (
	"bytes" // for equal function
	"testing"
	"github.com/SeanTolstoyevski/paket/pengine"
)

var globpaket *pengine.Paket

var testdata = []byte("this is paket test")
var testkey = []byte("b9e0a869abb532bfccd6bd6a8e624753")
var 	data = make(pengine.Datas)

func TestCreateRandomBytesLenght(t *testing.T) {
	bytes, _ := pengine.CreateRandomBytes(32)
	want := 32
	got := len(bytes)
	if got != want {
		t.Error("CreateRandomBytes(16) want 32.")
	}
}

func TestEncryptionandDecryption(t *testing.T) {
	encdata, err	 := pengine.Encrypt(testkey, testdata)
	if err != nil {
		t.Error("Encrypting error", err)
	}
	decdata, err := pengine.Decrypt(testkey, encdata)
	if err != nil {
		t.Error("Decrypting error", err)
	}
	if !bytes.Equal(decdata, testdata) {
		t.Error("Original data could not be found in decrypting:", string(decdata))
	}
}

func TestNew(t *testing.T) {

	datas, err := pengine.New(testkey, "testdata/data.pack", data)
	if err != nil {
		t.Error("Paket creation error", err)
	}
	enclen, orglen, err := datas.GetLen()
	if err != nil {
		t.Error("Error in GetLen.", err)
	}
	if enclen == 0 || orglen == 0  || orglen < 1 || enclen < 1 {
		t.Error("Error in lenght  query. The wrong value.", enclen, orglen)
	}
	_, sha, err := datas.GetFile("George Orwell - Animal Farm.pdf", true, true)
	if err != nil {
		t.Error("GetFile in error", err)
	}
	if !sha {
		t.Error("Third parameter is true. So the hash value must be true.", sha)
	}
	_, rsha, err := datas.GetFile("George Orwell - Animal Farm.pdf", true, false)
	if err != nil {
		t.Error("Error in test 'shaControl false'.", err)
	}
	if rsha {
		t.Error("sha should be false but true. Because third value is false.")
	}
	// I will add all the scenarios as I find time.
	cerr := datas.Close()
	if cerr != nil {
		t.Error(cerr )
	}
}

func TestNewwithShortKey(t *testing.T) {
	paket, err := pengine.New([]byte("kisabirtestaslinda"), "testdata/data.pack", data)
	if paket != nil {
		t.Error("Error should return for short key.", err)
	}
}

func init() {
	data["Comedy of Errors (complete text) - Shakespeare.txt"] = pengine.Values{"0", "91189", "91173", "91189", "2aa62dd2d930ed5d8e1c3a33fba4d8525e16448b12d567f5808452b94cacf693", "f5bac206eca4f1ecea0406063764663ff27f91b584186ad32e8eaa0ed4de4216"}
	data["George Orwell - Animal Farm.pdf"] = pengine.Values{"91189", "643961", "552756", "552772", "2d8d5810046a78daea56adcf73497b6f331023a0a2cb700db4bb029ca1425573", "f7a23c6bd96c61ff907d8f175e0268eb90c9e7cebb233090c9e656d8eba5d0e3"}
	data["openal_soft_readme.md"] = pengine.Values{"643961", "646858", "2881", "2897", "4034ec4242e7a700e2586f6520941599230e7bc8509ca60950e570df213c49ae", "762d323774b9c84882e3ad030e767dcbf4e1d92c26f7edb192cfbd2cdad5a214"}

	var err error
	globpaket, err = pengine.New(testkey, "testdata/data.pack", data)
	if err != nil {
		panic(err)
	}
}