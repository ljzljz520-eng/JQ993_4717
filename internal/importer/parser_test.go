package importer

import (
	"strings"
	"testing"
)

func TestParseCSVAndValidate(t *testing.T) {
	csv := "id,batch,name,scent,material,owner,tags\n1,993-27,Cedar,wood,wax,ops,safety|label\n2,993-28,Rose,floral,wax,ops,label\n"
	rows, invalid, err := New().ParseCSV(strings.NewReader(csv))
	if err != nil || len(rows) != 2 || len(invalid) != 0 {
		t.Fatal(err, invalid)
	}
	if len(ValidateBatch(rows)) != 0 {
		t.Fatal("valid rows rejected")
	}
}
