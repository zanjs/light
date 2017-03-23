// DO NOT EDIT THIS FILE!
// It is generated by `light` tool from source `model.go`.

package mapper

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/arstd/light/example/domain"
	"github.com/arstd/light/example/enum"
	"github.com/arstd/log"
	"github.com/lib/pq"
)

type ModelMapperImpl struct{}

func (*ModelMapperImpl) Insert(m *domain.Model, xtx ...*sql.Tx) (err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` insert into models(name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32) values (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s) returning id`)
	zmMap, _ := json.Marshal(m.Map)
	zmPointer, _ := json.Marshal(m.Pointer)
	zmStructSlice, _ := json.Marshal(m.StructSlice)
	xargs = append(xargs, m.Name, m.Flag, m.Score, zmMap, m.Time, pq.Array(m.Array), pq.Array(m.Slice), m.Status, zmPointer, zmStructSlice, m.Uint32)

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)
	log.Debug(xquery)
	log.Debug(xargs...)

	xdest := []interface{}{&m.Id}
	if len(xtx) > 0 {
		err = xtx[0].QueryRow(xquery, xargs...).Scan(xdest...)
	} else {
		err = db.QueryRow(xquery, xargs...).Scan(xdest...)
	}
	if err != nil {
		log.Error(err)
		log.Error(xquery)
		log.Error(xargs...)
	}
	return
}

func (*ModelMapperImpl) BatchInsert(ms []*domain.Model, xtx ...*sql.Tx) (xa int64, err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` insert into models(uint32, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice) values`)
	xargs = append(xargs)
	for i, m := range ms {
		if i != 0 {
			xbuf.WriteString(",")
		}
		xbuf.WriteString(` (%s+888, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)`)
		zmMap, _ := json.Marshal(m.Map)
		zmPointer, _ := json.Marshal(m.Pointer)
		zmStructSlice, _ := json.Marshal(m.StructSlice)
		xargs = append(xargs, i, m.Name, m.Flag, m.Score, zmMap, m.Time, pq.Array(m.Array), pq.Array(m.Slice), m.Status, zmPointer, zmStructSlice)
	}

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)
	log.Debug(xquery)
	log.Debug(xargs...)

	var xres sql.Result
	if len(xtx) > 0 {
		xres, err = xtx[0].Exec(xquery, xargs...)
	} else {
		xres, err = db.Exec(xquery, xargs...)
	}
	if err != nil {
		log.Error(xquery)
		log.Error(xargs...)
		log.Error(err)
	}
	return xres.RowsAffected()
}

