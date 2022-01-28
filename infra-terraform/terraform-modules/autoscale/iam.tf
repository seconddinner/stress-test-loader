resource "aws_iam_role" "stress_test_client_read_role" {
  name = "stress_test_client_read_role"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF

  tags = {
      tag-key = "stress_test"
  }
}

resource "aws_iam_instance_profile" "stress_test_client_read_profile" {
  name = "stress_test_client_read_profile"
  role = "${aws_iam_role.stress_test_client_read_role.name}"
}

resource "aws_iam_role_policy" "stress_test_client_read" {
  name = "stress_test_client_read_policy"
  role = "${aws_iam_role.stress_test_client_read_role.id}"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:GetObject"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:s3:::cube-loadtest/*"
    }
  ]
}
EOF
}