/*
Copyright 2019 Google, Inc. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package resource

const (
	deploymentTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: %s
spec:
  replicas: 3
  selector:
    matchLabels:
      %s: %s
  template:
    metadata:
      labels:
        %s: %s
    spec:
      containers:
      - name: %s
        image: %s
`

	horizontalPodAutoscalerTemplate = `apiVersion: autoscaling/v2beta1
kind: HorizontalPodAutoscaler
metadata:
  name: %s
spec:
  scaleTargetRef:
    kind: Deployment
    name: %s
    apiVersion: apps/v1
  minReplicas: 1
  maxReplicas: 5
  metrics:
  - type: Resource
    resource:
      name: cpu
      targetAverageUtilization: 80
`

	namespaceTemplate = `apiVersion: v1
kind: Namespace
metadata:
  name: %s
`

	serviceTemplate = `apiVersion: v1
kind: Service
metadata:
  name: %s
spec:
  selector:
    %s: %s
  ports:
  - protocol: TCP
    port: %d
    targetPort: %d
  type: LoadBalancer
`
)
