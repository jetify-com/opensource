package auth

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/pkg/errors"
)

const dirName = ".jetpack"
const jwksFileName = "jwks.json"
const cacheDuration = 1 * time.Hour

func (a *Authenticator) fetchJWKSWithCache() (*keyfunc.JWKS, error) {
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", a.Domain)
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	path := filepath.Join(wd, dirName, jwksFileName)
	// check Cache if miss, jwksJSON will be empty
	jwksJSON, err := readJWKSCache(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if jwksJSON == nil { // cache miss
		// save new keys to cache
		jwksJSON, err = saveJWKSCache(jwksURL, path)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	jwks, err := keyfunc.NewJSON(jwksJSON)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return jwks, nil
}

func readJWKSCache(path string) ([]byte, error) {
	fileInfo, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.WithStack(err)
	}
	modificationTime := fileInfo.ModTime()
	current := time.Now()
	if current.After(modificationTime.Add(cacheDuration)) {
		return nil, nil
	}
	byteContent, err := os.ReadFile(path)
	if err != nil {
		return nil, nil
	}
	return byteContent, nil
}

func saveJWKSCache(url string, path string) ([]byte, error) {
	var client http.Client
	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	out, err := os.Create(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer out.Close()
	jwks, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	fmt.Printf("resp.body: %s \n\n", jwks)
	_, err = out.Write(jwks)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return jwks, nil
}
