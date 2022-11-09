package json2prorobuf

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/EAHITechnology/json2prorobuf/proto"
	"github.com/EAHITechnology/json2prorobuf/utils"

	pbproto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	ErrProtoNotExists     = errors.New("proto not exists")
	ErrProtoAlreadyExists = errors.New("proto already exists")

	ErrPbTypeNoExists = errors.New("proto type no exists")

	ErrMapKeyIsMessageType = errors.New("key should not message type")
)

type Json2PbParserManager struct {
	globalProtoMap map[string]pref.FileDescriptor
	lock           sync.RWMutex
}

func NewJson2PbParserManager() *Json2PbParserManager {
	return &Json2PbParserManager{
		globalProtoMap: make(map[string]pref.FileDescriptor),
	}
}

func appendField(msdIdx int, field proto.Field, number int32, kind pref.Kind, pb *descriptorpb.FileDescriptorProto) {
	pb.MessageType[msdIdx].Field = append(pb.MessageType[msdIdx].Field, &descriptorpb.FieldDescriptorProto{
		Name:     pbproto.String(field.Name),
		JsonName: pbproto.String(field.Name),
		Number:   pbproto.Int32(number),
		Type:     descriptorpb.FieldDescriptorProto_Type(kind).Enum(),
	})
}

func appendNestedTypeField(msdIdx int, name string, number int32, kind pref.Kind, pb *descriptorpb.FileDescriptorProto) {
	pb.MessageType[msdIdx].NestedType[len(pb.MessageType[msdIdx].NestedType)-1].Field = append(pb.MessageType[msdIdx].NestedType[len(pb.MessageType[msdIdx].NestedType)-1].Field, &descriptorpb.FieldDescriptorProto{
		Name:     pbproto.String(name),
		JsonName: pbproto.String(name),
		Number:   pbproto.Int32(number),
		Label:    descriptorpb.FieldDescriptorProto_Label(pref.Optional).Enum(),
		Type:     descriptorpb.FieldDescriptorProto_Type(kind).Enum(),
	})
}

