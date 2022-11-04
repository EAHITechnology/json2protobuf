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
    j2pManager := json2prorobuf.NewJson2PbParserManager()
    
    schema := []proto.JsonSchema{}
    schema = append(schema, proto.JsonSchema{
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
    
    schema = append(schema, proto.JsonSchema{
        MsgName: "Test3",
		    Fields: []proto.Field{
			      {
				        Name: "test3basetype",
				        Type: proto.Typ {
					          Name: proto.Int64Kind,
				        },
			      },
		    },
    })
    
    schema = append(schema, proto.JsonSchema{
		    MsgName: "Test1",
   		  Fields: []proto.Field{
			      {
				        Name: "basetype",
			 	        Type: proto.Typ {
					          Name: proto.Int64Kind,
				        },
			      },
			      {
				        Name: "arraytype",
				        Type: proto.Typ {
					          Name:        proto.ArrayKind,
					          ElementType: proto.StringKind,
				        },
			      },
			      {
				        Name: "mapType",
				        Type: proto.Typ {
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
}
```
