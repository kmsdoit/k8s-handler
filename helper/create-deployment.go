package helper

import (
	"dev-ops-in-golang/models"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func ConvertToDeployment(req models.DeploymentRequest) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &req.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": req.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": req.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  req.Name,
						Image: req.Image,
						Ports: func() []corev1.ContainerPort {
							if req.Service != nil {
								return []corev1.ContainerPort{{
									ContainerPort: req.Service.TargetPort,
								}}
							}
							return nil
						}(),
					}},
				},
			},
		},
	}

	return deployment
}

func ConvertToService(req models.DeploymentRequest) *corev1.Service {
	if req.Service == nil {
		return nil
	}
	serviceType := corev1.ServiceType(req.Service.Type)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: req.Name,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": req.Name},
			Ports: []corev1.ServicePort{{
				Port:       req.Service.Port,
				TargetPort: intstrFromInt32(req.Service.TargetPort),
			}},
			Type: serviceType,
		},
	}
	if serviceType == corev1.ServiceTypeNodePort && req.Service.NodePort != nil {
		service.Spec.Ports[0].NodePort = *req.Service.NodePort
	}
	return service
}

func intstrFromInt32(i int32) intstr.IntOrString {
	return intstr.FromInt(int(i))
}
