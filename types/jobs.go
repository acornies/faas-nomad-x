package types

import "github.com/hashicorp/nomad/api"

type Jobs interface {
	Register(*api.Job, *api.WriteOptions) (*api.JobRegisterResponse, *api.WriteMeta, error)
	Info(jobID string, q *api.QueryOptions) (*api.Job, *api.QueryMeta, error)
	List(q *api.QueryOptions) ([]*api.JobListStub, *api.QueryMeta, error)
	Deregister(jobID string, purge bool, q *api.WriteOptions) (string, *api.WriteMeta, error)
	Allocations(jobID string, allAllocs bool, q *api.QueryOptions) ([]*api.AllocationListStub, *api.QueryMeta, error)
}
