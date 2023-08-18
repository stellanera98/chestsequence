package main

import (
	"encoding/json"
	"os"
)

func GetJson(filename string, obj interface{}) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, &obj)
	return err
}

func SaveJson(filename string, obj interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	cont, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(cont)
	if err != nil {
		return err
	}
	return file.Close()
}
