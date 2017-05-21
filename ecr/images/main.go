package main

import (
	"log"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
)

const keep = 500

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
	// repositories = []string{"portal-app"}
	for _, repository := range repositories {
		log.Println(strings.Repeat("=", 10), repository)
		// log.Println(repository)
		deleteUntaggedImages(repository)
		deleteOldImages(repository)
	}
}

func deleteOldImages(repo string) {
	log.Println("get images except latest", keep)
	// imageDetails := make([]ecr.ImageDetail, 0)
	paramsDII := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repo),
	}
	respDII, err := svc.DescribeImages(paramsDII)
	imageDetails := respDII.ImageDetails
	if err != nil {
		log.Println(err.Error())
	}
	for respDII.NextToken != nil {
		paramsDII = &ecr.DescribeImagesInput{
			RepositoryName: aws.String(repo),
			NextToken:      respDII.NextToken,
		}
		respDII, err = svc.DescribeImages(paramsDII)
		imageDetails = append(imageDetails, respDII.ImageDetails...)
	}
	log.Println("total number of images:", len(imageDetails))
	sort.SliceStable(imageDetails, func(i, j int) bool {
		return imageDetails[i].ImagePushedAt.Unix() > imageDetails[j].ImagePushedAt.Unix()
	})
	if len(imageDetails) <= keep {
		log.Println("image number under threshold")
	} else {
		log.Println("delete old images")
		imageIds := make([]*ecr.ImageIdentifier, 0)
		for _, imageDetail := range imageDetails[keep:] {
			log.Println(*imageDetail.ImagePushedAt)
			imageID := &ecr.ImageIdentifier{ImageDigest: imageDetail.ImageDigest, ImageTag: imageDetail.ImageTags[0]}
			imageIds = append(imageIds, imageID)
		}
		paramsBDI := &ecr.BatchDeleteImageInput{
			ImageIds:       imageIds,
			RepositoryName: aws.String(repo),
		}
		_, err = svc.BatchDeleteImage(paramsBDI)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func deleteUntaggedImages(repo string) {
	log.Println("get untagged images")
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
	if len(respLI.ImageIds) > 0 {
		log.Println(respLI.ImageIds)
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
	} else {
		log.Println("all images tagged")
	}
}
