package models

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/unidoc/unioffice/presentation"

	"github.com/uadmin/uadmin"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/spreadsheet"
)

type Format int

func (Format) PDF() Format {
	return 1
}

func (Format) TXT() Format {
	return 2
}

func (Format) HTML() Format {
	return 3
}

func (Format) DOCX() Format {
	return 4
}

func (Format) XLSX() Format {
	return 5
}

func (Format) PPTX() Format {
	return 6
}

func (f *Format) GetRawText(doc *DocumentVersion) string {
	if *f == f.PDF() { // PDF
		// OCR
		// First convert the PDF to TIFF image format
		// get base path of file
		filePath := strings.TrimPrefix(doc.File, "/")
		basePath := filepath.Dir(filePath)
		cmd := exec.Command("convert", "-density", "300", filePath, "-depth", "8", "-strip", "-background", "white", "-alpha", "off", filepath.Join(basePath, "doc.tiff"))
		rawCmd := strings.Join([]string{"convert", "-density 300", filePath, "-depth 8", "-strip", "-background white", "-alpha off", filepath.Join(basePath, "doc.tiff")}, " ")
		out, err := cmd.CombinedOutput()
		uadmin.Trail(uadmin.DEBUG, "%s output: %s", rawCmd, string(out))
		if err != nil {
			uadmin.Trail(uadmin.ERROR, "Unable to convert PDF to TIFF: %s", err)
			return ""
		}
		cmd = exec.Command("tesseract", filepath.Join(basePath, "doc.tiff"), filepath.Join(basePath, "rawtext"))
		rawCmd = strings.Join([]string{"tesseract", filepath.Join(basePath, "doc.tiff"), filepath.Join(basePath, "rawtext")}, " ")
		out, err = cmd.CombinedOutput()
		uadmin.Trail(uadmin.DEBUG, "%s output: %s", rawCmd, string(out))
		if err != nil {
			uadmin.Trail(uadmin.ERROR, "Unable to OCR TIFF document: %s", err)
			return ""
		}
		buf, err := ioutil.ReadFile(filepath.Join(basePath, "rawtext.txt"))
		if err != nil {
			uadmin.Trail(uadmin.ERROR, "Error reading raw text file: %s", err)
			return ""
		}
		return string(buf)
	}

	if *f == f.TXT() { // Text
		// Read file
		filePath := strings.TrimPrefix(doc.File, "/")
		buf, err := ioutil.ReadFile(filePath)
		if err != nil {
			uadmin.Trail(uadmin.ERROR, "Error reading text file: %s", err)
			return ""
		}
		return string(buf)
	}

	if *f == f.HTML() {
		// Strip tags
		return ""
	}

	if *f == f.DOCX() {
		doc, _ := document.Open(strings.TrimPrefix(doc.File, "/"))
		buf := ""
		for _, para := range doc.Paragraphs() {
			for _, run := range para.Runs() {
				buf += run.Text() + " "
			}
		}
		return buf
	}

	if *f == f.XLSX() {
		excel, _ := spreadsheet.Open(strings.TrimPrefix(doc.File, "/"))
		buf := ""
		for _, sheet := range excel.Sheets() {
			for col := 'A'; col <= 'Z'; col++ {
				for row := 1; row < 100; row++ {
					ref := string(col) + fmt.Sprint(row)
					if sheet.Cell(ref).GetString() != "" {
						buf += sheet.Cell(ref).GetString() + " "
					}
				}

			}
			return buf
		}
	}

	if *f == f.PPTX() {
		// Using Python but cannot pass the value to Go
		// python.Initialize()
		// python.PyRun_SimpleString("from pptx import Presentation\nprs = Presentation('." + doc.File + "')\ntext_runs = []\nfor slide in prs.slides:\n\tfor shape in slide.shapes:\n\t\tif not shape.has_text_frame: continue\n\t\tfor paragraph in shape.text_frame.paragraphs:\n\t\t\tfor run in paragraph.runs: text_runs.append(run.text)\nprint text_runs")
		// python.Finalize()

		// For .pptx only. Not working with other powerpoint extension.
		pres, _ := presentation.Open(strings.TrimPrefix(doc.File, "/"))
		for _, slide := range pres.Slides() {
			for _, c := range slide.X().CSld.SpTree.Choice {
				// Shape
				for _, sp := range c.Sp {
					if sp.TxBody == nil {
						continue
					}
					// Paragraph
					for _, para := range sp.TxBody.P {
						if para == nil {
							continue
						}
						// Run
						for _, run := range para.EG_TextRun {
							uadmin.Trail(uadmin.DEBUG, run.R.T)
						}
					}
				}
			}
		}
	}

	return ""

}
