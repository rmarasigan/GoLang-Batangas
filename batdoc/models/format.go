package models

import (
	"github.com/uadmin/uadmin"
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
	return ""
}
