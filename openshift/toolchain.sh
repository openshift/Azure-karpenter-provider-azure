#!/usr/bin/env bash
set -exuo pipefail

K8S_VERSION="${K8S_VERSION:=1.32.x}"
export GOFLAGS=""

TOOL_DEST="$(go env GOPATH)/bin"
KUBEBUILDER_ASSETS="/usr/local/kubebuilder/bin"

main() {
    tools
    kubebuilder
    gettrivy
    golangci_lint_custom
}

tools() {
    go install -mod=mod -modfile=./openshift/go.tools.mod tool

    # These tools have transitive dependency conflicts with other tools in the
    # shared go.tools.mod (e.g. go-openapi/swag, viper, charmbracelet/x/ansi),
    # so install them separately with their own isolated dep graphs.
    go install -mod=mod github.com/google/ko@v0.17.1
    go install -mod=mod github.com/Azure/aks-node-viewer/cmd/aks-node-viewer@v0.0.2-alpha
}

golangci_lint_custom() {
    TOOL_DEST="$(go env GOPATH)/bin"
    SCRIPT_DIR="$(dirname "$(realpath "$0")")"
    sed "s|\$TOOL_DEST|$TOOL_DEST|g" "$SCRIPT_DIR/../hack/custom-gcl.template.yml" > .custom-gcl.yml
    "$TOOL_DEST/golangci-lint" custom -v
    rm .custom-gcl.yml
}

 kubebuilder() {
    echo "[INF] Setting up kubebuilder binaries for Kubernetes ${K8S_VERSION}"
    mkdir -p "${TOOL_DEST}" "${KUBEBUILDER_ASSETS}"

    arch="$(go env GOARCH)"
    os="$(go env GOOS)"
    curl -fsSL "https://github.com/kubernetes-sigs/controller-runtime/releases/download/v0.22.3/setup-envtest-${os}-${arch}" --output "${TOOL_DEST}/setup-envtest"
    chmod +x "${TOOL_DEST}/setup-envtest"

    envtest_path="$("${TOOL_DEST}/setup-envtest" use -p path "${K8S_VERSION}" --arch="${arch}" --bin-dir="${KUBEBUILDER_ASSETS}")"
    ln -sf "${envtest_path}"/* "${KUBEBUILDER_ASSETS}"
    find "${KUBEBUILDER_ASSETS}"

    if [[ "${K8S_VERSION}" = "1.25.x" ]] && [[ "${os}" = "linux" ]]; then
        for binary in kube-apiserver kubectl; do
            rm -f "${KUBEBUILDER_ASSETS}/${binary}"
            curl -fsSL "https://dl.k8s.io/v1.25.16/bin/linux/${arch}/${binary}" --output "${KUBEBUILDER_ASSETS}/${binary}"
            chmod +x "${KUBEBUILDER_ASSETS}/${binary}"
        done
    fi
}

gettrivy() {
    if [[ -x "${TOOL_DEST}/trivy" ]] || command -v trivy >/dev/null 2>&1; then
        return
    fi

    trivy_version="0.69.3"
    case "$(go env GOARCH)" in
        amd64) trivy_platform="64bit" ;;
        arm64) trivy_platform="ARM64" ;;
        *) echo "[ERR] Unsupported architecture for Trivy: $(go env GOARCH)" >&2; return 1 ;;
    esac

    trivy_name="trivy_${trivy_version}_Linux-${trivy_platform}.tar.gz"
    trivy_archive="/tmp/${trivy_name}"
    trivy_checksums="/tmp/trivy_checksums.txt"
    curl -fsSL "https://github.com/aquasecurity/trivy/releases/download/v${trivy_version}/${trivy_name}" --output "${trivy_archive}"
    curl -fsSL "https://github.com/aquasecurity/trivy/releases/download/v${trivy_version}/trivy_${trivy_version}_checksums.txt" --output "${trivy_checksums}"
    (cd "$(dirname "${trivy_archive}")" && sha256sum --ignore-missing --check "${trivy_checksums}")
    tar -xzf "${trivy_archive}" -C "${TOOL_DEST}" trivy
    chmod +x "${TOOL_DEST}/trivy"
    rm -f "${trivy_archive}" "${trivy_checksums}"
}

main "$@"
