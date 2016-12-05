// DO NOT EDIT THIS FILE !
// It is generated by yan tool, source from model.go.
package mapper

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/arstd/light/example/domain"
	"github.com/arstd/light/example/enum"
	"github.com/arstd/log"
	"github.com/lib/pq"
)

type ModelMapperImpl struct{}

func (*ModelMapperImpl) Insert(tx *sql.Tx, m *domain.Model) (err error) {
	var (
		stmt string
		buf  bytes.Buffer
		args []interface{}
	)

	stmt = `insert into model(name, flag, score, map, time, slice, status, pointer, struct_slice, uint32) values (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s) returning id `
	args = append(args, m.Name)
	args = append(args, m.Flag)
	args = append(args, m.Score)
	var xm_Map []byte
	xm_Map, err = json.Marshal(m.Map)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_Map)
	args = append(args, m.Time)
	var xm_Slice []byte
	xm_Slice, err = json.Marshal(m.Slice)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_Slice)
	args = append(args, int32(m.Status))
	var xm_Pointer []byte
	xm_Pointer, err = json.Marshal(m.Pointer)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_Pointer)
	var xm_StructSlice []byte
	xm_StructSlice, err = json.Marshal(m.StructSlice)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_StructSlice)
	args = append(args, time.Unix(int64(m.Uint32), 0))
	buf.WriteString(stmt)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)
	var dest []interface{}

	dest = append(dest, &m.Id)

	err = tx.QueryRow(query, args...).Scan(dest...)
	if err != nil {
		log.Error(err)
		log.Error(query)
		log.Error(args...)
		return
	}

	return nil
}

func (*ModelMapperImpl) Update(tx *sql.Tx, m *domain.Model) (xi int64, err error) {
	var (
		stmt string
		buf  bytes.Buffer
		args []interface{}
	)

	stmt = `update model set name=%s, flag=%s, score=%s, map=%s, time=%s, slice=%s, status=%s, pointer=%s, struct_slice=%s, uint32=%s where id=%s `
	args = append(args, m.Name)
	args = append(args, m.Flag)
	args = append(args, m.Score)
	var xm_Map []byte
	xm_Map, err = json.Marshal(m.Map)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_Map)
	args = append(args, m.Time)
	var xm_Slice []byte
	xm_Slice, err = json.Marshal(m.Slice)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_Slice)
	args = append(args, int32(m.Status))
	var xm_Pointer []byte
	xm_Pointer, err = json.Marshal(m.Pointer)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_Pointer)
	var xm_StructSlice []byte
	xm_StructSlice, err = json.Marshal(m.StructSlice)
	if err != nil {
		log.Error(err)
		return
	}
	args = append(args, xm_StructSlice)
	args = append(args, time.Unix(int64(m.Uint32), 0))
	args = append(args, m.Id)
	buf.WriteString(stmt)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)
	res, err := db.Exec(query, args...)
	if err != nil {
		log.Error(err)
		log.Error(query)
		log.Error(args...)
		return 0, err
	}
	return res.RowsAffected()
}

func (*ModelMapperImpl) Delete(tx *sql.Tx, id int) (xi int64, err error) {
	var (
		stmt string
		buf  bytes.Buffer
		args []interface{}
	)

	stmt = `delete from model where id=%s `
	args = append(args, id)
	buf.WriteString(stmt)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)
	res, err := db.Exec(query, args...)
	if err != nil {
		log.Error(err)
		log.Error(query)
		log.Error(args...)
		return 0, err
	}
	return res.RowsAffected()
}

func (*ModelMapperImpl) Get(tx *sql.Tx, id int) (xm *domain.Model, err error) {
	var (
		stmt string
		buf  bytes.Buffer
		args []interface{}
	)

	stmt = `select id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32 from model where id=%s `
	args = append(args, id)
	buf.WriteString(stmt)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)
	var dest []interface{}
	xm = &domain.Model{}
	dest = append(dest, &xm.Id)
	dest = append(dest, &xm.Name)
	dest = append(dest, &xm.Flag)
	dest = append(dest, &xm.Score)
	var xxm_Map []byte
	dest = append(dest, &xxm_Map)
	dest = append(dest, &xm.Time)
	var xxm_Slice []byte
	dest = append(dest, &xxm_Slice)
	dest = append(dest, &xm.Status)
	var xxm_Pointer []byte
	dest = append(dest, &xxm_Pointer)
	var xxm_StructSlice []byte
	dest = append(dest, &xxm_StructSlice)
	var xxm_Uint32 pq.NullTime
	dest = append(dest, &xxm_Uint32)
	err = db.QueryRow(query, args...).Scan(dest...)
	if err != nil {
		log.Error(err)
		log.Error(query)
		log.Error(args...)
		return
	}
	xm.Map = map[string]interface{}{}
	err = json.Unmarshal(xxm_Map, &xm.Map)
	if err != nil {
		log.Error(err)
		return
	}
	xm.Slice = []string{}
	err = json.Unmarshal(xxm_Slice, &xm.Slice)
	if err != nil {
		log.Error(err)
		return
	}
	xm.Pointer = &domain.Model{}
	err = json.Unmarshal(xxm_Pointer, xm.Pointer)
	if err != nil {
		log.Error(err)
		return
	}
	xm.StructSlice = []*domain.Model{}
	err = json.Unmarshal(xxm_StructSlice, &xm.StructSlice)
	if err != nil {
		log.Error(err)
		return
	}
	if xxm_Uint32.Valid {
		xm.Uint32 = uint32(xxm_Uint32.Time.Unix())
	}
	return
}

