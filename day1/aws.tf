provider "aws" {
  region = "us-east-1" # Replace with your desired region
}

resource "random_id" "bucket_suffix" {
  byte_length = 8
}

resource "aws_s3_bucket" "my_bucket" {
  bucket = "ceesay-multi27"

  tags = {
    Name        = "My bucket"
    Environment = "Dev"
  }
}

resource "aws_vpc" "multi-cloud" {
  cidr_block = "172.16.0.0/16"

  tags = {
    Name = "tf-multi-cloud"
  }
}

resource "aws_subnet" "multi-cloud" {
  vpc_id            = aws_vpc.multi-cloud.id
  cidr_block        = "172.16.10.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name = "tf-multi-cloud"
  }
}


# IAM role for EC2
resource "aws_iam_role" "ec2_admin_role" {
  name = "ec2-admin-role"
  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = "sts:AssumeRole",
        Principal = {
          Service = "ec2.amazonaws.com"
        },
        Effect = "Allow",
        Sid    = "AllowMultiCloud"
      }
    ]
  })

  tags = {
    Name = "EC2 Admin Role"
  }
}

# IAM policy attachment (attaches the AdministratorAccess policy)
resource "aws_iam_role_policy_attachment" "ec2_admin_policy_attachment" {
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess" # AWS managed policy for full admin access
  role       = aws_iam_role.ec2_admin_role.name
}

resource "aws_instance" "multi-cloud" {
  ami                  = "ami-05b10e08d247fb927"
  instance_type        = "t2.micro"
  subnet_id            = aws_subnet.multi-cloud.id
  key_name             = "ec2_log"
  iam_instance_profile = aws_iam_instance_profile.ec2_admin_instance_profile.name


  tags = {
    Name = "tf-multi-cloud"
  }
}

# IAM instance profile
resource "aws_iam_instance_profile" "ec2_admin_instance_profile" {
  name = "ec2-admin-instance-profile"
  role = aws_iam_role.ec2_admin_role.name
}
