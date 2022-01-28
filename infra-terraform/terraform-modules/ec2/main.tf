
data "aws_ami" "centos7" {
  owners      = ["679593333241"]
  most_recent = true

  filter {
    name   = "name"
    values = ["CentOS Linux 7 x86_64 HVM EBS *"]
  }

  filter {
    name   = "architecture"
    values = ["x86_64"]
  }

  filter {
    name   = "root-device-type"
    values = ["ebs"]
  }
}

resource "aws_instance" "stress_test_loader" {
  count                       = var.instance_count
  ami                         = data.aws_ami.centos7.id
  instance_type               = var.stress_test_loader_instance_type
  key_name                    = var.key_name
  subnet_id                   = var.aws_subnet[0][0]
  associate_public_ip_address = true
  root_block_device {
    volume_size = 30
  }

  tags = {
    Name = "hmg"
  }
  vpc_security_group_ids = [aws_security_group.instance_sg.id]
}

resource "aws_security_group" "instance_sg" {
  description = "controls direct access to application instances"
  vpc_id      = var.aws_vpc_id
  name        = "ec2-instsg"

  ingress {
    protocol  = "tcp"
    from_port = 22
    to_port   = 22

    cidr_blocks = [
      "${var.admin_cidr_ingress}",
    ]
  }

  ingress {
    protocol  = "tcp"
    from_port = 3128
    to_port   = 3128

    cidr_blocks = [
      "${var.admin_cidr_ingress}",
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

}
