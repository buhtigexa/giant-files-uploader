package test

import (
	"bugtigexa.giantfilesuploader.com/cmd"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
	"time"
)

func TestCreateDirectory(t *testing.T) {
	clean := func(fname string) {
		entries, err := os.ReadDir(fname)
		if err != nil {
			log.Fatal(err)
		}
		for _, e := range entries {
			filename := fname + "/" + e.Name()
			if _, err := os.Stat(filename); err != nil {
				continue
			}
			if err := os.Remove(filename); err != nil {
				log.Fatal(err)
			}
		}
		if err := os.Remove(fname); err != nil {
			log.Fatal(err)
		}
	}

	fm := cmd.NewFileManager()
	data := cmd.Data{
		FileName: "bigfile",
		Part:     2,
		Value:    []byte("this is a big file"),
		Time:     time.Now(),
	}
	dirname := data.GetDirname()
	defer clean(dirname)

	by, err := json.Marshal(&data)
	require.Nil(t, err)
	assert.True(t, len(by) >= 1)
	n, err := fm.Store(by)

	require.Nil(t, err)
	require.True(t, n > 1)
	finfo, err := os.Stat(dirname)
	require.Nil(t, err)
	assert.True(t, finfo.IsDir())
	fileName := data.GetFileName()
	finfo, err = os.Stat(fileName)
	require.Nil(t, err)
	assert.False(t, finfo.IsDir())
}