func fillingFieldsPbSchema(msdIdx int, typ proto.Kind, field proto.Field, number *int32, arrayFlag, mapKeyFlag, mapValFlag bool, pb *descriptorpb.FileDescriptorProto) error {
	switch typ {
	case proto.Int64Kind:
		if mapKeyFlag {
			appendNestedTypeField(msdIdx, "key", 1, pref.Int64Kind, pb)
			return nil
		}
		if mapValFlag {
			appendNestedTypeField(msdIdx, "value", 2, pref.Int64Kind, pb)
			return nil
		}
		appendField(msdIdx, field, *number, pref.Int64Kind, pb)
		if arrayFlag {
			pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].Label = descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum()
		}
		(*number)++
	case proto.Int32Kind:
		if mapKeyFlag {
			appendNestedTypeField(msdIdx, "key", 1, pref.Int32Kind, pb)
			return nil
		}
		if mapValFlag {
			appendNestedTypeField(msdIdx, "value", 2, pref.Int32Kind, pb)
			return nil
		}
		appendField(msdIdx, field, *number, pref.Int32Kind, pb)
		if arrayFlag {
			pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].Label = descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum()
		}
		(*number)++
	case proto.FloatKind:
		if mapKeyFlag {
			appendNestedTypeField(msdIdx, "key", 1, pref.FloatKind, pb)
			return nil
		}
		if mapValFlag {
			appendNestedTypeField(msdIdx, "value", 2, pref.FloatKind, pb)
			return nil
		}
		appendField(msdIdx, field, *number, pref.FloatKind, pb)
		if arrayFlag {
			pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].Label = descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum()
		}
		(*number)++
	case proto.StringKind:
		if mapKeyFlag {
			appendNestedTypeField(msdIdx, "key", 1, pref.StringKind, pb)
			return nil
		}
		if mapValFlag {
			appendNestedTypeField(msdIdx, "value", 2, pref.StringKind, pb)
			return nil
		}
		appendField(msdIdx, field, *number, pref.StringKind, pb)
		if arrayFlag {
			pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].Label = descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum()
		}
		(*number)++
	case proto.BoolKind:
		if mapKeyFlag {
			appendNestedTypeField(msdIdx, "key", 1, pref.BoolKind, pb)
			return nil
		}
		if mapValFlag {
			appendNestedTypeField(msdIdx, "value", 2, pref.BoolKind, pb)
			return nil
		}
		appendField(msdIdx, field, *number, pref.BoolKind, pb)
		if arrayFlag {
			pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].Label = descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum()
		}
		(*number)++
	case proto.ArrayKind:
		return fillingFieldsPbSchema(msdIdx, field.Type.ElementType, field, number, true, false, false, pb)
	case proto.MapKind:
		if !strings.HasSuffix(field.Name, "_map") {
			field.Name = field.Name + "_map"
		}

		entryName := ""
		fns := strings.Split(field.Name, "_")
		for _, fn := range fns {
			entryName += utils.FirstUpper(fn)
		}
		entryName += "Entry"

		appendField(msdIdx, field, *number, pref.MessageKind, pb)
		pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].Label = descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum()
		pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].TypeName = pbproto.String("." + *(pb.Package) + "." + *(pb.MessageType[msdIdx].Name) + "." + entryName)
		pb.MessageType[msdIdx].NestedType = append(pb.MessageType[msdIdx].NestedType, &descriptorpb.DescriptorProto{
			Name:  pbproto.String(entryName),
			Field: []*descriptorpb.FieldDescriptorProto{},
			Options: &descriptorpb.MessageOptions{
				MapEntry: pbproto.Bool(true),
			},
		})
		if err := fillingFieldsPbSchema(msdIdx, field.Type.KeyType, field, number, false, true, false, pb); err != nil {
			return err
		}

		if err := fillingFieldsPbSchema(msdIdx, field.Type.ValueType, field, number, false, false, true, pb); err != nil {
			return err
		}
		(*number)++
	default:
		if mapKeyFlag {
			return ErrMapKeyIsMessageType
		}
		if mapValFlag {
			appendNestedTypeField(msdIdx, "value", 2, pref.MessageKind, pb)
			pb.MessageType[msdIdx].NestedType[len(pb.MessageType[msdIdx].NestedType)-1].Field[1].TypeName = pbproto.String("." + *(pb.Package) + "." + field.Type.ValueType.String())
			return nil
		}
		appendField(msdIdx, field, *number, pref.MessageKind, pb)
		pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].TypeName = pbproto.String("." + *(pb.Package) + "." + field.Type.Name.String())
		if arrayFlag {
			pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].TypeName = pbproto.String("." + *(pb.Package) + "." + field.Type.ElementType.String())
			pb.MessageType[msdIdx].Field[len(pb.MessageType[msdIdx].Field)-1].Label = descriptorpb.FieldDescriptorProto_Label(pref.Repeated).Enum()
		}
		(*number)++
	}
	return nil
}

func fillingServicesPbSchema(service proto.JsonServiceSchema, pb *descriptorpb.FileDescriptorProto) {
	pb.Service = append(pb.Service, &descriptorpb.ServiceDescriptorProto{
		Name:   pbproto.String(service.Name),
		Method: []*descriptorpb.MethodDescriptorProto{},
	})

	for _, sd := range service.ServiceDescs {
		pb.Service[len(pb.Service)-1].Method = append(pb.Service[len(pb.Service)-1].Method, &descriptorpb.MethodDescriptorProto{
			Name:            pbproto.String(sd.Name),
			InputType:       pbproto.String(sd.Input),
			OutputType:      pbproto.String(sd.Output),
			ClientStreaming: pbproto.Bool(sd.ClientStreaming),
			ServerStreaming: pbproto.Bool(sd.ServerStreaming),
		})
	}
}

