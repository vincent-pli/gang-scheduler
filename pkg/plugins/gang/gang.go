/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gang

import (
	"context"
	"fmt"

	batchv1alpha1 "github.ibm.com/panpxpx/klsf/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clientx "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	// "sigs.k8s.io/controller-runtime/pkg/config"
)

// Name is the name of the plugin used in the plugin registry and configurations.
const (
	Name          = "Gang"
	jobNamek      = "job-name"
	jobNamespacek = "job-namespace"
)

// var (
// 	client clientX.Client
// )

// WeightSort is a plugin that implements Priority based sorting.
type Gang struct {
	args   *Args
	handle framework.FrameworkHandle
	// client clientX.Client
}

type Args struct {
	KubeConfig string `json:"kubeconfig,omitempty"`
	Master     string `json:"master,omitempty"`
}

var _ framework.PreBindPlugin = &Gang{}

// Name returns name of the plugin.
func (pl *Gang) Name() string {
	return Name
}

func (g *Gang) PreBind(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) *framework.Status {
	var jobName string
	var jobNamespace string

	labels := pod.Labels
	for k, v := range labels {
		switch k {
		case jobNamek:
			jobName = v
		case jobNamespacek:
			jobNamespace = v
		}
	}

	if jobName == "" || jobNamespace == "" {
		return framework.NewStatus(framework.Success, "")
	}

	job := &batchv1alpha1.LSFJob{}
	namespacename := types.NamespacedName{Name: jobName, Namespace: jobNamespace}

	c, err := clientx.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		framework.NewStatus(framework.Error, err.Error())
	}

	err = c.Get(context.TODO(), namespacename, job)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("pod related job: <%s/%s> not found", jobNamespace, jobName)
		}
		framework.NewStatus(framework.Error, err.Error())
	}

	// labelSet := labelsx.Set{
	// 	jobNamespacek: jobNamespace,
	// 	jobNamek:      jobName,
	// }

	minAvailable := job.Spec.MinAvailable

	numWaiting := 1
	var waitingPods []framework.WaitingPod
	g.handle.IterateOverWaitingPods(func(wpod framework.WaitingPod) {
		labels := wpod.GetPod().Labels
		if labels[jobNamespacek] == jobNamespace && labels[jobNamek] == jobName {
			numWaiting++
			waitingPods = append(waitingPods, wpod)
		}
	})

	if int32(numWaiting) >= minAvailable {
		for _, waitingPod := range waitingPods {
			waitingPod.Allow(Name)
		}
	} else {
		return framework.NewStatus(framework.Wait, "")
	}

	// pods, err := g.handle.SnapshotSharedLister().Pods().List(labelSet.AsSelector())
	// if err != nil {
	// 	framework.NewStatus(framework.Error, err.Error())
	// }

	// for index, pod := range pods {

	// }

	// nodeInfo, err := g.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	// if err != nil {
	// 	return framework.NewStatus(framework.Error, err.Error())
	// }

	return framework.NewStatus(framework.Success, "")
}

// func Config(client clientX.Client) {
// 	client = client
// }

// New initializes a new plugin and returns it.
func New(plArgs *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
	args := &Args{}
	if err := framework.DecodeInto(plArgs, args); err != nil {
		return nil, err
	}

	return &Gang{
		args:   args,
		handle: handle,
		// client: client,
	}, nil
}
