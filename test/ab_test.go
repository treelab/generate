package test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"github.com/treelab/generate"
)

func Output(t *testing.T, filename string) {
	inputFiles := []string{filename}
	schemas, err := generate.ReadInputFiles(inputFiles, false)
	assert.Nil(t, err)

	g := generate.New(schemas...)

	err = g.CreateTypes()
	assert.Nil(t, err)

	codeBuf := new(bytes.Buffer)

	generate.Output(codeBuf, g, "test")
	cupaloy.SnapshotT(t, codeBuf.String())
}

func TestAbandoned(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "abandoned")
	Output(t, path)
}

func TestAdditionalProperties(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "additionalProperties")
	Output(t, path)
}

func TestAdditionalProperties2(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "additionalProperties2")
	Output(t, path)
}

func TestAnonarrayitems(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "anonarrayitems")
	Output(t, path)
}

func TestAdditionalPropertiesMarshal(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "additionalPropertiesMarshal")
	Output(t, path)
}

func TestAprefnoprop(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "aprefnoprop")
	Output(t, path)
}

func TestArray(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "array")
	Output(t, path)
}

func TestArrayofref(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "arrayofref")
	Output(t, path)
}

func TestCustomer(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "customer")
	Output(t, path)
}

func TestExample1(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "example1")
	Output(t, path)
}

func TestExample1a(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "example1a")
	Output(t, path)
}

func TestIssue6(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "issue6")
	Output(t, path)
}

func TestIssue14(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "issue14")
	Output(t, path)
}

func TestIssue35(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "issue35")
	Output(t, path)
}

func TestIssue39(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "issue39")
	Output(t, path)
}

func TestMultiple(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "multiple")
	Output(t, path)
}

func TestNestedarrayofref(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "nestedarrayofref")
	Output(t, path)
}

func TestRecursion(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "recursion")
	Output(t, path)
}

func TestSchemaid(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "schemaid")
	Output(t, path)
}

func TestSimple(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "simple")
	Output(t, path)
}

func TestTest(t *testing.T) {
	path := fmt.Sprintf("./_mock/%s.json", "test")
	Output(t, path)
}
