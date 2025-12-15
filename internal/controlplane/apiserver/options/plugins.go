package options

// This file exists to force the desired plugin implementations to be linked.
// This should probably be part of some configuration fed into the build for a
// given binary target.
import (
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"

	// Admission policies.
	"github.com/LiangNing7/minerx/internal/controlplane/admission/plugin/admit"
	"github.com/LiangNing7/minerx/internal/controlplane/admission/plugin/deny"
	"github.com/LiangNing7/minerx/internal/controlplane/admission/plugin/namespace/autoprovision"
	"github.com/LiangNing7/minerx/internal/controlplane/admission/plugin/namespace/exists"
	"github.com/LiangNing7/minerx/internal/controlplane/admission/plugin/namespace/lifecycle"
)

// AllOrderedPlugins is the list of all the plugins in order.
var AllOrderedPlugins = []string{
	admit.PluginName,         // AlwaysAdmit
	autoprovision.PluginName, // NamespaceAutoProvision
	lifecycle.PluginName,     // NamespaceLifecycle
	exists.PluginName,        // NamespaceExists

	// new admission plugins should generally be inserted above here
	// webhook, resourcequota, and deny plugins must go at the end
	// minerset.PluginName, // MinerSet

	// mutatingwebhook.PluginName,   // MutatingAdmissionWebhook
	// validatingwebhook.PluginName, // ValidatingAdmissionWebhook
	// resourcequota.PluginName, // ResourceQuota
	deny.PluginName, // AlwaysDeny
}

// RegisterAllAdmissionPlugins registers all admission plugins.
// The order of registration is irrelevant, see AllOrderedPlugins for execution order.
func RegisterAllAdmissionPlugins(plugins *admission.Plugins) {
	admit.Register(plugins) // DEPRECATED as no real meaning
	autoprovision.Register(plugins)
	lifecycle.Register(plugins)
	exists.Register(plugins)
	// minerset.Register(plugins)
	deny.Register(plugins) // DEPRECATED as no real meaning
}

// DefaultOffAdmissionPlugins get admission plugins off by default for kube-apiserver.
func DefaultOffAdmissionPlugins() sets.Set[string] {
	defaultOnPlugins := sets.New(
		autoprovision.PluginName, // NamespaceAutoProvision
		lifecycle.PluginName,     // NamespaceLifecycle
	)

	return sets.New(AllOrderedPlugins...).Difference(defaultOnPlugins)
}
