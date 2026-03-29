package printers

import (
	"bytes"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestPrintSecrets_NoSensitiveDataInOutput(t *testing.T) {
	list := &v1.SecretList{
		Items: []v1.Secret{{
			ObjectMeta: metav1.ObjectMeta{Name: "my-secret"},
			Type:       v1.SecretTypeOpaque,
			Data: map[string][]byte{
				"password": []byte("super-secret-value"),
				"token":    []byte("my-api-token"),
			},
		}},
	}

	var output bytes.Buffer
	Printer(&output, list, "")

	if strings.Contains(output.String(), "super-secret-value") {
		t.Fatalf("secret value leaked in output: %q", output.String())
	}

	if strings.Contains(output.String(), "my-api-token") {
		t.Fatalf("token leaked in output: %q", output.String())
	}

	if !strings.Contains(output.String(), "2") {
		t.Fatalf("expected data count in output, got %q", output.String())
	}
}

func TestPrinter_EmptyList(t *testing.T) {
	testCases := []struct {
		name     string
		resource interface{}
	}{
		{name: "pods", resource: &v1.PodList{}},
		{name: "secrets", resource: &v1.SecretList{}},
		{name: "configmaps", resource: &v1.ConfigMapList{}},
		{name: "deployments", resource: &appsv1.DeploymentList{}},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			var output bytes.Buffer
			Printer(&output, testCase.resource, "")

			if output.String() != "" {
				t.Fatalf("expected no output, got %q", output.String())
			}
		})
	}
}

func TestPrinter_PatternFilter(t *testing.T) {
	list := &v1.ConfigMapList{
		Items: []v1.ConfigMap{
			{ObjectMeta: metav1.ObjectMeta{Name: "nginx-config"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "redis-config"}},
		},
	}

	var output bytes.Buffer
	Printer(&output, list, "nginx")

	if !strings.Contains(output.String(), "nginx-config") {
		t.Fatalf("expected matching configmap in output, got %q", output.String())
	}

	if strings.Contains(output.String(), "redis-config") {
		t.Fatalf("expected non-matching configmap to be filtered out, got %q", output.String())
	}
}

func TestMatchesPattern(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		pattern string
		match   bool
	}{
		{
			name:    "empty pattern matches everything",
			value:   "configmap",
			pattern: "",
			match:   true,
		},
		{
			name:    "substring matches",
			value:   "nginx-config",
			pattern: "nginx",
			match:   true,
		},
		{
			name:    "non matching pattern returns false",
			value:   "redis-config",
			pattern: "nginx",
			match:   false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			if got := matchesPattern(testCase.value, testCase.pattern); got != testCase.match {
				t.Fatalf("matchesPattern(%q, %q) = %v, want %v", testCase.value, testCase.pattern, got, testCase.match)
			}
		})
	}
}

func TestPrinter_UnstructuredList(t *testing.T) {
	list := &unstructured.UnstructuredList{
		Items: []unstructured.Unstructured{
			{
				Object: map[string]interface{}{
					"apiVersion": "example.com/v1",
					"kind":       "Widget",
					"metadata": map[string]interface{}{
						"name":      "demo-widget",
						"namespace": "default",
					},
				},
			},
		},
	}
	list.SetKind("Widget")

	var output bytes.Buffer
	Printer(&output, list, "demo")

	if !strings.Contains(output.String(), "Widget") {
		t.Fatalf("expected kind header in output, got %q", output.String())
	}

	if !strings.Contains(output.String(), "demo-widget") {
		t.Fatalf("expected resource name in output, got %q", output.String())
	}
}
