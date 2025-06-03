ARG BUILDER_IMAGE

FROM ${BUILDER_IMAGE} AS builder

WORKDIR /build
COPY . .

RUN mkdir /artifacts
RUN make -C validator/ validator
RUN cp /build/validator/validator /artifacts/nvidia-validator
RUN cp /build/validator/plugin-workload-validation.yaml /artifacts/
RUN cp /build/validator/cuda-workload-validation.yaml /artifacts/

FROM registry.ddbuild.io/images/nvidia-cuda-sample:11.6

LABEL maintainers="Compute"

USER root

RUN ln -s /cuda-samples/vectorAdd /usr/bin/vectorAdd
RUN mkdir -p /var/nvidia/manifests

COPY --from=builder /artifacts/nvidia-validator  /usr/bin/nvidia-validator
COPY --from=builder /artifacts/plugin-workload-validation.yaml /var/nvidia/manifests
COPY --from=builder /artifacts/cuda-workload-validation.yaml /var/nvidia/manifests

ENV NVIDIA_DISABLE_REQUIRE="true"
ENV NVIDIA_VISIBLE_DEVICES=void
ENV NVIDIA_DRIVER_CAPABILITIES=utility,compute

ENTRYPOINT ["/bin/bash"]
