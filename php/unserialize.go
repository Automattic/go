package php

import (
	"fmt"
	"strconv"
)

// Null represents a PHP null
type Null interface{}

var (
	// ErrWrongType indicates that the PHP value is not compatible with the way you requested it
	ErrWrongType = fmt.Errorf("Value is not of requested type")
	// ErrMalformedInput indicates that, while parsing, something unexpected was found
	ErrMalformedInput = fmt.Errorf("Input is malformed")
	// ErrUnsupportedType indicates that the value represents a type that is incompatible with the requested operation
	ErrUnsupportedType = fmt.Errorf("The supplied type is not supported for this operation")
)

// http://www.phpinternalsbook.com/classes_objects/serialization.html
const (
	idArray  = 'a'
	idString = 's'
	idInt    = 'i'
	idFloat  = 'd'
	idObject = 'O'
	idNull   = 'N'
	idBool   = 'b'
	idRef    = 'R'
	idOref   = 'r'
)

const (
	synSep = ':'
	synDq  = '"'
	synSq  = '\''
	synEnd = ';'
	synCbo = '{'
	synCbc = '}'
	synNil = '\000'
)

func vLength(body []byte) (length int, skip int, err error) {
	err = ErrMalformedInput
	if body[0] != synSep {
		return
	}
	for i := 1; i < len(body); i++ {
		skip = i + 2
		if body[i] == synSep || body[i] == synEnd {
			length, err = strconv.Atoi(string(body[1:i]))
			return
		}
	}
	return
}

func vRaw(body []byte) (data []byte, skip int, err error) {
	err = ErrMalformedInput
	if body[0] != synSep {
		return
	}
	for i := 1; i < len(body); i++ {
		skip = i + 2
		if body[i] == synEnd {
			data = body[1:i]
			err = nil
			return
		}
	}
	return
}

func unpackInt(body []byte, position int) (*Value, int, error) {
	oldPosition := position
	raw, skip, err := vRaw(body[position+1:])
	if err != nil {
		return nil, 0, err
	}
	position = position + skip
	val, err := strconv.Atoi(string(raw))
	if err != nil {
		return nil, position, err
	}
	return &Value{
		kind:    KindInt,
		content: val,
		bytes:   body[oldPosition:position],
	}, position, nil
}

func unpackFloat(body []byte, position int) (*Value, int, error) {
	oldPosition := position
	raw, skip, err := vRaw(body[position+1:])
	if err != nil {
		return nil, 0, err
	}
	position = position + skip
	val, err := strconv.ParseFloat(string(raw), 64)
	if err != nil {
		return nil, position, err
	}
	return &Value{
		kind:    KindFloat,
		content: val,
		bytes:   body[oldPosition:position],
	}, position, nil
}

func unpackString(body []byte, position int) (*Value, int, error) {
	oldPosition := position
	strLen, skip, err := vLength(body[position+1:])
	if err != nil {
		return nil, 0, err
	}
	position = position + skip
	if body[position] != synDq || body[position+strLen+1] != synDq || body[position+strLen+2] != synEnd {
		return nil, position, ErrMalformedInput
	}
	return &Value{
		kind:    KindString,
		content: body[position+1 : position+strLen+1],
		bytes:   body[oldPosition : position+strLen+3],
	}, position + strLen + 3, nil
}

func unpackArray(body []byte, position int) (*Value, int, error) {
	oldPosition := position
	aLen, skip, err := vLength(body[position+1:])
	if err != nil {
		return nil, 0, err
	}
	position = position + skip + 1
	if body[position-1] != synCbo {
		return nil, position, ErrMalformedInput
	}
	var content = []*Row{}
	for i := 0; i < aLen; i++ {
		k, skip, err := unmarshal(body[position:])
		position = position + skip
		if err != nil {
			return nil, position, ErrMalformedInput
		}
		v, skip, err := unmarshal(body[position:])
		position = position + skip
		if err != nil {
			return nil, position, ErrMalformedInput
		}
		if k == nil || v == nil {
			return nil, position, ErrMalformedInput
		}
		content = append(content, &Row{Key: k, Val: v})
	}
	if len(content) != aLen {
		return nil, position, ErrMalformedInput
	}
	return &Value{
		kind:    KindArray,
		content: content,
		bytes:   body[oldPosition : position+1],
	}, position + 1, nil
}

