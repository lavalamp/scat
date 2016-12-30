package cpprocs_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"

	"secsplit/checksum"
	"secsplit/cpprocs"
)

func TestCat(t *testing.T) {
	const (
		hashData = "x"
		hashStr  = "2d711642b726b04401627ca9fbac32f5c8530fb1903cc4db02258717921a4881"
		data     = "xxx"
	)
	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	cat := cpprocs.NewCat(dir)

	// Ls() empty
	hashes, err := cat.Ls()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(hashes))

	// NewCmd()
	hash := checksum.Sum([]byte(hashData))
	cmd, err := cat.NewCmd(hash)
	assert.NoError(t, err)
	cmd.Stdin = strings.NewReader(data)
	err = cmd.Run()
	assert.NoError(t, err)
	path := filepath.Join(dir, hashStr)
	fdata, err := ioutil.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, data, string(fdata))

	// Ls()
	hashes, err = cat.Ls()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(hashes))
	assert.Equal(t, []checksum.Hash{hash}, hashes)
}
