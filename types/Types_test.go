package types

import (
	"testing"
)

func TestSimpleDeser(t *testing.T) {
	typ := CharType{}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	if deserialize.typId() != typ.typId() {
		t.Fail()
		return
	}
}

func TestEmptyStructDeser(t *testing.T) {
	typ := StructType{
		Fields: []Type{},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserStr, ok := deserialize.(StructType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if len(deserStr.Fields) != 0 {
		t.Errorf("Deser too many fields (%d)", deserStr.Fields)
		return
	}
}

func TestSimplStructDeser(t *testing.T) {
	typ := StructType{
		Fields: []Type{CharType{}},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserStr, ok := deserialize.(StructType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if len(deserStr.Fields) != 1 {
		t.Errorf("Deser not right no of fields (got %d, expected %d)", len(deserStr.Fields), 1)
		return
	}
	if (deserStr.Fields[0].typId() != CharType{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Fields[0].typId(), CharType{}.typId())
		return

	}
}
func TestSimplStructDeser2(t *testing.T) {
	typ := StructType{
		Fields: []Type{CharType{},
			Int64Type{}},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserStr, ok := deserialize.(StructType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if len(deserStr.Fields) != 2 {
		t.Errorf("Deser not right no of fields (got %d, expected %d)", len(deserStr.Fields), 2)
		return
	}
	if (deserStr.Fields[0].typId() != CharType{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Fields[0].typId(), CharType{}.typId())
		return

	}
	if (deserStr.Fields[1].typId() != Int64Type{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Fields[0].typId(), Int64Type{}.typId())
		return

	}
}

func TestNesedStructDeser(t *testing.T) {
	typ := StructType{
		Fields: []Type{
			CharType{},
			StructType{
				Fields: []Type{Int64Type{}},
			},
			Int64Type{},
		},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserStr, ok := deserialize.(StructType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if len(deserStr.Fields) != 3 {
		t.Errorf("Deser not right no of fields (got %d, expected %d)", len(deserStr.Fields), 3)
		return
	}
	if (deserStr.Fields[0].typId() != CharType{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Fields[0].typId(), CharType{}.typId())
		return
	}
	if (deserStr.Fields[2].typId() != Int64Type{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Fields[0].typId(), Int64Type{}.typId())
		return
	}
	nestedField, ok := deserStr.Fields[1].(StructType)
	if !ok {
		t.Error("Wrong nested type")
		return
	}
	if len(nestedField.Fields) != 1 {
		t.Errorf("Deser not right no of fields (got %d, expected %d)", len(nestedField.Fields), 1)
		return
	}
	if (nestedField.Fields[0].typId() != Int64Type{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", nestedField.Fields[0].typId(), Int64Type{}.typId())
		return
	}
}

func TestSimplUnionDeser(t *testing.T) {
	typ := UnionType{
		Members: []Type{CharType{}},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserStr, ok := deserialize.(UnionType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if len(deserStr.Members) != 1 {
		t.Errorf("Deser not right no of Members (got %d, expected %d)", len(deserStr.Members), 1)
		return
	}
	if (deserStr.Members[0].typId() != CharType{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Members[0].typId(), CharType{}.typId())
		return

	}
}
func TestSimplUnionDeser2(t *testing.T) {
	typ := UnionType{
		Members: []Type{CharType{},
			Int64Type{}},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserStr, ok := deserialize.(UnionType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if len(deserStr.Members) != 2 {
		t.Errorf("Deser not right no of Members (got %d, expected %d)", len(deserStr.Members), 2)
		return
	}
	if (deserStr.Members[0].typId() != CharType{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Members[0].typId(), CharType{}.typId())
		return

	}
	if (deserStr.Members[1].typId() != Int64Type{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Members[0].typId(), Int64Type{}.typId())
		return

	}
}

func TestNesedUnionDeser(t *testing.T) {
	typ := UnionType{
		Members: []Type{
			CharType{},
			UnionType{
				Members: []Type{Int64Type{}},
			},
			Int64Type{},
		},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserStr, ok := deserialize.(UnionType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if len(deserStr.Members) != 3 {
		t.Errorf("Deser not right no of Members (got %d, expected %d)", len(deserStr.Members), 3)
		return
	}
	if (deserStr.Members[0].typId() != CharType{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Members[0].typId(), CharType{}.typId())
		return
	}
	if (deserStr.Members[2].typId() != Int64Type{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", deserStr.Members[0].typId(), Int64Type{}.typId())
		return
	}
	nestedField, ok := deserStr.Members[1].(UnionType)
	if !ok {
		t.Error("Wrong nested type")
		return
	}
	if len(nestedField.Members) != 1 {
		t.Errorf("Deser not right no of Members (got %d, expected %d)", len(nestedField.Members), 1)
		return
	}
	if (nestedField.Members[0].typId() != Int64Type{}.typId()) {
		t.Errorf("Wrong typ id (got %d, expected %d)", nestedField.Members[0].typId(), Int64Type{}.typId())
		return
	}
}

func TestArrDeserSimple(t *testing.T) {
	typ := ArrayType{
		Length: 42,
		Typ:    CharType{},
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserArr, ok := deserialize.(ArrayType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if deserArr.Length != typ.Length {
		t.Errorf("Mistmatch! Expected length %d, got %d", deserArr.Length, typ.Length)
		return
	}
	if deserArr.Typ.typId() != typ.Typ.typId() {
		t.Errorf("Mistmatch! Expected length %d, got %d", deserArr.Typ.typId(), typ.Typ.typId())
		return
	}
}

func TestArrDeserNested(t *testing.T) {
	nestedOrig := ArrayType{
		Typ:    CharType{},
		Length: 17,
	}
	typ := ArrayType{
		Length: 42,
		Typ:    nestedOrig,
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserArr, ok := deserialize.(ArrayType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if deserArr.Length != typ.Length {
		t.Errorf("Mistmatch! Expected length %d, got %d", deserArr.Length, typ.Length)
		return
	}
	if deserArr.Typ.typId() != typ.Typ.typId() {
		t.Errorf("Mistmatch! Expected length %d, got %d", deserArr.Typ.typId(), typ.Typ.typId())
		return
	}
	nestedArr, ok := deserArr.Typ.(ArrayType)
	if !ok {
		t.Error("Wrong nested type")
		return
	}
	if nestedArr.Length != nestedOrig.Length {
		t.Errorf("Mistmatch! Expected length %d, got %d", nestedArr.Length, nestedOrig.Length)
		return
	}
	if nestedArr.Typ.typId() != nestedOrig.Typ.typId() {
		t.Errorf("Mistmatch! Expected type %d, got %d", nestedArr.Typ.typId(), nestedOrig.Typ.typId())
		return
	}

}

func TestArrDeserNestedStruct(t *testing.T) {
	nestedOrig := StructType{
		Fields: []Type{CharType{}},
	}
	typ := ArrayType{
		Length: 42,
		Typ:    nestedOrig,
	}
	serialize := typ.Serialize()
	deserialize, err := Deserialize(serialize)
	if err != nil {
		t.Error(err)
		return
	}
	deserArr, ok := deserialize.(ArrayType)
	if !ok {
		t.Error("Wrong type")
		return
	}
	if deserArr.Length != typ.Length {
		t.Errorf("Mistmatch! Expected length %d, got %d", deserArr.Length, typ.Length)
		return
	}
	if deserArr.Typ.typId() != typ.Typ.typId() {
		t.Errorf("Mistmatch! Expected length %d, got %d", deserArr.Typ.typId(), typ.Typ.typId())
		return
	}
	nestedArr, ok := deserArr.Typ.(StructType)
	if !ok {
		t.Error("Wrong nested type")
		return
	}
	if len(nestedArr.Fields) != len(nestedOrig.Fields) {
		t.Errorf("Mistmatch! Expected length %d, got %d", len(nestedArr.Fields), len(nestedOrig.Fields))
		return
	}
	if nestedArr.Fields[0].typId() != nestedOrig.Fields[0].typId() {
		t.Errorf("Mistmatch! Expected type %d, got %d", nestedArr.Fields[0].typId(), nestedOrig.Fields[0].typId())
		return
	}
}
