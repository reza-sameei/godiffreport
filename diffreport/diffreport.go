/**
 * Copyright (c) 2019 TRIALBLAZE PTY. LTD. All rights reserved.
 *
 * Created by Reza Same'ei (reza.sameei@trialblaze.com).
 * User: Reza Same'ei
 * Date: 2019-05-26
 * Time: 10:44
 *
 * Description: diffreport.go
 */
package diffreport

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

type (
	Diff struct {
		Flag string
		Type string
		Desc string
		Key  map[string]string
		From string
		To   string
	}

	DiffReportContext struct {
		Key  map[string]string
		List []*Diff
	}
)

func DiffReportContextNew() *DiffReportContext {
	return &DiffReportContext{
		map[string]string{},
		[]*Diff{},
	}
}

func (rc *DiffReportContext) Add(
	flag string, tpe string,
	desc string, key string,
	from string, to string,
) *DiffReportContext {

	keys := make(map[string]string)
	for k, v := range rc.Key {
		keys[k] = v
	}
	keys[tpe] = key

	rc.List = append(
		rc.List,
		&Diff{flag, tpe, desc, keys, from, to},
	)

	return rc
}

// ============================

func (rc *DiffReportContext) ValNew(tpe, key, current string) *DiffReportContext {
	return rc.Add("New", tpe, "", key, "", current)
}

func (rc *DiffReportContext) ValDel(tpe, key, past string) *DiffReportContext {
	return rc.Add("Delete", tpe, "", key, past, "")
}

func (rc *DiffReportContext) ValUpdate(tpe, key, past, current string) *DiffReportContext {
	return rc.Add("Update", tpe, "", key, past, current)
}

func (rc *DiffReportContext) StructNew(tpe, key string) *DiffReportContext {
	return rc.Add("New", tpe, "", key, "", "{...}")
}

func (rc *DiffReportContext) StructDel(tpe, key string) *DiffReportContext {
	return rc.Add("Delete", tpe, "", key, "{...}", "")
}

func (rc *DiffReportContext) StructUpdate(tpe, key string) *DiffReportContext {
	return rc.Add("Update", tpe, "", key, "{...}", "{...}")
}

func (rc *DiffReportContext) ValCheckBool(tpc, key string, past, current bool) *DiffReportContext {
	if past != current {
		return rc.ValUpdate(tpc, key, strconv.FormatBool(past), strconv.FormatBool(current))
	}
	return rc
}

func (rc *DiffReportContext) ValCheckInt(tpe, key string, past, current int64) *DiffReportContext {
	if past != current {
		return rc.ValUpdate(tpe, key, strconv.FormatInt(past, 10), strconv.FormatInt(current, 10))
	}
	return rc
}

func (rc *DiffReportContext) ValCheckStr(tpe, key, past, current string) *DiffReportContext {
	if past != current {
		return rc.ValUpdate(tpe, key, past, current)
	}
	return rc
}

func (rc *DiffReportContext) ValCheckFloat64(tpe, key string, past, current float64) *DiffReportContext {
	if past != current {
		return rc.ValUpdate(tpe, key,
			strconv.FormatFloat(past, 'f', 2, 64),
			strconv.FormatFloat(current, 'f', 2, 64),
			)
	}
	return rc
}
func (rc *DiffReportContext) ValCheckFloat32(tpe, key string, past, current float32) *DiffReportContext {
	if past != current {
		return rc.ValUpdate(tpe, key,
			strconv.FormatFloat(float64(past), 'f', 2, 32),
			strconv.FormatFloat(float64(current), 'f', 2, 32),
		)
	}
	return rc
}

// ============================

func (rc *DiffReportContext) AddNew(tpe, key string, current interface{}) *DiffReportContext {
	switch v := current.(type) {
	case string:
		return rc.ValNew(tpe, key, v)
	case int64:
		return rc.ValNew(tpe, key, strconv.FormatInt(v, 10))
	case bool:
		return rc.ValNew(tpe, key, strconv.FormatBool(v))
	case struct{}:
		return rc.StructNew(tpe, key)
	default:
		log.Panicf("Unspoorted Type: %T, %s", current, current)
	}
	return rc // Never
}

func (rc *DiffReportContext) AddDel(tpe, key string, past interface{}) *DiffReportContext {
	switch v := past.(type) {
	case string:
		return rc.ValDel(tpe, key, v)
	case int64:
		return rc.ValDel(tpe, key, strconv.FormatInt(v, 10))
	case bool:
		return rc.ValDel(tpe, key, strconv.FormatBool(v))
	case struct{}:
		return rc.StructDel(tpe, key)
	default:
		log.Panicf("Unspoorted Type: %T, %s", past, past)
	}
	return rc // Never
}

func (rc *DiffReportContext) AddChangeType(tpe, key string, past, current reflect.Type) *DiffReportContext {
	return rc.Add("ChangeType", tpe, "", key, past.Name(), current.Name())
}

func (rc *DiffReportContext) KeysOf(maps ...map[string]interface{}) []string {

	unique := make(map[string]string)

	for _, m := range maps {
		for k, _ := range m {
			unique[k] = ""
		}
	}

	ls := []string{}
	for k, _ := range unique {
		ls = append(ls, k)
	}

	return ls
}

func (rc *DiffReportContext) MapSIOf(obj interface{}) (map[string]interface{}, error) {

	js, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var mp map[string]interface{}
	err = json.Unmarshal(js, &mp)
	if err != nil {
		return nil, err
	}

	return mp, nil
}

