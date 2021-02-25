package enhancedmaps

import (
	"encoding/gob"
	"errors"
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

// DataType is the type of data stored in a map
type DataType map[string]interface{}

// Map is the map strcture that holds the data and the mutex
type Map struct {
	mutex sync.Mutex
	data  DataType
}

// New creates a new instance of the map.
func New() *Map {
	return &Map{
		mutex: sync.Mutex{},
		data:  DataType{},
	}
}

// ReadFile reads a file and loads its data in the map
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

// WriteFile writes the data of the map to a file
func (db *Map) WriteFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return ErrorFilePermission
		}
		return ErrorFile
	}

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(db.data)

	if err != nil {
		return ErrorData
	}

	return nil
}

// Get returns the elements of the map at the given key
// If an error occurs, it returns nil and the error
func (db *Map) Get(key string) (interface{}, error) {
	if db.data != nil {
		if val, ok := db.data[key]; ok {
			return val, nil
		}
		return 0, ErrorNotExist
	}
	return nil, ErrorNotInit
}

// GetSafe is like Get, but lock and unlock using mutex
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

// GetInt get the element at the given key, and type assert it as an int
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

// GetString get the element at the given key, and type assert it as a string
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

// Set set an element at the given key
func (db *Map) Set(key string, value interface{}) error {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		db.data[key] = value
		return nil
	}

	return ErrorNotInit
}

// Remove removes an element at a given key
func (db *Map) Remove(key string) error {
	if db.data != nil {
		db.mutex.Lock()
		defer db.mutex.Unlock()
		delete(db.data, key)
		return nil
	}

	return ErrorNotInit
}

// Exists tells if there is an element at the given key
func (db *Map) Exists(key string) (bool, error) {
	if db.data != nil {
		_, ok := db.data[key]
		return ok, nil
	}

	return false, ErrorNotInit
}

// GetUnsafeData returns the raw DataType under the map
func (db *Map) GetUnsafeData() DataType {
	return db.data
}

// Keys returns all the keys init in the map, if error empty slice
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

// Values returns all the values of the elements in the map, if error empty slice
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

// Entries returns all the pairs key/values of the map, if error empty slices
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
