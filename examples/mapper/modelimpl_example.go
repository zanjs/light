package mapper

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/arstd/gobatis/examples/domain"
	"github.com/arstd/gobatis/examples/enums"
	"github.com/wothing/log"
)

type ModelImplExample struct{}

func (*ModelImplExample) Insert(tx *sql.Tx, m *domain.Model) error {
	var (
		buf  bytes.Buffer
		args []interface{}
		err  error
	)

	buf.WriteString(`insert into models(buildin_bool, buildin_byte,
		buildin_float32, buildin_float64,
	  buildin_int, buildin_int16, buildin_int32, buildin_int64, buildin_int8,
	  buildin_rune, buildin_string, buildin_uint, buildin_uint16, buildin_uint32,
	  buildin_uint64, buildin_uint8, buildin_map, enum_status, ptr_model)
	  values (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)
	  returning id`)

	m_BuildinMap, err := json.Marshal(m.BuildinMap)
	if err != nil {
		log.Error(err)
		return err
	}
	m_PtrModel, err := json.Marshal(m.PtrModel)
	if err != nil {
		log.Error(err)
		return err
	}
	args = append(args, m.BuildinBool, m.BuildinByte, m.BuildinFloat32,
		m.BuildinFloat64, m.BuildinInt, m.BuildinInt16, m.BuildinInt32,
		m.BuildinInt64, m.BuildinInt8, m.BuildinRune, m.BuildinString,
		m.BuildinUint, m.BuildinUint16, m.BuildinUint32, m.BuildinUint64,
		m.BuildinUint8, m_BuildinMap, m.EnumStatus, m_PtrModel)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)

	err = tx.QueryRow(query, args...).Scan(&m.Id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// update models
// set buildin_bool=${m.BuildinBool}, buildin_byte=${m.BuildinByte},
//   buildin_float32=${m.BuildinFloat32}, buildin_float64=${m.BuildinFloat64},
//   buildin_int=${m.BuildinInt}, buildin_int16=${m.BuildinInt16},
//   buildin_int32=${m.BuildinInt32}, buildin_int64=${m.BuildinInt64},
//   buildin_int8=${m.BuildinInt8}, buildin_rune=${m.BuildinRune},
//   buildin_string=${m.BuildinString}, buildin_uint=${m.BuildinUint},
//   buildin_uint16=${m.BuildinUint16}, buildinuint32=${m.BuildinUint32},
//   buildin_uint64=${m.BuildinUint64}, buildin_uint8=${m.BuildinUint8,
//   buildin_map=${m.BuildinMap}, enum_status=${m.EnumStatus},
//   ptr_model=${m.PtrModel})
// where id=${m.Id}
func (*ModelImplExample) Update(tx *sql.Tx, m *domain.Model) (int64, error) {
	var (
		buf  bytes.Buffer
		args []interface{}
		err  error
	)

	buf.WriteString(`update models
	set buildin_bool=%s, buildin_byte=%s,
	  buildin_float32=%s, buildin_float64=%s,
	  buildin_int=%s, buildin_int16=%s,
	  buildin_int32=%s, buildin_int64=%s,
	  buildin_int8=%s, buildin_rune=%s,
	  buildin_string=%s, buildin_uint=%s,
	  buildin_uint16=%s, buildin_uint32=%s,
	  buildin_uint64=%s, buildin_uint8=%s,
	  buildin_map=%s, enum_status=%s,
	  ptr_model=%s
	where id=%s`)

	m_BuildinMap, err := json.Marshal(m.BuildinMap)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	m_PtrModel, err := json.Marshal(m.PtrModel)
	if err != nil {
		log.Error(err)
		return 0, err
	}
	args = append(args, m.BuildinBool, m.BuildinByte, m.BuildinFloat32,
		m.BuildinFloat64, m.BuildinInt, m.BuildinInt16, m.BuildinInt32,
		m.BuildinInt64, m.BuildinInt8, m.BuildinRune, m.BuildinString,
		m.BuildinUint, m.BuildinUint16, m.BuildinUint32, m.BuildinUint64,
		m.BuildinUint8, m_BuildinMap, m.EnumStatus, m_PtrModel, m.Id)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)

	res, err := tx.Exec(query, args...)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return res.RowsAffected()
}

// delete from models
// where id=${m.Id}
func (*ModelImplExample) Delete(tx *sql.Tx, id int) (int64, error) {
	var (
		buf  bytes.Buffer
		args []interface{}
		err  error
	)

	buf.WriteString(`delete from models
	where id=%s`)
	args = append(args, id)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)

	res, err := tx.Exec(query, args...)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return res.RowsAffected()
}

// select count(*)
// from models
// where buildin_bool=${d.Name}
//   [?{m.BuildinInt != false} and buildin_int=${m.BuildinInt} ]
//   [?{len(ss) != 0} and enum_status in (${ss}) ]
// order by id
// offset ${offset} limit ${limit}
func (*ModelImplExample) Count(tx *sql.Tx, m *domain.Model, ss []enums.Status) (int64, error) {
	var (
		buf  bytes.Buffer
		args []interface{}
	)

	buf.WriteString(`select count(*)
	from models
	where buildin_bool=%s`)
	args = append(args, m.BuildinBool)

	if m.BuildinInt != 0 {
		buf.WriteString(` and buildin_int=%s`)
		args = append(args, m.BuildinInt)
	}
	if len(ss) != 0 {
		var stmt = ` and enum_status in (${ss})`
		stmt = strings.Replace(stmt, "${ss}", strings.Repeat(",%s", len(ss))[1:], -1)
		buf.WriteString(stmt)
		for _, s := range ss {
			args = append(args, int32(s))
		}
	}

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)

	var count int64

	err := tx.QueryRow(query, args...).Scan(&count)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return count, nil
}

