package php

import (
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type TestSuite struct{}

var _ = Suite(&TestSuite{})

func (t *TestSuite) TestReferences(c *C) {
	val, err := Unmarshal([]byte(`a:12:{i:0;i:1;i:1;i:2;i:2;O:8:"stdClass":2:{s:2:"id";s:9:"testClass";s:4:"some";s:5:"thing";}i:3;i:4;i:4;i:5;i:5;a:2:{i:0;s:3:"foo";i:1;s:3:"bar";}i:6;i:7;i:7;i:8;i:8;r:4;i:9;i:10;i:10;i:11;i:11;R:9;}`))

	c.Assert(err, IsNil)
	c.Assert(val.IsArray(), Equals, true)
	row8, err := val.GetKey(8)
	c.Assert(err, IsNil)
	c.Assert(row8.kind, Equals, KindObjReference)
	row8resolve, err := row8.GetKey("some")
	c.Assert(err, IsNil)
	c.Assert(row8resolve.kind, Equals, KindString)
	s, err := row8resolve.String()
	c.Assert(err, IsNil)
	c.Assert(s, Equals, "thing")

	// Variable references
	row11, err := val.GetKey(11)
	c.Assert(err, IsNil)
	c.Assert(row11.kind, Equals, KindVarReference)
	row11resolve, err := row11.GetKey(1)
	c.Assert(err, IsNil)
	c.Assert(row11resolve.kind, Equals, KindString)
	s, err = row11resolve.String()
	c.Assert(err, IsNil)
	c.Assert(s, Equals, "bar")
}

func (t *TestSuite) TestObject(c *C) {
	val, err := Unmarshal(
		[]byte("O:3:\"Foo\":4:{s:3:\"one\";s:3:\"aaa\";s:6:\"\000*\000two\";s:3:\"bbb\";s:10:\"\000Foothree\000\";s:3:\"ccc\";s:4:\"four\";s:3:\"ddd\";}"),
	)
	c.Assert(err, IsNil)
	j, _ := val.JSON()
	// need better things here. this will fail at some point. maybe just on a different machine. due to map randomization....
	c.Assert(string(j), Equals, `{"four":"ddd","one":"aaa","three":"ccc","two":"bbb"}`)
}

func (t *TestSuite) TestNullBoolOffsets(c *C) {
	a, err := Unmarshal([]byte("a:7:{i:0;b:1;i:1;b:0;i:2;b:0;i:3;b:1;i:4;N;i:5;N;i:6;s:4:\"addd\";}"))
	c.Assert(err, IsNil)
	c.Assert(a.IsArray(), Equals, true)
	array, _ := a.Rows()
	c.Assert(len(array), Equals, 7)
}

func (t *TestSuite) TestNull(c *C) {
	n, err := Unmarshal([]byte("N;"))
	c.Assert(err, IsNil)
	c.Assert(n.IsNull(), Equals, true)
}

func (t *TestSuite) TestBool(c *C) {
	bt, err := Unmarshal([]byte("b:1;"))
	c.Assert(err, IsNil)
	c.Assert(bt.IsBool(), Equals, true)
	b, _ := bt.Bool()
	c.Assert(b, Equals, true)
	f, err := Unmarshal([]byte("b:0;"))
	c.Assert(err, IsNil)
	c.Assert(f.IsBool(), Equals, true)
	b, _ = f.Bool()
	c.Assert(b, Equals, false)
}

func (t *TestSuite) TestNestedArray(c *C) {
	val, err := Unmarshal([]byte("a:3:{i:0;i:1;i:1;a:2:{i:0;i:2;i:1;i:3;}i:2;i:4;}"))
	c.Assert(err, IsNil)
	c.Assert(val.IsArray(), Equals, true)
	shouldBeArray, err := val.GetKey(1)
	c.Assert(err, IsNil)
	c.Assert(shouldBeArray.IsArray(), Equals, true)
	v, _ := shouldBeArray.GetKey(1)
	i, _ := v.Int()
	c.Assert(i, Equals, 3)
}

func (t *TestSuite) TestSimpleArray(c *C) {
	val, err := Unmarshal([]byte("a:3:{i:0;i:1;i:1;i:2;i:2;i:3;}"))
	c.Assert(err, IsNil)
	array, err := val.Rows()
	c.Assert(err, IsNil)
	i, err := array[1].Key.Int()
	c.Assert(err, IsNil)
	c.Assert(i, Equals, 1)
	i, err = array[2].Val.Int()
	c.Assert(err, IsNil)
	c.Assert(i, Equals, 3)
}

func (t *TestSuite) TestFloat(c *C) {
	tests := map[float64][]byte{
		10.99:     []byte("d:10.99;"),
		999999.99: []byte("d:999999.99;"),
	}
	for want, input := range tests {
		v, err := Unmarshal(input)
		c.Assert(err, IsNil)
		c.Assert(v.IsFloat(), Equals, true)
		s, err := v.Float()
		c.Assert(err, IsNil)
		c.Assert(s, Equals, want)
	}
}

func (t *TestSuite) TestInt(c *C) {
	tests := map[int][]byte{
		1:      []byte("i:1;"),
		100009: []byte("i:100009;"),
	}
	for want, input := range tests {
		v, err := Unmarshal(input)
		c.Assert(err, IsNil)
		c.Assert(v.IsInt(), Equals, true)
		s, err := v.Int()
		c.Assert(err, IsNil)
		c.Assert(s, Equals, want)
	}
}

func (t *TestSuite) TestString(c *C) {
	tests := map[string][]byte{
		"foobarbazboo":            []byte(`s:12:"foobarbazboo";`),
		"fo\to\nb\ra\"r'b;aüzboo": []byte("s:20:\"fo\to\nb\ra\"r'b;aüzboo\";"),
	}
	for want, input := range tests {
		v, err := Unmarshal(input)
		c.Assert(err, IsNil)
		c.Assert(v.IsString(), Equals, true)
		s, err := v.String()
		c.Assert(err, IsNil)
		c.Assert(s, Equals, want)
	}
}
