# 🔎 ksearch

> Search your Kubernetes resources like you mean it.

[![CI](https://img.shields.io/github/actions/workflow/status/arush-sal/ksearch/ci.yml?branch=master&label=CI)](https://github.com/arush-sal/ksearch/actions)
[![Release](https://img.shields.io/github/v/release/arush-sal/ksearch)](https://github.com/arush-sal/ksearch/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/arush-sal/ksearch)](./go.mod)
[![License](https://img.shields.io/github/license/arush-sal/ksearch)](./LICENSE)

[Install with Krew](./README.md#installation) • [Download a release](ttps://github.com/arush-sal/ksearch/releases) • [See examples](./README.md#examples)

## ✨ Features

- 🔍 Search resources by name pattern
- 📦 Filter results by selected kinds
- 🧭 Scope to a namespace or search broadly
- 🧠 Discover supported resources dynamically from the cluster
- ⚡ Reuse cached output with TTL control
- 🚀 Install via Krew, GitHub Releases, or source build

## 🚀 Installation

### Krew

```bash
kubectl krew install ksearch
```

### GitHub Releases

Download the latest archive from:
`https://github.com/arush-sal/ksearch/releases`

### Build from source

```bash
git clone https://github.com/arush-sal/ksearch.git
cd ksearch
make build
./ksearch --help
```

## ⚙️ Quick Start

```bash
kubectl ksearch
kubectl ksearch -n kube-system
kubectl ksearch -n kube-system -p nginx
```

## 🎯 Examples

Search within one namespace:

```bash
kubectl ksearch -n default
```

Find resources related to one workload:

```bash
kubectl ksearch -n prod -p nginx
```

Limit output to selected kinds:

```bash
kubectl ksearch -n prod -k deployment,service,configmap,secret
```

Skip cache for a fresh read:

```bash
kubectl ksearch --no-cache
```

## 🛠 Flags

| Flag              | Description                                   |
|-------------------|-----------------------------------------------|
| `-n, --namespace` | Namespace to search                           |
| `-p, --pattern`   | Match resource names by substring             |
| `-k, --kinds`     | Comma-separated kinds or resources to include |
| `--cache-ttl`     | Cache TTL, defaults to `1m`                   |
| `--no-cache`      | Skip cached output                            |

Environment override:

```bash
export KSEARCH_CACHE_TTL=30s
```

## 📝 Notes

- Resource discovery depends on what the current cluster exposes.
- Cached output is stored locally and reused until the TTL expires.
- Use `--no-cache` when you need a fully fresh read.

## 🤝 Contributing

Contributor-facing documentation lives in `CONTRIBUTION.md`.
