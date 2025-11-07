package collector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemInfoParseProcFile(t *testing.T) {
	// happy path
	m := MemInfo{ProcFileName: "../fixtures/meminfo_partial"}

	err := m.ParseProcFile()
	assert.Nil(t, err)
	assert.Equal(t, map[string]uint64{"memfree": 2906963968, "memtotal": 3981893632}, m.data)

	// invalid file format
	m.ProcFileName = "../fixtures/meminfo_invalid"
	err = m.ParseProcFile()
	assert.EqualError(t, err, "could not parse meminfo. could not convert \"bar\" to uint64 for metric \"foo\". strconv.ParseUint: parsing \"bar\": invalid syntax")

	// error opening the file
	m.ProcFileName = "this_does_not_exist"
	err = m.ParseProcFile()
	assert.EqualError(t, err, "could not open meminfo proc file \"this_does_not_exist\". open this_does_not_exist: no such file or directory")

}

func TestMemInfoparseMemInfo(t *testing.T) {
	// happy path - meminfo file format
	r := strings.NewReader("MemTotal:        3888568 kB\nMemFree:         2838832 kB")
	d, err := parseMemInfo(r)
	assert.Nil(t, err)
	assert.Equal(t, map[string]uint64{"memtotal": 3888568, "memfree": 2838832}, d)

	// incorrect units
	r = strings.NewReader("MemTotal:        foo kB\nMemFree:         bar kB")
	d, err = parseMemInfo(r)
	assert.EqualError(t, err, "could not convert \"foo\" to uint64 for metric \"MemTotal\". strconv.ParseUint: parsing \"foo\": invalid syntax")
	assert.Equal(t, map[string]uint64{}, d)

	// unexpected individual metric format (missing ':')
	r = strings.NewReader("MemTotal        3888568 kB\nMemFree:         2838832 kB")
	d, err = parseMemInfo(r)
	assert.Nil(t, err)
	assert.Equal(t, map[string]uint64{"memfree": 2838832}, d)

	// empty file passed (io.Reader)
	r = strings.NewReader("")
	d, err = parseMemInfo(r)
	assert.Nil(t, err)
	assert.Equal(t, map[string]uint64{}, d)
}

func TestOutputPromMetrics(t *testing.T) {
	assert.True(t, true)
}
