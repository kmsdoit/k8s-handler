package api

import (
	"context"
	"net/http"

	"dev-ops-in-golang/helper"
	"dev-ops-in-golang/k8s"
	"dev-ops-in-golang/models"

	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentAPI struct {
	K8s *k8s.K8sClient
}

func (api *DeploymentAPI) CreateDeploymentHandler(c *gin.Context) {
	var req models.DeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deployment := helper.ConvertToDeployment(req)
	result, err := api.K8s.CreateDeployment(req.Namespace, deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Service 생성
	if req.Service != nil {
		service := helper.ConvertToService(req)
		if service != nil {
			_, err := api.K8s.Clientset.CoreV1().Services(req.Namespace).Create(context.Background(), service, metav1.CreateOptions{})
			if err != nil {
				// Deployment는 성공했지만 Service는 실패한 경우
				c.JSON(http.StatusCreated, gin.H{"message": "Deployment created, but failed to create service", "deployment": result.Name, "service_error": err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment created", "deployment": result.Name})
}

func (api *DeploymentAPI) DeleteDeploymentHandler(c *gin.Context) {
	var req models.DeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Deployment 삭제
	err := api.K8s.Clientset.AppsV1().Deployments(req.Namespace).Delete(context.Background(), req.Name, metav1.DeleteOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete deployment: " + err.Error()})
		return
	}

	svcErr := api.K8s.Clientset.CoreV1().Services(req.Namespace).Delete(context.Background(), req.Name, metav1.DeleteOptions{})
	if svcErr != nil {
		// Service가 없어서 NotFound가 나는 경우는 무시
		if !errors.IsNotFound(svcErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Deployment deleted, but failed to delete service: " + svcErr.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment and service deleted"})
}

func (api *DeploymentAPI) GetRolloutStatusHandler(c *gin.Context) {
	ns := c.Param("namespace")
	name := c.Param("name")
	deployment, err := api.K8s.Clientset.AppsV1().Deployments(ns).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	available := deployment.Status.AvailableReplicas
	desired := *deployment.Spec.Replicas
	c.JSON(200, gin.H{
		"available": available,
		"desired":   desired,
		"updated":   deployment.Status.UpdatedReplicas,
		"ready":     deployment.Status.ReadyReplicas,
	})
}

func (api *DeploymentAPI) GetDeploymentHandler(c *gin.Context) {
	ns := c.Param("namespace")
	deployment, err := api.K8s.Clientset.AppsV1().Deployments(ns).List(context.Background(), v1.ListOptions{})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"deployment.items.spec": deployment.Items[0].Spec})
}

func (api *DeploymentAPI) UpdateDeploymentHandler(c *gin.Context) {
	var req models.DeploymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deployment := helper.ConvertToDeployment(req)

	existing, err := api.K8s.Clientset.AppsV1().Deployments(req.Namespace).Get(context.Background(), req.Name, v1.GetOptions{})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found: " + err.Error()})
		return
	}
	deployment.ResourceVersion = existing.ResourceVersion

	result, err := api.K8s.Clientset.AppsV1().Deployments(req.Namespace).Update(context.Background(), deployment, v1.UpdateOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deployment updated", "deployment": result.Name})
}
