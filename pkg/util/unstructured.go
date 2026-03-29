package util

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/cbor"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/rest"
)

var (
	unstructuredScheme          = runtime.NewScheme()
	unstructuredParameterScheme = runtime.NewScheme()
	unstructuredParameterCodec  = runtime.NewParameterCodec(unstructuredParameterScheme)
	versionV1                   = schema.GroupVersion{Version: "v1"}
)

func init() {
	metav1.AddToGroupVersion(unstructuredScheme, versionV1)
	metav1.AddToGroupVersion(unstructuredParameterScheme, versionV1)
}

func listUnstructuredResource(ctx context.Context, cfg *rest.Config, namespace string, meta ResourceMeta) (*unstructured.UnstructuredList, error) {
	if cfg == nil {
		return nil, fmt.Errorf("rest config is required for dynamic resource %s", meta.Resource)
	}

	httpClient, err := rest.HTTPClientFor(cfg)
	if err != nil {
		return nil, err
	}

	client, err := rest.UnversionedRESTClientForConfigAndClient(configForUnstructured(cfg), httpClient)
	if err != nil {
		return nil, err
	}

	segments := make([]string, 0, 6)
	if meta.APIGroup == "" {
		segments = append(segments, "api")
	} else {
		segments = append(segments, "apis", meta.APIGroup)
	}
	segments = append(segments, meta.APIVersion)
	if meta.Namespaced && namespace != "" {
		segments = append(segments, "namespaces", namespace)
	}
	segments = append(segments, meta.Resource)

	var out unstructured.UnstructuredList
	if err := client.Get().
		AbsPath(segments...).
		SpecificallyVersionedParams(&metav1.ListOptions{}, unstructuredParameterCodec, versionV1).
		Do(ctx).
		Into(&out); err != nil {
		return nil, err
	}

	out.SetKind(meta.Kind)
	out.SetAPIVersion(joinGroupVersion(meta.APIGroup, meta.APIVersion))

	return &out, nil
}

func configForUnstructured(inConfig *rest.Config) *rest.Config {
	config := rest.CopyConfig(inConfig)
	config.ContentType = "application/json"
	config.AcceptContentTypes = "application/json"
	config.GroupVersion = nil
	config.APIPath = "/"
	config.NegotiatedSerializer = basicUnstructuredNegotiatedSerializer{
		supportedMediaTypes: []runtime.SerializerInfo{
			{
				MediaType:        "application/json",
				MediaTypeType:    "application",
				MediaTypeSubType: "json",
				EncodesAsText:    true,
				Serializer:       json.NewSerializerWithOptions(json.DefaultMetaFactory, unstructuredCreater{nested: unstructuredScheme}, unstructuredTyper{nested: unstructuredScheme}, json.SerializerOptions{}),
				PrettySerializer: json.NewSerializerWithOptions(json.DefaultMetaFactory, unstructuredCreater{nested: unstructuredScheme}, unstructuredTyper{nested: unstructuredScheme}, json.SerializerOptions{Pretty: true}),
				StreamSerializer: &runtime.StreamSerializerInfo{
					EncodesAsText: true,
					Serializer:    json.NewSerializerWithOptions(json.DefaultMetaFactory, unstructuredScheme, unstructuredScheme, json.SerializerOptions{}),
					Framer:        json.Framer,
				},
			},
			{
				MediaType:        "application/cbor",
				MediaTypeType:    "application",
				MediaTypeSubType: "cbor",
				Serializer:       cbor.NewSerializer(unstructuredCreater{nested: unstructuredScheme}, unstructuredTyper{nested: unstructuredScheme}),
				StreamSerializer: &runtime.StreamSerializerInfo{
					Serializer: cbor.NewSerializer(unstructuredScheme, unstructuredScheme, cbor.Transcode(false)),
					Framer:     cbor.NewFramer(),
				},
			},
		},
	}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return config
}

func joinGroupVersion(group, version string) string {
	if group == "" {
		return version
	}
	return group + "/" + version
}

type basicUnstructuredNegotiatedSerializer struct {
	supportedMediaTypes []runtime.SerializerInfo
}

func (s basicUnstructuredNegotiatedSerializer) SupportedMediaTypes() []runtime.SerializerInfo {
	return s.supportedMediaTypes
}

func (s basicUnstructuredNegotiatedSerializer) EncoderForVersion(encoder runtime.Encoder, gv runtime.GroupVersioner) runtime.Encoder {
	return runtime.WithVersionEncoder{
		Version:     gv,
		Encoder:     encoder,
		ObjectTyper: permissiveTyper{nested: unstructuredScheme},
	}
}

func (s basicUnstructuredNegotiatedSerializer) DecoderToVersion(decoder runtime.Decoder, gv runtime.GroupVersioner) runtime.Decoder {
	return decoder
}

type unstructuredCreater struct {
	nested runtime.ObjectCreater
}

func (c unstructuredCreater) New(kind schema.GroupVersionKind) (runtime.Object, error) {
	out, err := c.nested.New(kind)
	if err == nil {
		return out, nil
	}

	out = &unstructured.Unstructured{}
	out.GetObjectKind().SetGroupVersionKind(kind)
	return out, nil
}

type unstructuredTyper struct {
	nested runtime.ObjectTyper
}

func (t unstructuredTyper) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	kinds, unversioned, err := t.nested.ObjectKinds(obj)
	if err == nil {
		return kinds, unversioned, nil
	}
	if _, ok := obj.(runtime.Unstructured); ok && !obj.GetObjectKind().GroupVersionKind().Empty() {
		return []schema.GroupVersionKind{obj.GetObjectKind().GroupVersionKind()}, false, nil
	}
	return nil, false, err
}

func (t unstructuredTyper) Recognizes(gvk schema.GroupVersionKind) bool {
	return true
}

type permissiveTyper struct {
	nested runtime.ObjectTyper
}

func (t permissiveTyper) ObjectKinds(obj runtime.Object) ([]schema.GroupVersionKind, bool, error) {
	kinds, unversioned, err := t.nested.ObjectKinds(obj)
	if err == nil {
		return kinds, unversioned, nil
	}
	if _, ok := obj.(runtime.Unstructured); ok {
		return []schema.GroupVersionKind{obj.GetObjectKind().GroupVersionKind()}, false, nil
	}
	return nil, false, err
}

func (t permissiveTyper) Recognizes(gvk schema.GroupVersionKind) bool {
	return true
}
