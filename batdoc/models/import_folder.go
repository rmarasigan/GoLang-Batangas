package models

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/uadmin/uadmin"
)

type ImportFolder struct {
	uadmin.Model
	Path string
}

func (i ImportFolder) String() string {
	return i.Path
}

func (i *ImportFolder) Save() {
	// Some magic code
	const BASE_PATH = "./media/"
	ROOT := i.Path

	if !strings.Contains(ROOT, ".") {
		fmt.Println("Saving...")
		folderParent := filepath.Join(BASE_PATH + filepath.Base(ROOT))
		os.MkdirAll(folderParent, os.ModePerm)
		Copy(ROOT, folderParent)
	}

	uadmin.Save(i)
}

// Copy copies src to dest, doesn't matter if src
// is a directory or a file.
func Copy(source, destination string) error {
	// Lstat returns FileInfo describing the named file.
	info, _ := os.Lstat(source)
	fmt.Println("Saving... Copy")

	// Check if it is a directory
	if info.IsDir() {
		return copy(source, destination, info)
	}

	return nil
}

// copy dispatches copy-funcs according to the mode.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func copy(source, destination string, info os.FileInfo) error {

	fmt.Println("Saving... copy")
	if info.Mode()&os.ModeSymlink != 0 {
		source, _ := os.Readlink(source)
		return os.Symlink(source, destination)
	}

	// Return Directory
	if info.IsDir() {
		return destinationCopy(source, destination, info)
	}

	// Creating the destination of file
	fileToCreate, _ := os.Create(destination)
	// Opening the Source of file
	sourceFile, _ := os.Open(source)
	// Copying the source file
	io.Copy(fileToCreate, sourceFile)
	return nil
}

func destinationCopy(sourceDirectory, destinationDirectory string, info os.FileInfo) error {
	fmt.Println("Saving... destination copy")
	// Making Directory from parent to subdirectories
	if err := os.MkdirAll(destinationDirectory, os.ModePerm); err != nil {
		return err
	}

	// Reading the directory of our source file/folder
	contents, err := ioutil.ReadDir(sourceDirectory)
	if err != nil {
		return err
	}

	// Folder Saving Data
	if !strings.Contains(destinationDirectory, ".") {
		folder := Folder{}
		uadmin.Get(&folder, "path=?", destinationDirectory)

		folderPath := folder.Path

		if destinationDirectory != folderPath {
			folderName := folder.Name
			currentFolderName := filepath.Base(destinationDirectory)
			currentFolderDirectory := filepath.Dir(destinationDirectory)

			uadmin.Get(&folder, "name=?", currentFolderName)
			if currentFolderName != folderName {
				folder.Name = currentFolderName
				folder.Path = currentFolderDirectory

				folder.Save()
			}
		}
	}

	for _, content := range contents {
		currentSrcDir, currentDestDir := filepath.Join(sourceDirectory, content.Name()), filepath.Join(destinationDirectory, content.Name())
		copy(currentSrcDir, currentDestDir, content)

		// Document
		folderCurrentDirectory := filepath.Base(destinationDirectory)
		documentName := filepath.Base(currentDestDir)

		// currentFolderDirectory
		documentFileDirectory := "/" + currentDestDir

		// Folder Settings Update
		if !strings.Contains(documentFileDirectory, ".") {
			folderTrimParent := strings.TrimRight(strings.TrimLeft(strings.TrimRight(documentFileDirectory, filepath.Base(documentFileDirectory)), "/"), "/")

			folder := Folder{}
			folderParent := filepath.Base(folderTrimParent)
			uadmin.Get(&folder, "name=?", folderParent)

			// Folder Parent ID
			folderParentID := folder.ID
			UINTPARENT_ID := strconv.FormatUint(uint64(folderParentID), 10)

			// Getting Subdirectory
			folderSubdirectory := filepath.Base(documentFileDirectory)

			// Connect to DB and execute SQL
			db := uadmin.GetDB()
			sqlUPDATE := fmt.Sprintf("UPDATE folders SET parent_id = " + UINTPARENT_ID + " WHERE name IN (SELECT name FROM folders WHERE name = '" + folderSubdirectory + "')")
			db.Exec(sqlUPDATE)
		}

		// Document
		document := Document{}
		if strings.Contains(documentName, ".") {
			folder := Folder{}

			uadmin.Get(&folder, "name=?", folderCurrentDirectory)
			uadmin.Get(&document, "name=?", documentName)

			if documentName != document.Name {
				uadmin.Get(&document, "file=?", documentFileDirectory)

				if documentFileDirectory != document.File {
					document.Name = documentName
					document.File = documentFileDirectory
					document.FolderID = folder.ID
					document.CreatedDate = time.Now()
					document.Save()
				}
			}
		}
	}
	return nil
}
