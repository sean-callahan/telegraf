# Kubernete_State Plugin
[kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) is an open source project designed to generate metrics derived from the state of Kubernetes objects – the abstractions Kubernetes uses to represent your cluster. With this information you can monitor details such as:

State of nodes, pods, and jobs
Compliance with replicaSets specs
Resource requests and min/max limits

The Kubernete State Plugin gathers information based on [kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)

### Configuration:

### Metrics:

#### kube_configmap
- ts: configmap created
- measurement & fields:
  - gauge (always 1)
- tags:
  - namespace
  - configmap
  - resource_version

#### kube_cronjob
- ts: now
- measurement & fields:
  - status_active ( # of active jobs)
  - spec_starting_deadline_seconds
  - next_schedule_time
  - created
  - status_last_schedule_time
- tags:
  - namespace
  - cronjob
  - schedule
  - concurrency_policy
  - label_* 

#### kube_daemonset
- ts: daemonset created
- measurement & fields:
  - metadata_generation
- tags:
  - namespace
  - daemonset
  - label_*

#### kube_daemonset_status
- ts: now
- measurement & fields:
  - current_number_scheduled
  - desired_number_scheduled
  - number_available
  - number_misscheduled
  - number_ready
  - number_unavailable
  - updated_number_scheduled
- tags:
  - namespace
  - daemonset
  - label_*

#### kube_deployment
- ts: deployment created
- measurement & fields:
  - spec_replicas
  - metadata_generation
  - spec_strategy_rollingupdate_max_unavailable
  - spec_strategy_rollingupdate_max_surge
- tags:
  - namespace
  - deployment
  - spec_paused ("true", "false")

#### kube_deployment_status
- ts: now
- measurement & fields:
  - replicas
  - replicas_available
  - replicas_unavailable
  - replicas_updated
  - observed_generation
- tags:
  - namespace
  - deployment

#### kube_endpoint
- ts: endpoint created
- measurement & fields:
  - address_available
  - address_not_ready
- tags:
  - namespace
  - endpoint

#### kube_hpa
- ts: hpa created
- measurement & fields:
  - metadata_generation
  - spec_max_replicas
  - spec_min_replicas
- tags:
  - namespace
  - hpa
  - label_*

#### kube_hpa_status
- ts: now
- measurement & fields:
  - current_replicas
  - desired_replicas
  - condition_true (1 = "true", 0 = "false")
  - condition_false (1 = "true", 0 = "false")
  - condition_unkown
- tags:
  - namespace
  - hpa
  - condition ("true", "false", "unkown")
  
#### kube_job
- ts: job completion time
- measurement & fields:
  - status_succeeded
  - status_failed
  - status_active
  - spec_parallelism
  - spec_completions
  - created
  - spec_active_deadline_seconds
  - status_start_time
- tags:
  - namespace
  - job_name
  - label_*

#### kube_job_condition
- ts: job condition last transition time
- measurement & fields:
  - completed (1 = "true", 0 = "false")
  - failed (1 = "true", 0 = "false")
- tags:
  - namespace
  - job_name
  - condition

#### kube_limitrange
- ts: limit range creation time
- measurement & fields:
  - min_pod_cpu
  - min_pod_memory
  - min_pod_storage
  - min_pod_ephemeral_storage
  - min_container_cpu
  - min_container_memory
  - min_container_storage
  - min_container_ephemeral_storage
  - min_persistentvolumeclaim_cpu
  - min_persistentvolumeclaim_memory
  - min_persistentvolumeclaim_storage
  - min_persistentvolumeclaim_ephemeral_storage
  - max_pod_cpu
  - max_pod_memory
  - max_pod_storage
  - max_pod_ephemeral_storage
  - max_container_cpu
  - max_container_memory
  - max_container_storage
  - max_container_ephemeral_storage
  - max_persistentvolumeclaim_cpu
  - max_persistentvolumeclaim_memory
  - max_persistentvolumeclaim_storage
  - max_persistentvolumeclaim_ephemeral_storage
  - default_pod_cpu
  - default_pod_memory
  - default_pod_storage
  - default_pod_ephemeral_storage
  - default_container_cpu
  - default_container_memory
  - default_container_storage
  - default_container_ephemeral_storage
  - default_persistentvolumeclaim_cpu
  - default_persistentvolumeclaim_memory
  - default_persistentvolumeclaim_storage
  - default_persistentvolumeclaim_ephemeral_storage
  - default_request_pod_cpu
  - default_request_pod_memory
  - default_request_pod_storage
  - default_request_pod_ephemeral_storage
  - default_request_container_cpu
  - default_request_container_memory
  - default_request_container_storage
  - default_request_container_ephemeral_storage
  - default_request_persistentvolumeclaim_cpu
  - default_request_persistentvolumeclaim_memory
  - default_request_persistentvolumeclaim_storage
  - default_request_persistentvolumeclaim_ephemeral_storage
  - max_limit_request_ratio_pod_cpu
  - max_limit_request_ratio_pod_memory
  - max_limit_request_ratio_pod_storage
  - max_limit_request_ratio_pod_ephemeral_storage
  - max_limit_request_ratio_container_cpu
  - max_limit_request_ratio_container_memory
  - max_limit_request_ratio_container_storage
  - max_limit_request_ratio_container_ephemeral_storage
  - max_limit_request_ratio_persistentvolumeclaim_cpu
  - max_limit_request_ratio_persistentvolumeclaim_memory
  - max_limit_request_ratio_persistentvolumeclaim_storage
  - max_limit_request_ratio_persistentvolumeclaim_ephemeral_storage
- tags:
  - namespace
  - limitrange

#### kube_namespace
- ts: namespace creation time
- measurement & fields:
  - gauge (always 1)
- tags
  - namespace
  - status_phase
  - label_*
  - annotation_*

#### kube_node
- ts: now
- measurement & fields:
  - created
  - status_capacity_cpu_cores
  - status_capacity_ephemera_storage_bytes
  - status_capacity_memory_bytes
  - status_capacity_pods
  - status_capacity_*
  - status_allocatable_cpu_cores
  - status_allocatable_ephemera_storage_bytes
  - status_allocatable_memory_bytes
  - status_allocatable_pods
  - status_allocatable_*
- tags:
  - node
  - kernel_version
  - os_image
  - container_runtime_version
  - kubelet_version
  - kubeproxy_version
  - status_phase
  - provider_id
  - spec_unschedulable
  - label_*

#### kube_node_spec_taint
- ts: now
- measurement & fields:
  - gauge (always 1)
- tags:
  - node
  - key
  - value
  - effect

#### kube_node_status_conditions  
- ts: condition last transition time
- measurement & fields:
  - gauge (always 1)
- tags:
  - node
  - condition
  - status

#### kube_persistentvolume
- ts: now
- measurement & fields:
  - status_pending (1 = "true", 0 = "false")
  - status_available (1 = "true", 0 = "false")
  - status_bound (1 = "true", 0 = "false")
  - status_released (1 = "true", 0 = "false")
  - status_failed (1 = "true", 0 = "false")
- tags:
  - persistentvolume
  - storageclass
  - status
  - label_*

#### kube_persistentvolumeclaim
- ts: now
- measurement & fields:
  - status_lost (1 = "true", 0 = "false")
  - status_bound (1 = "true", 0 = "false")
  - status_failed (1 = "true", 0 = "false")
  - resource_requests_storage_bytes
- tags:
  - namespace
  - persistentvolumeclaim
  - storageclass
  - volumename
  - status
  - label_*

#### kube_pod

- ts: pod created
- measurement & fields:
  - gauge (always 1)
- tags:
  - namesapce
  - pod
  - node
  - created_by_kind
  - created_by_name
  - owner_kind
  - owner_name
  - owner_is_controller ("true", "false")
  - label_*

#### kube_pod_status

- ts: now
- measurement & fields:
  - status_phase_pending (1 = "true", 0 = "false")
  - status_phase_succeeded (1 = "true", 0 = "false")
  - status_phase_failed (1 = "true", 0 = "false")
  - status_phase_running (1 = "true", 0 = "false")
  - status_phase_unknown (1 = "true", 0 = "false")
  - completion_time
  - scheduled_time
- tags:
  - namesapce
  - pod
  - node
  - host_ip
  - pod_ip
  - status_phase ("pending", "succeeded", "failed", "running", "unknown")
  - ready ("true", "false")
  - scheduled ("true", "false")

#### kube_pod_container

- ts: now
- measurement & fields:
  - status_restarts_total
  - status_waiting (1 = "true", 0 = "false")
  - status_running (1 = "true", 0 = "false")
  - status_terminated (1 = "true", 0 = "false")
  - status_ready (1 = "true", 0 = "false")
  - resource_requests_cpu_cores
  - resource_requests_memory_bytes
  - resource_requests_storage_bytes
  - resource_requests_ephemeral_storage_bytes
  - resource_limits_cpu_cores
  - resource_limits_memory_bytes
  - resource_limits_storage_bytes
  - resource_limits_ephemeral_storage_bytes
- tags:
  - namespace
  - pod_name
  - node_name
  - container
  - image
  - image_id
  - container_id
  - status_waiting_reason
  - status_terminated_reason


#### kube_pod_volume
- ts: now
- measurement & fields:
  - read_only (1 = "true", 0 = "false")
- tags:
  - namespace
  - pod
  - volume
  - persistentvolumeclaim

#### kube_replicasets  
- ts: replicaset creation time
- measurement & fields:
  - metadata_generation
  - spec_replicas
- tags
  - namespace
  - replicaset 

#### kube_replicasets_status
- ts: now
- measurement & fields:
  - replicas
  - fully_labeled_replicas
  - ready_replicas
  - observed_generation
- tags:
  - namespace
  - replicaset

#### kube_replicationcontroller
- ts: replicaset creation time
- measurement & fields:
  - metadata_generation
  - spec_replicas
- tags
  - namespace
  - replicationcontroller

#### kube_replicationcontroller_status
- ts: now
- measurement & fields:
  - replicas
  - fully_labeled_replicas
  - ready_replicas
  - available_replicas
  - observed_generation
- tags:
  - namespace
  - replicationcontroller

#### kube_resourcequota
- ts: resourcequota creation time
- measurement & fields:
  - gauge
- tags:
  - namespace
  - resourcequota
  - resource
  - type ("hard", "used")

#### kube_secret
- ts: secret creation time
- measurement & fields:
  - gauge (always 1)
- tags:
  - namespace
  - secret
  - resource_version
  - label_*

#### kube_service
- ts: service creation time
- measurement & fields:
  - gauge (always 1)
- tags:
  - namespace
  - service
  - type
  - cluster_ip
  - label_*

#### kube_statefulset
- ts: statefulset creation time
- measurement & fields:
  - metadata_generation
  - replicas
- tags:
  - namespace
  - statefulset
  - label_*

#### kube_statefulset_status
- ts: now
- measurement & fields:
  - replicas
  - replicas_current
  - replicas_ready
  - replicas_updated
  - observed_generation
- tags:
  - namespace
  - statefulset
