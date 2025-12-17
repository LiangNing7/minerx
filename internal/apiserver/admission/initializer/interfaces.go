package initializer

import (
	"k8s.io/apiserver/pkg/admission"

	clientset "github.com/LiangNing7/minerx/pkg/generated/clientset/versioned"
	"github.com/LiangNing7/minerx/pkg/generated/informers"
)

// WantsExternalInformerFactory defines a function which sets InformerFactory for admission plugins that need it.
type WantsExternalInformerFactory interface {
	admission.InitializationValidator
	SetExternalInformerFactory(informers.SharedInformerFactory)
}

// WantsExternalClientSet defines a function which sets external ClientSet for admission plugins that need it.
type WantsExternalClientSet interface {
	admission.InitializationValidator
	SetExternalClientSet(clientset.Interface)
}
