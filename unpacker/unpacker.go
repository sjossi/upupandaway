/*
unpacker unpacks a folder of interest with lots of subfolders, according to
it's instructions in main_instructions.ini.
*/
package unpacker

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	// "errors"
	"encoding/csv"
	"strconv"

	// "regexp"
	ini "github.com/ochinchina/go-ini"
)

func check(err error) {
	if err != nil {
		log.Printf("Error: %q", err)
	}
}

func ExtractFiles(ini *Ini, toBase string) {
	// ExtractFiles extracts all the files according to the files.ini
	// instructions
	//
	// Relative unpacking/copying without implicit directory will be placed in
	// ./tmp, since the updater also runs in /tmp. This creates the closest
	// representation to the actual file system.

	// Generate a new folder every time to avoid conflicts
	os.Mkdir(toBase, 0755)

	for _, instruction := range ini.Instructions.Instructions {
		switch instruction.InstructionStep {
		case Copy:
			// args: from, to
			var folder string
			if strings.HasPrefix(ini.Filename, "files.ini") {
				// TODO: log.Printf("folder: %s rootdir %s arg0 %s",
				//           ini.Folder, ini.RootDir, instruction.Arguments[0])
				folder = ""
			} else {
				folder = ini.Folder
			}
			from := filepath.Join(ini.RootDir, folder, instruction.Arguments[0])

			// If a relative path is defined, save in /tmp.  This is because
			// the updater scripts are ran in /tmp
			var to string
			if filepath.Dir(instruction.Arguments[1]) == "." {
				to = filepath.Join(toBase, "/tmp/", instruction.Arguments[1])
			} else {
				to = filepath.Join(toBase, instruction.Arguments[1])
			}

			err := os.MkdirAll(filepath.Dir(to), 0755)
			if err != nil && !os.IsExist(err) {
				// Only fail if it's an error other than existing directory
				log.Fatal(err)
			}

			var reader io.Reader

			reader, err = os.Open(from)
			// handle implicit .gz
			if os.IsNotExist(err) {
				reader, err = os.Open(from + ".gz")
				check(err)
				reader, err = gzip.NewReader(reader)
				check(err)
			}

			writer, err := os.Create(to)
			check(err)

			_, err = io.Copy(writer, reader)
			// TODO: check(err)
			// TODO: log.Printf("copied %s to %s, bytes written: %d"
			//         from, to, bytesWritten)
		}
	}
}

func SimulateSteps(ini *Ini) []string {
	// SimulateExecute simulates an execute.ini instructions file
	//
	// Currently lists all the target files for Copy. In the future this should
	// become a representation of the file system that represents the target
	// state after all operations have executed (potentially missing details
	// because shell scripts were not executed)

	base := "/tmp"
	files := make([]string, 0)

	if !strings.HasPrefix(ini.Filename, "execute.ini") && !strings.HasPrefix(ini.Filename, "files.ini") {
		log.Printf("Not an execute.ini or files.ini: %s", ini.Filename)
		return files
	}

	for _, instruction := range ini.Instructions.Instructions {
		switch instruction.InstructionStep {
		case Copy:
			// args: from, to
			filename := instruction.Arguments[1]
			if filepath.Dir(filename) == "." {
				filename = filepath.Join(base, filename)
			}
			files = append(files, filename)
		case Create:
			// args: filepath
			filename := instruction.Arguments[0]
			if filepath.Dir(filename) == "." {
				filename = filepath.Join(base, filename)
			}
			files = append(files, filename)
		}
	}

	return files
}

func ParseIniTree(filename string) []*Ini {
	dir := filepath.Dir(filename)

	main := ParseMainIni(ini.Load(filename))
	main.RootDir = dir
	main.Filename = filename

	tree := make([]*Ini, main.Instructions.Count)

	tree[0] = main

	for i := 1; i < main.Instructions.Count; i++ {
		instruction := main.Instructions.Instructions[i]

		candidate := filepath.Join(dir, instruction.Arguments[0], instruction.Arguments[1])

		// Check if it's a normal file or compressed
		var reader io.Reader
		file, err := os.Open(candidate)
		if !os.IsNotExist(err) {
			reader = file
		} else {
			candidate += ".gz"

			file, err = os.Open(candidate)
			if err != nil {
				log.Printf("[!] could not open file %s: %q", candidate, err)
				continue
			}
			reader, err = gzip.NewReader(file)
			if err != nil {
				log.Printf("[!] could not decompress file")
				continue
			}
		}

		subini := ini.Load(reader)
		subini_ini := ParseSubIni(subini)

		subini_ini.RootDir = dir
		subini_ini.Folder = instruction.Arguments[0]
		subini_ini.Filename = instruction.Arguments[1]

		tree[i] = subini_ini
	}

	return tree
}