func (jm *Json2PbParserManager) parser(fileName string, desc proto.JsonProtoDesc) (pref.FileDescriptor, error) {
	pb := &descriptorpb.FileDescriptorProto{
		Syntax:      pbproto.String("proto3"),
		Name:        pbproto.String(fileName),
		Package:     pbproto.String(desc.Pkg),
		MessageType: []*descriptorpb.DescriptorProto{},
		Service:     []*descriptorpb.ServiceDescriptorProto{},
	}

	for msdIdx, s := range desc.FieldSchemas {
		pb.MessageType = append(pb.MessageType, &descriptorpb.DescriptorProto{
			Name:  pbproto.String(s.MsgName),
			Field: []*descriptorpb.FieldDescriptorProto{},
		})

		var number int32 = 1
		for _, field := range s.Fields {
			if err := fillingFieldsPbSchema(msdIdx, field.Type.Name, field, &number, false, false, false, pb); err != nil {
				return nil, err
			}
		}
	}

	for _, service := range desc.ServiceSchemas {
		fillingServicesPbSchema(service, pb)
	}

	pb.Options = &descriptorpb.FileOptions{
		GoPackage: pbproto.String(desc.GoPkg),
	}

	return protodesc.NewFile(pb, nil)
}

func (jm *Json2PbParserManager) Dump(fileName string) (string, error) {
	jm.lock.RLock()
	defer jm.lock.RUnlock()

	fd, ok := jm.globalProtoMap[fileName]
	if !ok {
		return "", ErrProtoNotExists
	}

	dumpStr := ""
	dumpStr += "syntax = \"proto3\";\n"
	dumpStr += "package " + string(fd.Package().Name()) + ";\n\n"
	dumpStr += "option go_package =\"" + fd.Options().(*descriptorpb.FileOptions).GetGoPackage() + "\";\n\n"

	for idx := 0; idx < fd.Messages().Len(); idx++ {
		msgName := fd.Messages().Get(idx).FullName()
		dumpStr += "message " + string(msgName.Name()) + " {\n"
		var number int64 = 0
		for fieldsIdx := 0; fieldsIdx < fd.Messages().Get(idx).Fields().Len(); fieldsIdx++ {
			number++
			field := fd.Messages().Get(idx).Fields().Get(fieldsIdx)
			value := field.Kind().String()
			if field.IsMap() {
				value = field.MapValue().Kind().String()
				if field.MapValue().Kind() == pref.MessageKind {
					value = string(field.MapValue().Message().Name())
				}
				dumpStr += "\tmap<" + field.MapKey().Kind().String() + "," + value + "> " + strings.Replace(string(field.FullName().Name()), "_map", "", -1) + " = " + strconv.FormatInt(number, 10) + ";\n"
				continue
			}

			if field.IsList() {
				if field.Kind() == pref.MessageKind {
					value = string(field.Message().Name())
				}
				dumpStr += "\trepeated " + value + " " + string(field.FullName().Name()) + " = " + strconv.FormatInt(number, 10) + ";\n"
				continue
			}

			if field.Kind() == pref.MessageKind {
				value = string(field.Message().Name())
			}
			dumpStr += "\t" + value + " " + string(field.FullName().Name()) + " = " + strconv.FormatInt(number, 10) + ";\n"
		}
		dumpStr += "}\n\n"
	}

	for idx := 0; idx < fd.Services().Len(); idx++ {
		dumpStr += "service " + string(fd.Services().Get(idx).Name()) + " {\n"
		for midx := 0; midx < fd.Services().Get(idx).Methods().Len(); midx++ {
			method := fd.Services().Get(idx).Methods().Get(midx)
			dumpStr += "\trpc " + string(method.Name()) + " (" + string(method.Input().Name()) + ") returns (" + string(method.Output().Name()) + ");\n"
		}
		dumpStr += "}\n"
	}

	return dumpStr, nil
}

func (jm *Json2PbParserManager) AddItem(fileName string, desc proto.JsonProtoDesc) error {
	jm.lock.Lock()
	defer jm.lock.Unlock()

	fileDescriptor, err := jm.parser(fileName, desc)
	if err != nil {
		return err
	}
	jm.globalProtoMap[fileName] = fileDescriptor
	return nil
}

func (jm *Json2PbParserManager) GetPBSchema(fileName string) (pref.FileDescriptor, error) {
	jm.lock.RLock()
	defer jm.lock.RUnlock()

	fd, ok := jm.globalProtoMap[fileName]
	if !ok {
		return nil, ErrProtoNotExists
	}

	return fd, nil
}
