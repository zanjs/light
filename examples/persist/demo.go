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
	// where [?(d.Id!="")id<=${d.Id} and] demo_name=${d.DemoName}
	// [?(d.DemoStatus!="")and demo_status=${d.DemoStatus}] and demo_struct=${d.DemoStruct}
	List(d *domain.Demo) (ds []*domain.Demo, err error)

	// delete from demos where id=${id}
	Delete(id string) (err error)

	// select id,code,source,user_id,mobile,total,reduced,amount,deposit,remain,
	// 	add_on_fee,status,pay_status,status_migrate,service_code,goods,
	// 	add_on,backend,discounts,extras,creator_id,reserve_time,
	// 	hospital_score,doctor_score,evaluation
	// from orders
	// where 1=1
	// [?(userId != "") and user_id=${userId}]
	// [?(start != 0 && end != 0) and create_at between ${start} and ${end}]
	// [?(start != 0 && end == 0) and create_at >= ${start}]
	// [?(start == 0 && end != 0) and create_at <= ${end}]
	// order by create_at desc
	// List2(userId string, start, end int64, page, size int64) (data []*pb.Order, err error)
}
