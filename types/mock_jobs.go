package types

import (
	"github.com/hashicorp/nomad/api"
	"github.com/stretchr/testify/mock"
)

type MockJobs struct {
	mock.Mock
}

func (m *MockJobs) Register(job *api.Job, options *api.WriteOptions) (
	*api.JobRegisterResponse,
	*api.WriteMeta, error) {

	args := m.Called(job, options)

	var resp *api.JobRegisterResponse
	if r := args.Get(0); r != nil {
		resp = r.(*api.JobRegisterResponse)
	}

	var meta *api.WriteMeta
	if r := args.Get(1); r != nil {
		meta = r.(*api.WriteMeta)
	}

	return resp, meta, args.Error(2)
}

func (m *MockJobs) Allocations(jobID string, allAllocs bool, q *api.QueryOptions) ([]*api.AllocationListStub, *api.QueryMeta, error) {
	args := m.Called(jobID, allAllocs, q)

	var allocs []*api.AllocationListStub
	if a := args.Get(0); a != nil {
		allocs = a.([]*api.AllocationListStub)
	}

	var meta *api.QueryMeta
	if r := args.Get(1); r != nil {
		meta = r.(*api.QueryMeta)
	}

	return allocs, meta, args.Error(2)
}

func (m *MockJobs) Info(jobID string, q *api.QueryOptions) (*api.Job, *api.QueryMeta, error) {
	args := m.Called(jobID, q)

	var job *api.Job
	if j := args.Get(0); j != nil {
		job = j.(*api.Job)
	}

	var meta *api.QueryMeta
	if r := args.Get(1); r != nil {
		meta = r.(*api.QueryMeta)
	}

	return job, meta, args.Error(2)
}

// List returns mock info from the job API
func (m *MockJobs) List(q *api.QueryOptions) ([]*api.JobListStub, *api.QueryMeta, error) {
	args := m.Called(q)

	var jobs []*api.JobListStub
	if j := args.Get(0); j != nil {
		jobs = j.([]*api.JobListStub)
	}

	var meta *api.QueryMeta
	if r := args.Get(1); r != nil {
		meta = r.(*api.QueryMeta)
	}

	return jobs, meta, args.Error(2)
}

func (m *MockJobs) Deregister(jobID string, purge bool, q *api.WriteOptions) (
	string, *api.WriteMeta, error) {

	args := m.Called(jobID, purge, q)

	return "", nil, args.Error(2)
}
