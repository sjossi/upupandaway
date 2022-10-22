package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sjossi/upupandaway/unpacker"
)

func main() {
	// Currently runs a hacky "extract all" for the provided directory.
	//
	// Cleaner steps are in place in the library, but I couldn't put in the
	// time to create a proper CLI yet, so it is what it is now.

	log.Print("[+] Welcome to .up .up and away")

	if len(os.Args) < 2 {
		log.Printf("[!] Usage: %s %s", filepath.Base(os.Args[0]), "<path>")
		os.Exit(1)
	}

	upDir := os.Args[1]

	mainInstructions := filepath.Join(upDir, "main_instructions.ini")
	iniTree := unpacker.ParseIniTree(mainInstructions)

	// Extracts to extracted_ folder with timestamp within the supplied root folder
	// TODO: keep track of it in a configuration object
	toBase := filepath.Join(upDir, "./extracted_"+time.Now().Format("20060102150405")+"")

	// TODO: add logging configuration to configuration object

	log.Printf("[+] Extracting to %s", toBase)

	for _, ini := range iniTree {
		if strings.HasPrefix(ini.Filename, "files.ini") || strings.HasPrefix(ini.Filename, "execute.ini") {
			log.Printf("[+] Extracting %s", ini.Filename)
			unpacker.ExtractFiles(ini, toBase)
		}
	}
}
