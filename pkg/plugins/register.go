package plugins

import (
	"github.com/vincent-pli/gang-scheduler/pkg/plugins/gang"
	"github.com/vincent-pli/gang-scheduler/pkg/plugins/sample"
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

func Register() *cobra.Command {
	return app.NewSchedulerCommand(
		app.WithPlugin(sample.Name, sample.New),
		app.WithPlugin(gang.Name, gang.New),
	)
}
