package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/iancoleman/strcase"
	"github.com/jhump/protoreflect/desc/protoparse"
)

//go:embed "all:*"
var templates embed.FS

const (
	ProtoSourceFile = "protobufs/api.proto"
)

type destinationDatum struct {
	destFile       string
	sourceFile     string
	templateName   string
	path           string
	pkg            string
	excludedFields []string
}

var destinationData = []destinationDatum{
	{
		destFile:       "api.proto.datadog_formatter.go",
		sourceFile:     "templates/datadog_formatter.go.tmpl",
		templateName:   "datadog_formatter.go.tmpl",
		path:           "internal/writer/datadog",
		pkg:            "datadog",
		excludedFields: []string{"running"},
	},
	{
		destFile:       "api.proto.influxdb_formatter.go",
		sourceFile:     "templates/influxdb_formatter.go.tmpl",
		templateName:   "influx_db_formatter.go.tmpl",
		path:           "internal/writer/influxdb",
		pkg:            "influxdb",
		excludedFields: []string{"unix_time", "running"},
	},
	{
		destFile:       "api.proto.publisher_formatter.go",
		sourceFile:     "templates/publisher_formatter.go.tmpl",
		templateName:   "publisher_formatter.go.tmpl",
		path:           "internal/collector/publisher",
		pkg:            "publisher",
		excludedFields: []string{"unix_time", "running"},
	},
}

type RowData struct {
	RawFieldName        string
	CamelCasedFieldName string
}

type FileData struct {
	Rows []RowData
}

// This could be written with the purpose of a protobuf file in mind, but given the difficulty of using it in
// that way and the simplicity of the generated files, opting to use this mechanic instead.
func main() {

	// marshalled with Protobuf
	log.Println("Starting to generate golang proto generated code")

	loadedPb, err := protoparse.Parser{}.ParseFiles(ProtoSourceFile)
	if err != nil {
		panic(err)
	}

	if len(loadedPb) == 0 {
		panic(fmt.Sprintf(
			"Unable to properly parse data from raw source. Check raw data source: %s",
			ProtoSourceFile,
		))
	}

	data := loadedPb[0]

	fieldRowsToWrite := make([]RowData, 0)

	for _, messageType := range data.GetMessageTypes() {
		if messageType.GetName() == "PublishRequest" {
			names := messageType.GetFields()
			for _, name := range names {
				fieldRowsToWrite = append(fieldRowsToWrite, RowData{
					RawFieldName:        name.GetName(),
					CamelCasedFieldName: strcase.ToCamel(name.GetName()),
				})
			}
		}
	}

	if os.Getenv("DEBUG") == "true" {
		log.Println("PROTOBUF MARSHALED DATA - START DEBUG")
		spew.Dump(data)
		log.Println("Writing these rows to generated files:")
		spew.Dump(fieldRowsToWrite)
		log.Println("PROTOBUF MARSHALED DATA - END DEBUG")
	}

	for _, dataToWrite := range destinationData {
		log.Printf("Processing %s\n", dataToWrite.sourceFile)

		filteredRows := []RowData{}
		for _, row := range fieldRowsToWrite {
			shouldExcludeField := false
			for _, excludedField := range dataToWrite.excludedFields {
				if excludedField == row.RawFieldName {
					shouldExcludeField = true
				}
			}
			if shouldExcludeField {
				continue
			}

			filteredRows = append(filteredRows, row)
		}

		fileToWrite, err := os.Create(fmt.Sprintf("%s/%s", dataToWrite.path, dataToWrite.destFile))
		if err != nil {
			log.Println("create file: ", err)
			return
		}

		// load template from path, and generate strings to write to destinations
		tmpl := template.Must(template.ParseFS(templates, dataToWrite.sourceFile))
		err = tmpl.Execute(fileToWrite, FileData{
			filteredRows,
		})
		if err != nil {
			panic(err)
		}

	}
}
