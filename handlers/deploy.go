package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/acornies/faas-nomad-x/types"
	"github.com/hashicorp/nomad/api"
	"github.com/openfaas/faas/gateway/requests"
)

func MakeDeploy(config *types.ProviderConfig, client types.Jobs) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		// Parse the OpenFaaS create function request and map it
		// to Nomad job service type
		req := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		job := createJob(config.Nomad, req)

		// Use the Nomad API client to register the job
		jrRes, _, err := client.Register(&job, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Print("Failed to register Nomad job from function request", err)
			return
		}

		log.Printf("Registered Nomad job with eval index: %d", jrRes.EvalCreateIndex)
	}
}

// Use NomadConfig to set default job config
// then the mapped values from functionRequest as job overrides
func createJob(config types.NomadConfig, r requests.CreateFunctionRequest) api.Job {
	t := &config.Scheduling.JobType
	priority := &config.Scheduling.Priority
	jobName := fmt.Sprintf("%s-%s", config.Scheduling.JobPrefix, r.Service)

	return api.Job{
		Name:        &jobName,
		Type:        t,
		Datacenters: []string{config.Datacenter}, // TODO: map constraints
		Priority:    priority,
		TaskGroups:  createTaskGroups(config, r),
	}
}

// Used to create the job task groups
func createTaskGroups(config types.NomadConfig, r requests.CreateFunctionRequest) []*api.TaskGroup {

	taskCount := config.Scheduling.Count
	restartDelay, _ := time.ParseDuration(config.Scheduling.RestartDelay)
	restartMode := config.Scheduling.RestartMode
	restartAttempts := config.Scheduling.RestartAttempts
	diskSize := config.Scheduling.DiskSize

	return []*api.TaskGroup{
		&api.TaskGroup{
			Name:   &r.Service,
			Count:  &taskCount,
			Update: &api.UpdateStrategy{},
			RestartPolicy: &api.RestartPolicy{
				Delay:    &restartDelay,
				Mode:     &restartMode,
				Attempts: &restartAttempts,
			},
			EphemeralDisk: &api.EphemeralDisk{
				SizeMB: &diskSize,
			},
		},
	}
}
