package types

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

type Type interface {
	Name() string
	typId() byte
	Size() uint64
	Serialize() []byte
	Deserialize([]byte) (Type, error)
	// GetSuperTypes Gets all types that such that the current type is equal to or a subtype of
	// Ordered from same type up to highest subtype
	GetSuperTypes() []Type
}

var typIdMap = createTypeIdMap()

func createTypeIdMap() map[byte]Type {
	idMap := map[byte]Type{}
	types := []Type{CharType{}, Int32Type{}, Int64Type{}, Float32Type{}, Float64Type{}, BoolType{}, StructType{}, UnionType{}, ArrayType{}}
	for _, typ := range types {
		idMap[typ.typId()] = typ
	}
	return idMap
}

const (
	charTypeId    = 0
	int32TypeId   = 1
	int64TypeId   = 2
	float32TypeId = 3
	float64TypeId = 4
	boolTypeId    = 5
	structTypeId  = 6
	unionTypeId   = 7
	arrayTypeId   = 8
)

func Deserialize(raw []byte) (Type, error) {
	if len(raw) < 5 {
		return nil, errors.New("too short")
	}
	typId := raw[4]
	typ := typIdMap[typId]
	desertyp, err := typ.Deserialize(raw)
	if err != nil {
		return nil, err
	}
	return desertyp, nil
}

type CharType struct {
	Type
}

func (typ CharType) Name() string          { return "Char" }
func (typ CharType) typId() byte           { return charTypeId }
func (typ CharType) Size() uint64          { return 2 }
func (typ CharType) GetSuperTypes() []Type { return []Type{typ} }
func (typ CharType) Serialize() []byte {
	ser := make([]byte, 4)
	binary.BigEndian.PutUint32(ser, 5)
	return append(ser, typ.typId())
}
func (typ CharType) Deserialize(data []byte) (Type, error) {
	if len(data) != 5 {
		return typ, errors.New("invalid length")
	}
	return typ, nil
}

type Int32Type struct {
	Type
}

func (typ Int32Type) Name() string          { return "Int32" }
func (typ Int32Type) typId() byte           { return int32TypeId }
func (typ Int32Type) Size() uint64          { return 4 }
func (typ Int32Type) GetSuperTypes() []Type { return []Type{typ} }
func (typ Int32Type) Serialize() []byte {
	ser := make([]byte, 4)
	binary.BigEndian.PutUint32(ser, 5)

	return append(ser, typ.typId())
}
func (typ Int32Type) Deserialize(data []byte) (Type, error) {
	if len(data) != 5 {
		return typ, errors.New("invalid length")
	}
	return typ, nil
}

type Int64Type struct {
	Type
}

func (typ Int64Type) Name() string          { return "Int64" }
func (typ Int64Type) typId() byte           { return int64TypeId }
func (typ Int64Type) Size() uint64          { return 8 }
func (typ Int64Type) GetSuperTypes() []Type { return []Type{typ} }
func (typ Int64Type) Serialize() []byte {
	ser := make([]byte, 4)
	binary.BigEndian.PutUint32(ser, 5)

	return append(ser, typ.typId())
}
func (typ Int64Type) Deserialize(data []byte) (Type, error) {
	if len(data) != 5 {
		return typ, errors.New("invalid length")
	}
	return typ, nil
}

type Float32Type struct {
	Type
}

func (typ Float32Type) Name() string          { return "Float32" }
func (typ Float32Type) typId() byte           { return float32TypeId }
func (typ Float32Type) Size() uint64          { return 4 }
func (typ Float32Type) GetSuperTypes() []Type { return []Type{typ} }
func (typ Float32Type) Serialize() []byte {
	ser := make([]byte, 4)
	binary.BigEndian.PutUint32(ser, 5)

	return append(ser, typ.typId())
}
func (typ Float32Type) Deserialize(data []byte) (Type, error) {
	if len(data) != 5 {
		return typ, errors.New("invalid length")
	}
	return typ, nil
}

type Float64Type struct {
	Type
}