func (*ModelMapperImpl) Get(id int, xtx ...*sql.Tx) (xobj *domain.Model, err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32 from models where id=%s`)
	xargs = append(xargs, id)

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)
	log.Debug(xquery)
	log.Debug(xargs...)

	xobj = &domain.Model{}
	var xobjzMap []byte
	var xobjzTime pq.NullTime
	var xobjzPointer []byte
	var xobjzStructSlice []byte
	xdest := []interface{}{&xobj.Id, &xobj.Name, &xobj.Flag, &xobj.Score, &xobjzMap, &xobjzTime, pq.Array(&xobj.Array), pq.Array(&xobj.Slice), &xobj.Status, &xobjzPointer, &xobjzStructSlice, &xobj.Uint32}
	if len(xtx) > 0 {
		err = xtx[0].QueryRow(xquery, xargs...).Scan(xdest...)
	} else {
		err = db.QueryRow(xquery, xargs...).Scan(xdest...)
	}
	if err != nil {
		log.Error(err)
		log.Error(xquery)
		log.Error(xargs...)
	}
	xobj.Map = map[string]interface{}{}
	json.Unmarshal(xobjzMap, xobj.Map)
	xobj.Time = xobjzTime.Time
	xobj.Pointer = &domain.Model{}
	json.Unmarshal(xobjzPointer, &xobj.Pointer)
	xobj.StructSlice = []*domain.Model{}
	json.Unmarshal(xobjzStructSlice, &xobj.StructSlice)
	return
}

func (*ModelMapperImpl) Update(m *domain.Model, xtx ...*sql.Tx) (xa int64, err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` update models set name=%s, flag=%s, score=%s, map=%s, time=%s, slice=%s, status=%s, pointer=%s, struct_slice=%s, uint32=%s where id=%s`)
	zmMap, _ := json.Marshal(m.Map)
	zmPointer, _ := json.Marshal(m.Pointer)
	zmStructSlice, _ := json.Marshal(m.StructSlice)
	xargs = append(xargs, m.Name, m.Flag, m.Score, zmMap, m.Time, pq.Array(m.Slice), m.Status, zmPointer, zmStructSlice, m.Uint32, m.Id)

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)
	log.Debug(xquery)
	log.Debug(xargs...)

	var xres sql.Result
	if len(xtx) > 0 {
		xres, err = xtx[0].Exec(xquery, xargs...)
	} else {
		xres, err = db.Exec(xquery, xargs...)
	}
	if err != nil {
		log.Error(xquery)
		log.Error(xargs...)
		log.Error(err)
	}
	return xres.RowsAffected()
}

func (*ModelMapperImpl) Delete(id int, xtx ...*sql.Tx) (xa int64, err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` delete from models where id=%s`)
	xargs = append(xargs, id)

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)
	log.Debug(xquery)
	log.Debug(xargs...)

	var xres sql.Result
	if len(xtx) > 0 {
		xres, err = xtx[0].Exec(xquery, xargs...)
	} else {
		xres, err = db.Exec(xquery, xargs...)
	}
	if err != nil {
		log.Error(xquery)
		log.Error(xargs...)
		log.Error(err)
	}
	return xres.RowsAffected()
}

func (*ModelMapperImpl) Count(m *domain.Model, ss []enum.Status, xtx ...*sql.Tx) (xcnt int64, err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` select count(*) from models where name like %s`)
	xargs = append(xargs, m.Name)
	if m.Flag != false {
		xbuf.WriteString(` and flag=%s`)
		xargs = append(xargs, m.Flag)
	}
	if len(ss) != 0 {
		xargs = append(xargs)
		xbuf.WriteString(` and status in (`)
		xargs = append(xargs)
		for i, v := range ss {
			if i != 0 {
				xbuf.WriteString(",")
			}
			xbuf.WriteString(` %s`)
			xargs = append(xargs, v)
		}
		xbuf.WriteString(` )`)
		xargs = append(xargs)
	}
	if len(m.Slice) != 0 {
		xbuf.WriteString(` and slice && %s`)
		xargs = append(xargs, pq.Array(m.Slice))
	}

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)
	log.Debug(xquery)
	log.Debug(xargs...)

	if len(xtx) > 0 {
		err = xtx[0].QueryRow(xquery, xargs...).Scan(&xcnt)
	} else {
		err = db.QueryRow(xquery, xargs...).Scan(&xcnt)
	}
	if err != nil {
		log.Error(err)
		log.Error(xquery)
		log.Error(xargs...)
	}
	return
}

func (*ModelMapperImpl) List(m *domain.Model, ss []enum.Status, offset int, limit int, xtx ...*sql.Tx) (xdata []*domain.Model, err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32 from models where name like %s`)
	xargs = append(xargs, m.Name)
	if m.Flag != false {
		xargs = append(xargs)
		if len(ss) != 0 {
			xargs = append(xargs)
			xbuf.WriteString(` and status in (`)
			xargs = append(xargs)
			for i, v := range ss {
				if i != 0 {
					xbuf.WriteString(",")
				}
				xbuf.WriteString(` %s`)
				xargs = append(xargs, v)
			}
			xbuf.WriteString(` )`)
			xargs = append(xargs)
		}
		xbuf.WriteString(` and flag=%s`)
		xargs = append(xargs, m.Flag)
	}
	if len(m.Array) != 0 {
		xargs = append(xargs)
		xbuf.WriteString(` and xarray && array[`)
		xargs = append(xargs)
		for i, v := range m.Array {
			if i != 0 {
				xbuf.WriteString(",")
			}
			xbuf.WriteString(` %s`)
			zv, _ := json.Marshal(v)
			xargs = append(xargs, zv)
		}
		xbuf.WriteString(` ]`)
		xargs = append(xargs)
	}
	if len(m.Slice) != 0 {
		xbuf.WriteString(` and slice && %s`)
		xargs = append(xargs, pq.Array(m.Slice))
	}
	xbuf.WriteString(` order by id offset %s limit %s`)
	xargs = append(xargs, offset, limit)

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)
	log.Debug(xquery)
	log.Debug(xargs...)

	var xrows *sql.Rows
	if len(xtx) > 0 {
		xrows, err = xtx[0].Query(xquery, xargs...)
	} else {
		xrows, err = db.Query(xquery, xargs...)
	}
	if err != nil {
		log.Error(err)
		log.Error(xquery)
		log.Error(xargs...)
		return
	}
	defer xrows.Close()

	xdata = []*domain.Model{}
	for xrows.Next() {
		xe := &domain.Model{}
		xdata = append(xdata, xe)
		var xezMap []byte
		var xezTime pq.NullTime
		var xezPointer []byte
		var xezStructSlice []byte
		xdest := []interface{}{&xe.Id, &xe.Name, &xe.Flag, &xe.Score, &xezMap, &xezTime, pq.Array(&xe.Array), pq.Array(&xe.Slice), &xe.Status, &xezPointer, &xezStructSlice, &xe.Uint32}
		err = xrows.Scan(xdest...)
		if err != nil {
			log.Error(err)
			return
		}
		xe.Map = map[string]interface{}{}
		json.Unmarshal(xezMap, xe.Map)
		xe.Time = xezTime.Time
		xe.Pointer = &domain.Model{}
		json.Unmarshal(xezPointer, &xe.Pointer)
		xe.StructSlice = []*domain.Model{}
		json.Unmarshal(xezStructSlice, &xe.StructSlice)
	}
	if err = xrows.Err(); err != nil {
		log.Error(err)
	}
	return
}

