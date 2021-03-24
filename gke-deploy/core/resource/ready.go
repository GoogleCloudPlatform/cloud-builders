package resource

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// IsReady returns true if a deployed object is ready. Please check the comments of each kind's
// implementation for a description of what is considered to be ready for that kind of object.
func IsReady(ctx context.Context, obj *Object) (bool, error) {
	kind := ObjectKind(obj)
	switch kind {
	case "DaemonSet":
		return daemonSetIsReady(ctx, obj)
	case "Deployment":
		return deploymentIsReady(ctx, obj)
	case "PersistentVolumeClaim":
		return persistentVolumeClaimIsReady(ctx, obj)
	case "Pod":
		return podIsReady(ctx, obj)
	case "PodDisruptionBudget":
		return podDisruptionBudgetIsReady(ctx, obj)
	case "ReplicaSet":
		return replicaSetIsReady(ctx, obj)
	case "ReplicationController":
		return replicationControllerIsReady(ctx, obj)
	case "Service":
		return serviceIsReady(ctx, obj)
	case "StatefulSet":
		return statefulSetIsReady(ctx, obj)
	default:
		return true, nil
	}
}

// daemonSetIsReady returns true if a deployed object with kind "DaemonSet" is ready.
// This returns true if the following bullets are true:
// * status.observedGeneration == metadata.generation
// * status.numberAvailable == status.desiredNumberScheduled
// * status.numberReady == status.desiredNumberScheduled
func daemonSetIsReady(ctx context.Context, obj *Object) (bool, error) {
	generation, ok, err := unstructured.NestedInt64(obj.Object, "metadata", "generation")
	if err != nil {
		return false, fmt.Errorf("failed to get metadata.generation field: %v", err)
	}
	if !ok {
		return false, nil
	}

	observedGeneration, ok, err := unstructured.NestedInt64(obj.Object, "status", "observedGeneration")
	if err != nil {
		return false, fmt.Errorf("failed to get status.observedGeneration field: %v", err)
	}
	if !ok || observedGeneration != generation {
		return false, nil
	}

	desiredNumberScheduled, ok, err := unstructured.NestedInt64(obj.Object, "status", "desiredNumberScheduled")
	if err != nil {
		return false, fmt.Errorf("failed to get status.desiredNumberScheduled field: %v", err)
	}
	if !ok {
		return false, nil
	}

	numberAvailable, ok, err := unstructured.NestedInt64(obj.Object, "status", "numberAvailable")
	if err != nil {
		return false, fmt.Errorf("failed to get status.numberAvailable field: %v", err)
	}
	if !ok || numberAvailable != desiredNumberScheduled {
		return false, nil
	}

	numberReady, ok, err := unstructured.NestedInt64(obj.Object, "status", "numberReady")
	if err != nil {
		return false, fmt.Errorf("failed to get status.numberReady field: %v", err)
	}
	if !ok || numberReady != desiredNumberScheduled {
		return false, nil
	}

	return true, nil
}

