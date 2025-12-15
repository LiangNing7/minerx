package initializer

import (
	"k8s.io/apiserver/pkg/admission"

	kubeinformers "k8s.io/client-go/informers"
	kubeclientset "k8s.io/client-go/kubernetes"
)

// WantsInternalInformerFactory defines a function which sets InformerFactory for admission plugins that need it.
type WantsInternalInformerFactory interface {
	admission.InitializationValidator
	SetInternalInformerFactory(kubeinformers.SharedInformerFactory)
}

// WantsInternalClientSet defines a function which sets external ClientSet for admission plugins that need it.
type WantsInternalClientSet interface {
	admission.InitializationValidator
	SetInternalClientSet(kubeclientset.Interface)
}
