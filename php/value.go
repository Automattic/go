package php

import (
	"encoding/json"
	"fmt"
)

const (
	KindString       = 1
	KindInt          = 2
	KindFloat        = 3
	KindArray        = 4
	KindObject       = 5
	KindNull         = 6
	KindBool         = 7
	KindVarReference = 8
	KindObjReference = 9
)

// Row represents a key/value pair for PHP objects and PHP arrays
type Row struct {
	Key *Value
	Val *Value
}

// Value represents a PHP value
type Value struct {
	refs      []*Value
	kind      int
	content   interface{}
	bytes     []byte
	className []byte
	prefix    []byte
	suffix    []byte
}

func (v *Value) isRef() bool {
	if v.kind == KindObjReference || v.kind == KindVarReference {
		return true
	}
	return false
}

func (v *Value) findRef(id int) *Value {
	wantId := id - 1
	if wantId > len(v.refs) {
		return nil
	}
	return v.refs[wantId]
}

// Kind returns the kind of value that this represents
func (v *Value) Kind() int {
	if v.isRef() {
		return v.findRef(v.content.(int)).kind
	}
	return v.kind
}

// KindString returns a string representation of the kind of PHP value being represented
// "unknown" is returned if this function is not kept up to date with new types of values
func (v *Value) KindString() string {
	if v.isRef() {
		return v.findRef(v.content.(int)).KindString()
	}
	switch v.kind {
	case KindString:
		return "string"
	case KindInt:
		return "int"
	case KindBool:
		return "boolean"
	case KindNull:
		return "null"
	case KindFloat:
		return "float"
	case KindArray:
		return "array"
	case KindObject:
		return "object"
	case KindObjReference:
		return "object-reference"
	case KindVarReference:
		return "reference"
	}
	return "unknown"
}

// IsPrivate will tell you whether the member part of a key value pair on an object is private
// This data is lost when converted to JSON
func (v *Value) IsPrivate() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsPrivate()
	}
	if v.suffix != nil {
		return true
	}
	return false
}

// IsProtected will tell you whether the member part of a key value pair on an object is protected
// This data is lost when converted to JSON
func (v *Value) IsProtected() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsProtected()
	}
	if v.prefix != nil && v.suffix == nil {
		return true
	}
	return false
}

// ClassName will give you the class name of an object
// This data is lost when converted to JSON
func (v *Value) ClassName() ([]byte, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).ClassName()
	}
	if v.kind != KindObject {
		return nil, ErrWrongType
	}
	return v.className, nil
}

// resolve recursively returns the real value of the value, so to speak
//
// For PHP arrays and strings it returns a map[string]interface{}
// This conversion looses: class names, and member visibility.
//
// For string, int, float, bool types it returns those types
//
// For PHP null values it would return a nil interface{}
//
// References and duplicate objects are also lost because this would lead to infinite output in many cases
func (v *Value) resolve() interface{} {
	switch v.kind {
	case KindArray, KindObject:
		var rval = map[string]interface{}{}
		for _, row := range v.content.([]*Row) {
			if row.Val.isRef() {
				continue
			}
			switch row.Key.kind {
			case KindString:
				key, _ := row.Key.String()
				rval[key] = row.Val.resolve()
			case KindInt:
				key, _ := row.Key.Int()
				rval[fmt.Sprintf("%d", key)] = row.Val.resolve()
			default:
				panic("Unknown key data!")
			}
		}
		return rval
	case KindString:
		s, _ := v.String()
		return s
	}
	return v.content
}

// JSON will return JSON for the value. This operation looses class names, and member visibility
// as JSON does not support these concepts.
func (v *Value) JSON() ([]byte, error) {
	return json.Marshal(v.resolve())
}

// IsNull tells you whether the PHP type was a NULL
func (v *Value) IsNull() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsNull()
	}
	if v.kind != KindNull {
		return false
	}
	return true
}

// Null returns the null interface{} that this value represents
func (v *Value) Null() (Null, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).Null()
	}
	if v.kind != KindNull {
		return false, ErrWrongType
	}
	return v.content, nil
}

// IsBool tells you whether the PHP type was a boolean
func (v *Value) IsBool() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsBool()
	}
	if v.kind != KindBool {
		return false
	}
	return true
}

