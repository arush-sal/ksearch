#!/usr/bin/env bash

set -euo pipefail

cleanup() {
	kubectl krew uninstall ksearch >/dev/null 2>&1 || true
}

trap cleanup EXIT

install_krew() {
	if kubectl krew version >/dev/null 2>&1; then
		return
	fi

	local temp_dir os arch krew
	temp_dir="$(mktemp -d)"
	os="$(uname | tr '[:upper:]' '[:lower:]')"
	arch="$(uname -m | sed -e 's/x86_64/amd64/' -e 's/\(arm\)\(64\)\?.*/\1\2/' -e 's/aarch64$/arm64/')"
	krew="krew-${os}_${arch}"

	curl -fsSL "https://github.com/kubernetes-sigs/krew/releases/latest/download/${krew}.tar.gz" -o "${temp_dir}/${krew}.tar.gz"
	tar -C "${temp_dir}" -zxf "${temp_dir}/${krew}.tar.gz"
	"${temp_dir}/${krew}" install krew
	rm -rf "${temp_dir}"
}

main() {
	export KREW_ROOT="${KREW_ROOT:-${HOME}/.krew}"
	export PATH="${KREW_ROOT}/bin:${PATH}"

	install_krew

	rm -rf dist
	goreleaser release --snapshot --clean

	local manifest archive plugin_path help_output
	manifest="dist/krew/ksearch.yaml"
	archive="dist/ksearch_linux_amd64.tar.gz"

	if [[ ! -f "${manifest}" ]]; then
		echo "missing manifest: ${manifest}" >&2
		exit 1
	fi

	if [[ ! -f "${archive}" ]]; then
		echo "missing archive: ${archive}" >&2
		exit 1
	fi

	kubectl krew uninstall ksearch >/dev/null 2>&1 || true
	kubectl krew install --manifest="${manifest}" --archive="${archive}"

	plugin_path="$(command -v kubectl-ksearch)"
	if [[ -z "${plugin_path}" ]]; then
		echo "kubectl-ksearch not found on PATH after install" >&2
		exit 1
	fi

	if ! kubectl plugin list | grep -q 'kubectl-ksearch'; then
		echo "kubectl plugin list did not include kubectl-ksearch" >&2
		exit 1
	fi

	help_output="$(kubectl ksearch --help)"
	if [[ "${help_output}" != *"kubectl ksearch [flags]"* ]]; then
		echo "unexpected help output:" >&2
		printf '%s\n' "${help_output}" >&2
		exit 1
	fi

	printf 'Validated Krew install via %s\n' "${archive}"
	printf 'Plugin command path: %s\n' "${plugin_path}"
	printf '%s\n' "${help_output}" | sed -n '1,10p'
}

main "$@"
