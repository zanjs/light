package persist

import (
	"fmt"
	"math/rand"
	"testing"

	m "github.com/arstd/persist/examples/domain"
)

var x DemoPersist
var id string

func TestDemoDataAdd(t *testing.T) {
	d := &m.Demo{
		DemoName: fmt.Sprintf("U_%d", rand.Intn(9999)),
	}

	err := x.Add(d)
	if err != nil {
		t.Fatalf("insert %#v error: %s", d, err)
	}

	id = d.Id
}

func TestDemoDataUpdate(t *testing.T) {
	d := &m.Demo{
		Id:       id,
		DemoName: fmt.Sprintf("Ux_%d", rand.Intn(9999)),
	}

	err := x.Update(d)
	if err != nil {
		t.Fatalf("update %#v error: %s", d, err)
	}
}

func TestDemoDataList(t *testing.T) {
	d := &m.Demo{
		Id:       id,
		DemoName: fmt.Sprintf("Ux_%d", rand.Intn(9999)),
	}
	mts, err := x.List(d)

	if err != nil {
		t.Fatalf("list(%+v) error: %s", d, err)
	}

	t.Logf("%v", mts)
}

func TestDemoDataDelete(t *testing.T) {
	err := x.Delete(id)
	if err != nil {
		t.Fatalf("delete by Id=%s error: %s", id, err)
	}
}
