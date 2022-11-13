package utils

import "testing"

func TestGenChineseValue(t *testing.T) {
	l := 10
	t.Log(GenChineseValue(l))
}

func TestGenDatetimeValue(t *testing.T) {
	t.Log(GenDatetimeValue())
}

func TestGenFloatValue(t *testing.T) {
	t.Log(GenFloatValue())
}

func TestGenIntValue(t *testing.T) {
	l := 100
	t.Log(GenIntValue(l))
	l = 1000
	t.Log(GenIntValue(l))
}

func TestGenStringValue(t *testing.T) {
	l := 10
	t.Log(GenStringValue(l))
}

func TestUUID(t *testing.T) {
	t.Log(UUID())
}

func TestTakeParentheses(t *testing.T) {
	str := "select * from table where id = (select id from table2 where id = 1)"
	t.Log(TakeParentheses(str))
}
