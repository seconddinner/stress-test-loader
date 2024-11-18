FROM golang:1.23.3
LABEL maintainer="Seconddinner Engineering Team"

ENV PACKER_VERSION=1.9.1
ENV PACKER_SHA256SUM=793ed62255b9e572eda0c77d2a770f5fde501314b7598320786f1e51feb260d6

RUN apt-get update; \ 
    apt-get install -y --no-install-recommends git bash wget openssl protobuf-compiler unzip zip jq

ADD https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_linux_amd64.zip ./
ADD https://releases.hashicorp.com/packer/${PACKER_VERSION}/packer_${PACKER_VERSION}_SHA256SUMS ./

RUN sed -i '/.*linux_amd64.zip/!d' packer_${PACKER_VERSION}_SHA256SUMS
RUN unzip packer_${PACKER_VERSION}_linux_amd64.zip -d /bin
RUN rm -f packer_${PACKER_VERSION}_linux_amd64.zip

ENV TERRAFORM_VERSION=1.4.6
ENV TERRAFORM_SHA256SUM=e079db1a8945e39b1f8ba4e513946b3ab9f32bd5a2bdf19b9b186d22c5a3d53b

ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip ./
ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_SHA256SUMS ./

RUN sed -i '/.*linux_amd64.zip/!d' terraform_${TERRAFORM_VERSION}_SHA256SUMS
RUN unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /bin
RUN rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
ENV PATH="/go/bin:${PATH}"

ENV PULUMI_VERSION=v3.139.0
ENV PULUMI_B3SUMS=3f28f7e315525ebab264097b5f8be473c0d28fffa93c11299bddf5b90ffa147b

ADD https://github.com/pulumi/pulumi/releases/download/${PULUMI_VERSION}/pulumi-${PULUMI_VERSION}-linux-x64.tar.gz ./
RUN tar xzf pulumi-${PULUMI_VERSION}-linux-x64.tar.gz -C /bin
RUN rm -f pulumi-${PULUMI_VERSION}-linux-x64.tar.gz

ENV PATH="/bin/pulumi:${PATH}"

ADD https://packages.microsoft.com/config/ubuntu/20.04/packages-microsoft-prod.deb ./
RUN dpkg -i packages-microsoft-prod.deb ; rm packages-microsoft-prod.deb
RUN apt-get update && apt-get install -y dotnet-sdk-9.0

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install