package csv

import (
	"bytes"
	"github.com/gocarina/gocsv"
	"io"
	"strings"
)

func Parse(data string, out interface{}) error {
	return gocsv.UnmarshalDecoder(&trimDecoder{
		gocsv.LazyCSVReader(bytes.NewBufferString(data)),
	}, out)
}

type trimDecoder struct {
	csvReader gocsv.CSVReader
}

func (c *trimDecoder) GetCSVRow() ([]string, error) {
	recoder, err := c.csvReader.Read()
	for i, r := range recoder {
		recoder[i] = strings.TrimSpace(r)
	}
	return recoder, err
}

func (c *trimDecoder) GetCSVRows() ([][]string, error) {
	var records [][]string
	for {
		record, err := c.GetCSVRow()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}
