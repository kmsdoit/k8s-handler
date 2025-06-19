package models

type ServiceSpec struct {
	Port       int32  `json:"port" binding:"required"`
	TargetPort int32  `json:"targetPort" binding:"required"`
	Type       string `json:"type" binding:"required"` // NodePort, ClusterIP, LoadBalancer
	NodePort   *int32 `json:"nodePort,omitempty"`      // NodePort 타입일 때만 사용
}

type DeploymentRequest struct {
	Name      string       `json:"name" binding:"required"`
	Namespace string       `json:"namespace" binding:"required"`
	Image     string       `json:"image" binding:"required"`
	Replicas  int32        `json:"replicas" binding:"required"`
	Service   *ServiceSpec `json:"service,omitempty"`
}

type DeleteDeploymentRequest struct {
	Namespace string `json:"namespace" binding:"required"`
	Name      string `json:"name" binding:"required"`
}
