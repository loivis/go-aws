package main

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var env = "test"
var services = []string{
	"charge-drive-mongodb",
	"gireve-emp-mongodb",
	"share-charge-adapter-mongodb",
	"smarthome-mongodb",
	"stripe-payment-service-mongodb",
}

func main() {
	session, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
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
		}}
	dvOutput, err := svc.DescribeVolumes(&dvInput)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(dvOutput.Volumes)

	// take snapshot for volumes
	for _, volume := range dvOutput.Volumes[:] {
		isTarget := false
		// get tags from volume
		tags := []*ec2.Tag{}
		for _, tag := range volume.Tags {
			for _, service := range services {
				if *tag.Value == service {
					fmt.Println("# " + env + ": " + service)
					isTarget = true
				}
			}
			if match, _ := regexp.MatchString("aws:.*", *tag.Key); !match {
				tags = append(tags, tag)
			}
		}

		// fmt.Println(tags)

		if isTarget {
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
}
