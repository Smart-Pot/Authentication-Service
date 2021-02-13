package oauth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	c := m.Run()
	os.Exit(c)
}



func TestReadConfig(t *testing.T) {
	wd,_ := os.Getwd()
	err :=ReadConfig(filepath.Join(wd,"..","..","test_files"))
	/*
	Test file looks like this:
	Google:
		ClientID: testclient
		ClientSecret: testsecret

	*/
	assert.Nil(t,err)
	assert.Equal(t,"testclient",Config.Google.ClientID)
	assert.Equal(t,"testsecret",Config.Google.ClientSecret)
}