func ParseMainIni(in *ini.Ini) *Ini {
	// Parses a main.ini file
	//
	// The main.ini points to other .ini files.
	// TODO: check steps and other verifications

	// TODO: debug: log.Print("ParseMainIni()")

	ini := new(Ini)

	settings := ParseSettings(in)
	ini.Settings = settings

	instructions := ParseInstructions(in, "Instructions", true)
	ini.Instructions = instructions

	instructions_ext := ParseInstructions(in, "Instructions_Ext", true)
	ini.Instructions_Ext = instructions_ext

	datastorage := ParseDataStorage(in)
	ini.DataStorage = datastorage

	return ini
}

func ParseSubIni(in *ini.Ini) *Ini {
	// Parses sub ini files
	//
	// Parses ini files that are called by main.ini.

	// TODO: debug log.Print("ParseExecuteIni")

	ini := new(Ini)

	settings := ParseSettings(in)
	ini.Settings = settings

	instructions := ParseInstructions(in, "Instructions", false)
	ini.Instructions = instructions

	return ini
}

func ParseSettings(in *ini.Ini) Settings {
	// Parses a settings.ini
	//

	// TODO: debug log.Print("ParseInstructions()")

	section := "Settings"

	res := Settings{}

	res.Packageid = in.GetInt64WithDefault(section, "PackageID", 0)
	res.TotalStepsCount = in.GetInt64WithDefault(section, "TotalStepsCount", 0)

	var compressiontype CompressionType
	compression := in.GetValueWithDefault(section, "CompressionType", "")
	switch compression {
	case "GZIP":
		compressiontype = GZIP
	default:
		compressiontype = UNDEFINED
	}
	res.CompressionType = compressiontype

	return res
}

func ParseDataStorage(in *ini.Ini) DataStorage {
	// Parses a data_storage.ini file

	// log.Print("ParseInstructions()")

	section := "DataStorage"

	res := DataStorage{}

	res.Count = in.GetIntWithDefault(section, "Count", 0)
	res.UPType = in.GetValueWithDefault(section, "UPType", "")
	res.SubUPType = in.GetValueWithDefault(section, "SubUPType", "")
	res.ReTransmit = in.GetValueWithDefault(section, "ReTransmit", "")
	res.NewPackage = in.GetValueWithDefault(section, "NewPackage", "")

	return res
}

func ParseInstructions(in *ini.Ini, section string, has_steps bool) Instructions {
	// Parses an instructions.ini file

	// log.Print("ParseInstructions()")

	ins := Instructions{}

	// section := "Instructions"

	ins.Count = in.GetIntWithDefault(section, "Count", 0)

	instructions := make([]Instruction, 0)

	// TODO: is this not ini parsed?!
	for i := 1; i <= ins.Count; i++ {
		line, err := in.GetValue(section, strconv.FormatInt(int64(i), 10))
		if err != nil {
			log.Printf("Could not get step %d: %q", i, err)
		}

		r := csv.NewReader(strings.NewReader(line))
		r.FieldsPerRecord = -1
		r.LazyQuotes = true
		r.TrimLeadingSpace = true

		tokens, err := r.Read()
		if err != nil {
			log.Printf("Could not read line %s: %q", line, err)
		}

		var step InstructionStep

		switch tokens[0] {
		case "Execute":
			step = Execute
		case "ImageUpdate":
			step = ImageUpdate
		case "FileUpdate":
			step = FileUpdate
		case "BreakPoint":
			step = BreakPoint
		case "Copy":
			step = Copy
		case "Remove":
			step = Remove
		case "Create":
			step = Create
		case "RemoveFolderContent":
			step = RemoveFolderContent
		default:
			log.Printf("Could not parse instructionstep: %s", tokens[0])
		}

		var args []string
		var steps int

		if has_steps {
			args = tokens[1 : len(tokens)-1]
			steps, err = strconv.Atoi(tokens[len(tokens)-1])
			if err != nil {
				log.Printf("could not parse steps: %q", err)
			}
		} else {
			args = tokens[1:]
			steps = 0
		}

		instruction := Instruction{
			StepNo:          i,
			InstructionStep: step,
			Arguments:       args,
			Steps:           steps,
		}

		instructions = append(instructions, instruction)
	}

	ins.Instructions = instructions

	return ins
}
