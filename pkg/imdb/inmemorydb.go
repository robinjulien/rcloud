package main

import (
	"encoding/gob"
	"errors"
	"os"
	"sync"
)

var (
	ErrorNotExists error = errors.New("No value found for the key given.")
	ErrorType      error = errors.New("Invalid type.")
	ErrorFile      error = errors.New("Error when reading/writing file")
	ErrorNotInit   error = errors.New("DB not init.")
	ErrorData      error = errors.New("Invalid data.")
)

type DataType map[string]interface{}

type DB struct {
	mutex sync.Mutex
	data  DataType
}

func New() *DB {
	return &DB{
		mutex: sync.Mutex{},
		data:  DataType{},
	}
}

func (db *DB) ReadFile(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		return ErrorFile
	}

	decoder := gob.NewDecoder(file)
	var data DataType
	err = decoder.Decode(&data)

	if err != nil {
		return ErrorData
	}

	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.data = data
	return nil
}

func (db *DB) WriteFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	if err != nil {
		return ErrorFile
	}

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(db.data)

	if err != nil {
		return ErrorData
	}

	return nil
}

func (db *DB) Get(key string) (interface{}, error) {
	if db.data != nil {
		if val, ok := db.data[key]; ok {
			return val, nil
		}
		return 0, ErrorNotExists
	}
	return nil, ErrorNotInit
}

func (db *DB) GetSafe(key string) (interface{}, error) {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		if val, ok := db.data[key]; ok {
			return val, nil
		}
		return 0, ErrorNotExists
	}
	return nil, ErrorNotInit
}

func (db *DB) GetInt(key string) (int, error) {
	tmp, err := db.Get(key)

	if err != nil {
		return 0, err
	}

	val, ok := tmp.(int)

	if ok {
		return val, nil
	} else {
		return 0, ErrorType
	}
}

func (db *DB) GetString(key string) (string, error) {
	tmp, err := db.Get(key)

	if err != nil {
		return "", err
	}

	val, ok := tmp.(string)

	if ok {
		return val, nil
	} else {
		return "", ErrorType
	}
}

func (db *DB) Set(key string, value interface{}) error {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		db.data[key] = value
		return nil
	}

	return ErrorNotInit
}

func (db *DB) Remove(key string) error {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		delete(db.data, key)
		return nil
	}

	return ErrorNotInit
}

func (db *DB) Exists(key string) (bool, error) {
	if db.data != nil {
		_, ok := db.data[key]
		return ok, nil
	}

	return false, ErrorNotInit
}
