package types

import (
	"fmt"
	"reflect"
)

type Type struct{}

func (typ Type) Name() string { return "" }
func (typ Type) Size() uint64 { return 0 }

type CharType struct {
	Type
}

func (typ CharType) Name() string { return "Char" }
func (typ CharType) Size() uint64 { return 2 }

type Int32Type struct {
	Type
}

func (typ Int32Type) Name() string { return "Int32" }
func (typ Int32Type) Size() uint64 { return 4 }

type Int64Type struct {
	Type
}

func (typ Int64Type) Name() string { return "Int64" }
func (typ Int64Type) Size() uint64 { return 8 }

type Float32Type struct {
	Type
}

func (typ Float32Type) Name() string { return "Float32" }
func (typ Float32Type) Size() uint64 { return 4 }

type Float64Type struct {
	Type
}

func (typ Float64Type) Name() string { return "Float64" }
func (typ Float64Type) Size() uint64 { return 8 }

type BoolType struct {
	Type
}

func (typ BoolType) Name() string { return "Bool" }
func (typ BoolType) Size() uint64 { return 1 }

type StructType struct {
	Type
	Fields []Type
}

func (typ StructType) Name() string {
	name := "Struct"
	for _, fieldTyp := range typ.Fields {
		name += "-" + fieldTyp.Name()
	}
	return name
}
func (typ StructType) Size() uint64 {
	size := uint64(0)
	for _, fieldTyp := range typ.Fields {
		size += fieldTyp.Size()
	}
	return size
}
func (typ StructType) IsSubtype(superTyp StructType) bool {
	if len(superTyp.Fields) >= len(typ.Fields) {
		return false
	}
	for i, superField := range superTyp.Fields {
		if !reflect.DeepEqual(typ.Fields[i], superField) {
			return false
		}
	}
	return true
}

type UnionType struct {
	Type
	Members []Type
}

func (typ UnionType) Name() string {
	name := "Union"
	for _, memberTyp := range typ.Members {
		name += "-" + memberTyp.Name()
	}
	return name
}
func (typ UnionType) Size() uint64 {
	size := uint64(0)
	for _, memberTyp := range typ.Members {
		if memberTyp.Size() > size {
			size = memberTyp.Size()
		}
	}
	return size
}

type ArrayType struct {
	Type
	Length uint64
	Typ    Type
}

func (typ ArrayType) Name() string {
	return fmt.Sprintf("Array-%s-%d", typ.Typ.Name(), typ.Length)
}
func (typ ArrayType) Size() uint64 {
	return typ.Length * typ.Typ.Size()
}
