package php

import (
	"bytes"
	"log"
	"testing"
)

func TestReferences(t *testing.T) {
	val, err := Unmarshal([]byte(`a:12:{i:0;i:1;i:1;i:2;i:2;O:8:"stdClass":2:{s:2:"id";s:9:"testClass";s:4:"some";s:5:"thing";}i:3;i:4;i:4;i:5;i:5;a:2:{i:0;s:3:"foo";i:1;s:3:"bar";}i:6;i:7;i:7;i:8;i:8;r:4;i:9;i:10;i:10;i:11;i:11;R:9;}`))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !val.IsArray() {
		t.Error("not an array?")
		return
	}

	// Object references
	row8, err := val.GetKey(8)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if row8.kind != KindObjReference {
		t.Error("should be object reference")
		return
	}
	row8resolve, err := row8.GetKey("some")
	if err != nil {
		t.Error(err.Error())
		return
	}
	if row8resolve.kind != KindString {
		t.Error("should resolve to a string!")
		return
	}
	if s, _ := row8resolve.String(); s != "thing" {
		t.Error("should have resolved to 'thing'")
	}

	// Variable references
	row11, err := val.GetKey(11)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if row11.kind != KindVarReference {
		t.Error("should be variable reference")
		return
	}
	row11resolve, err := row11.GetKey(1)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if row11resolve.kind != KindString {
		t.Error("should resolve to a string!")
		return
	}
	if s, _ := row11resolve.String(); s != "bar" {
		t.Error("should have resolved to 'bar'")
	}

}

func TestObject(t *testing.T) {
	val, err := Unmarshal(
		[]byte("O:3:\"Foo\":4:{s:3:\"one\";s:3:\"aaa\";s:6:\"\000*\000two\";s:3:\"bbb\";s:10:\"\000Foothree\000\";s:3:\"ccc\";s:4:\"four\";s:3:\"ddd\";}"),
	)
	if err != nil {
		t.Error(err.Error())
		return
	}
	j, _ := val.JSON()
	// need better things here. this will fail at some point. maybe just on a different machine. due to map randomization....
	if !bytes.Equal(j, []byte(`{"four":"ddd","one":"aaa","three":"ccc","two":"bbb"}`)) {
		t.Error("did not marshal to expected json...")
	}
}

func TestNullBoolOffsets(t *testing.T) {
	a, err := Unmarshal([]byte("a:7:{i:0;b:1;i:1;b:0;i:2;b:0;i:3;b:1;i:4;N;i:5;N;i:6;s:4:\"addd\";}"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !a.IsArray() {
		t.Error("something went wrong")
	}
	if array, _ := a.Rows(); len(array) != 7 {
		log.Printf("%#v", array)
		t.Error("something else went wrong")
	}
}

func TestNull(t *testing.T) {
	n, err := Unmarshal([]byte("N;"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !n.IsNull() {
		t.Error("is not null")
		return
	}
}

func TestBool(t *testing.T) {
	bt, err := Unmarshal([]byte("b:1;"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !bt.IsBool() {
		t.Error("is not boolean")
		return
	}
	if b, _ := bt.Bool(); !b {
		t.Error("shoulf be true")
	}

	f, err := Unmarshal([]byte("b:0;"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !f.IsBool() {
		t.Error("is not boolean")
		return
	}
	if b, _ := f.Bool(); b {
		t.Error("shoulf be false")
	}
}

func TestNestedArray(t *testing.T) {
	val, err := Unmarshal([]byte("a:3:{i:0;i:1;i:1;a:2:{i:0;i:2;i:1;i:3;}i:2;i:4;}"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !val.IsArray() {
		t.Error("not an array")
		return
	}
	shouldBeArray, err := val.GetKey(1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if !shouldBeArray.IsArray() {
		t.Error("sub-array not an array")
		return
	}

	v, _ := shouldBeArray.GetKey(1)
	if i, _ := v.Int(); i != 3 {
		t.Error("wrong value")
		return
	}
}

func TestSimpleArray(t *testing.T) {
	val, err := Unmarshal([]byte("a:3:{i:0;i:1;i:1;i:2;i:2;i:3;}"))
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !val.IsArray() {
		t.Error("not an array")
		return
	}
	array, err := val.Rows()
	if err != nil {
		t.Error(err.Error())
		return
	}
	array[0].Key.Int()
	if i, err := array[1].Key.Int(); err != nil {
		t.Error(err.Error())
		return
	} else if i != 1 {
		t.Error("wrong index value")
		return
	}
	if i, err := array[2].Val.Int(); err != nil {
		t.Error(err.Error())
		return
	} else if i != 3 {
		t.Error("wrong value")
		return
	}
}

func TestFloat(t *testing.T) {
	tests := map[float64][]byte{
		10.99:     []byte("d:10.99;"),
		999999.99: []byte("d:999999.99;"),
	}
	for want, input := range tests {
		v, err := Unmarshal(input)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !v.IsFloat() {
			t.Error("is not a float")
			return
		}
		s, err := v.Float()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if s != want {
			t.Error("Content decoded to the wrong string")
		}
	}
}

func TestInt(t *testing.T) {
	tests := map[int][]byte{
		1:      []byte("i:1;"),
		100009: []byte("i:100009;"),
	}
	for want, input := range tests {
		v, err := Unmarshal(input)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !v.IsInt() {
			t.Error("not an integer")
			return
		}
		s, err := v.Int()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if s != want {
			t.Error("Content decoded to the wrong string")
		}
	}
}

func TestString(t *testing.T) {
	tests := map[string][]byte{
		"foobarbazboo":            []byte(`s:12:"foobarbazboo";`),
		"fo\to\nb\ra\"r'b;aüzboo": []byte("s:20:\"fo\to\nb\ra\"r'b;aüzboo\";"),
	}
	for want, input := range tests {
		t.Logf("Testing for %#v", want)
		v, err := Unmarshal(input)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if !v.IsString() {
			t.Error("is not a string")
			return
		}
		s, err := v.String()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if s != want {
			t.Error("Content decoded to the wrong string")
		}
	}
}
