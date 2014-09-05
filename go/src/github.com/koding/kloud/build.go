package kloud

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/koding/kloud/eventer"
	"github.com/koding/kloud/machinestate"
	"github.com/koding/kloud/protocol"

	"github.com/koding/kite"
)

type Build struct {
	*Kloud
	deployer kite.Handler
}

// prepare prepares the steps to initialize the build. The build is done
// async, therefore if there is anything that needs to be checked it needs to
// be done. Any error here is passed directly to the client.
func (b *Build) prepare(r *kite.Request, c *Controller) (interface{}, error) {
	if c.CurrenState == machinestate.Building {
		return nil, NewError(ErrMachineIsBuilding)
	}

	if c.CurrenState == machinestate.Unknown {
		return nil, NewError(ErrMachineUnknownState)
	}

	// if it's something else (stopped, runnning, ...) it's been already built
	if !c.CurrenState.In(machinestate.Terminated, machinestate.NotInitialized) {
		return nil, NewError(ErrMachineInitialized)
	}

	b.Storage.UpdateState(c.MachineId, machinestate.Building)
	c.Eventer = b.NewEventer(r.Method + "-" + c.MachineId)

	instanceName := r.Username + "-" + strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	i, ok := c.Machine.Builder["instanceName"]
	if !ok || i == "" {
		// if it's empty we use the instance name that was generated above
		c.Machine.Builder["instanceName"] = instanceName
	} else {
		instanceName, ok = i.(string)
		if !ok {
			return nil, fmt.Errorf("instanceName is malformed: %v", i)
		}
	}

	// start our build process in async way
	go b.start(r, c)

	// but let the user know thay they can track us via the given event id
	return ControlResult{
		EventId: c.Eventer.Id(),
		State:   machinestate.Building,
	}, nil
}

func (b *Build) start(r *kite.Request, c *Controller) (resp interface{}, err error) {
	b.idlock.Get(c.MachineId).Lock()
	defer b.idlock.Get(c.MachineId).Unlock()

	// This is executed as the final step which stops the eventer and updates
	// the state in the storage.
	defer func() {
		status := machinestate.Running
		msg := "Build is finished successfully."
		eventErr := ""

		if err != nil {
			b.Log.Error("[%s] building failed. err %s.", c.MachineId, err.Error())

			status = c.CurrenState
			msg = ""
			eventErr = fmt.Sprintf("Building failed. Please contact support.")
		}

		// update final status in storage
		b.Storage.UpdateState(c.MachineId, status)

		// unlock distributed lock
		b.Locker.Unlock(c.MachineId)

		// let them know we are finished with our work
		c.Eventer.Push(&eventer.Event{
			Message:    msg,
			Status:     status,
			Percentage: 100,
			Error:      eventErr,
		})
	}()

	machOptions := &protocol.Machine{
		MachineId:   c.MachineId,
		Eventer:     c.Eventer,
		Credential:  c.Machine.Credential,
		Builder:     c.Machine.Builder,
		CurrentData: c.Machine.CurrentData,
	}

	msg := fmt.Sprintf("Building process started. Provider '%s'. MachineId: %+v",
		c.ProviderName, c.MachineId)

	c.Eventer.Push(&eventer.Event{Message: msg, Status: machinestate.Building})

	buildStub := `
provider      : %s
machineId     : %s
username      : %s
instanceName  : %s
`

	buildInfo := fmt.Sprintf(buildStub,
		c.ProviderName,
		c.MachineId,
		r.Username,
		c.Machine.Builder["instanceName"].(string),
	)

	b.Log.Info("[%s] building machine with following data: %s", c.MachineId, buildInfo)

	var artifact *protocol.Artifact

	artifact, err = c.Builder.Build(machOptions)
	if err != nil {
		return nil, err
	}

	if artifact == nil {
		return nil, NewError(ErrBadResponse)
	}

	// if the username is not explicit changed, assign the original username to it
	if artifact.Username == "" {
		artifact.Username = r.Username
	}

	// update if we somehow updated in build process
	c.Machine.Builder["instanceName"] = artifact.InstanceName

	r.Context.Set("buildArtifact", artifact)

	// Start the canceller for the build if something goes wrong. Like deleting
	// the terminate.
	defer func() {
		if err == nil || c.Canceller == nil {
			return
		}

		b.Log.Info("[%s] building machine failed. Starting canceller.", c.MachineId)

		if err := c.Canceller.Cancel(machOptions, artifact); err != nil {
			b.Log.Debug("[%s] couldn't run canceller. err: %s", c.MachineId, err)
		}
	}()

	deployArtifact, err := b.deployer.ServeKite(r)
	if err != nil {
		return nil, err
	}

	// garbage collect it
	r.Context = nil

	// b.Log.Debug("[controller]: building machine finished, result artifact is: %# v",
	// 	pretty.Formatter(artifact))

	resultStub := `
username   : %s
domain     : %s
ip address : %s
instance   : %s
kite query : %s
`

	resultInfo := fmt.Sprintf(resultStub,
		artifact.Username,
		artifact.DomainName,
		artifact.IpAddress,
		artifact.InstanceName,
		artifact.KiteQuery,
	)

	b.Log.Info("[%s] building machine was successfull. Artifact data: %s",
		c.MachineId, resultInfo)

	storageData := map[string]interface{}{
		"ipAddress":    artifact.IpAddress,
		"domainName":   artifact.DomainName,
		"instanceId":   artifact.InstanceId,
		"instanceName": artifact.InstanceName,
	}

	b.Log.Info("[%s] ========== %s finished ==========", c.MachineId, strings.ToUpper(r.Method))

	storageData["queryString"] = deployArtifact.(*protocol.Artifact).KiteQuery

	return true, b.Storage.Update(c.MachineId, &StorageData{
		Type: "build",
		Data: storageData,
	})
}