func unpackObject(body []byte, position int) (*Value, int, error) {
	var className []byte
	oldPosition := position
	strLen, skip, err := vLength(body[position+1:])
	if err != nil {
		return nil, 0, err
	}
	className = body[position+skip+1 : position+skip+1+strLen]
	position = position + skip + strLen + 1
	cLen, skip, err := vLength(body[position+1:])
	position = position + skip
	if err != nil {
		return nil, 0, err
	}
	if body[position] != synCbo {
		return nil, position, ErrMalformedInput
	}
	position++
	var content = []*Row{}
	for i := 0; i < cLen; i++ {
		k, skip, err := unmarshal(body[position:])
		position = position + skip
		if err != nil {
			return nil, position, ErrMalformedInput
		}

		s, _ := k.String()
		if s[0] == synNil {
			if s[1] == '*' && s[2] == synNil {
				k.prefix = []byte{0, 42, 0}
				k.content = []byte(s[3:])
			} else {
				k.prefix = append([]byte{0}, className...)
				k.suffix = []byte{0}
				k.content = []byte(s[1+len(className) : len(s)-1])
			}
		}
		v, skip, err := unmarshal(body[position:])
		position = position + skip
		if err != nil {
			return nil, position, ErrMalformedInput
		}
		if k == nil || v == nil {
			return nil, position, ErrMalformedInput
		}
		content = append(content, &Row{Key: k, Val: v})
	}
	if len(content) != cLen {
		return nil, position, ErrMalformedInput
	}
	return &Value{
		kind:      KindObject,
		content:   content,
		bytes:     body[oldPosition : position+1],
		className: className,
	}, position + 1, nil
}

func unpackNull(body []byte, position int) (*Value, int, error) {
	if body[position+1] != ';' {
		return nil, 0, ErrMalformedInput
	}
	return &Value{kind: KindNull, content: nil}, 2, nil
}

func unpackBool(body []byte, position int) (*Value, int, error) {
	if body[position+1] != ':' || body[position+3] != ';' {
		return nil, 0, ErrMalformedInput
	}
	switch body[position+2] {
	case '0':
		return &Value{kind: KindBool, content: false}, 4, nil
	case '1':
		return &Value{kind: KindBool, content: true}, 4, nil
	}
	return nil, 4, ErrMalformedInput
}

func unmarshal(body []byte) (*Value, int, error) {
	if len(body) < 1 {
		return nil, 0, nil
	}
	var position = 0
	switch body[position] {
	case idObject:
		return unpackObject(body, position)
	case idArray:
		return unpackArray(body, position)
	case idFloat:
		return unpackFloat(body, position)
	case idInt:
		return unpackInt(body, position)
	case idString:
		return unpackString(body, position)
	case idNull:
		return unpackNull(body, position)
	case idBool:
		return unpackBool(body, position)
	case idRef:
		referenceTo, skip, err := vLength(body[position+1:])
		position = position + skip
		if err != nil {
			return nil, position, err
		}
		return &Value{kind: KindVarReference, content: referenceTo}, position, nil
	case idOref:
		referenceTo, skip, err := vLength(body[position+1:])
		position = position + skip
		if err != nil {
			return nil, position, err
		}
		return &Value{kind: KindObjReference, content: referenceTo}, position, nil
	}
	return nil, 0, ErrMalformedInput
}

func resolveRefs(v *Value, list []*Value) []*Value {
	list = append(list, v)
	switch v.kind {
	case KindObject:
		for _, row := range v.content.([]*Row) {
			list = resolveRefs(row.Val, list)
		}
	case KindArray:
		for _, row := range v.content.([]*Row) {
			list = resolveRefs(row.Val, list)
		}
	}
	return list
}

func applyRefs(v *Value, list []*Value) {
	v.refs = list
	switch v.kind {
	case KindObject:
		for _, row := range v.content.([]*Row) {
			applyRefs(row.Val, list)
		}
	case KindArray:
		for _, row := range v.content.([]*Row) {
			applyRefs(row.Val, list)
		}
	}
}

// Unmarshal the serialized PHP data into a Value
func Unmarshal(body []byte) (*Value, error) {
	v, _, e := unmarshal(body)
	if e == nil {
		applyRefs(v, resolveRefs(v, nil))
	}
	return v, e
}
