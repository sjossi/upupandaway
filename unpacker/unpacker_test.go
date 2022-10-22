package unpacker

import (
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	ini "github.com/ochinchina/go-ini"
)

const MAIN_INSTRUCTIONS = "./test/main_instructions.ini"
const EXECUTE_INSTRUCTIONS = "./test/execute.ini.gz"

func TestExtractFiles(t *testing.T) {
	// bit silly, turn off later
	log.Print("TestExtractFiles")

	testsDir, exists := os.LookupEnv("TESTS_DIR")

	if !exists {
		log.Panic("Set TEST_DIR to a valid update folder")
	}

	mainInstructions := filepath.Join(testsDir, "main_instructions.ini")

	got := ParseIniTree(mainInstructions)

	toBase := "./extracted_" + time.Now().Format("20060102150405") + ""

	for _, ini := range got {
		if strings.HasPrefix(ini.Filename, "files.ini") || strings.HasPrefix(ini.Filename, "execute.ini") {
			ExtractFiles(ini, toBase)
		}
	}
}

func TestSimulateFullTree(t *testing.T) {
	// bit silly, turn off later
	log.Print("TestSimulateFullTree")

	testsDir, exists := os.LookupEnv("TESTS_DIR")

	if !exists {
		log.Panic("Set TEST_DIR to a valid update folder")
	}

	mainInstructions := filepath.Join(testsDir, "main_instructions.ini")

	got := ParseIniTree(mainInstructions)

	files := make([]string, 0)

	for _, v := range got {
		files = append(files, SimulateSteps(v)...)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	// pretty_files, err := json.MarshalIndent(files, "", "\t")
	// check(err)
	// log.Printf("%s", pretty_files)

	outfile, err := os.Create("file_list.txt")
	check(err)
	defer outfile.Close()

	for _, file := range files {
		outfile.WriteString(file + "\n")
	}

	// log.Printf("tree: %#v", got)
}

func TestSimulateExecute(t *testing.T) {
	file, err := os.Open(EXECUTE_INSTRUCTIONS)
	if err != nil {
		log.Panicf("could not open file: %q", err)
	}

	reader, err := gzip.NewReader(file)
	if err != nil {
		log.Panicf("could not decompress file: %q", err)
	}

	ini := ini.Load(reader)

	in := ParseSubIni(ini)
	got := SimulateSteps(in)

	want := []string{"/tmp/compactwnn_dictionary.sh"}

	if !reflect.DeepEqual(got, want) {
		t.Error("SimulateExecute: do not match")
		log.Printf("got: %#v\nwant: %#v", got, want)
	}
}

func TestParseIniTree(t *testing.T) {
	log.Print("TestParseIniTree")

	testsDir, exists := os.LookupEnv("TESTS_DIR")

	if !exists {
		log.Panic("Set TEST_DIR to a valid update folder")
	}

	mainInstructions := filepath.Join(testsDir, "main_instructions.ini")

	got := ParseIniTree(mainInstructions)

	log.Printf("tree: %#v", got)
}

func TestParseMainIni(t *testing.T) {
	log.Print("TestParseMainIni")

	have := ini.Load(MAIN_INSTRUCTIONS)

	got := ParseMainIni(have)
	want := &Ini{
		Instructions: Instructions{
			Count: 21,
			Instructions: []Instruction{
				{StepNo: 1, InstructionStep: 0, Arguments: []string{"cleandatapersist", "execute.ini"}, Steps: 4},
				{StepNo: 2, InstructionStep: 0, Arguments: []string{"bootstrap", "execute.ini"}, Steps: 7},
				{StepNo: 3, InstructionStep: 1, Arguments: []string{"ibc2", "binary.ini"}, Steps: 2},
				{StepNo: 4, InstructionStep: 1, Arguments: []string{"fail-safe", "binary.ini"}, Steps: 2},
				{StepNo: 5, InstructionStep: 0, Arguments: []string{"checksumoption", "execute.ini"}, Steps: 5},
				{StepNo: 6, InstructionStep: 1, Arguments: []string{"ibc1", "binary.ini"}, Steps: 3},
				{StepNo: 7, InstructionStep: 0, Arguments: []string{"linux1", "execute.ini"}, Steps: 6},
				{StepNo: 8, InstructionStep: 0, Arguments: []string{"getoldflavor", "execute.ini"}, Steps: 4},
				{StepNo: 9, InstructionStep: 0, Arguments: []string{"rootfs1upd", "execute.ini"}, Steps: 8},
				{StepNo: 10, InstructionStep: 0, Arguments: []string{"getnewflavor", "execute.ini"}, Steps: 4},
				{StepNo: 11, InstructionStep: 0, Arguments: []string{"passwdupdate", "execute.ini"}, Steps: 12},
				{StepNo: 12, InstructionStep: 0, Arguments: []string{"gps", "execute.ini"}, Steps: 6},
				{StepNo: 13, InstructionStep: 2, Arguments: []string{"resources", "files.ini"}, Steps: 801},
				{StepNo: 14, InstructionStep: 0, Arguments: []string{"usersettingsbackup", "execute.ini"}, Steps: 4},
				{StepNo: 15, InstructionStep: 0, Arguments: []string{"usersettingsrestore", "execute.ini"}, Steps: 4},
				{StepNo: 16, InstructionStep: 0, Arguments: []string{"usersettingscleanup", "execute.ini"}, Steps: 4},
				{StepNo: 17, InstructionStep: 0, Arguments: []string{"preloaddata", "execute.ini"}, Steps: 8},
				{StepNo: 18, InstructionStep: 0, Arguments: []string{"compactwnn", "execute.ini"}, Steps: 5},
				{StepNo: 19, InstructionStep: 0, Arguments: []string{"neutralizeid7", "execute.ini"}, Steps: 4},
				{StepNo: 20, InstructionStep: 0, Arguments: []string{"systemupdateid", "execute.ini"}, Steps: 2},
				{StepNo: 21, InstructionStep: 0, Arguments: []string{"vip", "execute.ini"}, Steps: 7}}},
		Settings: Settings{Packageid: 1587449549, CompressionType: 1, TotalStepsCount: 902},
		Instructions_Ext: Instructions{
			Count: 25,
			Instructions: []Instruction{
				{StepNo: 1, InstructionStep: 0, Arguments: []string{"cleandatapersist", "execute.ini"}, Steps: 4},
				{StepNo: 2, InstructionStep: 0, Arguments: []string{"bootstrap", "execute.ini"}, Steps: 7},
				{StepNo: 3, InstructionStep: 3, Arguments: []string{"failsafeos", "Start"}, Steps: 0},
				{StepNo: 4, InstructionStep: 1, Arguments: []string{"ibc2", "binary.ini"}, Steps: 2},
				{StepNo: 5, InstructionStep: 1, Arguments: []string{"fail-safe", "binary.ini"}, Steps: 2},
				{StepNo: 6, InstructionStep: 0, Arguments: []string{"checksumoption", "execute.ini"}, Steps: 5},
				{StepNo: 7, InstructionStep: 3, Arguments: []string{"failsafeos", "End"}, Steps: 0},
				{StepNo: 8, InstructionStep: 3, Arguments: []string{"reinstall", "Start"}, Steps: 0},
				{StepNo: 9, InstructionStep: 1, Arguments: []string{"ibc1", "binary.ini"}, Steps: 3},
				{StepNo: 10, InstructionStep: 0, Arguments: []string{"linux1", "execute.ini"}, Steps: 6},
				{StepNo: 11, InstructionStep: 0, Arguments: []string{"getoldflavor", "execute.ini"}, Steps: 4},
				{StepNo: 12, InstructionStep: 0, Arguments: []string{"rootfs1upd", "execute.ini"}, Steps: 8},
				{StepNo: 13, InstructionStep: 0, Arguments: []string{"getnewflavor", "execute.ini"}, Steps: 4},
				{StepNo: 14, InstructionStep: 0, Arguments: []string{"passwdupdate", "execute.ini"}, Steps: 12},
				{StepNo: 15, InstructionStep: 0, Arguments: []string{"gps", "execute.ini"}, Steps: 6},
				{StepNo: 16, InstructionStep: 2, Arguments: []string{"resources", "files.ini"}, Steps: 801},
				{StepNo: 17, InstructionStep: 0, Arguments: []string{"usersettingsbackup", "execute.ini"}, Steps: 4},
				{StepNo: 18, InstructionStep: 0, Arguments: []string{"usersettingsrestore", "execute.ini"}, Steps: 4},
				{StepNo: 19, InstructionStep: 0, Arguments: []string{"usersettingscleanup", "execute.ini"}, Steps: 4},
				{StepNo: 20, InstructionStep: 0, Arguments: []string{"preloaddata", "execute.ini"}, Steps: 8},
				{StepNo: 21, InstructionStep: 0, Arguments: []string{"compactwnn", "execute.ini"}, Steps: 5},
				{StepNo: 22, InstructionStep: 0, Arguments: []string{"neutralizeid7", "execute.ini"}, Steps: 4},
				{StepNo: 23, InstructionStep: 0, Arguments: []string{"systemupdateid", "execute.ini"}, Steps: 2},
				{StepNo: 24, InstructionStep: 0, Arguments: []string{"vip", "execute.ini"}, Steps: 7},
				{StepNo: 25, InstructionStep: 3, Arguments: []string{"reinstall", "End"}, Steps: 0}}},
		DataStorage: DataStorage{
			Count:      4,
			UPType:     "\"Reinstall\"",
			SubUPType:  "\"Mass\"",
			ReTransmit: "\"1\"",
			NewPackage: "\"1\""},
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("ParseMainIni: do not match")
		log.Printf("got: %#v\nwant: %#v", got, want)
	}
}

func TestParseSettingsIni(t *testing.T) {
	log.Print("TestParseSettingsIni()")

	have := ini.Load(MAIN_INSTRUCTIONS)
	got := ParseSettings(have)
	want := Settings{
		Packageid:       1587449549,
		CompressionType: 1,
		TotalStepsCount: 902,
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("ParseMainIni: do not match")
		log.Printf("got: %#v\nwant: %#v", got, want)
	}
}

func TestParseDataStorage(t *testing.T) {
	in := ini.Load(MAIN_INSTRUCTIONS)

	got := ParseDataStorage(in)
	want := DataStorage{
		Count:      4,
		UPType:     "\"Reinstall\"",
		SubUPType:  "\"Mass\"",
		ReTransmit: "\"1\"",
		NewPackage: "\"1\""}

	if !reflect.DeepEqual(got, want) {
		t.Error("ParseMainIni: do not match")
		log.Printf("got: %#v\nwant: %#v", got, want)
	}

}

func TestParseExecuteIni(t *testing.T) {
	file, err := os.Open(EXECUTE_INSTRUCTIONS)
	check(err)
	gz, err := gzip.NewReader(file)
	check(err)

	in := ini.Load(gz)

	got := ParseSubIni(in)

	want := &Ini{
		Instructions: Instructions{
			Count: 5,
			Instructions: []Instruction{
				{StepNo: 1, InstructionStep: 4,
					Arguments: []string{"e0000000001.dat", "compactwnn_dictionary.sh"}, Steps: 0},
				{StepNo: 2, InstructionStep: 0,
					Arguments: []string{"echo ========== Copy compactwnn dictionary to data_persist =========="}, Steps: 0},
				{StepNo: 3, InstructionStep: 0,
					Arguments: []string{"/tmp/compactwnn_dictionary.sh"}, Steps: 0},
				{StepNo: 4, InstructionStep: 0,
					Arguments: []string{"echo ========== Finish executing Custom Package =========="}, Steps: 0},
				{StepNo: 5, InstructionStep: 5,
					Arguments: []string{"compactwnn_dictionary.sh"}, Steps: 0}},
		},
		Settings: Settings{
			Packageid:       0,
			CompressionType: 0,
			TotalStepsCount: 0,
		},
		Instructions_Ext: Instructions{Count: 0, Instructions: nil},
		DataStorage: DataStorage{
			Count:      0,
			UPType:     "",
			SubUPType:  "",
			ReTransmit: "",
			NewPackage: "",
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("ParseExecuteIni: do not match")
		log.Printf("got: %#v\nwant: %#v", got, want)
	}
}

func TestParseInstructions(t *testing.T) {
	in := ini.Load(MAIN_INSTRUCTIONS)

	got := ParseInstructions(in, "Instructions", true)

	want := Instructions{
		Count: 21,
		Instructions: []Instruction{
			{StepNo: 1, InstructionStep: 0, Arguments: []string{"cleandatapersist", "execute.ini"}, Steps: 4},
			{StepNo: 2, InstructionStep: 0, Arguments: []string{"bootstrap", "execute.ini"}, Steps: 7},
			{StepNo: 3, InstructionStep: 1, Arguments: []string{"ibc2", "binary.ini"}, Steps: 2},
			{StepNo: 4, InstructionStep: 1, Arguments: []string{"fail-safe", "binary.ini"}, Steps: 2},
			{StepNo: 5, InstructionStep: 0, Arguments: []string{"checksumoption", "execute.ini"}, Steps: 5},
			{StepNo: 6, InstructionStep: 1, Arguments: []string{"ibc1", "binary.ini"}, Steps: 3},
			{StepNo: 7, InstructionStep: 0, Arguments: []string{"linux1", "execute.ini"}, Steps: 6},
			{StepNo: 8, InstructionStep: 0, Arguments: []string{"getoldflavor", "execute.ini"}, Steps: 4},
			{StepNo: 9, InstructionStep: 0, Arguments: []string{"rootfs1upd", "execute.ini"}, Steps: 8},
			{StepNo: 10, InstructionStep: 0, Arguments: []string{"getnewflavor", "execute.ini"}, Steps: 4},
			{StepNo: 11, InstructionStep: 0, Arguments: []string{"passwdupdate", "execute.ini"}, Steps: 12},
			{StepNo: 12, InstructionStep: 0, Arguments: []string{"gps", "execute.ini"}, Steps: 6},
			{StepNo: 13, InstructionStep: 2, Arguments: []string{"resources", "files.ini"}, Steps: 801},
			{StepNo: 14, InstructionStep: 0, Arguments: []string{"usersettingsbackup", "execute.ini"}, Steps: 4},
			{StepNo: 15, InstructionStep: 0, Arguments: []string{"usersettingsrestore", "execute.ini"}, Steps: 4},
			{StepNo: 16, InstructionStep: 0, Arguments: []string{"usersettingscleanup", "execute.ini"}, Steps: 4},
			{StepNo: 17, InstructionStep: 0, Arguments: []string{"preloaddata", "execute.ini"}, Steps: 8},
			{StepNo: 18, InstructionStep: 0, Arguments: []string{"compactwnn", "execute.ini"}, Steps: 5},
			{StepNo: 19, InstructionStep: 0, Arguments: []string{"neutralizeid7", "execute.ini"}, Steps: 4},
			{StepNo: 20, InstructionStep: 0, Arguments: []string{"systemupdateid", "execute.ini"}, Steps: 2},
			{StepNo: 21, InstructionStep: 0, Arguments: []string{"vip", "execute.ini"}, Steps: 7}}}

	if !reflect.DeepEqual(got, want) {
		t.Error("ParseMainIni: do not match")
		log.Printf("got: %#v", got)
		log.Printf("want: %#v", want)
	}

	got = ParseInstructions(in, "Instructions_Ext", true)
	want = Instructions{
		Count: 25,
		Instructions: []Instruction{
			{StepNo: 1, InstructionStep: 0, Arguments: []string{"cleandatapersist", "execute.ini"}, Steps: 4},
			{StepNo: 2, InstructionStep: 0, Arguments: []string{"bootstrap", "execute.ini"}, Steps: 7},
			{StepNo: 3, InstructionStep: 3, Arguments: []string{"failsafeos", "Start"}, Steps: 0},
			{StepNo: 4, InstructionStep: 1, Arguments: []string{"ibc2", "binary.ini"}, Steps: 2},
			{StepNo: 5, InstructionStep: 1, Arguments: []string{"fail-safe", "binary.ini"}, Steps: 2},
			{StepNo: 6, InstructionStep: 0, Arguments: []string{"checksumoption", "execute.ini"}, Steps: 5},
			{StepNo: 7, InstructionStep: 3, Arguments: []string{"failsafeos", "End"}, Steps: 0},
			{StepNo: 8, InstructionStep: 3, Arguments: []string{"reinstall", "Start"}, Steps: 0},
			{StepNo: 9, InstructionStep: 1, Arguments: []string{"ibc1", "binary.ini"}, Steps: 3},
			{StepNo: 10, InstructionStep: 0, Arguments: []string{"linux1", "execute.ini"}, Steps: 6},
			{StepNo: 11, InstructionStep: 0, Arguments: []string{"getoldflavor", "execute.ini"}, Steps: 4},
			{StepNo: 12, InstructionStep: 0, Arguments: []string{"rootfs1upd", "execute.ini"}, Steps: 8},
			{StepNo: 13, InstructionStep: 0, Arguments: []string{"getnewflavor", "execute.ini"}, Steps: 4},
			{StepNo: 14, InstructionStep: 0, Arguments: []string{"passwdupdate", "execute.ini"}, Steps: 12},
			{StepNo: 15, InstructionStep: 0, Arguments: []string{"gps", "execute.ini"}, Steps: 6},
			{StepNo: 16, InstructionStep: 2, Arguments: []string{"resources", "files.ini"}, Steps: 801},
			{StepNo: 17, InstructionStep: 0, Arguments: []string{"usersettingsbackup", "execute.ini"}, Steps: 4},
			{StepNo: 18, InstructionStep: 0, Arguments: []string{"usersettingsrestore", "execute.ini"}, Steps: 4},
			{StepNo: 19, InstructionStep: 0, Arguments: []string{"usersettingscleanup", "execute.ini"}, Steps: 4},
			{StepNo: 20, InstructionStep: 0, Arguments: []string{"preloaddata", "execute.ini"}, Steps: 8},
			{StepNo: 21, InstructionStep: 0, Arguments: []string{"compactwnn", "execute.ini"}, Steps: 5},
			{StepNo: 22, InstructionStep: 0, Arguments: []string{"neutralizeid7", "execute.ini"}, Steps: 4},
			{StepNo: 23, InstructionStep: 0, Arguments: []string{"systemupdateid", "execute.ini"}, Steps: 2},
			{StepNo: 24, InstructionStep: 0, Arguments: []string{"vip", "execute.ini"}, Steps: 7},
			{StepNo: 25, InstructionStep: 3, Arguments: []string{"reinstall", "End"}, Steps: 0}}}

	if !reflect.DeepEqual(got, want) {
		t.Error("ParseMainIni: do not match")
		log.Printf("got: %#v", got)
		log.Printf("want: %#v", want)
	}
}

// func TestParseDataStorage(t *testing.T) {
//     log.Print("test parse datastorage")
//
//     want := DataStorage {
//         Count: 4,
//         UPType: "Reinstall",
//         SubUPType: "Mass",
//         ReTransmit: "1",
//         NewPackage: "1",
//     }
//
//     have := []string {
//         "[DataStorage]",
//         "Count = 4",
//         "UPType = \"Reinstall\"",
//         "SubUPType = \"Mass\"",
//         "ReTransmit = \"1\"",
//         "NewPackage = \"1\"",
//     }
//
//     got, err := ParseDataStorage(have)
//
//     if err != nil {
//         log.Printf("Error parsing DataStorage: %q", err)
//     }
//
//     if !reflect.DeepEqual(want, got) {
//         t.Errorf("Mismatch. got: %+v, want: %+v", got, want)
//     }
// }
//
// func compareSlices[T comparable](t *testing.T, got []T, want []T) bool {
//     log.Print("compareSlices()")
//
//     if len(got) != len(want) {
//         t.Errorf("mismatched length! got %d\nwant %d", len(got), len(want))
//         return false
//     }
//
//     if !reflect.DeepEqual(got, want) {
//         t.Error("not equal")
//
//         for i, v := range got {
//             if v != want[i] {
//                 t.Errorf("mismatched contents! got %v\nwant %v", v, want[i])
//             }
//         }
//         return false
//     }
//
//     return true
// }
