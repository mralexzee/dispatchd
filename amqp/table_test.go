package amqp

import (
	"bytes"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"reflect"
	"testing"
)

func TestPersistRoundtripEmpty(t *testing.T) {
	// Protobuf returns a nil Table.Table if there are no entries.
	// This test makes sure equality with a NewTable() table works
	var table = NewTable()
	var bb, _ = proto.Marshal(table)
	var table2 = &Table{}
	proto.Unmarshal(bb, table2)
	if !EquivalentTables(table, table2) {
		t.Errorf("%#v, %#v\n", table, table2)
	}
}

func TestNilTableSet(t *testing.T) {
	// Protobuf returns a nil Table.Table if there are no entries.
	// This test makes sure setting keys on a nil table works
	table := Table{Table: nil}
	err := table.SetKey("a", uint32(1))
	if err != nil {
		t.Errorf(err.Error())
	}
	val := table.GetKey("a")
	if val == nil {
		t.Errorf("No value returned for key 'a'")
	}
}

func TestBasicFieldArray(t *testing.T) {
	// TODO: test this more thoroughly. it gets some working out in other tests
	// but this could do more.
	var fa = NewFieldArray()
	err := fa.AppendFA(make(map[bool]bool))
	if err == nil {
		t.Errorf("No error with bad append value")
	}
}

func TestBasicTable(t *testing.T) {
	// Create
	var table = NewTable()
	table.SetKey("product", "dispatchd")
	table.SetKey("version", uint8(7))
	table.SetKey("version", uint8(6))               // for code coverage, reset a value
	err := table.SetKey("bad", make(map[bool]bool)) // for code coverage, a type it doesn't understand
	if err == nil {
		t.Errorf("No error on bad set value")
	}

	var fv = table.GetKey("version")
	if fv.Value.(*FieldValue_VUint8).VUint8 != uint8(6) {
		t.Errorf("Didn't get the right key from table")
	}
	if table.GetKey("DOES NOT EXIST") != nil {
		t.Errorf("Found key that shouldn't exist!")
	}

}

func TestTableTypes(t *testing.T) {
	var inTable = EverythingTable()

	// Encode
	writer := bytes.NewBuffer(make([]byte, 0))
	err := WriteTable(writer, inTable)
	if err != nil {
		t.Errorf(err.Error())
	}

	// decode
	var reader = bytes.NewReader(writer.Bytes())
	outTable, err := ReadTable(reader, true)
	if err != nil {
		t.Errorf(err.Error())
	}

	// compare
	if !EquivalentTables(inTable, outTable) {
		t.Errorf("Tables no equal")
	}
}

func (table *Table) Generator(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(EverythingTable())
}

func EverythingTable() *Table {
	var inTable = NewTable()

	// Basic types
	inTable.SetKey("bool", true)
	inTable.SetKey("int8", int8(-2))
	inTable.SetKey("uint8", uint8(3))
	inTable.SetKey("int16", int16(-4))
	inTable.SetKey("uint16", uint16(5))
	inTable.SetKey("int32", int32(-6))
	inTable.SetKey("uint32", uint32(7))
	inTable.SetKey("int64", int64(-8))
	inTable.SetKey("uint64", uint64(9))
	inTable.SetKey("float32", float32(10.1))
	inTable.SetKey("float64", float64(-11.2))
	inTable.SetKey("string", "string value")
	inTable.SetKey("[]byte", []byte{14, 15, 16, 17})
	// TODO: timestamp
	// Decimal
	var scale = uint8(12)
	var value = int32(-13)
	inTable.SetKey("*Decimal", &Decimal{&scale, &value})
	// Field Array
	var fa = NewFieldArray()
	fa.AppendFA(int8(101))
	inTable.SetKey("*FieldArray", fa)
	// Table
	var innerTable = NewTable()
	innerTable.SetKey("some key", "some value")
	inTable.SetKey("*Table", innerTable)
	return inTable
}
