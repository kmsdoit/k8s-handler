package k8s

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	Clientset *kubernetes.Clientset
}

func NewClient() (*K8sClient, error) {
	kubeconfig := clientcmd.RecommendedHomeFile
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8sClient{Clientset: clientset}, nil
}

func (k *K8sClient) CreateDeployment(namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	return k.Clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
}

func (k *K8sClient) DeleteDeployment(namespace string, name string) error {
	return k.Clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}