func (*ModelMapperImpl) Count(tx *sql.Tx, m *domain.Model, ss []enum.Status) (xi int64, err error) {
	var (
		stmt string
		buf  bytes.Buffer
		args []interface{}
	)

	stmt = `select count(*) from model where name like %s `
	args = append(args, m.Name)
	buf.WriteString(stmt)

	if m.Flag != false {
		stmt = `and flag=%s  `
		args = append(args, m.Flag)
		buf.WriteString(stmt)
	}
	if len(ss) != 0 {
		stmt = `and status in (${ss})  `
		stmt = strings.Replace(stmt, "${"+"ss"+"}",
			strings.Repeat(",%s", len(ss))[1:], -1)
		for _, s := range ss {
			args = append(args, s)
		}
		buf.WriteString(stmt)
	}

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)
	var dest []interface{}
	dest = append(dest, &xi)
	err = db.QueryRow(query, args...).Scan(dest...)
	if err != nil {
		log.Error(err)
		log.Error(query)
		log.Error(args...)
		return
	}
	return
}

func (*ModelMapperImpl) List(tx *sql.Tx, m *domain.Model, ss []enum.Status, offset int, limit int) (xms []*domain.Model, err error) {
	var (
		stmt string
		buf  bytes.Buffer
		args []interface{}
	)

	stmt = `select id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32 from model where name like %s `
	args = append(args, m.Name)
	buf.WriteString(stmt)

	if m.Flag != false {
		stmt = `and flag=%s  `
		args = append(args, m.Flag)
		buf.WriteString(stmt)
	}
	if len(ss) != 0 {
		stmt = `and status in (${ss})  `
		stmt = strings.Replace(stmt, "${"+"ss"+"}",
			strings.Repeat(",%s", len(ss))[1:], -1)
		for _, s := range ss {
			args = append(args, s)
		}
		buf.WriteString(stmt)
	}
	if len(m.Slice) != 0 {
		stmt = `and slice ?| array[${m.Slice}]  `
		stmt = strings.Replace(stmt, "${"+"m.Slice"+"}",
			strings.Repeat(",%s", len(m.Slice))[1:], -1)
		for _, s := range m.Slice {
			args = append(args, s)
		}
		buf.WriteString(stmt)
	}

	stmt = `order by id offset %s limit %s `
	args = append(args, offset)
	args = append(args, limit)
	buf.WriteString(stmt)

	var ph []interface{}
	for i := range args {
		ph = append(ph, "$"+strconv.Itoa(i+1))
	}

	query := fmt.Sprintf(buf.String(), ph...)

	log.Debug(query)
	log.Debug(args...)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var data []*domain.Model
	for rows.Next() {
		x := &domain.Model{}
		data = append(data, x)

		var dest []interface{}
		dest = append(dest, &x.Id)
		dest = append(dest, &x.Name)
		dest = append(dest, &x.Flag)
		dest = append(dest, &x.Score)
		var xx_Map []byte
		dest = append(dest, &xx_Map)
		dest = append(dest, &x.Time)
		var xx_Slice []byte
		dest = append(dest, &xx_Slice)
		dest = append(dest, &x.Status)
		var xx_Pointer []byte
		dest = append(dest, &xx_Pointer)
		var xx_StructSlice []byte
		dest = append(dest, &xx_StructSlice)
		var xx_Uint32 pq.NullTime
		dest = append(dest, &xx_Uint32)
		err = rows.Scan(dest...)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		x.Map = map[string]interface{}{}
		err = json.Unmarshal(xx_Map, &x.Map)
		if err != nil {
			log.Error(err)
			return
		}
		x.Slice = []string{}
		err = json.Unmarshal(xx_Slice, &x.Slice)
		if err != nil {
			log.Error(err)
			return
		}
		x.Pointer = &domain.Model{}
		err = json.Unmarshal(xx_Pointer, x.Pointer)
		if err != nil {
			log.Error(err)
			return
		}
		x.StructSlice = []*domain.Model{}
		err = json.Unmarshal(xx_StructSlice, &x.StructSlice)
		if err != nil {
			log.Error(err)
			return
		}
		if xx_Uint32.Valid {
			x.Uint32 = uint32(xx_Uint32.Time.Unix())
		}
	}
	if err = rows.Err(); err != nil {
		log.Error(err)
		return nil, err
	}

	return data, nil
}