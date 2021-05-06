package types

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Type interface {
	Name() string
	Size() uint64
	// GetSuperTypes Gets all types that such that the current type is equal to or a subtype of
	// Ordered from same type up to highest subtype
	GetSuperTypes() []Type
}

type CharType struct {
	Type
}

func (typ CharType) Name() string          { return "Char" }
func (typ CharType) Size() uint64          { return 2 }
func (typ CharType) GetSuperTypes() []Type { return []Type{typ} }

type Int32Type struct {
	Type
}

func (typ Int32Type) Name() string          { return "Int32" }
func (typ Int32Type) Size() uint64          { return 4 }
func (typ Int32Type) GetSuperTypes() []Type { return []Type{typ} }

type Int64Type struct {
	Type
}

func (typ Int64Type) Name() string          { return "Int64" }
func (typ Int64Type) Size() uint64          { return 8 }
func (typ Int64Type) GetSuperTypes() []Type { return []Type{typ} }

type Float32Type struct {
	Type
}

func (typ Float32Type) Name() string          { return "Float32" }
func (typ Float32Type) Size() uint64          { return 4 }
func (typ Float32Type) GetSuperTypes() []Type { return []Type{typ} }

type Float64Type struct {
	Type
}

func (typ Float64Type) Name() string          { return "Float64" }
func (typ Float64Type) Size() uint64          { return 8 }
func (typ Float64Type) GetSuperTypes() []Type { return []Type{typ} }

type BoolType struct {
	Type
}

func (typ BoolType) Name() string          { return "Bool" }
func (typ BoolType) Size() uint64          { return 1 }
func (typ BoolType) GetSuperTypes() []Type { return []Type{typ} }

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

var globalSuperTypeCache = createSuperTypeCache()

// GetSuperTypes Get all super types and the type itself
func (typ StructType) GetSuperTypes() []Type {
	if cachedSuperTypes := globalSuperTypeCache.get(typ); cachedSuperTypes != nil {
		return cachedSuperTypes
	}
	noFields := len(typ.Fields)
	superTypes := make([]Type, noFields)
	for i := range typ.Fields {
		superTypes[i] = StructType{
			Fields: typ.Fields[:(noFields - i)],
		}
	}
	go globalSuperTypeCache.put(typ, superTypes) // So we don't need to wait for write-back
	return superTypes
}

func (typ StructType) TrimToSubType(subType StructType, data []byte) ([]byte, error) {
	if uint64(len(data)) != typ.Size() {
		return nil, errors.New("invalid data length")
	}
	if !subType.IsSubtype(typ) {
		return nil, errors.New("not actually a subtype")
	}
	return data[:subType.Size()], nil
}

func Trim(typ Type, subType Type, data []byte) ([]byte, error) {
	if typ == subType {
		return data, nil
	}
	structType, isStruct := typ.(StructType)
	subStructType, isSubStruct := subType.(StructType)
	if !isStruct || !isSubStruct {
		return nil, errors.New("subtyping only exists between structs")
	}
	return structType.TrimToSubType(subStructType, data)
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
func (typ UnionType) GetSuperTypes() []Type { return []Type{typ} }

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
func (typ ArrayType) GetSuperTypes() []Type { return []Type{typ} }

func createSuperTypeCache() superTypeCache {
	return superTypeCache{
		types:      map[Type][]Type{},
		writeMutex: &sync.Mutex{},
	}
}

type superTypeCache struct {
	types      map[Type][]Type
	writeMutex *sync.Mutex
}

func (cache *superTypeCache) get(typ Type) []Type {
	return cache.types[typ]
}
func (cache *superTypeCache) put(typ Type, superTypes []Type) {
	cache.writeMutex.Lock()
	defer cache.writeMutex.Unlock()
	if cache.get(typ) == nil {
		cache.types[typ] = superTypes
	}
}
