package main

import (
	"encoding/json"
	"fmt"
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
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}
}

// Write struct method of driver struct type
func (d *Driver) Write() error {

}

func (d *Driver) Read() error {

}

func (d *Driver) ReadAll() error {

}

func (d *Driver) Delete() error {

}

func (d *Driver) GetOrCreateMutex() *sync.Mutex {

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
	Age     json.Number // if will bw stringfor golang vut number for json
	Contact string
	Address Address
}

func main() {

	// the collection folder will be created in this folder
	dir := "./"

	db, err := New()(dir, nil)

	if err != nil {
		fmt.Println("Error", err)
	}

	empolyees := []User{
		{"Aron", "25", "12sefse", Address{"sef", "WB", "ind", "70998"}},
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
	records, err := db.readAll("users")
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

	// Delete

	if err := db.Delete("user", "Aron"); err != nil {
		fmt.Println("Error", err)
	}

	// Delete all

	if err := db.Delete("user", ""); err != nil {
		fmt.Println("Error", err)
	}
}
