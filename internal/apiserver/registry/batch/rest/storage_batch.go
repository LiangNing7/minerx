package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"

	cronjobstore "github.com/LiangNing7/minerx/internal/apiserver/registry/batch/cronjob/storage"
	jobstore "github.com/LiangNing7/minerx/internal/apiserver/registry/batch/job/storage"
	serializerutil "github.com/LiangNing7/minerx/internal/pkg/util/serializer"
	"github.com/LiangNing7/minerx/pkg/apis/batch"
	"github.com/LiangNing7/minerx/pkg/apis/batch/v1beta1"
	"github.com/LiangNing7/minerx/pkg/apiserver/storage"
)

// RESTStorageProvider is a struct for batch REST storage.
type RESTStorageProvider struct{}

// Implement RESTStorageProvider.
var _ storage.RESTStorageProvider = &RESTStorageProvider{}

// NewRESTStorage returns APIGroupInfo object.
func (p RESTStorageProvider) NewRESTStorage(
	apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter,
) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(batch.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	apiGroupInfo.NegotiatedSerializer = serializerutil.NewProtocolShieldSerializers(&legacyscheme.Codecs)

	storageMap, err := p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
	if err != nil {
		return genericapiserver.APIGroupInfo{}, err
	}
	apiGroupInfo.VersionedResourcesStorageMap[v1beta1.SchemeGroupVersion.Version] = storageMap

	return apiGroupInfo, nil
}

func (p RESTStorageProvider) v1beta1Storage(
	apiResourceConfigSource serverstorage.APIResourceConfigSource,
	restOptionsGetter generic.RESTOptionsGetter,
) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	//nolint:goconst
	// cronjobs
	if resource := "cronjobs"; apiResourceConfigSource.ResourceEnabled(v1beta1.SchemeGroupVersion.WithResource(resource)) {
		cronJobStorage, err := cronjobstore.NewStorage(restOptionsGetter)
		if err != nil {
			return storage, err
		}

		storage[resource] = cronJobStorage.CronJob
		storage[resource+"/status"] = cronJobStorage.Status
	}

	// jobs
	if resource := "jobs"; apiResourceConfigSource.ResourceEnabled(v1beta1.SchemeGroupVersion.WithResource(resource)) {
		jobStorage, err := jobstore.NewStorage(restOptionsGetter)
		if err != nil {
			return storage, err
		}

		storage[resource] = jobStorage.Job
		storage[resource+"/status"] = jobStorage.Status
	}

	return storage, nil
}

// GroupName return the api group name.
func (p RESTStorageProvider) GroupName() string {
	return batch.GroupName
}
