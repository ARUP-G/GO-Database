package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	// con
	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)

type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)

	opts := Options{}

	if options != nil {
		opts = *options
	}

	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	// if the databse ale=reasy exists
	if _, err := os.Stat(dir); err != nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", dir)
		return &driver, nil
	}

	// If database not exits
	opts.Logger.Debug("Creating the databse at '%s'... \n", dir)
	return &driver, os.Mkdir(dir, 0755)
}

// Write struct method of driver struct type
func (d *Driver) Write(collection, resouce string, v interface{}) error {
	if collection == "" {
		return fmt.Errorf("missing collection ")
	}

	if resouce == "" {
		return fmt.Errorf("missing Resouce ")
	}

	mutex := d.GetOrCreateMutex(collection)
	mutex.Lock()

	defer mutex.Unlock() // at the end after write function completed

	dir := filepath.Join(d.dir, collection)
	fnlPath := filepath.Join(dir, resouce+".josn")
	tmpPath := fnlPath + ".tmp"

	if err := os.Mkdir(dir, 0755); err != nil {
		return err
	}

	// converting collected data in json
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	b = append(b, byte('\n'))

	// Writing to the created file
	if err := os.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	return os.Rename(tmpPath, fnlPath)

}

func (d *Driver) Read(collection, resource string, v interface{}) error {

	if collection == "" {
		return fmt.Errorf("missing collection ")
	}

	if resource == "" {
		return fmt.Errorf("missing resouce,  Unable to save ")
	}

	record := filepath.Join(d.dir, collection, resource)

	if _, err := stat(record); err != nil {
		return err
	}

	b, err := os.ReadFile(record + ".json")
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &v)

}

func (d *Driver) ReadAll(collection string) ([]string, error) {
	// []string -> data
	if collection == "" {
		return nil, fmt.Errorf("missing collection !! Unable to read")
	}

	// Going to folder and joning the collection
	dir := filepath.Join(d.dir, collection)

	if _, err := stat(dir); err != nil {
		return nil, err
	}

	// Read everything '_' -> err
	files, _ := os.ReadDir(dir)

	var recoads []string

	// Accessing the files accorfing to names
	for _, file := range files {
		b, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		recoads = append(recoads, string(b))
	}
	return recoads, nil
}

func (d *Driver) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)

	mutex := d.GetOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	dir := filepath.Join(d.dir, path)

	switch fi, err := stat(dir); {
	case fi == nil, err != nil:
		return fmt.Errorf("unable to find file or directory %v", path)

		// Delete whole folder
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

		// Delete file
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir + ".json")
	}
	return nil
}

func (d *Driver) GetOrCreateMutex(collection string) *sync.Mutex {

	d.mutex.Lock()
	defer d.mutex.Unlock()
	m, present := d.mutexes[collection]

	if !present {
		m = &sync.Mutex{}
		d.mutexes[collection] = m
	}
	return m
}

func stat(path string) (fi os.FileInfo, err error) {
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path + ".json") // created dabase will have .json file
	}
	return
}

// Address
type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

// User details
type User struct {
	Name    string
	Age     json.Number // if will be string for golang but number for json
	Contact string
	Address Address
}

func main() {

	// the collection folder will be created in this folder
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	empolyees := []User{
		{"ron", "25", "12sefse", Address{"sef", "WB", "ind", "70998"}},
		{"Jon", "25", "sddse", Address{"feO", "WB", "ind", "700071"}},
		{"Vince", "31", "13fse", Address{"tyO", "WB", "ind", "733188"}},
		{"Leo", "35", "445e", Address{"GrO", "WB", "ind", "73314"}},
	}

	for _, value := range empolyees {

		// Write function
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Address: value.Address,
		})
	}

	// read all
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}

	fmt.Println(records)

	// Find
	alluser := []User{}

	for _, data := range records {
		empolyeeFound := User{}

		// As the data retrived is in json so convert
		if err := json.Unmarshal([]byte(data), &empolyeeFound); err != nil {
			fmt.Println("Error", err)
		}

		alluser = append(alluser, empolyeeFound)
	}
	fmt.Println(alluser)

	// // Delete

	// if err := db.Delete("user", "Aron"); err != nil {
	// 	fmt.Println("Error", err)
	// }

	// // Delete all

	// if err := db.Delete("user", ""); err != nil {
	// 	fmt.Println("Error", err)
	// }
}
