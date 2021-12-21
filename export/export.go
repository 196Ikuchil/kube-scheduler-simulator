package export

import (
	"context"

	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	v1 "k8s.io/client-go/applyconfigurations/core/v1"
	confstoragev1 "k8s.io/client-go/applyconfigurations/storage/v1"
	clientset "k8s.io/client-go/kubernetes"
	v1beta2config "k8s.io/kube-scheduler/config/v1beta2"
)

type Service struct {
	client              clientset.Interface
	podService          PodService
	nodeService         NodeService
	pvService           PersistentVolumeService
	pvcService          PersistentVolumeClaimService
	storageClassService StorageClassService
	schedulerService    SchedulerService
}

// Resources  will used to compile all the resources for export.
type Resources struct {
	Pods            []corev1.Pod                              `json:"podList"`
	Nodes           []corev1.Node                             `json:"nodeList"`
	Pvs             []corev1.PersistentVolume                 `json:"pvList"`
	Pvcs            []corev1.PersistentVolumeClaim            `json:"pvcList"`
	StorageClasses  []storagev1.StorageClass                  `json:"storageClassList"`
	SchedulerConfig *v1beta2config.KubeSchedulerConfiguration `json:"schedulerConfig"`
}

type ResourcesApplyConfiguration struct {
	Pods            []v1.PodApplyConfiguration                     `json:"podList"`
	Nodes           []v1.NodeApplyConfiguration                    `json:"nodeList"`
	Pvs             []v1.PersistentVolumeApplyConfiguration        `json:"pvList"`
	Pvcs            []v1.PersistentVolumeClaimApplyConfiguration   `json:"pvcList"`
	StorageClasses  []confstoragev1.StorageClassApplyConfiguration `json:"storageClassList"`
	SchedulerConfig *v1beta2config.KubeSchedulerConfiguration      `json:"schedulerConfig"`
}

type PodService interface {
	List(ctx context.Context) (*corev1.PodList, error)
	Apply(ctx context.Context, pod *v1.PodApplyConfiguration) (*corev1.Pod, error)
}

type NodeService interface {
	List(ctx context.Context) (*corev1.NodeList, error)
	Apply(ctx context.Context, nac *v1.NodeApplyConfiguration) (*corev1.Node, error)
}

type PersistentVolumeService interface {
	List(ctx context.Context) (*corev1.PersistentVolumeList, error)
	Apply(ctx context.Context, persistentVolume *v1.PersistentVolumeApplyConfiguration) (*corev1.PersistentVolume, error)
}

type PersistentVolumeClaimService interface {
	Get(ctx context.Context, name string) (*corev1.PersistentVolumeClaim, error)
	List(ctx context.Context) (*corev1.PersistentVolumeClaimList, error)
	Apply(ctx context.Context, persistentVolumeClaime *v1.PersistentVolumeClaimApplyConfiguration) (*corev1.PersistentVolumeClaim, error)
}

type StorageClassService interface {
	List(ctx context.Context) (*storagev1.StorageClassList, error)
	Apply(ctx context.Context, storageClass *confstoragev1.StorageClassApplyConfiguration) (*storagev1.StorageClass, error)
}

type SchedulerService interface {
	GetSchedulerConfig() *v1beta2config.KubeSchedulerConfiguration
	RestartScheduler(cfg *v1beta2config.KubeSchedulerConfiguration) error
}

func NewResourcesService(client clientset.Interface, pods PodService, nodes NodeService, pvs PersistentVolumeService, pvcs PersistentVolumeClaimService, storageClasss StorageClassService, schedulers SchedulerService) *Service {
	return &Service{
		client:              client,
		podService:          pods,
		nodeService:         nodes,
		pvService:           pvs,
		pvcService:          pvcs,
		storageClassService: storageClasss,
		schedulerService:    schedulers,
	}
}

