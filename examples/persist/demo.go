package persist

import "github.com/arstd/persist/examples/domain"

//go:generate go run ../../main.go

// DemoPersister xxx
type DemoPersister interface {

	// insert into demos(demo_name, demo_status, demo_struct)
	// values(${d.DemoName},${d.DemoStatus},${d.DemoStruct})
	// returning id
	Add(d *domain.Demo) (err error)

	// update demos
	// set demo_name=${d.DemoName}
	// where id=${d.Id}
	Update(d *domain.Demo) (err error)

	// select id, demo_name, demo_status, demo_struct
	// from demos where id=${id}
	Get(id string) (d *domain.Demo, err error)

	// select id, demo_name, demo_status, demo_struct
	// from demos
	// where id<=${d.Id} and demo_name=${d.DemoName}
	// and demo_status=${d.DemoStatus} and demo_struct=${d.DemoStruct}
	List(d *domain.Demo) (ds []*domain.Demo, err error)

	// delete from demos where id=${id}
	Delete(id string) (err error)
}
