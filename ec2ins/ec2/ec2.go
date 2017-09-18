package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Run ...
func Run(s *session.Session) {
	svc := ec2.New(s)
	instances, _ := describeInstances(svc)
	fmt.Println(len(instances.Reservations))
}

func describeInstances(client *ec2.EC2) (*ec2.DescribeInstancesOutput, error) {
	param := ec2.DescribeInstancesInput{Filters: []*ec2.Filter{
		&ec2.Filter{
			Name:   aws.String("instance-state-name"),
			Values: []*string{aws.String("stopped")}},
	}}
	return client.DescribeInstances(&param)
}
