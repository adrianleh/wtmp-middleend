package command

import (
	"errors"
	"strings"

	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type AcceptTypeCommandHandler struct{}

func (AcceptTypeCommandHandler) Handle(frame CommandFrame) error {
	content, err := getType(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetById(frame.ClientId)
	return cl.RegisterType(content.typ)
}

type acceptTypeCommandContent struct {
	typ types.Type
}


func getType(data []byte) (getCommandContent, error) {
	// TODO: below is placeholder: depends on type serialization here
	nullNameIdx := strings.Index(string(data), "\x00")
	if nullNameIdx < 0 {
		return getCommandContent{}, errors.New("invalid format: no string found for type name")
	}
	name := string(data[:nullNameIdx]) // convert to string

	var typ types.Type
	if strings.HasPrefix(name, "Array") {
		typ = types.ArrayType{}
	} else if strings.HasPrefix(name, "Union") {
		typ = types.UnionType{}
	} else if strings.HasPrefix(name, "Struct") {
		typ = types.StructType{}
	} else {
		switch name {
		case "Char":
			typ = types.CharType{} 
		case "Int32":
			typ = types.Int32Type{}
		case "Int64":
			typ = types.Int32Type{}
		case "Float32":
			typ = types.Float32Type{}
		case "Float64":
			typ = types.Float64Type{}
		case "Bool":
			typ = types.BoolType{}
		default:
			return getCommandContent{}, errors.New("invalid type: what type is this??")
		}
	}

	return getCommandContent{
		typ: typ,
	}, nil
}
