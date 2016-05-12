package persist

import (
	"testing"

	"github.com/gotips/log"

	m "github.com/arstd/persist/examples/domain"
)

var x DemoPersist
var id string

//
// func TestDemoDataAdd(t *testing.T) {
// 	d := &m.Demo{
// 		Name:       "Name",
// 		ThirdField: true,
// 		Status:     m.StatusNormal,
// 		Content:    bytes.NewBufferString("test content"),
// 	}
//
// 	err := x.Add(d)
// 	if err != nil {
// 		t.Fatalf("insert %#v error: %s", d, err)
// 	}
//
// 	id = d.Id
// }
//
// func TestDemoDataUpdate(t *testing.T) {
// 	d := &m.Demo{
// 		Id:         id,
// 		DemoName:   "demo_name",
// 		DemoStatus: "demo_status",
// 		DemoStruct: bytes.NewBufferString("{}"),
// 	}
//
// 	err := x.Update(d)
// 	if err != nil {
// 		t.Fatalf("update %#v error: %s", d, err)
// 	}
// }
//
// func TestDemoGet(t *testing.T) {
// 	d, err := x.Get(id)
//
// 	if err != nil {
// 		t.Fatalf("get(%+v) error: %s", d, err)
// 	}
//
// 	log.SetLevel(log.TraceLevel)
// 	log.JSON(d)
// }

func TestDemoDataList(t *testing.T) {
	d := &m.Demo{
		Name:       "Name",
		ThirdField: true,
		Status:     m.StatusNormal,
		Content:    &m.Demo{},
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	mts, err := x.List(tx, d, []m.Status{m.StatusNormal}, 1, 9999)
	CommitTx(tx)

	if err != nil {
		t.Fatalf("list(%+v) error: %s", d, err)
	}

	log.JSON(mts)
}

// func TestDemoDataDelete(t *testing.T) {
// 	err := x.Delete(id)
// 	if err != nil {
// 		t.Fatalf("delete by Id=%s error: %s", id, err)
// 	}
// }
