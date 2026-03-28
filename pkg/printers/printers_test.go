package printers

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func captureOutput(t *testing.T, f func()) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create pipe: %v", err)
	}

	originalStdout := os.Stdout
	os.Stdout = writer

	defer func() {
		os.Stdout = originalStdout
	}()

	f()

	if err := writer.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	var buffer bytes.Buffer
	if _, err := io.Copy(&buffer, reader); err != nil {
		t.Fatalf("copy captured output: %v", err)
	}

	if err := reader.Close(); err != nil {
		t.Fatalf("close reader: %v", err)
	}

	return buffer.String()
}

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

	output := captureOutput(t, func() {
		Printer(list, "")
	})

	if strings.Contains(output, "super-secret-value") {
		t.Fatalf("secret value leaked in output: %q", output)
	}

	if strings.Contains(output, "my-api-token") {
		t.Fatalf("token leaked in output: %q", output)
	}

	if !strings.Contains(output, "\t2\t") && !strings.Contains(output, " 2 ") && !strings.Contains(output, "2") {
		t.Fatalf("expected data count in output, got %q", output)
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
			output := captureOutput(t, func() {
				Printer(testCase.resource, "")
			})

			if output != "" {
				t.Fatalf("expected no output, got %q", output)
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

	output := captureOutput(t, func() {
		Printer(list, "nginx")
	})

	if !strings.Contains(output, "nginx-config") {
		t.Fatalf("expected matching configmap in output, got %q", output)
	}

	if strings.Contains(output, "redis-config") {
		t.Fatalf("expected non-matching configmap to be filtered out, got %q", output)
	}
}
