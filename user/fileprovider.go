package user

import (
	"encoding/csv"
	"fmt"
	"github.com/g-leon/go-data-collector/prn"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type reader interface {
	Read() ([]string, error)
}

// fileProvider is a prn based implementation
// of the provider.UserFetcher interface
type fileProvider struct {
	path string
}

func NewFileProvider(path string) *fileProvider {
	return &fileProvider{path: path}
}

// TableNames returns the list of names
// of the tables that can be found at
// this specific endpoint.
func (fp *fileProvider) TableNames() []string {
	tn := make([]string, 0)

	files, err := ioutil.ReadDir(fp.path)
	if err != nil {
		log.Printf("Unable to open path \"%s\": %s", fp.path, err)
		return tn
	}

	for _, f := range files {
		if !f.IsDir() && fp.isSupportedFileType(fp.fileTypeOf(f.Name())) {
			tn = append(tn, f.Name())
		}
	}

	return tn
}

// GetTable returns the users data
// from the table identified by the
// name given as parameter or an
// error if the table could not be
// found.
func (fp *fileProvider) GetTable(fileName string) ([]*Model, error) {
	ft := fp.fileTypeOf(fileName)

	f, err := os.Open(fp.buildFilePath(fileName))
	defer f.Close()
	if err != nil {
		log.Printf("Unable to open file: %s, %s", fileName, err)
		return nil, err
	}

	var r reader
	switch ft {
	case "csv":
		r = csv.NewReader(f)
	case "prn":
		headers := []string{"Name", "Address", "Postcode", "Phone", "Credit Limit", "Birthday"}
		r = prn.NewReader(f, headers)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ft)
	}

	res, err := fp.loadFile(r)
	if err != nil {
		log.Printf("Unable to load file: %s, %s", fileName, err)
		return nil, err
	}

	return res, nil
}

// fileTypeOf takes a file name and returns it's type.
func (fp *fileProvider) fileTypeOf(fileName string) string {
	parts := strings.Split(fileName, ".")
	return strings.ToLower(parts[len(parts)-1])
}

// loadFile reads all records from r and loads them
// into a slice of *Model.
func (fp *fileProvider) loadFile(r reader) ([]*Model, error) {
	models := make([]*Model, 0)

	skipHeader := true
	for {
		row, err := r.Read()
		if skipHeader {
			skipHeader = false
			continue
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return nil, err
		}

		models = append(models, &Model{
			Name:        row[0],
			Address:     row[1],
			Postcode:    row[2],
			Phone:       row[3],
			CreditLimit: row[4],
			Birthday:    row[5],
		})
	}

	return models, nil
}

// buildFilePath returns system specific absolute path
// of the file identified by given name.
func (fp *fileProvider) buildFilePath(fileName string) string {
	return fp.path + string(os.PathSeparator) + fileName
}

// isSupportedFileType checks if file is of type
// csv or prn.
func (fp *fileProvider) isSupportedFileType(fileName string) bool {
	return fileName == "csv" || fileName == "prn"
}
