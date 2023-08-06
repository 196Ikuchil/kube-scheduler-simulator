// Deprecated: it was renamed to the debuggablescheduler plugin. We'll remove it soon.
package externalscheduler

import (
	"github.com/spf13/cobra"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
	"k8s.io/kubernetes/pkg/scheduler/framework/runtime"

	"sigs.k8s.io/kube-scheduler-simulator/simulator/scheduler/plugin"
)

// Deprecated: the externalscheduler package was renamed to the debuggablescheduler plugin. We'll remove it soon.
func NewSchedulerCommand(opts ...Option) (*cobra.Command, func(), error) {
	opt := &options{pluginExtender: map[string]plugin.PluginExtenderInitializer{}, outOfTreeRegistry: map[string]runtime.PluginFactory{}}
	for _, o := range opts {
		o(opt)
	}

	scheduleropts, cancelFn, err := CreateOptionForOutOfTreePlugin(opt.outOfTreeRegistry, opt.pluginExtender)
	if err != nil {
		return nil, cancelFn, err
	}

	command := app.NewSchedulerCommand(scheduleropts...)

	return command, cancelFn, nil
}

type options struct {
	outOfTreeRegistry runtime.Registry
	pluginExtender    map[string]plugin.PluginExtenderInitializer
}

// Deprecated: the externalscheduler package was renamed to the debuggablescheduler plugin. We'll remove it soon.
type Option func(opt *options)

// WithPlugin creates an Option based on plugin name and factory.
// Deprecated: the externalscheduler package was renamed to the debuggablescheduler plugin. We'll remove it soon.
func WithPlugin(pluginName string, factory runtime.PluginFactory) Option {
	return func(opt *options) {
		opt.outOfTreeRegistry[pluginName] = factory
	}
}

// WithPluginExtenders creates an Option based on plugin name and plugin extenders.
// Deprecated: the externalscheduler package was renamed to the debuggablescheduler plugin. We'll remove it soon.
func WithPluginExtenders(pluginName string, e plugin.PluginExtenderInitializer) Option {
	return func(opt *options) {
		opt.pluginExtender[pluginName] = e
	}
}
