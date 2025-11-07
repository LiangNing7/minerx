// Package install installs the experimental API group, making it available as
// an option to all of the API encoding/decoding machinery.
package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	"github.com/LiangNing7/minerx/pkg/apis/apps"
	appsv1beta1 "github.com/LiangNing7/minerx/pkg/apis/apps/v1beta1"
)

func init() {
	Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme.
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(apps.AddToScheme(scheme))
	utilruntime.Must(appsv1beta1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(appsv1beta1.SchemeGroupVersion))
}