// Bool returns the boolean that this value represents
func (v *Value) Bool() (bool, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).Bool()
	}
	if v.kind != KindBool {
		return false, ErrWrongType
	}
	content, ok := v.content.(bool)
	if !ok {
		return false, ErrWrongType
	}
	return content, nil
}

// IsArray tells you whether the PHP type was an array
func (v *Value) IsArray() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsArray()
	}
	if v.kind != KindArray {
		return false
	}
	return true
}

// IsFloat tells you whether the PHP type was a float
func (v *Value) IsFloat() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsFloat()
	}
	if v.kind != KindFloat {
		return false
	}
	return true
}

// Float returns the float64 represented by the value
func (v *Value) Float() (float64, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).Float()
	}
	if v.kind != KindFloat {
		return 0, ErrWrongType
	}
	content, ok := v.content.(float64)
	if !ok {
		return 0, ErrWrongType
	}
	return content, nil
}

// IsString tells you whether the PHP type was a string
func (v *Value) IsString() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsString()
	}
	if v.kind != KindString {
		return false
	}
	return true
}

// String returns the string that the value represents. This requires allocating
// memory, so Bytes is preferred when possible
func (v *Value) String() (string, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).String()
	}
	if v.kind != KindString {
		return "", ErrWrongType
	}
	content, ok := v.content.([]byte)
	if !ok {
		return "", ErrWrongType
	}
	return string(content), nil
}

// Bytes returns the byte slice which makes up the underlying value of a PHP string
func (v *Value) Bytes() ([]byte, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).Bytes()
	}
	if v.kind != KindString {
		return nil, ErrWrongType
	}
	content, ok := v.content.([]byte)
	if !ok {
		return nil, ErrWrongType
	}
	return content, nil
}

// IsInt tells you whether the PHP type was an int
func (v *Value) IsInt() bool {
	if v.isRef() {
		return v.findRef(v.content.(int)).IsInt()
	}
	if v.kind != KindInt {
		return false
	}
	return true
}

// Int returns the integer that the value represents
func (v *Value) Int() (int, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).Int()
	}
	if v.kind != KindInt {
		return 0, ErrWrongType
	}
	content, ok := v.content.(int)
	if !ok {
		return 0, ErrWrongType
	}
	return content, nil
}

func getStringKey(key string, content []*Row) *Value {
	for _, row := range content {
		k, err := row.Key.String()
		if err != nil {
			continue
		}
		if k == key {
			return row.Val
		}
	}
	return nil
}

func getIntKey(key int, content []*Row) *Value {
	for _, row := range content {
		k, err := row.Key.Int()
		if err != nil {
			continue
		}
		if k == key {
			return row.Val
		}
	}
	return nil
}

// GetKey will return the value for the associated key if the
// value being represented is an array or an object. key clan
// be a string or an int (which are the only types that PHP
// can serialize array and object keys into)
func (v *Value) GetKey(key interface{}) (*Value, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).GetKey(key)
	}
	if v.kind != KindArray && v.kind != KindObject {
		return nil, ErrWrongType
	}
	content, ok := v.content.([]*Row)
	if !ok {
		return nil, ErrWrongType
	}
	switch key := key.(type) {
	case int:
		return getIntKey(key, content), nil
	case string:
		return getStringKey(key, content), nil
	}
	return nil, ErrUnsupportedType
}

// ForEach allows you to skip some boiler plate and iterate over an object to an array
// with iterFunc.  If iterFunc returns an error then iteration is halted and the
// error is returned
func (v *Value) ForEach(iterFunc func(*Value, *Value) error) error {
	if v.isRef() {
		return v.findRef(v.content.(int)).ForEach(iterFunc)
	}
	if v.kind != KindArray {
		return ErrWrongType
	}
	content, ok := v.content.([][]*Value)
	if !ok {
		return ErrWrongType
	}
	for _, row := range content {
		if err := iterFunc(row[0], row[1]); err != nil {
			return err
		}
	}
	return nil
}

// Rows returns slices of [0] key [1] value sets which make up the array or object
func (v *Value) Rows() ([]*Row, error) {
	if v.isRef() {
		return v.findRef(v.content.(int)).Rows()
	}
	if v.kind != KindArray && v.kind != KindObject {
		return nil, ErrWrongType
	}
	content, ok := v.content.([]*Row)
	if !ok {
		return nil, ErrWrongType
	}
	return content, nil
}
