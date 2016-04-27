package persis

import (
	"fmt"
	"math/rand"
	"testing"
)

var testID int

func TestMyTestTableAdd(t *testing.T) {
	mt := &MyTest{
		TestName: fmt.Sprintf("U_%d", rand.Intn(9999)),
	}

	err := MyTestTable.Add(mt)
	if err != nil {
		t.Fatalf("insert %#v error: %s", mt, err)
	}

	testID = mt.TestID
}

func TestMyTestTableUpdate(t *testing.T) {
	mt := &MyTest{
		TestID:   testID,
		TestName: fmt.Sprintf("Ux_%d", rand.Intn(9999)),
	}

	err := MyTestTable.Update(mt)
	if err != nil {
		t.Fatalf("update %#v error: %s", mt, err)
	}
}

func TestMyTestTableList(t *testing.T) {
	mts, err := MyTestTable.List(testID)

	if err != nil {
		t.Fatalf("list by TestID=%d error: %s", testID, err)
	}

	t.Logf("%v", mts)
}

func TestMyTestTableDelete(t *testing.T) {
	err := MyTestTable.Delete(testID)
	if err != nil {
		t.Fatalf("delete by TestID=%d error: %s", testID, err)
	}
}
