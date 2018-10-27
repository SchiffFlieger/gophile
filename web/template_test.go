package web

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/SchiffFlieger/gophile/auth"
	"github.com/stretchr/testify/assert"
)

func TestCreateFileData(t *testing.T) {
	user := "testuser"
	tm := time.Now().Format("02.01.2006 15:04:05")
	os.MkdirAll(path.Join(auth.USERDIR, user, "testing", "UDH"), os.ModePerm)
	os.MkdirAll(path.Join(auth.USERDIR, user, "testing", "PAN"), os.ModePerm)
	os.MkdirAll(path.Join(auth.USERDIR, user, "testing", "MAP"), os.ModePerm)
	f, _ := os.Create(path.Join(auth.USERDIR, user, "testing", "testfile"))
	f.Close()

	result := CreateFileData(user, path.Join("testing"), "")
	fdata := []FileTemplateData{
		{IsDir: true, FileName: "MAP", FileDate: tm, FileSize: "-"},
		{IsDir: true, FileName: "PAN", FileDate: tm, FileSize: "-"},
		{IsDir: true, FileName: "UDH", FileDate: tm, FileSize: "-"},
		{IsDir: false, FileName: "testfile", FileDate: tm, FileSize: "0 B"},
	}
	expected := GophileTemplateData{CurrentPath: "testing", User: user, ActionText: "",
		FileTable: fdata}

	assert.Equal(t, expected, result)

	os.RemoveAll(auth.USERDIR)
}

func TestIntToBytes(t *testing.T) {
	var b, kb, mb, gb int64
	b = 347
	kb = 93 * 1024
	mb = 4 * 1024 * 1024
	gb = 74 * 1024 * 1024 * 1024

	assert.Equal(t, "347 B", intToBytes(b))
	assert.Equal(t, "93 KB", intToBytes(kb))
	assert.Equal(t, "4 MB", intToBytes(mb))
	assert.Equal(t, "74 GB", intToBytes(gb))
}
