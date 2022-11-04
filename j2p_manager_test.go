package json2prorobuf

import (
	"io/ioutil"
	"testing"

	"github.com/EAHITechnology/json2prorobuf/proto"

	"github.com/stretchr/testify/assert"
)

func TestNewJson2PbParserManager(t *testing.T) {
	j2pManager := NewJson2PbParserManager()

	jss := []proto.JsonSchema{}
	jss = append(jss, proto.JsonSchema{
		MsgName: "Test2",
		Fields: []proto.Field{
			{
				Name: "test2basetype",
				Type: proto.Typ{
					Name: proto.Int64Kind,
				},
			},
		},
	})
	jss = append(jss, proto.JsonSchema{
		MsgName: "Test3",
		Fields: []proto.Field{
			{
				Name: "test3basetype",
				Type: proto.Typ{
					Name: proto.Int64Kind,
				},
			},
		},
	})
	jss = append(jss, proto.JsonSchema{
		MsgName: "Test1",
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
					Name: proto.NewKind("Test2"),
				},
			},
			{
				Name: "messageArrayType",
				Type: proto.Typ{
					Name:        proto.ArrayKind,
					ElementType: proto.NewKind("Test2"),
				},
			},
			{
				Name: "messageMapType",
				Type: proto.Typ{
					Name:      proto.MapKind,
					KeyType:   proto.StringKind,
					ValueType: proto.NewKind("Test3"),
				},
			},
		},
	})
	err := j2pManager.AddItem("testproto.proto", jss)
	assert.Equal(t, nil, err)

	str, err := j2pManager.Dump("testproto.proto")
	assert.Equal(t, nil, err)

	err = ioutil.WriteFile("testproto.proto", []byte(str), 0644)
	assert.Equal(t, nil, err)
}
