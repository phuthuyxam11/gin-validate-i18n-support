package utils

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
	"path/filepath"
)

func I18nInit(messageFolderPath string) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	files, err := ioutil.ReadDir(messageFolderPath)
	path, _ := filepath.Abs(messageFolderPath)

	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if filepath.Ext(filepath.Join(path, f.Name())) != ".json" {
			continue
		}
		if _, err = bundle.LoadMessageFile(filepath.Join(path, f.Name())); err != nil {
			return nil, err
		}
	}
	return bundle, nil
}
