package main

import (
	"log"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

var sess = session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))
var svc = ecr.New(sess)

func main() {
	params := &ecr.DescribeRepositoriesInput{}
	resp, err := svc.DescribeRepositories(params)
	if err != nil {
		log.Println(err.Error())
		return
	}
	// log.Println(resp.Repositories)
	repositories := make([]string, 0)
	for _, repository := range resp.Repositories {
		repositories = append(repositories, *repository.RepositoryName)
	}
	sort.Sort(sort.StringSlice(repositories))
	log.Println("repositories:", repositories)
	for _, repository := range repositories {
		log.Println(repository)
		deleteUntaggedImages(repository)
	}
}

func deleteUntaggedImages(repo string) {
	log.Println("get untagged images")
	// imageIds := make([]string, 0)
	paramsLI := &ecr.ListImagesInput{
		RepositoryName: aws.String(repo),
		Filter: &ecr.ListImagesFilter{
			TagStatus: aws.String("UNTAGGED"),
		},
	}
	respLI, err := svc.ListImages(paramsLI)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(respLI.ImageIds)
	if len(respLI.ImageIds) > 0 {
		log.Println("delete untagged images")
		paramsBDI := &ecr.BatchDeleteImageInput{
			ImageIds:       respLI.ImageIds,
			RepositoryName: aws.String(repo),
		}
		respBDI, err := svc.BatchDeleteImage(paramsBDI)
		if err != nil {
			log.Println(err.Error())
		}
		log.Println(respBDI)
	}
}
