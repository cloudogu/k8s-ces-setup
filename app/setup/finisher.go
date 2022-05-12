package setup

import (
	context2 "context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cloudogu/k8s-ces-setup/app/context"

	"k8s.io/client-go/kubernetes"
)

// Finisher finishes the k8s-ces-setup by setting the installed flag and removing itself from the cluster.
type Finisher struct {
	Client    kubernetes.Interface
	Namespace string
}

// NewFinisher creates a new Finisher.
func NewFinisher(client kubernetes.Interface, targetNamespace string) *Finisher {
	return &Finisher{Client: client, Namespace: targetNamespace}
}

// FinishSetup writes the installed flag into the setup config map and implodes by removing itself from the cluster.
func (f *Finisher) FinishSetup() error {
	err := f.writeInformationToClusterState()
	if err != nil {
		return err
	}

	err = f.removeK8sCesSetupFromCluster()
	if err != nil {
		return fmt.Errorf("failed to remove k8s-ces-setup from the cluster: %w", err)
	}

	return nil
}

func (f *Finisher) removeK8sCesSetupFromCluster() error {
	// TODO: this does currently not work as it is forbidden to delete own resources.
	//gracePeriod := int64(5)
	//
	//// ------- Cluster Role Bindings
	//logrus.Debug("Remove cluster role binding: k8s-ces-setup-cluster-resources")
	//err := f.Client.RbacV1().ClusterRoleBindings().Delete(context2.Background(), "k8s-ces-setup-cluster-resources", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete cluster role binding [%s]: %w", "k8s-ces-setup-cluster-resources", err)
	//}
	//logrus.Debug("Remove cluster role binding: k8s-ces-setup-cluster-non-resources")
	//err = f.Client.RbacV1().ClusterRoleBindings().Delete(context2.Background(), "k8s-ces-setup-cluster-non-resources", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete cluster role binding [%s]: %w", "k8s-ces-setup-cluster-non-resources", err)
	//}
	//
	//// ------- Cluster Roles
	//logrus.Debug("Remove cluster role: k8s-ces-setup-cluster-resources")
	//err = f.Client.RbacV1().ClusterRoles().Delete(context2.Background(), "k8s-ces-setup-cluster-resources", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete cluster role [%s]: %w", "k8s-ces-setup-cluster-resources", err)
	//}
	//logrus.Debug("Remove cluster role: k8s-ces-setup-cluster-non-resources")
	//err = f.Client.RbacV1().ClusterRoles().Delete(context2.Background(), "k8s-ces-setup-cluster-non-resources", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete cluster role [%s]: %w", "k8s-ces-setup-cluster-non-resources", err)
	//}
	//
	//// ------- Role Binding
	//logrus.Debug("Remove role binding: k8s-ces-setup")
	//err = f.Client.RbacV1().RoleBindings(f.Namespace).Delete(context2.Background(), "k8s-ces-setup", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete role binding [%s]: %w", "k8s-ces-setup", err)
	//}
	//
	//// ------- Role
	//logrus.Debug("Remove role: k8s-ces-setup")
	//err = f.Client.RbacV1().Roles(f.Namespace).Delete(context2.Background(), "k8s-ces-setup", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete role [%s]: %w", "k8s-ces-setup", err)
	//}
	//
	//// ------- Service Account
	//logrus.Debug("Remove service account: k8s-ces-setup")
	//err = f.Client.CoreV1().ServiceAccounts(f.Namespace).Delete(context2.Background(), "k8s-ces-setup", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete service account [%s]: %w", "k8s-ces-setup", err)
	//}
	//
	//// ------- Service
	//logrus.Debug("Remove service: k8s-ces-setup")
	//err = f.Client.CoreV1().Services(f.Namespace).Delete(context2.Background(), "k8s-ces-setup", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete service [%s]: %w", "k8s-ces-setup", err)
	//}
	//
	//// ------- Deployment
	//logrus.Debug("Remove deployment: k8s-ces-setup")
	//err = f.Client.AppsV1().Deployments(f.Namespace).Delete(context2.Background(), "k8s-ces-setup", metav1.DeleteOptions{GracePeriodSeconds: &gracePeriod})
	//if err != nil && !errors.IsNotFound(err) {
	//	return fmt.Errorf("failed to delete deployment [%s]: %w", "k8s-ces-setup", err)
	//}
	return nil
}

func (f *Finisher) writeInformationToClusterState() error {
	setupConfigMap, err := context.GetSetupConfigMap(f.Client, f.Namespace)
	if err != nil {
		return fmt.Errorf("failed to get setup config map: %w", err)
	}

	setupConfigMap.Data[context.SetupStateKey] = context.SetupStateInstalled
	setupConfigMap, err = f.Client.CoreV1().ConfigMaps(f.Namespace).Update(context2.Background(), setupConfigMap, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}