func (typ Float64Type) Name() string          { return "Float64" }
func (typ Float64Type) typId() byte           { return float64TypeId }
func (typ Float64Type) Size() uint64          { return 8 }
func (typ Float64Type) GetSuperTypes() []Type { return []Type{typ} }
func (typ Float64Type) Serialize() []byte {
	ser := make([]byte, 4)
	binary.BigEndian.PutUint32(ser, 5)

	return append(ser, typ.typId())
}
func (typ Float64Type) Deserialize(data []byte) (Type, error) {
	if len(data) != 5 {
		return typ, errors.New("invalid length")
	}
	return typ, nil
}

type BoolType struct {
	Type
}

func (typ BoolType) Name() string          { return "Bool" }
func (typ BoolType) typId() byte           { return boolTypeId }
func (typ BoolType) Size() uint64          { return 1 }
func (typ BoolType) GetSuperTypes() []Type { return []Type{typ} }
func (typ BoolType) Serialize() []byte {
	ser := make([]byte, 4)
	binary.BigEndian.PutUint32(ser, 5)
	return append(ser, typ.typId())
}
func (typ BoolType) Deserialize(data []byte) (Type, error) {
	if len(data) != 5 {
		return typ, errors.New("invalid length")
	}
	return typ, nil
}

type StructType struct {
	Type   `json:"-"`
	Fields []Type `json:"fields"`
}

