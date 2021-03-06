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
	"time"

	batchv1alpha1 "github.com/vincent-pli/job-management/pkg/apis/job/v1alpha1"
	apis "github.com/vincent-pli/job-management/pkg/apis"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	clientx "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

// Name is the name of the plugin used in the plugin registry and configurations.
const (
	Name          = "gang"
	jobNamek      = "job-name"
	jobNamespacek = "job-namespace"
)

var ()

// WeightSort is a plugin that implements Priority based sorting.
type Gang struct {
	args         *Args
	handle       framework.FrameworkHandle
	waitDuration time.Duration
	client       clientx.Client
}

type Args struct {
	KubeConfig string `json:"kubeconfig,omitempty"`
	Master     string `json:"master,omitempty"`
}

var _ framework.PermitPlugin = &Gang{}

// Name returns name of the plugin.
func (g *Gang) Name() string {
	return Name
}

func (g *Gang) Permit(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (*framework.Status, time.Duration) {
	fmt.Println("Gang scheduler works in Permit plugin...")
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
		return framework.NewStatus(framework.Success, ""), g.waitDuration
	}

	job := &batchv1alpha1.XJob{}
	namespacename := types.NamespacedName{Name: jobName, Namespace: jobNamespace}

	err := g.client.Get(context.TODO(), namespacename, job)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("pod related job: <%s/%s> not found", jobNamespace, jobName)
		}
		klog.V(3).Infof("error happend: %+v", err)
		return framework.NewStatus(framework.Error, err.Error()), g.waitDuration
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
		return framework.NewStatus(framework.Wait, ""), g.waitDuration
	}

	return framework.NewStatus(framework.Success, ""), g.waitDuration
}

// New initializes a new plugin and returns it.
func New(plArgs *runtime.Unknown, handle framework.FrameworkHandle) (framework.Plugin, error) {
	fmt.Println("gang plugtin new--->")
	waitDuration, err := time.ParseDuration("1h")
	if err != nil {
		return nil, err
	}

	args := &Args{}
	if err := framework.DecodeInto(plArgs, args); err != nil {
		return nil, err
	}

	scheme := runtime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)
	_ = apis.AddToScheme(scheme)

	c, err := clientx.New(config.GetConfigOrDie(), client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return &Gang{
		args:         args,
		handle:       handle,
		waitDuration: waitDuration,
		client:       c,
	}, nil
}
