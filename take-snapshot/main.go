package main

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var region = "eu-west-1"
var env = "test"
var services = []*string{
	aws.String("charge-drive-mongodb"),
	aws.String("gireve-emp-mongodb"),
	aws.String("share-charge-adapter-mongodb"),
	aws.String("smarthome-mongodb"),
	aws.String("stripe-payment-service-mongodb"),
}

func main() {
	session, err := session.NewSession(&aws.Config{Region: &region})
	if err != nil {
		fmt.Println(err)
	}
	svc := ec2.New(session)

	// get volumes
	dvInput := ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Env"),
				Values: []*string{aws.String(env)},
			},
			{
				Name:   aws.String("tag:Service"),
				Values: services,
			},
		}}
	dvOutput, err := svc.DescribeVolumes(&dvInput)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(len(dvOutput.Volumes))

	// take snapshot for volumes
	for _, volume := range dvOutput.Volumes[:] {
		// exclude tags with tag:Key "^aws:"
		tags := []*ec2.Tag{}
		for _, tag := range volume.Tags {
			if match, _ := regexp.MatchString("^aws:.*", *tag.Key); !match {
				tags = append(tags, tag)
			}
			if *tag.Key == "Name" {
				fmt.Println("# " + *tag.Value)
			}
		}
		// fmt.Println(tags)

		// take snapshot
		fmt.Println(*volume.VolumeId)
		csInput := ec2.CreateSnapshotInput{VolumeId: volume.VolumeId}
		csOutput, err := svc.CreateSnapshot(&csInput)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(*csOutput.SnapshotId)

		// tag snapshot
		ctInput := ec2.CreateTagsInput{Resources: []*string{csOutput.SnapshotId}, Tags: tags}
		_, err = svc.CreateTags(&ctInput)
		if err != nil {
			fmt.Println(err)
		}

	}
}
