package domain

import (
	"time"

	"github.com/arstd/gobatis/examples/enums"
)

// Model 示例结构体
type Model struct {
	Id          int `json:"id"`
	BuildinBool bool
	BuildinByte byte
	// BuildinComplex128 complex128
	// BuildinComplex64  complex64
	// BuildinError   error
	BuildinFloat32 float32
	BuildinFloat64 float64
	BuildinInt     int
	BuildinInt16   int16
	BuildinInt32   int32
	BuildinInt64   int64
	BuildinInt8    int8
	BuildinRune    rune
	BuildinString  string
	BuildinUint    uint
	BuildinUint16  uint16
	BuildinUint32  uint32
	BuildinUint64  uint64
	BuildinUint8   uint8
	// BuildinUintptr uintptr
	BuildinMap map[string]interface{}

	EnumStatus enums.Status
	PtrModel   *Model
	Time       time.Time
}

/*
create table models (
	id serial primary key,
	buildin_bool bool,
	buildin_byte smallint,
	buildin_float32 real,
	buildin_float64 double precision,
	buildin_int     int8,
	buildin_int16   smallint,
	buildin_int32   integer,
	buildin_int64   bigint,
	buildin_int8    smallint,
	buildin_rune    smallint,
	buildin_string  text,
	buildin_uint    bigint,
	buildin_uint16  integer,
	buildin_uint32  bigint,
	buildin_uint64  bigint,
	buildin_uint8   smallint,
	buildin_map     jsonb,
	enum_status smallint,
	ptr_model   jsonb,
	time timestamptz
)
*/
