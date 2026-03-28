# Proposal: Dependency Upgrade

## Purpose
Bring all direct dependencies to current stable versions to close security gaps,
restore compatibility with Kubernetes 1.27–1.32 clusters, and unblock
changes 0002 and 0003 which rely on modern APIs.

## Requirements

### Requirement: Go toolchain upgrade
The module SHALL declare Go 1.22 and compile cleanly under it.

#### Scenario: Build succeeds on Go 1.22
- GIVEN go.mod declares `go 1.22`
- WHEN `go build ./...` is run with Go 1.22
- THEN the binary is produced with no errors

### Requirement: Kubernetes client-go upgrade
The module SHALL use `k8s.io/client-go` v0.32.x.

#### Scenario: All List calls accept context
- GIVEN client-go v0.32.x
- WHEN any resource `.List()` method is called
- THEN the call signature includes `context.Context` as the first argument

### Requirement: controller-runtime upgrade
The module SHALL use `sigs.k8s.io/controller-runtime` v0.20.x.

#### Scenario: Kubeconfig loading still works
- GIVEN controller-runtime v0.20.x
- WHEN ksearch is invoked with a valid kubeconfig
- THEN a REST config is obtained without error

### Requirement: Cobra and Logrus upgrade
The module SHALL use cobra v1.8.x and logrus v1.9.x.

#### Scenario: All existing flags remain functional
- GIVEN cobra v1.8.x
- WHEN `ksearch --help` is run
- THEN --pattern, --namespace, and --kinds flags are listed correctly
