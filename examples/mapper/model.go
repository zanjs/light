package mapper

import (
	"database/sql"

	"github.com/arstd/gobatis/examples/domain"
	"github.com/arstd/gobatis/examples/enums"
)

//go:generate go run ../../*.go

// ModelMapper 示例接口
type ModelMapper interface {

	// insert into models(buildin_bool, buildin_byte,
	// 	 buildin_float32, buildin_float64,
	//   buildin_int, buildin_int16, buildin_int32, buildin_int64, buildin_int8,
	//   buildin_rune, buildin_string, buildin_uint, buildin_uint16, buildinuint32,
	//   buildin_uint64, buildin_uint8, buildin_map, enum_status,
	//   ptr_model)
	// values (${m.BuildinBool}, ${m.BuildinByte}, ${m.BuildinFloat32},
	//   ${m.BuildinFloat64}, ${m.BuildinInt}, ${m.BuildinInt16}, ${m.BuildinInt32},
	//   ${m.BuildinInt64}, ${m.BuildinInt8}, ${m.BuildinRune}, ${m.BuildinString},
	//   ${m.BuildinUint}, ${m.BuildinUint16}, ${m.BuildinUint32}, ${m.BuildinUint64},
	//   ${m.BuildinUint8, ${m.BuildinMap}, ${m.EnumStatus}, ${m.PtrModel})
	// returning id
	Insert(tx *sql.Tx, m *domain.Model) error

	// update models
	// set buildin_bool=${m.BuildinBool}, buildin_byte=${m.BuildinByte},
	//   buildin_float32=${m.BuildinFloat32}, buildin_float64=${m.BuildinFloat64},
	//   buildin_int=${m.BuildinInt}, buildin_int16=${m.BuildinInt16},
	//   buildin_int32=${m.BuildinInt32}, buildin_int64=${m.BuildinInt64},
	//   buildin_int8=${m.BuildinInt8}, buildin_rune=${m.BuildinRune},
	//   buildin_string=${m.BuildinString}, buildin_uint=${m.BuildinUint},
	//   buildin_uint16=${m.BuildinUint16}, buildin_uint32=${m.BuildinUint32},
	//   buildin_uint64=${m.BuildinUint64}, buildin_uint8=${m.BuildinUint8},
	//   buildin_map=${m.BuildinMap}, enum_status=${m.EnumStatus},
	//   ptr_model=${m.PtrModel}
	// where id=${m.Id}
	Update(tx *sql.Tx, m *domain.Model) (int64, error)

	// delete from models
	// where id=${m.Id}
	Delete(tx *sql.Tx, id int) (int64, error)

	// select id, buildin_bool, buildin_byte, buildin_float32, buildin_float64,
	//   buildin_int, buildin_int16, buildin_int32, buildin_int64, buildin_int8,
	//   buildin_rune, buildin_string, buildin_uint, buildin_uint16, buildinuint32,
	//   buildin_uint64, buildin_uint8, buildin_map, enum_status,
	//   ptr_model
	// from models
	// where id=${m.Id}
	Get(tx *sql.Tx, id int) (*domain.Model, error)

	// select count(*)
	// from models
	// where buildin_bool=${d.Name}
	//   [?{m.BuildinInt != false} and buildin_int=${m.BuildinInt} ]
	//   [?{len(ss) != 0} and enum_status in (${ss}) ]
	// order by id
	// offset ${offset} limit ${limit}
	Count(tx *sql.Tx, m *domain.Model, ss []enums.Status) (int64, error)

	// select sum(buildin_float64)
	// from models
	// where buildin_bool=${d.Name}
	//   [?{m.BuildinInt != false} and buildin_int=${m.BuildinInt} ]
	//   [?{len(ss) != 0} and enum_status in (${ss}) ]
	// order by id
	// offset ${offset} limit ${limit}
	Sum(tx *sql.Tx, m *domain.Model, ss []enums.Status) (float64, error)

	// select id, buildin_bool, buildin_byte, buildin_float32, buildin_float64,
	//   buildin_int, buildin_int16, buildin_int32, buildin_int64, buildin_int8,
	//   buildin_rune, buildin_string, buildin_uint, buildin_uint16, buildinuint32,
	//   buildin_uint64, buildin_uint8, buildin_map, enum_status,
	//   ptr_model
	// from models
	// where buildin_bool=${d.Name}
	//   [?{m.BuildinInt != false} and buildin_int=${m.BuildinInt} ]
	//   [?{len(ss) != 0} and enum_status in (${ss}) ]
	// order by id
	// offset ${offset} limit ${limit}
	List(tx *sql.Tx, m *domain.Model, ss []enums.Status, offset, limit int) ([]*domain.Model, error)
}
