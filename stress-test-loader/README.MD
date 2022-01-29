* stress-test-loader 

golang grpc service load and execute actual stress-test

* tmp cmd 
`aws ec2 describe-instances --query 'Reservations[*].Instances[*].{"public_ip":PublicIpAddress}' --output json > test.json` 