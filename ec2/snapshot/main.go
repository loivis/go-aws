package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	ec2 := ec2.New(sess, &aws.Config{Region: aws.String("eu-west-1")})
	respV, err := ec2.DescribeVolumes(nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("> number of volumes:", len(respV.Volumes))

	volumes := make(map[string]bool, len(respV.Volumes))

	for _, volume := range respV.Volumes {
		// fmt.Println(*volume.VolumeId)
		volumes[*volume.VolumeId] = true
	}

	fmt.Println("> print volumes")
	for volume := range volumes {
		fmt.Println(volume)
	}

	respSS, err := ec2.DescribeSnapshots(nil)
	if err != nil {
		panic(err)
	}

	// snapshots := make([]string, len(respSS.Snapshots))

	for _, snapshot := range respSS.Snapshots {
		if volumes[*snapshot.VolumeId] {
			fmt.Println("> exists", *snapshot.VolumeId)
		} else {
			fmt.Println("< not exists", *snapshot.VolumeId)
		}
	}
}
