# json2protobuf
Based on golang, json schema converts protobuffer dynamic description and .proto file lib.

# Usage

```
package main

import (
    "io/ioutil"
    
    "github.com/EAHITechnology/json2prorobuf/proto"
    "github.com/EAHITechnology/json2prorobuf"
)

func main() {
    j2pManager := NewJson2PbParserManager()

	jss := []proto.JsonFieldSchema{}
	jss = append(jss, proto.JsonFieldSchema{
		MsgName: "TestField",
		Fields: []proto.Field{
			{
				Name: "testfieldtype",
				Type: proto.Typ{
					Name: proto.Int64Kind,
				},
			},
		},
	})
	jss = append(jss, proto.JsonFieldSchema{
		MsgName: "Test",
		Fields: []proto.Field{
			{
				Name: "basetype",
				Type: proto.Typ{
					Name: proto.Int64Kind,
				},
			},
			{
				Name: "arraytype",
				Type: proto.Typ{
					Name:        proto.ArrayKind,
					ElementType: proto.StringKind,
				},
			},
			{
				Name: "mapType",
				Type: proto.Typ{
					Name:      "map",
					KeyType:   proto.StringKind,
					ValueType: proto.FloatKind,
				},
			},
			{
				Name: "messageType",
				Type: proto.Typ{
					Name: proto.NewKind("TestField"),
				},
			},
			{
				Name: "messageArrayType",
				Type: proto.Typ{
					Name:        proto.ArrayKind,
					ElementType: proto.NewKind("TestField"),
				},
			},
			{
				Name: "messageMapType",
				Type: proto.Typ{
					Name:      proto.MapKind,
					KeyType:   proto.StringKind,
					ValueType: proto.NewKind("TestField"),
				},
			},
		},
	})

	desc := proto.JsonProtoDesc{
		Pkg:          "common",
		GoPkg:        "./json2prorobuf",
		FieldSchemas: jss,
	}
	if err := j2pManager.AddItem("testproto.proto", desc); err != nil {
		panic(err)
	}

	str, err := j2pManager.Dump("testproto.proto")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("testproto.proto", []byte(str), 0644)
	if err != nil {
		panic(err)
	}
}
```
