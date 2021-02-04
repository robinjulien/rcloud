package enhancedmaps

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"sync"
)

// Errors
var (
	ErrorNotExist       error = errors.New("no value found for the key given")
	ErrorType           error = errors.New("invalid type")
	ErrorFilePermission error = errors.New("error when reading/writing file : Permission")
	ErrorFileNotExist   error = errors.New("error when reading file : Not exist")
	ErrorFile           error = errors.New("error when reading/writing file : uknown")
	ErrorNotInit        error = errors.New("map not init")
	ErrorData           error = errors.New("invalid data")
)

type DataType map[string]interface{}

type Map struct {
	mutex sync.Mutex
	data  DataType
}

func New() *Map {
	return &Map{
		mutex: sync.Mutex{},
		data:  DataType{},
	}
}

func (db *Map) ReadFile(filename string) error {
	file, err := os.Open(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrorFileNotExist
		} else if errors.Is(err, os.ErrPermission) {
			return ErrorFilePermission
		}
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

func (db *Map) WriteFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)

	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return ErrorFilePermission
		}
		return ErrorFile
	}

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(db.data)

	if err != nil {
		fmt.Println(err)
		return ErrorData
	}

	return nil
}

func (db *Map) Get(key string) (interface{}, error) {
	if db.data != nil {
		if val, ok := db.data[key]; ok {
			return val, nil
		}
		return 0, ErrorNotExist
	}
	return nil, ErrorNotInit
}

func (db *Map) GetSafe(key string) (interface{}, error) {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		if val, ok := db.data[key]; ok {
			return val, nil
		}
		return 0, ErrorNotExist
	}
	return nil, ErrorNotInit
}

func (db *Map) GetInt(key string) (int, error) {
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

func (db *Map) GetString(key string) (string, error) {
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

func (db *Map) Set(key string, value interface{}) error {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		db.data[key] = value
		return nil
	}

	return ErrorNotInit
}

func (db *Map) Remove(key string) error {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		delete(db.data, key)
		return nil
	}

	return ErrorNotInit
}

func (db *Map) Exists(key string) (bool, error) {
	if db.data != nil {
		_, ok := db.data[key]
		return ok, nil
	}

	return false, ErrorNotInit
}

func (db *Map) GetUnsafeData() DataType {
	return db.data
}

func (db *Map) Keys() []string {
	if db.data == nil {
		return []string{}
	}

	keys := make([]string, 0, len(db.data))

	for k := range db.data {
		keys = append(keys, k)
	}

	return keys
}

func (db *Map) Values() []interface{} {
	if db.data == nil {
		tmp := make([]interface{}, 0)
		return tmp
	}

	values := make([]interface{}, 0, len(db.data))

	for _, v := range db.data {
		values = append(values, v)
	}

	return values
}

func (db *Map) Entries() ([]string, []interface{}) {
	if db.data == nil {
		tmp := make([]interface{}, 0)
		return []string{}, tmp
	}

	keys := make([]string, 0, len(db.data))
	values := make([]interface{}, 0, len(db.data))

	for k, v := range db.data {
		keys = append(keys, k)
		values = append(values, v)
	}

	return keys, values
}
