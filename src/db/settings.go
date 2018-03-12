package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Setting struct {
	StorePath string
	Data      map[string]string
}

var instance *Setting

func Settings() *Setting {
	if instance == nil {
		instance = &Setting{}
	}

	return instance
}

func (setting *Setting) Get(value string) string {
	setting.Default()
	return setting.Data[value]
}

func (setting *Setting) Set(key, value string) {
	setting.Default()

	setting.Data[key] = value

	jsonBytes, marshalError := json.Marshal(setting.Data)

	if marshalError != nil {
		panic(marshalError)
	}

	writeError := ioutil.WriteFile(setting.StorePath, jsonBytes, os.ModePerm)

	if writeError != nil {
		panic(writeError)
	}
}

func (setting *Setting) Default() {
	if setting.Data == nil {
		data, ioError := ioutil.ReadFile(setting.StorePath)

		if ioError == nil {
			settingsError := json.Unmarshal(data, &setting.Data)

			if settingsError != nil {
				panic(settingsError)
			}

		} else {

			setting.Data = map[string]string{}
			jsonBytes, marshalError := json.Marshal(setting.Data)

			if marshalError != nil {
				panic(marshalError)
			}

			writeError := ioutil.WriteFile(setting.StorePath, jsonBytes, os.ModePerm)

			if writeError != nil {
				panic(writeError)
			}
		}

	}

}
