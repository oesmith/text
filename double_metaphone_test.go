package text

import (
	"encoding/csv"
	"io"
	"os"
	"testing"
)

func TestDoubleMetaphone(t *testing.T) {
	f, err := os.Open("double_metaphone.csv")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()
	r := csv.NewReader(f)
	i := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err.Error())
		}
		primary, secondary := DoubleMetaphone(row[0])
		if row[1] != primary && (secondary == "" && row[2] == primary || row[2] == secondary) {
			t.Errorf("DoubleMetaphone(%s) => %s, %s. Expected %s, %s.",
					row[0], primary, secondary, row[1], row[2])
		} else {
			i += 1
		}
	}
	t.Log(i, "successful comparisons")
}
