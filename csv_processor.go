package csvutil

import (
	"encoding/csv"
	"io"
)

// CSVProcessor process CSV read from csv.Reader, and write to csv.Writer.
type CSVProcessor struct {
	reader        *csv.Reader
	writer        *csv.Writer
	headerHandler func([]string) ([]string, error)
	preBodyRead   func() error
	recordHandler func([]string) ([]string, error)
}

// NewCSVProcessor returns new processor from reader and writer.
func NewCSVProcessor(r *csv.Reader, w *csv.Writer) *CSVProcessor {
	return &CSVProcessor{
		reader: r,
		writer: w,
	}
}

// NewReadOnlyCSVProcessor returns new processor from reader.
// Processor does not write to CSV.
func NewReadOnlyCSVProcessor(r *csv.Reader) *CSVProcessor {
	return &CSVProcessor{reader: r}
}

// SetHeaderHanlder set function for calling on header line read.
func (csvp *CSVProcessor) SetHeaderHanlder(f func([]string) ([]string, error)) {
	csvp.headerHandler = f
}

// SetPreBodyRead set function for calling before first body line read.
func (csvp *CSVProcessor) SetPreBodyRead(f func() error) {
	csvp.preBodyRead = f
}

// SetRecordHandler set function for calling on each body line read.
func (csvp *CSVProcessor) SetRecordHandler(f func([]string) ([]string, error)) {
	csvp.recordHandler = f
}

// Process CSV.
// Read each line and apply each functions and write line.
func (csvp *CSVProcessor) Process() error {
	if csvp.headerHandler != nil {
		hdr, err := csvp.reader.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		pHdr, err := csvp.headerHandler(hdr)
		if err != nil {
			return err
		}
		if csvp.writer != nil && pHdr != nil {
			csvp.writer.Write(pHdr)
		}
	}

	if csvp.preBodyRead != nil {
		err := csvp.preBodyRead()
		if err != nil {
			return err
		}
	}

	for {
		rec, err := csvp.reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		pRec, err := csvp.recordHandler(rec)
		if err != nil {
			return err
		}
		if csvp.writer != nil && pRec != nil {
			csvp.writer.Write(pRec)
		}
	}

	return nil
}
