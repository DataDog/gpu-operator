// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	context "context"

	imagev1 "github.com/openshift/api/image/v1"
	applyconfigurationsimagev1 "github.com/openshift/client-go/image/applyconfigurations/image/v1"
	scheme "github.com/openshift/client-go/image/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// ImageStreamsGetter has a method to return a ImageStreamInterface.
// A group's client should implement this interface.
type ImageStreamsGetter interface {
	ImageStreams(namespace string) ImageStreamInterface
}

// ImageStreamInterface has methods to work with ImageStream resources.
type ImageStreamInterface interface {
	Create(ctx context.Context, imageStream *imagev1.ImageStream, opts metav1.CreateOptions) (*imagev1.ImageStream, error)
	Update(ctx context.Context, imageStream *imagev1.ImageStream, opts metav1.UpdateOptions) (*imagev1.ImageStream, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, imageStream *imagev1.ImageStream, opts metav1.UpdateOptions) (*imagev1.ImageStream, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*imagev1.ImageStream, error)
	List(ctx context.Context, opts metav1.ListOptions) (*imagev1.ImageStreamList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *imagev1.ImageStream, err error)
	Apply(ctx context.Context, imageStream *applyconfigurationsimagev1.ImageStreamApplyConfiguration, opts metav1.ApplyOptions) (result *imagev1.ImageStream, err error)
	// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
	ApplyStatus(ctx context.Context, imageStream *applyconfigurationsimagev1.ImageStreamApplyConfiguration, opts metav1.ApplyOptions) (result *imagev1.ImageStream, err error)
	Secrets(ctx context.Context, imageStreamName string, options metav1.GetOptions) (*imagev1.SecretList, error)
	Layers(ctx context.Context, imageStreamName string, options metav1.GetOptions) (*imagev1.ImageStreamLayers, error)

	ImageStreamExpansion
}

// imageStreams implements ImageStreamInterface
type imageStreams struct {
	*gentype.ClientWithListAndApply[*imagev1.ImageStream, *imagev1.ImageStreamList, *applyconfigurationsimagev1.ImageStreamApplyConfiguration]
}

// newImageStreams returns a ImageStreams
func newImageStreams(c *ImageV1Client, namespace string) *imageStreams {
	return &imageStreams{
		gentype.NewClientWithListAndApply[*imagev1.ImageStream, *imagev1.ImageStreamList, *applyconfigurationsimagev1.ImageStreamApplyConfiguration](
			"imagestreams",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *imagev1.ImageStream { return &imagev1.ImageStream{} },
			func() *imagev1.ImageStreamList { return &imagev1.ImageStreamList{} },
		),
	}
}

// Secrets takes name of the imageStream, and returns the corresponding imagev1.SecretList object, and an error if there is any.
func (c *imageStreams) Secrets(ctx context.Context, imageStreamName string, options metav1.GetOptions) (result *imagev1.SecretList, err error) {
	result = &imagev1.SecretList{}
	err = c.GetClient().Get().
		Namespace(c.GetNamespace()).
		Resource("imagestreams").
		Name(imageStreamName).
		SubResource("secrets").
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// Layers takes name of the imageStream, and returns the corresponding imagev1.ImageStreamLayers object, and an error if there is any.
func (c *imageStreams) Layers(ctx context.Context, imageStreamName string, options metav1.GetOptions) (result *imagev1.ImageStreamLayers, err error) {
	result = &imagev1.ImageStreamLayers{}
	err = c.GetClient().Get().
		Namespace(c.GetNamespace()).
		Resource("imagestreams").
		Name(imageStreamName).
		SubResource("layers").
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}
