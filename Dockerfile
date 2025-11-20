ARG BUILDER_IMAGE

FROM ${BUILDER_IMAGE} AS builder

ARG TARGETARCH

RUN curl -o /usr/bin/kubectl -L "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/${TARGETARCH}/kubectl";
RUN chmod a+x /usr/bin/kubectl

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
COPY --from=builder /usr/bin/kubectl /usr/bin/kubectl
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
