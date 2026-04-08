ARG BUILDER_IMAGE

FROM ${BUILDER_IMAGE} AS builder

ARG TARGETARCH
ARG KUBECTL_VERSION=v1.34.5-dd.1

RUN curl -Lfs https://github.com/DataDog/kubernetes/releases/download/${KUBECTL_VERSION}/kubernetes-server-linux-${TARGETARCH}.tar.gz -O
RUN tar -C /usr/local/bin/ --strip-components 3 --exclude '*.tar' --exclude '*.docker_tag' -xvzf kubernetes-server-linux-${TARGETARCH}.tar.gz kubernetes/server/bin/kubectl
RUN chmod 755 /usr/local/bin/kubectl
RUN go tool nm /usr/local/bin/kubectl | grep -E 'sig.FIPSOnly'

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor vendor

# Copy the go source
COPY cmd/ cmd/
COPY api/ api/
COPY controllers/ controllers/
COPY internal/ internal/

# Copy Makefile
COPY Makefile Makefile
COPY *.mk .

# Build
ARG VERSION="unknown"
ARG GIT_COMMIT="unknown"
RUN make cmds

FROM registry.ddbuild.io/images/nvidia-cuda-sample:12.9 AS sample-getter

RUN stat /cuda-samples/vectorAdd

FROM registry.ddbuild.io/images/nvidia-cuda-base:12.9.0

LABEL maintainers="Compute"

ENV NVIDIA_VISIBLE_DEVICES=void

ARG VERSION="unknown"
ARG GIT_COMMIT="unknown"

USER root

RUN apt-get update && apt-get install -y --no-install-recommends kmod

WORKDIR /
COPY --from=builder /workspace/gpu-operator /usr/bin/
COPY --from=builder /usr/local/bin/kubectl /usr/bin/kubectl
COPY --from=builder /workspace/nvidia-validator /usr/bin/
COPY --from=sample-getter /cuda-samples/vectorAdd /usr/bin/vectorAdd
COPY --from=sample-getter /usr/local/cuda/compat /usr/local/cuda/compat

COPY assets /opt/gpu-operator/
COPY manifests /opt/gpu-operator/manifests
COPY validator/manifests /opt/validator/manifests

COPY hack/must-gather.sh /usr/bin/gather

# Add CRD resource into the image for helm upgrades
COPY deployments/gpu-operator/crds/nvidia.com_clusterpolicies.yaml /opt/gpu-operator/nvidia.com_clusterpolicies.yaml
COPY deployments/gpu-operator/crds/nvidia.com_nvidiadrivers.yaml /opt/gpu-operator/nvidia.com_nvidiadrivers.yaml
COPY deployments/gpu-operator/charts/node-feature-discovery/crds/nfd-api-crds.yaml /opt/gpu-operator/nfd-api-crds.yaml

RUN useradd gpu-operator
USER gpu-operator

ENTRYPOINT ["/usr/bin/gpu-operator"]
