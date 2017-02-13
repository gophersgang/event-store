FROM golang:1.7
RUN curl -sSL https://get.docker.com/ | sh

ENV CLOUDSDK_CORE_DISABLE_PROMPTS 1
ENV PATH /opt/google-cloud-sdk/bin:$PATH

USER root

RUN apt-get update -y
RUN apt-get install -y jq
RUN curl https://sdk.cloud.google.com | bash
RUN mv /root/google-cloud-sdk /opt
RUN gcloud components install kubectl

COPY . /go/src/github.com/vendasta/gosdks
RUN go install github.com/vendasta/gosdks/tools/mscli
CMD ["/bin/bash"]