func (typ StructType) Name() string {
	name := "Struct"
	for _, fieldTyp := range typ.Fields {
		name += "-" + fieldTyp.Name()
	}
	return name
}
func (typ StructType) typId() byte {
	return structTypeId
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

func (typ StructType) TrimToSuperType(subType StructType, data []byte) ([]byte, error) {
	if uint64(len(data)) != typ.Size() {
		return nil, errors.New("invalid data length")
	}
	if !subType.IsSubtype(typ) {
		return nil, errors.New("not actually a subtype")
	}
	return data[:subType.Size()], nil
}

func Trim(typ Type, superType Type, data []byte) ([]byte, error) {
	if typ == superType {
		return data, nil
	}
	structType, isStruct := typ.(StructType)
	superStructType, isSuperStruct := superType.(StructType)
	if !isStruct || !isSuperStruct {
		return nil, errors.New("subtyping only exists between structs")
	}
	return structType.TrimToSuperType(superStructType, data)
}
func (typ StructType) Serialize() []byte {
	serFields := make([]byte, 0)
	length := uint32(9)
	noFields := uint32(0)
	for _, field := range typ.Fields {
		fieldSer := field.Serialize()
		fieldLen := binary.BigEndian.Uint32(fieldSer[0:4])
		length = length + fieldLen
		noFields = noFields + 1
		serFields = append(serFields, fieldSer...)
	}
	lenRaw := make([]byte, 4)
	binary.BigEndian.PutUint32(lenRaw, length)
	noFieldRaw := make([]byte, 4)
	binary.BigEndian.PutUint32(noFieldRaw, noFields)
	return append(append(lenRaw, typ.typId()), append(noFieldRaw, serFields...)...)
}

func (typ StructType) Deserialize(data []byte) (Type, error) {
	if len(data) < 8 {
		return typ, errors.New("too short")
	}
	noFields := binary.BigEndian.Uint32(data[5:9])
	startIdx := uint32(4 + 1 + 4)
	var fields []Type
	for i := uint32(0); i < noFields; i++ {
		lenMaxIdx := startIdx + 4
		typIdIdx := lenMaxIdx
		if uint32(len(data)) < typIdIdx {
			return typ, errors.New("too short")
		}
		fieldLen := binary.BigEndian.Uint32(data[startIdx:lenMaxIdx])
		endIdx := startIdx + fieldLen
		if uint32(len(data)) < endIdx {
			return typ, errors.New("too short")
		}
		fieldArr := data[startIdx:endIdx]
		fieldTypId := data[typIdIdx]
		fieldTyp := typIdMap[fieldTypId]
		fieldDeserTyp, err := fieldTyp.Deserialize(fieldArr)
		if err != nil {
			return typ, err
		}
		fields = append(fields, fieldDeserTyp)
		startIdx = endIdx
	}
	typ.Fields = fields
	return typ, nil
}

type UnionType struct {
	Type    `json:"-"`
	Members []Type `json:"members"`
}

func (typ UnionType) Name() string {
	name := "Union"
	for _, memberTyp := range typ.Members {
		name += "-" + memberTyp.Name()
	}
	return name
}
func (typ UnionType) typId() byte {
	return unionTypeId
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
func (typ UnionType) Serialize() []byte {
	serFields := make([]byte, 0)
	length := uint32(9)
	noFields := uint32(0)
	for _, field := range typ.Members {
		fieldSer := field.Serialize()
		fieldLen := binary.BigEndian.Uint32(fieldSer[0:4])
		length = length + fieldLen
		noFields = noFields + 1
		serFields = append(serFields, fieldSer...)
	}
	lenRaw := make([]byte, 4)
	binary.BigEndian.PutUint32(lenRaw, length)
	noFieldRaw := make([]byte, 4)
	binary.BigEndian.PutUint32(noFieldRaw, noFields)
	return append(append(lenRaw, typ.typId()), append(noFieldRaw, serFields...)...)
}

func (typ UnionType) Deserialize(data []byte) (Type, error) {
	if len(data) < 8 {
		return typ, errors.New("too short")
	}
	noFields := binary.BigEndian.Uint32(data[5:9])
	startIdx := uint32(4 + 1 + 4)
	var fields []Type
	for i := uint32(0); i < noFields; i++ {
		lenMaxIdx := startIdx + 4
		typIdIdx := lenMaxIdx
		if uint32(len(data)) < typIdIdx {
			return typ, errors.New("too short")
		}
		fieldLen := binary.BigEndian.Uint32(data[startIdx:lenMaxIdx])
		endIdx := startIdx + fieldLen
		if uint32(len(data)) < endIdx {
			return typ, errors.New("too short")
		}
		fieldArr := data[startIdx:endIdx]
		fieldTypId := data[typIdIdx]
		fieldTyp := typIdMap[fieldTypId]
		fieldDeserTyp, err := fieldTyp.Deserialize(fieldArr)
		if err != nil {
			return typ, err
		}
		fields = append(fields, fieldDeserTyp)
		startIdx = endIdx
	}
	typ.Members = fields
	return typ, nil
}

type ArrayType struct {
	Type   `json:"-"`
	Length uint64 `json:"length"`
	Typ    Type   `json:"typ"`
}

func (typ ArrayType) Name() string {
	return fmt.Sprintf("Array-%s-%d", typ.Typ.Name(), typ.Length)
}
func (typ ArrayType) typId() byte {
	return arrayTypeId
}
func (typ ArrayType) Size() uint64 {
	return typ.Length * typ.Typ.Size()
}
func (typ ArrayType) GetSuperTypes() []Type { return []Type{typ} }
func (typ ArrayType) Serialize() []byte {
	typSer := typ.Typ.Serialize()
	typSerLen := uint32(len(typSer))
	lenRaw := make([]byte, 4)
	binary.BigEndian.PutUint32(lenRaw, typSerLen+4+1+8)
	innerLenRaw := make([]byte, 8)
	binary.BigEndian.PutUint64(innerLenRaw, typ.Length)
	result := append(lenRaw, typ.typId())
	result = append(result, innerLenRaw...)
	result = append(result, typSer...)
	return result
}
func (typ ArrayType) Deserialize(data []byte) (Type, error) {
	nonForeignSize := 4 + 1 + 8
	minLength := nonForeignSize + 4 + 1
	if len(data) < minLength {
		return typ, errors.New("too short")
	}
	lenRaw := data[5 : 5+8]
	typ.Length = binary.BigEndian.Uint64(lenRaw)
	innerTypId := data[minLength-1]
	inner := typIdMap[innerTypId]
	innerDeser, err := inner.Deserialize(data[nonForeignSize:])
	if err != nil {
		return typ, err
	}
	typ.Typ = innerDeser
	return typ, nil
}

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
