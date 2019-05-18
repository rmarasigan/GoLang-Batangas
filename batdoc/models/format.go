package models

import (
	"fmt"
	"github.com/uadmin/uadmin"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/presentation"
	"github.com/unidoc/unioffice/spreadsheet"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
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
			//uadmin.Trail(uadmin.DEBUG, buf)
			return buf

		}
	}
	if *f == f.PPTX() {
		pres, _ := presentation.Open("/home/ubuntu/Documents/test.pptx")
		buf := ""
		_ = buf
		for _, slide := range pres.Slides() {
			fmt.Println("slide")
			for _, box := range slide.PlaceHolders() {
				if box.X() == nil || box.X().TxBody == nil {
					continue
				}
				for _, para := range box.X().TxBody.P {
					if para == nil {
						continue
					}
					for _, run := range para.EG_TextRun {
						uadmin.Trail(uadmin.DEBUG, run)
					}
				}
			}
		}
	}
	return ""
}
