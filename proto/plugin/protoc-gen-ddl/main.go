package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/gertd/go-pluralize"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
)

var (
	text = `CREATE TABLE {{ .Name }} (
  {{- range $_, $col := .Columns }}
  {{ $col.Name }} {{ $col.Type }} NOT NULL,
  {{- end }}
) PRIMARY KEY (id);
`
)

type Table struct {
	Name    string
	Columns []Column
}

type Column struct {
	Name string
	Type string
}

func readStdin() (*plugin.CodeGeneratorRequest, error) {
	r := os.Stdin
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request")
	}

	var req plugin.CodeGeneratorRequest
	if err = proto.Unmarshal(buf, &req); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal request")
	}
	return &req, nil
}

func writeStdout(resp *plugin.CodeGeneratorResponse) error {
	buf, err := proto.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "failed to marshal response")

	}
	_, err = os.Stdout.Write(buf)
	return err
}

func appendFile(codeGeneratorResponse *plugin.CodeGeneratorResponse, text, fileName string, data interface{}) error {
	tpl, err := template.New("").Parse(text)
	if err != nil {
		return errors.Wrap(err, "failed to new template")
	}
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return errors.Wrap(err, "failed to execute template")
	}
	codeGeneratorResponse.File = append(codeGeneratorResponse.File, &plugin.CodeGeneratorResponse_File{
		Name:    proto.String(fileName),
		Content: proto.String(buf.String()),
	})
	return nil
}

func processReq(req *plugin.CodeGeneratorRequest) (*plugin.CodeGeneratorResponse, error) {
	plu := pluralize.NewClient()
	var res plugin.CodeGeneratorResponse
	for _, f := range req.ProtoFile {
		if f.GetPackage() == "entity" {
			for _, message := range f.MessageType {
				snakedMessageName := strcase.ToSnake(message.GetName())
				tableName := plu.Plural(snakedMessageName)
				tabl := &Table{
					Name: tableName,
				}
				for _, field := range message.Field {
					column := Column{
						Name: strcase.ToSnake(field.GetName()),
						Type: field.GetTypeName(),
					}
					typ := field.GetType()
					switch typ {
					case descriptor.FieldDescriptorProto_TYPE_UINT32:
						column.Type = "INT64"
					case descriptor.FieldDescriptorProto_TYPE_BOOL:
						column.Type = "BOOL"
					case descriptor.FieldDescriptorProto_TYPE_STRING:
						column.Type = "STRING(MAX)"
					default:
						if column.Name == "user_items" {
							column.Name = "user_item_ids"
							column.Type = "ARRAY<INT64>"
							break
						}
						log.Panicf("column type not match: table:%s, column:%s, type:%s", tableName, column.Name, typ)
					}
					tabl.Columns = append(tabl.Columns, column)
				}
				if err := appendFile(&res, text, fmt.Sprintf("db/%s.sql", tableName), tabl); err != nil {
					return nil, errors.Wrap(err, "failed to append file")
				}
			}
		}
	}
	return &res, nil
}

func run() error {
	req, err := readStdin()
	if err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	resp, err := processReq(req)
	if err != nil {
		return errors.Wrap(err, "failed to process request")
	}

	return writeStdout(resp)
}

// build and compile
// cd $GOPATH/src/github.com/TatsuNet/ex_gin_pb/proto/plugin/protoc-gen-ddl&& go build -o protoc-gen-ddl main.go && cd ../.. && protoc entity/*.proto -I. --plugin=./plugin/protoc-gen-ddl/protoc-gen-ddl --ddl_out=../
func main() {
	if err := run(); err != nil {
		log.Fatalf("failed to run:%+v", err)
	}
}
