//+build integration

package integration

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kachan1208/something/processor/dao/data"
	"github.com/kachan1208/something/processor/model"
	"github.com/kachan1208/something/processor/processor"
	"github.com/kachan1208/something/processor/stream/buffer"
	"github.com/kachan1208/something/processor/stream/decoder"
	"github.com/kachan1208/something/processor/tests/fixtures"
)

type DataStorageMock struct{}

func (s *DataStorageMock) Load(data string) (io.ReadCloser, error) {
	return ioutil.NopCloser(strings.NewReader(data)), nil
}

func (s *DataStorageMock) Close(file io.ReadCloser) error {
	return file.Close()
}

type StorageClientMock struct{}

func (s *StorageClientMock) StoreBatchOfPorts(data []*model.Port) (uint64, error) {
	return uint64(len(data)), nil
}

type LazyStorageClientMock struct{}

func (s *LazyStorageClientMock) StoreBatchOfPorts(data []*model.Port) (uint64, error) {
	time.Sleep(time.Millisecond * 100)

	return uint64(len(data)), nil
}

func TestProcessorProcessGetJobID(t *testing.T) {
	proc := processor.NewPortProcessor(
		decoder.NewJSONStreamDecoder(),
		&DataStorageMock{},
		buffer.NewInterfaceStream(2),
		&StorageClientMock{})

	jobID, err := proc.Process("does_it_really_matter")

	assert.NoError(t, err)
	assert.NotNil(t, jobID)
}

func TestProcessorProcessNoFile(t *testing.T) {
	proc := processor.NewPortProcessor(
		decoder.NewJSONStreamDecoder(),
		data.NewLocalStorage("."),
		buffer.NewInterfaceStream(2),
		&StorageClientMock{})

	_, err := proc.Process("no_file")

	assert.Error(t, err)
}

func TestProcessorGetStatusProcessSmall(t *testing.T) {
	proc := processor.NewPortProcessor(
		decoder.NewJSONStreamDecoder(),
		&DataStorageMock{},
		buffer.NewInterfaceStream(2),
		&StorageClientMock{})

	jobID, err := proc.Process(fixtures.TestJSONDataSmall.Data)
	assert.NoError(t, err)
	assert.NotNil(t, jobID)

	time.Sleep(time.Millisecond * 100)

	status, num, err := proc.GetStatus(jobID)
	assert.NoError(t, err)
	assert.Equal(t, processor.StatusDone, status)
	assert.Equal(t, fixtures.TestJSONDataSmall.Len, num)
}

func TestProcessorProcessFile(t *testing.T) {
	proc := processor.NewPortProcessor(
		decoder.NewJSONStreamDecoder(),
		data.NewLocalStorage("../fixtures"),
		buffer.NewInterfaceStream(1024),
		&StorageClientMock{})

	jobID, err := proc.Process("ports.json")
	assert.NoError(t, err)
	assert.NotNil(t, jobID)

	time.Sleep(time.Second)

	status, num, err := proc.GetStatus(jobID)
	assert.NoError(t, err)
	assert.Equal(t, processor.StatusDone, status)
	assert.Equal(t, uint64(1632), num)
}

func TestProcessorGetStatusOneJobWrongJobID(t *testing.T) {
	proc := processor.NewPortProcessor(
		decoder.NewJSONStreamDecoder(),
		&DataStorageMock{},
		buffer.NewInterfaceStream(1),
		&StorageClientMock{})

	jobID, err := proc.Process("does_it_really_matter")
	assert.NoError(t, err)
	assert.NotNil(t, jobID)

	status, num, err := proc.GetStatus("job_id")
	assert.Error(t, err)
	assert.Zero(t, status, num)
}

func TestProcessorGetStatusNoJobs(t *testing.T) {
	proc := processor.NewPortProcessor(
		decoder.NewJSONStreamDecoder(),
		&DataStorageMock{},
		buffer.NewInterfaceStream(1),
		&StorageClientMock{})

	status, num, err := proc.GetStatus("job_id")
	assert.Error(t, err)
	assert.Zero(t, status, num)
}

func TestProcessorProcessQueueMultipleJobs(t *testing.T) {
	proc := processor.NewPortProcessor(
		decoder.NewJSONStreamDecoder(),
		data.NewLocalStorage("../fixtures"),
		buffer.NewInterfaceStream(1),
		&LazyStorageClientMock{})

	jobID, err := proc.Process("ports.json")
	time.Sleep(time.Millisecond * 50)

	firstStatus, _, err := proc.GetStatus(jobID)
	assert.NoError(t, err)
	assert.Equal(t, processor.StatusRunning, firstStatus)

	jobID2, err := proc.Process("ports.json")
	time.Sleep(time.Millisecond * 50)

	firstStatus, _, err = proc.GetStatus(jobID)
	assert.NoError(t, err)
	assert.Equal(t, processor.StatusRunning, firstStatus)

	secondStatus, _, err := proc.GetStatus(jobID2)
	assert.NoError(t, err)
	assert.Equal(t, processor.StatusQueued, secondStatus)
}

//Test stream error
//Test save error
