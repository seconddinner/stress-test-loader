FROM golang:1.19
LABEL maintainer="Seconddinner Engineering Team"

ENV PACKER_VERSION=1.7.10
ENV PACKER_SHA256SUM=1c8c176dd30f3b9ec3b418f8cb37822261ccebdaf0b01d9b8abf60213d1205cb

RUN apt-get update; \ 
    apt-get install -y --no-install-recommends git bash wget openssl protobuf-compiler unzip zip

ADD https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_amd64.zip ./
ADD https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_SHA256SUMS ./

RUN sed -i '/.*linux_amd64.zip/!d' packer_${PACKER_VERSION}_SHA256SUMS
RUN unzip packer_${PACKER_VERSION}_linux_amd64.zip -d /bin
RUN rm -f packer_${PACKER_VERSION}_linux_amd64.zip

ENV TERRAFORM_VERSION=1.3.2
ENV TERRAFORM_SHA256SUM=ce1a8770aaf27736a3352c5c31e95fb10d0944729b9d81013bf6848f8657da5f

ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip ./
ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_SHA256SUMS ./

RUN sed -i '/.*linux_amd64.zip/!d' terraform_${TERRAFORM_VERSION}_SHA256SUMS
RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /bin
RUN rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
ENV PATH="/go/bin:${PATH}"

ENV PULUMI_VERSION=v3.51.1
ENV PULUMI_SHA256SUM=24acebecc361f63f3884e96d0e2761c73c782561f600c58fbb440e88467d61f78337f4ebced571ab8e490138b768dac507de1a0cbd8481f427cf5e5ee42f1683

ADD https://github.com/pulumi/pulumi/releases/download/${PULUMI_VERSION}/pulumi-${PULUMI_VERSION}-linux-x64.tar.gz ./
RUN tar xzf pulumi-${PULUMI_VERSION}-linux-x64.tar.gz -C /bin
RUN rm -f pulumi-${PULUMI_VERSION}-linux-x64.tar.gz

ENV PATH="/bin/pulumi:${PATH}"

ADD https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb ./
RUN dpkg -i packages-microsoft-prod.deb ; rm packages-microsoft-prod.deb
RUN apt-get update && apt-get install -y dotnet-sdk-6.0

