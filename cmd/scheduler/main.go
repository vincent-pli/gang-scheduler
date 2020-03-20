package main

import (
	// "flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"

	// "sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/angao/scheduler-framework-sample/pkg/plugins"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	// _ = clientgoscheme.AddToScheme(scheme)

	// _ = batchv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	// var metricsAddr string
	// var enableLeaderElection bool
	// flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	// flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
	// 	"Enable leader election for controller manager. "+
	// 		"Enabling this will ensure there is only one active controller manager.")

	// ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	// mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
	// 	Scheme:             scheme,
	// 	MetricsBindAddress: metricsAddr,
	// 	Port:               9443,
	// 	LeaderElection:     enableLeaderElection,
	// 	LeaderElectionID:   "ba5c279d.klsf.ibm.com",
	// })
	// if err != nil {
	// 	setupLog.Error(err, "unable to start manager")
	// 	os.Exit(1)
	// }

	// setupLog.Info("starting manager")
	// if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
	// 	setupLog.Error(err, "problem running manager")
	// 	os.Exit(1)
	// }

	cmd := plugins.Register()
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