// Get all resources from each service.
func (s *Service) get(ctx context.Context) (*Resources, error) {
	pods, err := s.podService.List(ctx)
	if err != nil {
		return nil, xerrors.Errorf("call list pods: %w", err)
	}
	nodes, err := s.nodeService.List(ctx)
	if err != nil {
		return nil, xerrors.Errorf("call list nodes: %w", err)
	}
	pvs, err := s.pvService.List(ctx)
	if err != nil {
		return nil, xerrors.Errorf("call list PersistentVolumes: %w", err)
	}
	pvcs, err := s.pvcService.List(ctx)
	if err != nil {
		return nil, xerrors.Errorf("call list PersistentVolumeClaims: %w", err)
	}
	scs, err := s.storageClassService.List(ctx)
	if err != nil {
		return nil, xerrors.Errorf("to call list storageClasses")
	}
	ss := s.schedulerService.GetSchedulerConfig()

	return &Resources{
		Pods:            pods.Items,
		Nodes:           nodes.Items,
		Pvs:             pvs.Items,
		Pvcs:            pvcs.Items,
		StorageClasses:  scs.Items,
		SchedulerConfig: ss,
	}, nil
}

func (s *Service) Export(ctx context.Context) (*Resources, error) {
	resources, err := s.get(ctx)
	if err != nil {
		return nil, xerrors.Errorf("export resources all: %w", err)
	}
	return resources, nil
}

// Import all resources from posted data.
// (1) Restart scheduler based on the data.
// (2) Apply each resource to the scheduler.
//     * If UID is not nil, an error will occur. (try to find existing resource by UID)
// (3) Get all resources. (Separated the get function to unify the struct format.)
func (s *Service) Import(ctx context.Context, resources *ResourcesApplyConfiguration) (*Resources, error) {
	// TODO: Issue: #12 PR: #13
	// if err := s.schedulerService.RestartScheduler(resources.SchedulerConfig); err != nil {
	// 	klog.Warningf("failed to start scheduler with imported configuration: %v", err)
	// 	return nil, xerrors.Errorf("restart scheduler with imported configuration: %w", err)
	// }
	for i := range resources.StorageClasses {
		sc := resources.StorageClasses[i]
		sc.ObjectMetaApplyConfiguration.UID = nil
		_, err := s.storageClassService.Apply(ctx, &sc)
		if err != nil {
			return nil, xerrors.Errorf("apply StorageClass: %w", err)
		}
	}
	for i := range resources.Pvcs {
		pvc := resources.Pvcs[i]
		pvc.ObjectMetaApplyConfiguration.UID = nil
		_, err := s.pvcService.Apply(ctx, &pvc)
		if err != nil {
			return nil, xerrors.Errorf("apply PersistentVolumeClaims: %w", err)
		}
	}
	for i := range resources.Pvs {
		pv := resources.Pvs[i]
		pv.ObjectMetaApplyConfiguration.UID = nil
		if *pv.Status.Phase == "Bound" {
			// PersistentVolumeClaims's UID has been changed to a new value.
			pvc, err := s.pvcService.Get(ctx, *pv.Spec.ClaimRef.Name)
			if err == nil {
				pv.Spec.ClaimRef.UID = &pvc.UID
			} else {
				pv.Spec.ClaimRef.UID = nil
			}
		}
		_, err := s.pvService.Apply(ctx, &pv)
		if err != nil {
			return nil, xerrors.Errorf("apply PersistentVolume: %w", err)
		}
	}
	for i := range resources.Nodes {
		node := resources.Nodes[i]
		node.ObjectMetaApplyConfiguration.UID = nil
		_, err := s.nodeService.Apply(ctx, &node)
		if err != nil {
			return nil, xerrors.Errorf("apply Node: %w", err)
		}
	}
	for i := range resources.Pods {
		pod := resources.Pods[i]
		pod.ObjectMetaApplyConfiguration.UID = nil
		_, err := s.podService.Apply(ctx, &pod)
		if err != nil {
			return nil, xerrors.Errorf("apply Pod: %w", err)
		}
	}
	rs, err := s.get(ctx)
	if err != nil {
		return nil, xerrors.Errorf("load the resources after import: %w", err)
	}
	return rs, nil
}
