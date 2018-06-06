package session

import (
	"fmt"
	"io"
	"strconv"

	"bitbucket.org/linkernetworks/aurora/src/aurora/serviceprovider"
	"bitbucket.org/linkernetworks/aurora/src/deployment"
	"bitbucket.org/linkernetworks/aurora/src/net/http/response"
	restful "github.com/emicklei/go-restful"
	"github.com/linkernetworks/logger"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KubernetesDeploymentTargetHandler func(req *restful.Request, resp *restful.Response, sp *serviceprovider.Container, dt *deployment.KubeDeploymentTarget)

func NewKubernetesDeploymentTargetHandler(sp *serviceprovider.Container, h KubernetesDeploymentTargetHandler) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		if sp.Config.JobServer == nil {
			response.InternalServerError(req.Request, resp.ResponseWriter, fmt.Errorf("JobServer config is not defined."))
			return
		}

		var targetName = req.PathParameter("target")
		var targets = deployment.LoadDeploymentTargets(sp.Config.JobServer.DeploymentTargets, sp.Redis)
		dt, ok := targets[targetName]
		if !ok {
			response.Forbidden(req.Request, resp.ResponseWriter, fmt.Errorf("deployment target not found."))
			return
		}

		kdt, supported := dt.(*deployment.KubeDeploymentTarget)
		if !supported {
			response.Forbidden(req.Request, resp.ResponseWriter, fmt.Errorf("deployment target is not a kubernetes cluster."))
			return
		}

		h(req, resp, sp, kdt)
	}
}

func GetKubernetesJobPodNamesHandler(req *restful.Request, resp *restful.Response, sp *serviceprovider.Container, dt *deployment.KubeDeploymentTarget) {
	var namespace = req.PathParameter("namespace")
	var jobName = req.PathParameter("jobName")
	podList, err := dt.Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "job-name=" + jobName})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, fmt.Errorf("failed to get pods: %v", err))
		return
	}

	if len(podList.Items) == 0 {
		response.NotFound(req.Request, resp.ResponseWriter, fmt.Errorf("pod not found: %v", err))
		return
	}

	var names []string
	for _, pod := range podList.Items {
		names = append(names, pod.Name)
	}

	resp.WriteEntity(map[string]interface{}{
		"error": false,
		"names": names,
	})
}

func GetKubernetesJobContainerNamesHandler(req *restful.Request, resp *restful.Response, sp *serviceprovider.Container, dt *deployment.KubeDeploymentTarget) {
	var namespace = req.PathParameter("namespace")
	var jobName = req.PathParameter("jobName")

	job, err := dt.Clientset.BatchV1().Jobs(namespace).Get(jobName, metav1.GetOptions{})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, fmt.Errorf("failed to get job: %v", err))
		return
	}

	var names []string
	for _, container := range job.Spec.Template.Spec.Containers {
		names = append(names, container.Name)
	}

	resp.WriteEntity(map[string]interface{}{
		"error": false,
		"names": names,
	})
}

func GetKubernetesJobContainerLogHandler(req *restful.Request, resp *restful.Response, sp *serviceprovider.Container, dt *deployment.KubeDeploymentTarget) {
	var namespace = req.PathParameter("namespace")
	var jobName = req.PathParameter("jobName")
	var containerName = req.PathParameter("container")

	podList, err := dt.Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: "job-name=" + jobName})
	if err != nil {
		response.InternalServerError(req.Request, resp.ResponseWriter, fmt.Errorf("failed to get pods: %v", err))
		return
	}

	if len(podList.Items) == 0 {
		response.NotFound(req.Request, resp.ResponseWriter, fmt.Errorf("pod not found: %v", err))
		return
	}

	copyLogs := func(podName, containerName string) (err error) {
		var readCloser io.ReadCloser
		var options = v1.PodLogOptions{Container: containerName}
		if tail := req.QueryParameter("tail"); tail != "" {
			if tailLines, err := strconv.ParseInt(tail, 10, 64); err == nil {
				options.TailLines = &tailLines
			}
		}
		readCloser, err = dt.GetPodLogs(podName, &options)
		if err != nil {
			return err
		}
		_, err = io.Copy(resp, readCloser)
		return err
	}

	for _, pod := range podList.Items {
		if containerName == "" {
			for _, container := range pod.Spec.Containers {
				if err := copyLogs(pod.Name, container.Name); err != nil {
					logger.Errorf("failed to copy the container log: %v", err)
				}
			}
		} else {
			if err := copyLogs(pod.Name, containerName); err != nil {
				logger.Errorf("failed to copy the container log: %v", err)
			}
		}

	}
}
