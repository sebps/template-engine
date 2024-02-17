package cmd

// Basic imports
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type RenderTestSuite struct {
	suite.Suite
	wantDir      string
	wantFile     string
	haveDir      string
	haveFile     string
	inDir        string
	inFile       string
	dataCsvFile  string
	dataXlsxFile string
	dataJsonFile string
	cmd          *cobra.Command
}

// before each test
func (suite *RenderTestSuite) SetupTest() {
	test := suite.T().Name()
	fmt.Printf("Test : %s \n", test)

	suite.wantDir = fmt.Sprintf("./tests/%s/want", test)
	suite.wantFile = fmt.Sprintf("./tests/%s/want/want.json", test)
	suite.haveDir = fmt.Sprintf("./tests/%s/have", test)
	suite.haveFile = fmt.Sprintf("./tests/%s/have/have.json", test)
	suite.inDir = fmt.Sprintf("./tests/%s/in", test)
	suite.inFile = fmt.Sprintf("./tests/%s/in/in.json", test)
	suite.dataCsvFile = fmt.Sprintf("./tests/%s/data.csv", test)
	suite.dataXlsxFile = fmt.Sprintf("./tests/%s/data.xlsx", test)
	suite.dataJsonFile = fmt.Sprintf("./tests/%s/data.json", test)
	suite.cmd = rootCmd
}

func (suite *RenderTestSuite) BeforeTest(suiteName, testName string) {
	os.RemoveAll(suite.haveDir)
}

func (suite *RenderTestSuite) AfterTest(suiteName, testName string) {
	os.RemoveAll(suite.haveDir)
}

// All methods that begin with "Test" are run as tests within a suite.

// csv data
func (suite *RenderTestSuite) Test_Render_CsvData_FileInput_SingleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataCsvFile,
		"--in",
		suite.inFile,
		"--out",
		suite.haveFile,
		"--multiple-output",
		"false",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	want, err := os.ReadFile(suite.wantFile)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	have, err := os.ReadFile(suite.haveFile)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_CsvData_FileInput_MultipleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataCsvFile,
		"--in",
		suite.inFile,
		"--out",
		suite.haveFile,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_CsvData_DirInput_SingleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataCsvFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"false",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_CsvData_DirInput_MultipleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataCsvFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_CsvData_DirInput_MultipleOutput_Filtered() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataCsvFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"$[?(@.sku in (record1,record2))]",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

// json object data
func (suite *RenderTestSuite) Test_Render_JsonObjectData_FileInput_SingleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inFile,
		"--out",
		suite.haveFile,
		"--multiple-output",
		"false",
		"--data-filter",
		"$.records",
	})

	suite.cmd.Execute()

	want, err := os.ReadFile(suite.wantFile)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	have, err := os.ReadFile(suite.haveFile)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonObjectData_FileInput_MultipleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inFile,
		"--out",
		suite.haveFile,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"$.records",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonObjectData_DirInput_SingleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"false",
		"--data-filter",
		"$.records",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonObjectData_DirInput_MultipleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"$.records",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonObjectData_DirInput_MultipleOutput_Filtered() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"$.records[?(@.sku in (record1,record2))]",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

// json array data
func (suite *RenderTestSuite) Test_Render_JsonArrayData_FileInput_SingleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inFile,
		"--out",
		suite.haveFile,
		"--multiple-output",
		"false",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	want, err := os.ReadFile(suite.wantFile)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	have, err := os.ReadFile(suite.haveFile)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonArrayData_FileInput_MultipleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inFile,
		"--out",
		suite.haveFile,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonArrayData_DirInput_SingleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"false",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonArrayData_DirInput_MultipleOutput() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

func (suite *RenderTestSuite) Test_Render_JsonArrayData_DirInput_MultipleOutput_Filtered() {
	suite.cmd.SetArgs([]string{
		"render",
		"--data",
		suite.dataJsonFile,
		"--in",
		suite.inDir,
		"--out",
		suite.haveDir,
		"--multiple-output",
		"true",
		"--multiple-output-filename-pattern",
		"{i}",
		"--data-filter",
		"$[?(@.sku in (record1,record2))]",
	})

	suite.cmd.Execute()

	wantDir, err := os.ReadDir(suite.wantDir)
	if err != nil {
		log.Fatalf("unable to read want file: %v", err)
	}

	haveDir, err := os.ReadDir(suite.haveDir)
	if err != nil {
		log.Fatalf("unable to read have file: %v", err)
	}

	want := make([][]byte, len(wantDir))
	for i := 0; i < len(wantDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.wantDir, wantDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		want[i] = fileContent
	}

	have := make([][]byte, len(haveDir))
	for i := 0; i < len(haveDir); i++ {
		fileContent, err := os.ReadFile(filepath.Join(suite.haveDir, haveDir[i].Name()))
		if err != nil {
			log.Fatalf("unable to read want file: %v", err)
		}
		have[i] = fileContent
	}

	assert.Equal(suite.T(), reflect.DeepEqual(want, have), true)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRenderTestSuite(t *testing.T) {
	suite.Run(t, new(RenderTestSuite))
}