func (rc *DiffReportContext) Check(tpe, key string, past, current interface{}) (*DiffReportContext, error) {

	switch {

	case past == nil && current == nil:
		return rc, nil

	case past == nil && current != nil:
		rc.StructNew(tpe, key)

	case past != nil && current == nil:
		rc.StructDel(tpe, key)

	case past != nil && current != nil:

		pmap, err := rc.MapSIOf(past)
		if err != nil {
			return nil, err
		}
		cmap, err := rc.MapSIOf(current)
		if err != nil {
			return nil, err
		}

		rc.StepIn(tpe, key)
		fmt.Printf("#1 Tpe: %s, Key: %s, %s\n", tpe, key, rc.DebugJSON())
		rc.CheckStructProps(tpe, key, pmap, cmap)
		fmt.Println(rc.DebugJSON())
		rc.StepOut(tpe)
	}

	return rc, nil
}

func (rc *DiffReportContext) CheckStructProps(tpe, key string, past, current map[string]interface{}) *DiffReportContext {

	keys := rc.KeysOf(past, current)

	for _, k := range keys {

		pv, pf := past[k]
		cv, cf := current[k]

		switch {
		case !pf && !cf:
			continue // Nothing, Impossible
		case !pf && cf:
			rc.AddNew(key, k, cv)
		case pf && !cf:
			rc.AddDel(key, k, pv)
		case pf && cf:
			fmt.Printf("#2 Tpe: %s, Key: %s, Prop: %s, %s\n", tpe, key, k, rc.DebugJSON())
			rc.CheckVal(tpe, key, k, pv, cv)
			fmt.Printf("#3 Tpe: %s, Key: %s, Prop: %s, %s\n", tpe, key, k, rc.DebugJSON())
		}
	}

	return rc
}

func (rc *DiffReportContext) CheckVal(tpe, key, prop string, past, current interface{}) *DiffReportContext {

	pt := reflect.TypeOf(past)
	ct := reflect.TypeOf(current)

	/*if !IsSupportedByType(pt) || !IsSupportedByType(ct) {
		log.Panicf("Unsupported types: %T(%v), %T(%v)", past, past, current, current)
	}*/

	switch {
	case pt.Kind() != ct.Kind():
		return rc.AddChangeType(key, prop, pt, ct)
	case pt.Kind() == reflect.String:
		return rc.ValCheckStr(key, prop, past.(string), current.(string))
	case pt.Kind() == reflect.Int64:
		return rc.ValCheckInt(key, prop, past.(int64), current.(int64))
	case pt.Kind() == reflect.Int32:
		return rc.ValCheckInt(key, prop, int64(past.(int32)), int64(current.(int32)))
	case pt.Kind() == reflect.Int16:
		return rc.ValCheckInt(key, prop, int64(past.(int16)), int64(current.(int16)))
	case pt.Kind() == reflect.Int8:
		return rc.ValCheckInt(key, prop, int64(past.(int8)), int64(current.(int8)))
	case pt.Kind() == reflect.Int:
		return rc.ValCheckInt(key, prop, int64(past.(int)), int64(current.(int)))
	case pt.Kind() == reflect.Float32:
		return rc.ValCheckFloat32(key, prop, past.(float32), current.(float32))
	case pt.Kind() == reflect.Float64:
		return rc.ValCheckFloat64(key, prop, past.(float64), current.(float64))
	case pt.Kind() == reflect.Bool:
		return rc.ValCheckBool(key, prop, past.(bool), current.(bool))
	case pt.Kind() == reflect.Map:
		rc.StepIn(tpe, key)
		rc.CheckStructProps(key, prop, past.(map[string]interface{}), current.(map[string]interface{}))
		rc.StepOut(tpe)
		return rc
	default:
		log.Panicf("Unsupported types: %T(%v), %T(%v)", past, past, current, current)
	}
	return rc // Never
}

func (rc *DiffReportContext) StepIn(tpe string, key string) *DiffReportContext {
	rc.Key[tpe] = key
	return rc
}

func (rc *DiffReportContext) StepOut(tpe string) *DiffReportContext {
	delete(rc.Key, tpe)
	return rc
}

func (rc *DiffReportContext) String() string {
	return fmt.Sprintf(
		"%T{ Keys: %v, List: %v }",
		rc, rc.Key, rc.List,
	)
}

func (r *Diff) String() string {
	return fmt.Sprintf(
		"%T{ Flag: %s, Type: %s, Desc: %s, Key: %v }",
		r, r.Flag, r.Type, r.Desc, r.Key,
	)
}

func (rc *DiffReportContext) DebugJSON() string {
	bs, _ := json.MarshalIndent(rc, "", "\t")
	return string(bs)
}

//

var supportedTypes = map[reflect.Kind]reflect.Kind{
	reflect.Map:    reflect.Map,
	reflect.String: reflect.String,
	reflect.Bool:   reflect.Bool,
	reflect.Int:    reflect.Int,
	reflect.Int8:   reflect.Int8,
	reflect.Int16:  reflect.Int16,
	reflect.Int32:  reflect.Int32,
	reflect.Int64:  reflect.Int64,
	reflect.Float32: reflect.Float32,
	reflect.Float64: reflect.Float64,
}

func IsSupportedByType(t reflect.Type) bool {
	_, ok := supportedTypes[t.Kind()]
	return ok
}
func IsSupported(i interface{}) bool {
	return IsSupportedByType(reflect.TypeOf(i))
}
