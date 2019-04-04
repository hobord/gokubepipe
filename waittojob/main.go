package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/hobord/gokubepipe/kubeclient"
	batch "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	namespaceName := flag.String("n", "", "namespace")
	jobName := flag.String("j", "", "job name")
	flag.Parse()
	if *jobName == "" {
		fmt.Println("Job name parameter is mandatory! Use -p=jobname ")
		os.Exit(2)
	}

	clientset := kubeclient.GetClientset()

	jobs, err := clientset.BatchV1().Jobs(*namespaceName).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var myJob batch.Job
	for _, job := range jobs.Items {
		if *jobName == job.ObjectMeta.GetName() {
			myJob = job
		}
	}

	if myJob.Status.Succeeded != 0 {
		if myJob.Status.Failed == 1 {
			panic("Job failed")
		}
		for myJob.Status.Succeeded != 1 {
			if myJob.Status.Failed == 1 {
				panic("Job failed")
			}
			time.Sleep(10 * time.Second)
		}
	} else {
		fmt.Println("not found")
	}
}
