// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"flag"
	"os"

	debug "github.com/talos-systems/go-debug"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	cgrecord "k8s.io/client-go/tools/record"
	capiv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	infrav1alpha2 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha2"
	infrav1alpha3 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	"github.com/talos-systems/sidero/app/caps-controller-manager/controllers"
	metalv1alpha1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

const (
	debugAddr = ":9994"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

//nolint:wsl
func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = capiv1.AddToScheme(scheme)
	_ = infrav1alpha2.AddToScheme(scheme)
	_ = infrav1alpha3.AddToScheme(scheme)
	_ = metalv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
		webhookPort          int
	)

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", true,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.IntVar(&webhookPort, "webhook-port", 0, "Webhook Server port, disabled by default. When enabled, the manager will only work as webhook server, no reconcilers are installed.")
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	go func() {
		debugLogFunc := func(msg string) {
			setupLog.Info(msg)
		}
		if err := debug.ListenAndServe(context.TODO(), debugAddr, debugLogFunc); err != nil {
			setupLog.Error(err, "failed to start debug server")
			os.Exit(1)
		}
	}()

	// Machine and cluster operations can create enough events to trigger the event recorder spam filter
	// Setting the burst size higher ensures all events will be recorded and submitted to the API
	broadcaster := cgrecord.NewBroadcasterWithCorrelatorOptions(cgrecord.CorrelatorOptions{
		BurstSize: 100,
	})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "controller-leader-election-capm",
		Port:               webhookPort,
		EventBroadcaster:   broadcaster,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(mgr.GetConfig())
	if err != nil {
		setupLog.Error(err, "unable to create k8s client")
		os.Exit(1)
	}

	eventBroadcaster := cgrecord.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(
		&typedcorev1.EventSinkImpl{
			Interface: clientset.CoreV1().Events(""),
		})

	recorder := eventBroadcaster.NewRecorder(
		mgr.GetScheme(),
		corev1.EventSource{Component: "caps-controller-manager"})

	if webhookPort == 0 {
		if err = (&controllers.MetalClusterReconciler{
			Client: mgr.GetClient(),
			Log:    ctrl.Log.WithName("controllers").WithName("MetalCluster"),
			Scheme: mgr.GetScheme(),
		}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: 10}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "MetalCluster")
			os.Exit(1)
		}

		if err = (&controllers.MetalMachineReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("MetalMachine"),
			Scheme:   mgr.GetScheme(),
			Recorder: recorder,
		}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: 10}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "MetalMachine")
			os.Exit(1)
		}

		if err = (&controllers.ServerBindingReconciler{
			Client:   mgr.GetClient(),
			Log:      ctrl.Log.WithName("controllers").WithName("ServerBinding"),
			Scheme:   mgr.GetScheme(),
			Recorder: recorder,
		}).SetupWithManager(mgr, controller.Options{MaxConcurrentReconciles: 10}); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "ServerBinding")
			os.Exit(1)
		}
	} else {
		if err = (&infrav1alpha3.MetalCluster{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "MetalCluster")
			os.Exit(1)
		}

		if err = (&infrav1alpha3.MetalMachine{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "MetalMachine")
			os.Exit(1)
		}

		if err = (&infrav1alpha3.MetalMachineTemplate{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "MetalMachineTemplate")
			os.Exit(1)
		}

		if err = (&infrav1alpha3.ServerBinding{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "ServerBinding")
			os.Exit(1)
		}
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
