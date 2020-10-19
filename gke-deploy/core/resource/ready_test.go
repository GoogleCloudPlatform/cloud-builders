package resource

import (
	"context"
	"testing"
)

func TestIsReady(t *testing.T) {
	ctx := context.Background()

	testDaemonsetReadyFile := "testing/daemonset-ready.yaml"
	testDaemonsetUnreadyFile := "testing/daemonset-unready.yaml"
	testDaemonsetUnready2File := "testing/daemonset-unready-2.yaml"
	testDaemonsetUnready3File := "testing/daemonset-unready-3.yaml"
	testDaemonsetUnready4File := "testing/daemonset-unready-4.yaml"
	testDeploymentReadyFile := "testing/deployment-ready.yaml"
	testDeploymentUnreadyFile := "testing/deployment-unready.yaml"
	testDeploymentUnready2File := "testing/deployment-unready-2.yaml"
	testDeploymentUnready3File := "testing/deployment-unready-3.yaml"
	testDeploymentUnready4File := "testing/deployment-unready-4.yaml"
	testDeploymentUnready5File := "testing/deployment-unready-5.yaml"
	testDeploymentUnready6File := "testing/deployment-unready-6.yaml"
	testDeploymentUnready7File := "testing/deployment-unready-7.yaml"
	testDeploymentUnready8File := "testing/deployment-unready-8.yaml"
	testDeploymentUnready9File := "testing/deployment-unready-9.yaml"
	testPvcReadyFile := "testing/pvc-ready.yaml"
	testPvcUnreadyFile := "testing/pvc-unready.yaml"
	testPodReadyFile := "testing/pod-ready.yaml"
	testPodReady2File := "testing/pod-ready-2.yaml"
	testPodUnreadyFile := "testing/pod-unready.yaml"
	testPodUnready2File := "testing/pod-unready-2.yaml"
	testPdbReadyFile := "testing/pdb-ready.yaml"
	testPdbReady2File := "testing/pdb-ready-2.yaml"
	testPdbUnreadyFile := "testing/pdb-unready.yaml"
	testPdbUnready2File := "testing/pdb-unready-2.yaml"
	testReplicasetReadyFile := "testing/replicaset-ready.yaml"
	testReplicasetUnreadyFile := "testing/replicaset-unready.yaml"
	testReplicasetUnready2File := "testing/replicaset-unready-2.yaml"
	testReplicasetUnready3File := "testing/replicaset-unready-3.yaml"
	testReplicasetUnready4File := "testing/replicaset-unready-4.yaml"
	testReplicasetUnready5File := "testing/replicaset-unready-5.yaml"
	testReplicationcontrollerReadyFile := "testing/replicationcontroller-ready.yaml"
	testReplicationcontrollerUnreadyFile := "testing/replicationcontroller-unready.yaml"
	testReplicationcontrollerUnready2File := "testing/replicationcontroller-unready-2.yaml"
	testReplicationcontrollerUnready3File := "testing/replicationcontroller-unready-3.yaml"
	testReplicationcontrollerUnready4File := "testing/replicationcontroller-unready-4.yaml"
	testReplicationcontrollerUnready5File := "testing/replicationcontroller-unready-5.yaml"
	testServiceReadyFile := "testing/service-ready.yaml"
	testServiceReady2File := "testing/service-ready-2.yaml"
	testServiceReady3File := "testing/service-ready-3.yaml"
	testServiceReady4File := "testing/service-ready-4.yaml"
	testServiceUnreadyFile := "testing/service-unready.yaml"
	testServiceUnready2File := "testing/service-unready-2.yaml"
	testServiceUnready3File := "testing/service-unready-3.yaml"
	testStatefulsetReadyFile := "testing/statefulset-ready.yaml"
	testStatefulsetUnreadyFile := "testing/statefulset-unready.yaml"
	testStatefulsetUnready2File := "testing/statefulset-unready-2.yaml"
	testStatefulsetUnready3File := "testing/statefulset-unready-3.yaml"
	testStatefulsetUnready4File := "testing/statefulset-unready-4.yaml"
	testStatefulsetUnready5File := "testing/statefulset-unready-5.yaml"
	testHpaFile := "testing/hpa.yaml"

	tests := []struct {
		name string

		obj *Object

		want bool
	}{{
		name: "DaemonSet is ready",

		obj: newObjectFromFile(t, testDaemonsetReadyFile),

		want: true,
	}, {
		name: "DaemonSet is not ready, status.numberReady != status.desiredNumberScheduled",

		obj: newObjectFromFile(t, testDaemonsetUnreadyFile),

		want: false,
	}, {
		name: "DaemonSet is not ready, status.numberAvailable != status.desiredNumberScheduled",

		obj: newObjectFromFile(t, testDaemonsetUnready2File),

		want: false,
	}, {
		name: "DaemonSet is not ready, status.numberAvailable is empty",

		obj: newObjectFromFile(t, testDaemonsetUnready3File),

		want: false,
	}, {
		name: "DaemonSet is not ready, status.observedGeneration != metadata.generation",

		obj: newObjectFromFile(t, testDaemonsetUnready4File),

		want: false,
	}, {
		name: "Deployment is ready",

		obj: newObjectFromFile(t, testDeploymentReadyFile),

		want: true,
	}, {
		name: "Deployment is not ready, Available condition status is False",

		obj: newObjectFromFile(t, testDeploymentUnreadyFile),

		want: false,
	}, {
		name: "Deployment is not ready, ReplicaFailure condition status is True",

		obj: newObjectFromFile(t, testDeploymentUnready2File),

		want: false,
	}, {
		name: "Deployment is not ready, Progressing condition with ReplicaSetUpdated reason status is False",

		obj: newObjectFromFile(t, testDeploymentUnready3File),

		want: false,
	}, {
		name: "Deployment is not ready, status.readyReplicas != spec.replicas",

		obj: newObjectFromFile(t, testDeploymentUnready4File),

		want: false,
	}, {
		name: "Deployment is not ready, status.availableReplicas != spec.replicas",

		obj: newObjectFromFile(t, testDeploymentUnready5File),

		want: false,
	}, {
		name: "Deployment is not ready, status.conditions is empty",

		obj: newObjectFromFile(t, testDeploymentUnready6File),

		want: false,
	}, {
		name: "Deployment is not ready, status has no fields",

		obj: newObjectFromFile(t, testDeploymentUnready7File),

		want: false,
	}, {
		name: "Deployment is not ready, status.replicas != spec.replicas",

		obj: newObjectFromFile(t, testDeploymentUnready8File),

		want: false,
	}, {
		name: "Deployment is not ready, status.observedGeneration != metadata.generation",

		obj: newObjectFromFile(t, testDeploymentUnready9File),

		want: false,
	}, {
		name: "PersistentVolumeClaim is ready",

		obj: newObjectFromFile(t, testPvcReadyFile),

		want: true,
	}, {
		name: "PersistentVolumeClaim is not ready, status.phase is Pending",

		obj: newObjectFromFile(t, testPvcUnreadyFile),

		want: false,
	}, {
		name: "Pod is ready, Ready condition status is True",

		obj: newObjectFromFile(t, testPodReadyFile),

		want: true,
	}, {
		name: "Pod is ready, Ready condition reason is PodCompleted",

		obj: newObjectFromFile(t, testPodReady2File),

		want: true,
	}, {
		name: "Pod is not ready, No Ready condition with status True or reason PodCompleted",

		obj: newObjectFromFile(t, testPodUnreadyFile),

		want: false,
	}, {
		name: "Pod is not ready, status.conditions is empty",

		obj: newObjectFromFile(t, testPodUnready2File),

		want: false,
	}, {
		name: "PodDisruptionBudget is ready, status.currentHealthy > status.desiredHealthy",

		obj: newObjectFromFile(t, testPdbReadyFile),

		want: true,
	}, {
		name: "PodDisruptionBudget is ready, status.currentHealthy == status.desiredHealthy",

		obj: newObjectFromFile(t, testPdbReady2File),

		want: true,
	}, {
		name: "PodDisruptionBudget is not ready, status.currentHealthy < status.desiredHealthy",

		obj: newObjectFromFile(t, testPdbUnreadyFile),

		want: false,
	}, {
		name: "PodDisruptionBudget is not ready, status.observedGeneration != metadata.generation",

		obj: newObjectFromFile(t, testPdbUnready2File),

		want: false,
	}, {
		name: "ReplicaSet is ready",

		obj: newObjectFromFile(t, testReplicasetReadyFile),

		want: true,
	}, {
		name: "ReplicaSet is not ready, status.readyReplicas != spec.replicas",

		obj: newObjectFromFile(t, testReplicasetUnreadyFile),

		want: false,
	}, {
		name: "ReplicaSet is not ready, status.availableReplicas != spec.replicas",

		obj: newObjectFromFile(t, testReplicasetUnready2File),

		want: false,
	}, {
		name: "ReplicaSet is not ready, status.readyReplicas and status.availableReplicas are empty",

		obj: newObjectFromFile(t, testReplicasetUnready3File),

		want: false,
	}, {
		name: "ReplicaSet is not ready, status.replicas != spec.replicas",

		obj: newObjectFromFile(t, testReplicasetUnready4File),

		want: false,
	}, {
		name: "ReplicaSet is not ready, status.observedGeneration != metadata.generation",

		obj: newObjectFromFile(t, testReplicasetUnready5File),

		want: false,
	}, {
		name: "ReplicationController is ready",

		obj: newObjectFromFile(t, testReplicationcontrollerReadyFile),

		want: true,
	}, {
		name: "ReplicationController is not ready, status.replicas != spec.replicas",

		obj: newObjectFromFile(t, testReplicationcontrollerUnreadyFile),

		want: false,
	}, {
		name: "ReplicationController is not ready, status.availableReplicas != spec.replicas",

		obj: newObjectFromFile(t, testReplicationcontrollerUnready2File),

		want: false,
	}, {
		name: "ReplicationController is not ready, status.readyReplicas != spec.replicas",

		obj: newObjectFromFile(t, testReplicationcontrollerUnready3File),

		want: false,
	}, {
		name: "ReplicationController is not ready, status.readyReplicas and status.availableReplicas are empty",

		obj: newObjectFromFile(t, testReplicationcontrollerUnready4File),

		want: false,
	}, {
		name: "ReplicationController is not ready, status.observedGeneration != metadata.generation",

		obj: newObjectFromFile(t, testReplicationcontrollerUnready5File),

		want: false,
	}, {
		name: "Service with LoadBalancer type is ready",

		obj: newObjectFromFile(t, testServiceReadyFile),

		want: true,
	}, {
		name: "Service with ClusterIP type is ready",

		obj: newObjectFromFile(t, testServiceReady2File),

		want: true,
	}, {
		name: "Service with NodePort type is ready",

		obj: newObjectFromFile(t, testServiceReady3File),

		want: true,
	}, {
		name: "Service with ExternalName type is ready",

		obj: newObjectFromFile(t, testServiceReady4File),

		want: true,
	}, {
		name: "Service is not ready, LoadBalancer type and status.loadBalancer.ingress is empty",

		obj: newObjectFromFile(t, testServiceUnreadyFile),

		want: false,
	}, {
		name: "Service is not ready, LoadBalancer type and status.loadBalancer.ingress has item with empty ip",

		obj: newObjectFromFile(t, testServiceUnready2File),

		want: false,
	}, {
		name: "Service is not ready, LoadBalancer type and spec.clusterIP is empty",

		obj: newObjectFromFile(t, testServiceUnready3File),

		want: false,
	}, {
		name: "StatefulSet is ready",

		obj: newObjectFromFile(t, testStatefulsetReadyFile),

		want: true,
	}, {
		name: "StatefulSet is not ready, status.readyReplicas != spec.replicas",

		obj: newObjectFromFile(t, testStatefulsetUnreadyFile),

		want: false,
	}, {
		name: "StatefulSet is not ready, status.currentReplicas != spec.replicas",

		obj: newObjectFromFile(t, testStatefulsetUnready2File),

		want: false,
	}, {
		name: "StatefulSet is not ready, status.readyReplicas and status.currentReplicas are empty",

		obj: newObjectFromFile(t, testStatefulsetUnready3File),

		want: false,
	}, {
		name: "StatefulSet is not ready, status.replicas != spec.replicas",

		obj: newObjectFromFile(t, testStatefulsetUnready4File),

		want: false,
	}, {
		name: "StatefulSet is not ready, status.observedGeneration != metadata.generation",

		obj: newObjectFromFile(t, testStatefulsetUnready5File),

		want: false,
	}, {
		name: "Default kind is always ready",

		obj: newObjectFromFile(t, testHpaFile),

		want: true,
	}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got, err := IsReady(ctx, tc.obj); got != tc.want || err != nil {
				t.Errorf("IsReady(ctx, %v) = %t, %v; want %t, <nil>", tc.obj, got, err, tc.want)
			}
		})
	}
}
