# Testing Specification

## Purpose
Defines the test coverage requirements for ksearch. All packages SHALL have
accompanying test files. Tests MUST be runnable without a live Kubernetes
cluster unless explicitly annotated as integration tests.

## Requirements

### Requirement: Unit tests for all packages
Every package SHALL have a `*_test.go` file covering its core logic.
Tests SHALL use only the stdlib `testing` package and
`k8s.io/client-go/kubernetes/fake` for Kubernetes API interactions.

#### Scenario: Tests pass without a cluster
- GIVEN no kubeconfig or cluster is available
- WHEN `go test ./...` is run
- THEN all non-integration tests pass

### Requirement: Secret output safety test
The printer SHALL have a test that constructs a `*v1.SecretList` with
known values in `Data` and asserts those values never appear in output.

#### Scenario: Secret values absent from printer output
- GIVEN a SecretList with Data containing known sentinel strings
- WHEN Printer() is called
- THEN the output contains the data count but not the sentinel strings

### Requirement: Pattern filter correctness
The printer SHALL have tests verifying that the pattern filter includes
matching items and excludes non-matching items.

#### Scenario: Matching item included
- GIVEN a ConfigMapList with items "nginx-config" and "redis-config"
- WHEN Printer() is called with pattern="nginx"
- THEN "nginx-config" appears in output and "redis-config" does not

### Requirement: Empty list produces no output
The printer SHALL produce zero bytes for any resource type with an empty
item list.

#### Scenario: Empty list is silent
- GIVEN any resource list type with zero Items
- WHEN Printer() is called
- THEN no bytes are written

### Requirement: Getter channel always closes
`Getter()` SHALL always close the result channel on return, including when
an unknown kind is supplied.

#### Scenario: Channel closed on unknown kind
- GIVEN an unrecognised kind string
- WHEN Getter() is called
- THEN the channel is closed within 2 seconds without blocking

### Requirement: Cache security test
The cache package SHALL have a test asserting that no known sensitive
strings appear in a written cache file when fed printer-formatted output.

#### Scenario: Sensitive strings absent from cache file
- GIVEN a SectionEntry whose Output is known printer-formatted text
- WHEN Write() is called and the file is read back raw
- THEN no sensitive sentinel strings are present in the JSON file

### Requirement: Integration tests gated by build tag
Tests that require a live Kubernetes cluster SHALL use the `integration`
build tag and SHALL create their own cluster using `kind` if available.

#### Scenario: Integration tests skipped in default run
- GIVEN no `kind` binary and no live cluster
- WHEN `go test ./...` is run (no -tags flag)
- THEN integration tests are not executed and the run succeeds