// deploymentIsReady returns true if a deployed object with kind "Deployment" is ready.
// This returns true if the following bullets are true:
// * status.observedGeneration == metadata.generation
// * status.replicas == spec.replicas
// * status.readyReplicas == spec.replicas
// * status.availableReplicas == spec.replicas
// * status.conditions is not empty
// * All items in status.conditions matches any:
//   * type == "Progressing" AND status == "True" AND reason == "NewReplicaSetAvailable"
//   * type == "Available" AND status == "True"
// * All items in status.conditions do not match any:
//   * type == "ReplicaFailure" AND status ==" True"
func deploymentIsReady(ctx context.Context, obj *Object) (bool, error) {
	generation, ok, err := unstructured.NestedInt64(obj.Object, "metadata", "generation")
	if err != nil {
		return false, fmt.Errorf("failed to get metadata.generation field: %v", err)
	}
	if !ok {
		return false, nil
	}

	observedGeneration, ok, err := unstructured.NestedInt64(obj.Object, "status", "observedGeneration")
	if err != nil {
		return false, fmt.Errorf("failed to get status.observedGeneration field: %v", err)
	}
	if !ok || observedGeneration != generation {
		return false, nil
	}

	specReplicas, ok, err := unstructured.NestedInt64(obj.Object, "spec", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get spec.replicas field: %v", err)
	}
	if !ok {
		return false, nil
	}

	statusReplicas, _, err := unstructured.NestedInt64(obj.Object, "status", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.replicas field: %v", err)
	}
	if statusReplicas != specReplicas {
		return false, nil
	}

	readyReplicas, _, err := unstructured.NestedInt64(obj.Object, "status", "readyReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.readyReplicas field: %v", err)
	}
	if readyReplicas != specReplicas {
		return false, nil
	}

	availableReplicas, _, err := unstructured.NestedInt64(obj.Object, "status", "availableReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.availableReplicas field: %v", err)
	}
	if availableReplicas != specReplicas {
		return false, nil
	}

	conditions, ok, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return false, fmt.Errorf("failed to get status.conditions field: %v", err)
	}
	if !ok || len(conditions) == 0 {
		return false, nil
	}
	for _, c := range conditions {
		cMap, ok := c.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("failed to convert conditions to map")
		}
		cType, ok, err := unstructured.NestedString(cMap, "type")
		if err != nil {
			return false, fmt.Errorf("failed to get type field: %v", err)
		}
		if !ok || cType == "" {
			return false, nil
		}

		switch cType {
		case "Available":
			status, ok, err := unstructured.NestedString(cMap, "status")
			if err != nil {
				return false, fmt.Errorf("failed to get status field: %v", err)
			}
			if !ok || status != "True" {
				return false, nil
			}
		case "Progressing":
			status, ok, err := unstructured.NestedString(cMap, "status")
			if err != nil {
				return false, fmt.Errorf("failed to get status field: %v", err)
			}
			if !ok || status == "" {
				return false, nil
			}
			reason, ok, err := unstructured.NestedString(cMap, "reason")
			if err != nil {
				return false, fmt.Errorf("failed to get reason field: %v", err)
			}
			if !ok || status != "True" || reason != "NewReplicaSetAvailable" {
				return false, nil
			}
		case "ReplicaFailure":
			status, ok, err := unstructured.NestedString(cMap, "status")
			if err != nil {
				return false, fmt.Errorf("failed to get status field: %v", err)
			}
			if !ok || status == "True" {
				return false, nil
			}
		default:
			return false, nil
		}
	}
	return true, nil
}

// persistentVolumeClaimIsReady returns true if a deployed object with kind "PersistentVolumeClaim" is ready.
// This returns true if the following bullets are true:
// * status.phase == "Bound"
func persistentVolumeClaimIsReady(ctx context.Context, obj *Object) (bool, error) {
	phase, ok, err := unstructured.NestedString(obj.Object, "status", "phase")
	if err != nil {
		return false, fmt.Errorf("failed to get status.phase field: %v", err)
	}
	if !ok || phase == "" {
		return false, nil
	}

	return phase == "Bound", nil
}

// podIsReady returns true if a deployed object with kind "Pod" is ready.
// This returns true if the following bullets are true:
// * status.conditions contains at least one item that matches any:
//   * type == "Ready" AND status == "True"
//   * type == "Ready" AND reason == "PodCompleted"
func podIsReady(ctx context.Context, obj *Object) (bool, error) {
	conditions, ok, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return false, fmt.Errorf("failed to get status.conditions field: %v", err)
	}
	if !ok || len(conditions) == 0 {
		return false, nil
	}
	for _, c := range conditions {
		cMap, ok := c.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("failed to convert conditions to map")
		}
		cType, ok, err := unstructured.NestedString(cMap, "type")
		if err != nil {
			return false, fmt.Errorf("failed to get type field: %v", err)
		}
		if !ok || cType == "" {
			return false, nil
		}

		switch cType {
		case "Ready":
			status, ok, err := unstructured.NestedString(cMap, "status")
			if err != nil {
				return false, fmt.Errorf("failed to get status field: %v", err)
			}
			if !ok {
				return false, nil
			}
			if status == "True" {
				return true, nil
			}

			reason, ok, err := unstructured.NestedString(cMap, "reason")
			if err != nil {
				return false, fmt.Errorf("failed to get reason field: %v", err)
			}
			if !ok {
				return false, nil
			}
			if reason == "PodCompleted" {
				return true, nil
			}
		default:
			// Skip
		}
	}

	return false, nil
}

