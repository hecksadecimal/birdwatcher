package cherp_api

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"sync"

	cookiejar "github.com/juju/persistent-cookiejar"
)

var lock sync.Mutex

// Load loads the file at path into v.
// Use os.IsNotExist() to see if the returned error is due
// to the file being missing.
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(f, v)
}

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := json.Marshal(v)
	if err != nil {
		return err
	}
	rread := bytes.NewReader(r)
	_, err = io.Copy(f, rread)
	return err
}

func SaveJar() {
	lock.Lock()
	defer lock.Unlock()
	client.Jar.(*cookiejar.Jar).Save()
}
