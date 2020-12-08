// Copyright © 2020 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/banzaicloud/kafka-operator/api/v1beta1"
)

func createMinimalKafkaClusterCR(name, namespace string) *v1beta1.KafkaCluster {
	return &v1beta1.KafkaCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.KafkaClusterSpec{
			ListenersConfig: v1beta1.ListenersConfig{
				ExternalListeners: []v1beta1.ExternalListenerConfig{
					{
						CommonListenerSpec: v1beta1.CommonListenerSpec{
							Name:          "test",
							ContainerPort: 9733,
						},
						ExternalStartingPort: 11202,
						HostnameOverride:     "test-host",
						AccessMethod:         corev1.ServiceTypeLoadBalancer,
					},
				},
				InternalListeners: []v1beta1.InternalListenerConfig{
					{
						CommonListenerSpec: v1beta1.CommonListenerSpec{
							Type:          "plaintext",
							Name:          "internal",
							ContainerPort: 29092,
						},
						UsedForInnerBrokerCommunication: true,
					},
					{
						CommonListenerSpec: v1beta1.CommonListenerSpec{
							Type:          "plaintext",
							Name:          "controller",
							ContainerPort: 29093,
						},
						UsedForInnerBrokerCommunication: false,
						UsedForControllerCommunication:  true,
					},
				},
			},
			BrokerConfigGroups: map[string]v1beta1.BrokerConfig{
				"default": {
					StorageConfigs: []v1beta1.StorageConfig{
						{
							MountPath: "/kafka-logs",
							PvcSpec: &corev1.PersistentVolumeClaimSpec{
								AccessModes: []corev1.PersistentVolumeAccessMode{
									corev1.ReadWriteOnce,
								},
								Resources: corev1.ResourceRequirements{
									Requests: map[corev1.ResourceName]resource.Quantity{
										corev1.ResourceStorage: resource.MustParse("10Gi"),
									},
								},
							},
						},
					},
				},
			},
			Brokers: []v1beta1.Broker{
				{
					Id:                0,
					BrokerConfigGroup: "default",
				},
			},
			ClusterImage: "ghcr.io/banzaicloud/kafka:2.13-2.6.0-bzc.1",
			ZKAddresses:  []string{},
			MonitoringConfig: v1beta1.MonitoringConfig{
				CCJMXExporterConfig: "custom_property: custom_value",
			},
			CruiseControlConfig: v1beta1.CruiseControlConfig{
				TopicConfig: &v1beta1.TopicConfig{
					Partitions:        7,
					ReplicationFactor: 2,
				},
				Config: "some.config=value",
			},
		},
	}
}

func waitForClusterRunningState(kafkaCluster *v1beta1.KafkaCluster, namespace string) {
	Eventually(func() (v1beta1.ClusterState, error) {
		createdKafkaCluster := &v1beta1.KafkaCluster{}
		err := k8sClient.Get(context.TODO(), types.NamespacedName{Name: kafkaCluster.Name, Namespace: namespace}, createdKafkaCluster)
		if err != nil {
			return v1beta1.KafkaClusterReconciling, err
		}
		return createdKafkaCluster.Status.State, nil
	}, 5*time.Second, 100*time.Millisecond).Should(Equal(v1beta1.KafkaClusterRunning))
}

func waitForClusterDeletion(kafkaCluster *v1beta1.KafkaCluster) {
	Eventually(func() error {
		createdKafkaCluster := &v1beta1.KafkaCluster{}
		err := k8sClient.Get(context.TODO(), types.NamespacedName{Name: kafkaCluster.Name, Namespace: kafkaCluster.Namespace}, createdKafkaCluster)
		if err == nil {
			return errors.New("cluster should be deleted")
		}
		if apierrors.IsNotFound(err) {
			return nil
		} else {
			return err
		}
	}, 5*time.Second, 100*time.Millisecond).Should(Succeed())
}
