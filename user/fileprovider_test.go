package user

import (
	"encoding/csv"
	"os"
	"testing"
)

func TestFileProvider_TableNames(t *testing.T) {
	f1, err := os.Create("text.csv")
	if err != nil {
		t.Error(err)
	}
	f2, err := os.Create("text.prn")
	if err != nil {
		t.Error(err)
	}
	f3, err := os.Create("text.bla")
	if err != nil {
		t.Error(err)
	}
	defer func() {
		os.Remove(f1.Name())
		os.Remove(f2.Name())
		os.Remove(f3.Name())
	}()

	fp := NewFileProvider(".")
	if len(fp.TableNames()) != 2 {
		t.Errorf("Expected %d, got %d.", 2, len(fp.TableNames()))
	}
}

func TestFileProvider_fileTypeOf(t *testing.T) {
	fp := NewFileProvider(".")

	testTable := []struct {
		in  string
		out string
	}{
		{in: "test.csv", out: "csv"},
		{in: "test.PRN", out: "prn"},
	}

	for _, tt := range testTable {
		if ft := fp.fileTypeOf(tt.in); ft != tt.out {
			t.Errorf("Expected %v, got %v", tt.out, ft)
		}
	}

}

func TestNewFileProvider_loadFile(t *testing.T) {
	f, err := os.Create("test.csv")
	defer os.Remove(f.Name())
	if err != nil {
		t.Error(err)
	}

	u := &Model{
		Name:        "Johnson John",
		Address:     "Voorstraat 32",
		Postcode:    "3122gg",
		Phone:       "020 3849381",
		CreditLimit: "10000",
		Birthday:    "01/01/1987",
	}

	record := u.Name + "," +
		u.Address + "," +
		u.Postcode + "," +
		u.Phone + "," +
		u.CreditLimit + "," +
		u.Birthday
	_, err = f.WriteString(record + "\n")
	_, err = f.WriteString(record)
	if err != nil {
		t.Error(err)
	}
	f.Close()

	fp := NewFileProvider(".")
	fo, err := os.Open(f.Name())
	defer fo.Close()
	if err != nil {
		t.Error(err)
	}
	r := csv.NewReader(fo)
	us, err := fp.loadFile(r)
	if err != nil {
		t.Error(err)
	}

	if len(us) != 1 {
		t.Errorf("Expected %v, got %v.", 1, len(us))
	}
	if us[0].Name != u.Name {
		t.Errorf("Expected %v, got %v.", u.Name, us[0].Name)
	}
	if us[0].Address != u.Address {
		t.Errorf("Expected %v, got %v.", u.Address, us[0].Address)
	}
	if us[0].Postcode != u.Postcode {
		t.Errorf("Expected %v, got %v.", u.Postcode, us[0].Postcode)
	}
	if us[0].Phone != u.Phone {
		t.Errorf("Expected %v, got %v.", u.Phone, us[0].Phone)
	}
	if us[0].CreditLimit != u.CreditLimit {
		t.Errorf("Expected %v, got %v.", u.CreditLimit, us[0].CreditLimit)
	}
	if us[0].Birthday != u.Birthday {
		t.Errorf("Expected %v, got %v.", u.Birthday, us[0].Birthday)
	}
}
