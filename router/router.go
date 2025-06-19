package router

import (
	"dev-ops-in-golang/api"
	"dev-ops-in-golang/k8s"

	"github.com/gin-gonic/gin"
)

func SetupRouter(k8sClient *k8s.K8sClient) *gin.Engine {
	r := gin.Default()
	deploymentAPI := &api.DeploymentAPI{K8s: k8sClient}
	r.GET("/deployments/:namespace/list", deploymentAPI.GetDeploymentHandler)
	r.POST("/deployments", deploymentAPI.CreateDeploymentHandler)
	r.PATCH("/deployments", deploymentAPI.UpdateDeploymentHandler)
	r.DELETE("/deployments", deploymentAPI.DeleteDeploymentHandler)
	r.GET("/deployments/:namespace/:name/rollout-status", deploymentAPI.GetRolloutStatusHandler)
	return r
}
