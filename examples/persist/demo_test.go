package persist

import (
	"testing"

	m "github.com/arstd/persist/examples/domain"
)

var x DemoPersist
var id string

func TestDemoDataAdd(t *testing.T) {
	d := &m.Demo{
		DemoName:   "demo_name",
		DemoStatus: "demo_status",
		DemoStruct: &m.Demo{
			DemoName:   "demo_name",
			DemoStatus: "demo_status",
		},
	}

	err := x.Add(d)
	if err != nil {
		t.Fatalf("insert %#v error: %s", d, err)
	}

	id = d.Id
}

func TestDemoDataUpdate(t *testing.T) {
	d := &m.Demo{
		Id:         id,
		DemoName:   "demo_name",
		DemoStatus: "demo_status",
		DemoStruct: &m.Demo{
			DemoName:   "demo_name",
			DemoStatus: "demo_status",
		},
	}

	err := x.Update(d)
	if err != nil {
		t.Fatalf("update %#v error: %s", d, err)
	}
}

func TestDemoGet(t *testing.T) {
	d, err := x.Get(id)

	if err != nil {
		t.Fatalf("get(%+v) error: %s", d, err)
	}

	t.Logf("%v", d)
}

func TestDemoDataList(t *testing.T) {
	d := &m.Demo{
		Id:         id,
		DemoName:   "demo_name",
		DemoStatus: "demo_status",
		DemoStruct: &m.Demo{
			DemoName:   "demo_name",
			DemoStatus: "demo_status",
		},
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