// select sum(buildin_float64)
// from models
// where buildin_bool=${d.Name}
//   [?{m.BuildinInt != false} and buildin_int=${m.BuildinInt} ]
//   [?{len(ss) != 0} and enum_status in (${ss}) ]
// order by id
// offset ${offset} limit ${limit}
func (*ModelImplExample) Sum(tx *sql.Tx, m *domain.Model, ss []enums.Status) (float64, error) {
	var (
		buf  bytes.Buffer
		args []interface{}
	)

	buf.WriteString(`select sum(buildin_float64)
	from models
	where buildin_bool=%s`)
	args = append(args, m.BuildinBool)

	if m.BuildinInt != 0 {
		buf.WriteString(` and buildin_int=%s`)
		args = append(args, m.BuildinInt)
	}
	if len(ss) != 0 {
		var stmt = ` and enum_status in (${ss})`
		stmt = strings.Replace(stmt, "${ss}", strings.Repeat(",%s", len(ss))[1:], -1)
		buf.WriteString(stmt)
		for _, s := range ss {
			args = append(args, int32(s))
		}
	}

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)

	var sum float64
	err := tx.QueryRow(query, args...).Scan(&sum)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	return sum, nil
}

func (*ModelImplExample) List(tx *sql.Tx, m *domain.Model, ss []enums.Status,
	offset, limit int) ([]*domain.Model, error) {
	var (
		buf  bytes.Buffer
		args []interface{}
	)

	buf.WriteString(`select id, buildin_bool, buildin_byte,
		buildin_float32, buildin_float64,
	  buildin_int, buildin_int16, buildin_int32, buildin_int64, buildin_int8,
	  buildin_rune, buildin_string, buildin_uint, buildin_uint16, buildin_uint32,
	  buildin_uint64, buildin_uint8, buildin_map, enum_status, ptr_model
	from models
	where buildin_bool=%s`)
	args = append(args, m.BuildinBool)

	if m.BuildinInt != 0 {
		buf.WriteString(` and buildin_int=%s`)
		args = append(args, m.BuildinInt)
	}
	if len(ss) != 0 {
		var stmt = ` and enum_status in (${ss})`
		stmt = strings.Replace(stmt, "${ss}", strings.Repeat(",%s", len(ss))[1:], -1)
		buf.WriteString(stmt)
		for _, s := range ss {
			args = append(args, int32(s))
		}
	}

	buf.WriteString(` order by id offset %s limit %s`)
	args = append(args, offset, limit)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)

	rows, err := tx.Query(query, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var data []*domain.Model
	for rows.Next() {
		var x domain.Model
		data = append(data, &x)
		var x_BuildinMap, x_PtrModel []byte
		err = rows.Scan(&x.Id, &x.BuildinBool, &x.BuildinByte, &x.BuildinFloat32,
			&x.BuildinFloat64, &x.BuildinInt, &x.BuildinInt16, &x.BuildinInt32,
			&x.BuildinInt64, &x.BuildinInt8, &x.BuildinRune, &x.BuildinString,
			&x.BuildinUint, &x.BuildinUint16, &x.BuildinUint32, &x.BuildinUint64,
			&x.BuildinUint8, &x_BuildinMap, &x.EnumStatus, &x_PtrModel)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = json.Unmarshal(x_BuildinMap, &x.BuildinMap)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		err = json.Unmarshal(x_PtrModel, &x.PtrModel)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	return data, nil
}

// select id, buildin_bool, buildin_byte, buildin_float32, buildin_float64,
//   buildin_int, buildin_int16, buildin_int32, buildin_int64, buildin_int8,
//   buildin_rune, buildin_string, buildin_uint, buildin_uint16, buildinuint32,
//   buildin_uint64, buildin_uint8, buildin_map, enum_status,
//   ptr_model
// from models
// where id=${m.Id}
func (*ModelImplExample) Get(tx *sql.Tx, id int) (*domain.Model, error) {
	var (
		buf  bytes.Buffer
		args []interface{}
		err  error
	)

	buf.WriteString(`select id, buildin_bool, buildin_byte, buildin_float32, buildin_float64,
	  buildin_int, buildin_int16, buildin_int32, buildin_int64, buildin_int8,
	  buildin_rune, buildin_string, buildin_uint, buildin_uint16, buildin_uint32,
	  buildin_uint64, buildin_uint8, buildin_map, enum_status, ptr_model
	from models
	where id=%s`)
	args = append(args, id)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)

	x := &domain.Model{}
	var x_BuildinMap, x_PtrModel []byte
	err = db.QueryRow(query, args...).Scan(
		&x.Id, &x.BuildinBool, &x.BuildinByte, &x.BuildinFloat32,
		&x.BuildinFloat64, &x.BuildinInt, &x.BuildinInt16, &x.BuildinInt32,
		&x.BuildinInt64, &x.BuildinInt8, &x.BuildinRune, &x.BuildinString,
		&x.BuildinUint, &x.BuildinUint16, &x.BuildinUint32, &x.BuildinUint64,
		&x.BuildinUint8, &x_BuildinMap, &x.EnumStatus, &x_PtrModel)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = json.Unmarshal(x_BuildinMap, &x.BuildinMap)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = json.Unmarshal(x_PtrModel, &x.PtrModel)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return x, nil
}
