packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "PNS_version" {
  type = string
}

source "amazon-ebs" "ubuntu-us-west-2" {
  ami_name              = var.PNS_version
  instance_type         = "t4g.medium"
  force_deregister      = "true"
  force_delete_snapshot = "true"
  region                = "us-west-2"
  source_ami_filter {
    filters = {
      name                = "ubuntu/images/hvm-ssd/ubuntu-focal-20.04-arm64-server-*"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
      architecture        = "arm64"
    }
    most_recent = true
    owners      = ["099720109477"]
  }
  ssh_username = "ubuntu"
}
build {
  name    = "add-files-and-update-packer"
  sources = ["source.amazon-ebs.ubuntu-us-west-2"]
  provisioner "file" {
    source      = "stress-test-loader-linux.gz"
    destination = "/tmp/stress-test-loader-linux.gz"
  }
  provisioner "file" {
    source      = "cicd/config"
    destination = "/tmp/config"
  }
  provisioner "shell" {
    inline = [
      "while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for cloud-init...'; sleep 1; done",
      "sudo rm -rf /var/lib/apt/lists/partial/",
      "sudo apt-get -y -q update",
      "sudo apt-get -y -q upgrade ",
      "sudo apt-get install -y prometheus-node-exporter",

      "sudo mkdir -p /usr/local/stress-test-loader",
      "sudo gzip -d -f /tmp/stress-test-loader-linux.gz",
      "sudo mv /tmp/stress-test-loader-linux /usr/local/stress-test-loader/",
      "sudo mv /tmp/config/systemd/stress-test-loader.service /etc/systemd/system",
      "sudo mv /tmp/config/config.json /usr/local/stress-test-loader",

      "sudo systemctl daemon-reload",
      "sudo systemctl enable prometheus-node-exporter",
      "sudo systemctl enable stress-test-loader.service",
      "sudo apt-get clean"
    ]
  }
}