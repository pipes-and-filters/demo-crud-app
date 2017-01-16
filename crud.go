package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"os"
	"regexp"

	"github.com/pipes-and-filters/filters"
	"github.com/ugorji/go/codec"
)

var (
	chainsFile       string
	chains           filters.Chains
	ErrInvalidAuthor = errors.New("Missing or invalid data for author.")
)

func main() {
	flag.StringVar(
		&chainsFile,
		"chains",
		os.Getenv("CRUD_DEMO_CHAIN_FILE"),
		"Chains file for setting data.",
	)
	flag.Parse()
	d := codec.NewDecoder(os.Stdin, &codec.MsgpackHandle{})
	w := Wrapper{}
	err := d.Decode(&w)
	if err != nil {
		log.Fatal(err)
	}
	chains, err = filters.ChainsFile(chainsFile)
	if err != nil {
		log.Fatal(err)
	}
	if w.Method == "Create" {
		if !w.Author.Validate() {
			log.Fatal(ErrInvalidAuthor)
		}
		err := w.Author.NewId()
		if err != nil {
			log.Fatal(err)
		}
	}
	var dbop filters.Chain
	dbop, err = chains.Get(w.Method)
	if err != nil {
		log.Fatal(err)
	}
	dbe, err := dbop.Exec()
	if err != nil {
		log.Fatal(err)
	}
	var in bytes.Buffer
	dbe.SetInput(&in)
	dbe.SetOutput(os.Stdout)
	enc := codec.NewEncoder(&in, &codec.MsgpackHandle{})
	//go func() {
	err = enc.Encode(w.Author)
	if err != nil {
		log.Fatal(err)
	}
	//}()
	err = dbe.Run()
	if err != nil {
		log.Fatal(err)
	}

}

type Wrapper struct {
	Method string
	Author Author
}

type Author struct {
	FirstName string
	LastName  string
	Email     string
	Books     []Book
	Id        string
}

func (a *Author) NewId() error {
	id, err := chains.Get("Id")
	if err != nil {
		return err
	}
	ide, err := id.Exec()
	if err != nil {
		return err
	}
	var ids bytes.Buffer
	ide.SetOutput(&ids)
	err = ide.Run()
	if err != nil {
		return err
	}
	a.Id = ids.String()
	return nil
}

func (a Author) Validate() bool {
	if a.validateBooks() && a.validateEmail() && a.validateNames() {
		return true
	}
	return false
}

func (a Author) validateBooks() bool {
	if len(a.Books) <= 0 {
		return false
	}
	return true
}

func (a Author) validateNames() bool {
	if a.FirstName == "" || a.LastName == "" {
		return false
	}
	return true
}

func (a Author) validateEmail() bool {
	exp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return exp.MatchString(a.Email)
}

type Book struct {
	Title string
}
