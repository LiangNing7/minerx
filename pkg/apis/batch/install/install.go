package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	"github.com/LiangNing7/minerx/pkg/apis/batch"
	batchv1beta1 "github.com/LiangNing7/minerx/pkg/apis/batch/v1beta1"
)

func init() {
	Install(legacyscheme.Scheme)
}

// Install registers the API group and adds types to a scheme.
func Install(scheme *runtime.Scheme) {
	utilruntime.Must(batch.AddToScheme(scheme))
	utilruntime.Must(batchv1beta1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(batchv1beta1.SchemeGroupVersion))
}