// podDisruptionBudget returns true if a deployed object with kind "PodDisruptionBudget" is ready.
// This returns true if the following bullets are true:
// * status.observedGeneration == metadata.generation
// * status.currentHealthy >= status.desiredHealthy
func podDisruptionBudgetIsReady(ctx context.Context, obj *Object) (bool, error) {
	generation, ok, err := unstructured.NestedInt64(obj.Object, "metadata", "generation")
	if err != nil {
		return false, fmt.Errorf("failed to get metadata.generation field: %v", err)
	}
	if !ok {
		return false, nil
	}

	observedGeneration, ok, err := unstructured.NestedInt64(obj.Object, "status", "observedGeneration")
	if err != nil {
		return false, fmt.Errorf("failed to get status.observedGeneration field: %v", err)
	}
	if !ok || observedGeneration != generation {
		return false, nil
	}

	desiredHealthy, ok, err := unstructured.NestedInt64(obj.Object, "status", "desiredHealthy")
	if err != nil {
		return false, fmt.Errorf("failed to get status.desiredHealthy field: %v", err)
	}

	currentHealthy, ok, err := unstructured.NestedInt64(obj.Object, "status", "currentHealthy")
	if err != nil {
		return false, fmt.Errorf("failed to get status.currentHealthy field: %v", err)
	}
	if !ok || currentHealthy < desiredHealthy {
		return false, nil
	}

	return true, nil
}

// replicaSetIsReady returns true if a deployed object with kind "ReplicaSet" is ready.
// This returns true if the following bullets are true:
// * status.observedGeneration == metadata.generation
// * status.replicas == spec.replicas
// * status.readyReplicas == spec.replicas
// * status.availableReplicas == spec.replicas
// * All items in status.conditions do not match any:
//   * type == "ReplicaFailure" AND status == "True"
func replicaSetIsReady(ctx context.Context, obj *Object) (bool, error) {
	generation, ok, err := unstructured.NestedInt64(obj.Object, "metadata", "generation")
	if err != nil {
		return false, fmt.Errorf("failed to get metadata.generation field: %v", err)
	}
	if !ok {
		return false, nil
	}

	observedGeneration, ok, err := unstructured.NestedInt64(obj.Object, "status", "observedGeneration")
	if err != nil {
		return false, fmt.Errorf("failed to get status.observedGeneration field: %v", err)
	}
	if !ok || observedGeneration != generation {
		return false, nil
	}

	specReplicas, ok, err := unstructured.NestedInt64(obj.Object, "spec", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get spec.replicas field: %v", err)
	}
	if !ok {
		return false, nil
	}

	statusReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.replicas field: %v", err)
	}
	if !ok || statusReplicas != specReplicas {
		return false, nil
	}

	readyReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "readyReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.readyReplicas field: %v", err)
	}
	if !ok || readyReplicas != specReplicas {
		return false, nil
	}

	availableReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "availableReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.availableReplicas field: %v", err)
	}
	if !ok || availableReplicas != specReplicas {
		return false, nil
	}

	conditions, ok, err := unstructured.NestedSlice(obj.Object, "status", "conditions")
	if err != nil {
		return false, fmt.Errorf("failed to get status.conditions field: %v", err)
	}
	if !ok || len(conditions) == 0 {
		return true, nil
	}
	for _, c := range conditions {
		cMap, ok := c.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("failed to convert conditions to map")
		}
		cType, ok, err := unstructured.NestedString(cMap, "type")
		if err != nil {
			return false, fmt.Errorf("failed to get type field: %v", err)
		}
		if !ok || cType == "" {
			return false, nil
		}

		switch cType {
		case "ReplicaFailure":
			status, ok, err := unstructured.NestedString(cMap, "status")
			if err != nil {
				return false, fmt.Errorf("failed to get status field: %v", err)
			}
			if !ok || status == "True" {
				return false, nil
			}
		default:
			// Skip
		}
	}

	return true, nil
}

