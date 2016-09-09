package mapper

import (
	"database/sql"

	"github.com/arstd/yan/example/domain"
	"github.com/arstd/yan/example/enum"
	"github.com/wothing/17mei/pb"
)

//go:generate yan -force

// ModelMapper 示例接口
type ModelMapper interface {

	// select id, tags from goods where id in (${ids})
	ListGoodsTags(ids []string) ([]*pb.Goods, error)

	// update goods set tags=${req.Tags} where id in (${req.Ids})
	GoodsSetTags(req *pb.GoodsSetTagsReq) (int64, error)

	// insert into model(name, flag, score, map, time, slice, status, pointer, struct_slice)
	// values (${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Slice},
	//   ${m.Status}, ${m.Pointer}, ${m.StructSlice})
	// returning id
	Insert(tx *sql.Tx, m *domain.Model) error

	// update model
	// set name=${m.Name}, flag=${m.Flag}, score=${m.Score},
	//   map=${m.Map}, time=${m.Time}, slice=${m.Slice},
	//   status=${m.Status}, pointer=${m.Pointer}, struct_slice=${m.StructSlice}
	// where id=${m.Id}
	Update(tx *sql.Tx, m *domain.Model) (int64, error)

	// delete from model
	// where id=${id}
	Delete(tx *sql.Tx, id int) (int64, error)

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice
	// from model
	// where id=${id}
	Get(tx *sql.Tx, id int) (*domain.Model, error)

	// select count(*)
	// from model
	// where name like ${m.Name}
	//   [?{ m.Flag != false } and flag=${m.Flag} ]
	//   [?{ len(ss) != 0 } and status in (${ss}) ]
	Count(tx *sql.Tx, m *domain.Model, ss []enum.Status) (int64, error)

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice
	// from model
	// where name like ${m.Name}
	//   [?{ m.Flag != false } and flag=${m.Flag} ]
	//   [?{ len(ss) != 0 } and status in (${ss}) ]
	//   [?{ len(m.Slice) != 0 } and slice ?| array[${m.Slice}] ]
	// order by id
	// offset ${offset} limit ${limit}
	List(tx *sql.Tx, m *domain.Model, ss []enum.Status, offset, limit int) ([]*domain.Model, error)
}