func (*ModelMapperImpl) Page(m *domain.Model, ss []enum.Status, offset int, limit int, xtx ...*sql.Tx) (xcnt int64, xdata []*domain.Model, err error) {
	var (
		xbuf  bytes.Buffer
		xargs []interface{}
	)
	xbuf.WriteString(` select id, name, flag, score, map, time, slice, status, pointer, struct_slice from models where name like %s`)
	xargs = append(xargs, m.Name)
	if m.Flag != false {
		xargs = append(xargs)
		if len(ss) != 0 {
			xargs = append(xargs)
			xbuf.WriteString(` and status in (`)
			xargs = append(xargs)
			for i, v := range ss {
				if i != 0 {
					xbuf.WriteString(",")
				}
				xbuf.WriteString(` %s`)
				xargs = append(xargs, v)
			}
			xbuf.WriteString(` )`)
			xargs = append(xargs)
		}
		xbuf.WriteString(` and flag=%s`)
		xargs = append(xargs, m.Flag)
	}
	if len(m.Slice) != 0 {
		xbuf.WriteString(` and slice && %s`)
		xargs = append(xargs, pq.Array(m.Slice))
	}
	xbuf.WriteString(` order by id offset %s limit %s`)
	xargs = append(xargs, offset, limit)

	xholder := make([]interface{}, len(xargs))
	for i := range xargs {
		xholder[i] = fmt.Sprintf("$%d", i+1)
	}
	xquery := fmt.Sprintf(xbuf.String(), xholder...)

	xfindex := strings.LastIndex(xquery, " from ")
	xobindex := strings.LastIndex(xquery, "order by")
	xtquery := `select count(*)` + xquery[xfindex:xobindex]
	xdcnt := strings.Count(xquery[xobindex:], "$")
	xtargs := xargs[:len(xargs)-xdcnt]
	log.Debug(xtquery)
	log.Debug(xtargs...)

	if len(xtx) > 0 {
		err = xtx[0].QueryRow(xtquery, xtargs...).Scan(&xcnt)
	} else {
		err = db.QueryRow(xtquery, xtargs...).Scan(&xcnt)
	}
	if err != nil {
		log.Error(err)
		log.Error(xquery)
		log.Error(xargs...)
		return
	}
	if xcnt == 0 {
		return
	}

	log.Debug(xquery)
	log.Debug(xargs...)

	var xrows *sql.Rows

	if len(xtx) > 0 {
		xrows, err = xtx[0].Query(xquery, xargs...)
	} else {
		xrows, err = db.Query(xquery, xargs...)
	}
	if err != nil {
		log.Error(err)
		log.Error(xquery)
		log.Error(xargs...)
		return
	}
	defer xrows.Close()

	xdata = []*domain.Model{}
	for xrows.Next() {
		xe := &domain.Model{}
		xdata = append(xdata, xe)
		var xezMap []byte
		var xezTime pq.NullTime
		var xezPointer []byte
		var xezStructSlice []byte
		xdest := []interface{}{&xe.Id, &xe.Name, &xe.Flag, &xe.Score, &xezMap, &xezTime, pq.Array(&xe.Slice), &xe.Status, &xezPointer, &xezStructSlice}
		err = xrows.Scan(xdest...)
		if err != nil {
			log.Error(err)
			return
		}
		xe.Map = map[string]interface{}{}
		json.Unmarshal(xezMap, xe.Map)
		xe.Time = xezTime.Time
		xe.Pointer = &domain.Model{}
		json.Unmarshal(xezPointer, &xe.Pointer)
		xe.StructSlice = []*domain.Model{}
		json.Unmarshal(xezStructSlice, &xe.StructSlice)
	}
	if err = xrows.Err(); err != nil {
		log.Error(err)
	}
	return
}