// replicationControllerIsReady returns true if a deployed object with kind "ReplicaSet" is ready.
// This returns true if the following bullets are true:
// * status.observedGeneration == metadata.generation
// * status.replicas == spec.replicas
// * status.readyReplicas == spec.replicas
// * status.availableReplicas == spec.replicas
func replicationControllerIsReady(ctx context.Context, obj *Object) (bool, error) {
	generation, ok, err := unstructured.NestedInt64(obj.Object, "metadata", "generation")
	if err != nil {
		return false, fmt.Errorf("failed to get metadata.generation field: %v", err)
	}
	if !ok {
		return false, nil
	}

	observedGeneration, ok, err := unstructured.NestedInt64(obj.Object, "status", "observedGeneration")
	if err != nil {
		return false, fmt.Errorf("failed to get status.observedGeneration field: %v", err)
	}
	if !ok || observedGeneration != generation {
		return false, nil
	}

	specReplicas, ok, err := unstructured.NestedInt64(obj.Object, "spec", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get spec.replicas field: %v", err)
	}
	if !ok {
		return false, nil
	}

	statusReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.replicas field: %v", err)
	}
	if !ok || statusReplicas != specReplicas {
		return false, nil
	}

	readyReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "readyReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.readyReplicas field: %v", err)
	}
	if !ok || readyReplicas != specReplicas {
		return false, nil
	}

	availableReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "availableReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.availableReplicas field: %v", err)
	}
	if !ok || availableReplicas != specReplicas {
		return false, nil
	}

	return true, nil
}

// serviceIsReady returns true if a deployed object with kind "Service" is ready.
// This returns true if the following bullets are true:
// * Any of the following are true
//   * type == "ClusterIP" (default)
//   * type == "NodePort"
//   * type == "ExternalName"
//   * type == "LoadBalancer" AND "spec.clusterIP" is not empty AND "status.loadBalancer.ingress" is
//     not empty AND all objects in "status.loadBalancer.ingress" has an "ip" that is not empty
func serviceIsReady(ctx context.Context, obj *Object) (bool, error) {
	serviceType, ok, err := unstructured.NestedString(obj.Object, "spec", "type")
	if err != nil {
		return false, fmt.Errorf("failed to get spec.type field: %v", err)
	}
	if !ok || serviceType == "" {
		return false, nil
	}

	if serviceType == "ClusterIP" || serviceType == "NodePort" || serviceType == "ExternalName" {
		return true, nil
	}

	clusterIP, ok, err := unstructured.NestedString(obj.Object, "spec", "clusterIP")
	if err != nil {
		return false, fmt.Errorf("failed to get spec.clusterIP field: %v", err)
	}
	if !ok || clusterIP == "" {
		return false, nil
	}

	ingress, ok, err := unstructured.NestedSlice(obj.Object, "status", "loadBalancer", "ingress")
	if err != nil {
		return false, fmt.Errorf("failed to get status.loadBalancer.ingress field: %v", err)
	}
	if !ok || len(ingress) == 0 {
		return false, nil
	}
	for _, i := range ingress {
		iMap, ok := i.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("failed to convert ingress to map")
		}
		ip, ok, err := unstructured.NestedString(iMap, "ip")
		if err != nil {
			return false, fmt.Errorf("failed to get ip field: %v", err)
		}
		if !ok || ip == "" {
			return false, nil
		}
	}
	return true, nil
}

// statefulSetIsReady returns true if a deployed object with kind "Service" is ready.
// This returns true if the following bullets are true:
// * status.observedGeneration == metadata.generation
// * status.replicas == spec.replicas
// * status.readyReplicas == spec.replicas
// * status.currentReplicas == spec.replicas
func statefulSetIsReady(ctx context.Context, obj *Object) (bool, error) {
	generation, ok, err := unstructured.NestedInt64(obj.Object, "metadata", "generation")
	if err != nil {
		return false, fmt.Errorf("failed to get metadata.generation field: %v", err)
	}
	if !ok {
		return false, nil
	}

	observedGeneration, ok, err := unstructured.NestedInt64(obj.Object, "status", "observedGeneration")
	if err != nil {
		return false, fmt.Errorf("failed to get status.observedGeneration field: %v", err)
	}
	if !ok || observedGeneration != generation {
		return false, nil
	}

	specReplicas, ok, err := unstructured.NestedInt64(obj.Object, "spec", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get spec.replicas field: %v", err)
	}
	if !ok {
		return false, nil
	}

	statusReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "replicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.replicas field: %v", err)
	}
	if !ok || statusReplicas != specReplicas {
		return false, nil
	}

	readyReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "readyReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.readyReplicas field: %v", err)
	}
	if !ok || readyReplicas != specReplicas {
		return false, nil
	}

	currentReplicas, ok, err := unstructured.NestedInt64(obj.Object, "status", "currentReplicas")
	if err != nil {
		return false, fmt.Errorf("failed to get status.currentReplicas field: %v", err)
	}
	if !ok || currentReplicas != specReplicas {
		return false, nil
	}

	return true, nil
}
