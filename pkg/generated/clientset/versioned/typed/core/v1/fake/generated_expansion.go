package fake

import (
    "context"

    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Finalize takes the representation of a namespace to update.  Returns the server's representation of the namespace, and an error, if it occurs.
func (c *fakeNamespaces) Finalize(ctx context.Context, namespace *v1.Namespace, opts metav1.UpdateOptions) (*v1.Namespace, error) {
    return nil, nil
}
