package command

import (
	"errors"
	"strings"

	"github.com/adrianleh/WTMP-middleend/client"
	"github.com/adrianleh/WTMP-middleend/types"
)

type SendCommandHandler struct{}

func (SendCommandHandler) Handle(frame CommandFrame) error {
	content, err := sendData(frame.Data)
	if err != nil {
		return err
	}

	cl := client.Clients.GetByName(content.target)
	return cl.Push(content.typ, content.msg)
}

type sendCommandContent struct {
	typ types.Type
	target string
	msg []byte
}

func sendData(data []byte) (sendCommandContent, error) {
	// TODO: below is placeholder: depends on type serialization here
	// Spaghetti code to be refactored
	nullNameIdx := strings.Index(string(data), "\x00")
	if nullNameIdx < 0 {
		return sendCommandContent{}, errors.New("invalid format: no string found for type name")
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
			return sendCommandContent{}, errors.New("invalid type: what type is this??")
		}
	}

	restOfData := data[nullNameIdx+1:]
	nullTargetIdx := strings.Index(string(restOfData), "\x00")
	if nullTargetIdx < 0 {
		return sendCommandContent{}, errors.New("invalid format: no string found for target")
	}
	target := string(restOfData[:nullTargetIdx]) // convert to string
	msg := restOfData[nullTargetIdx+1:] // TODO: I think we said empty msgs are okay?

	return sendCommandContent{
		typ: typ,
		target: target,
		msg: msg,
	}, nil
}
