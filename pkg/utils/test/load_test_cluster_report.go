package test

import (
	"time"

	"github.com/fusor/cpma/pkg/api"
	O7tapiroute "github.com/openshift/api/route/v1"
	k8sapicore "k8s.io/api/core/v1"
	k8sapistorage "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	k8smachinery "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateTestPVList create test pv list
func CreateTestPVList() *k8sapicore.PersistentVolumeList {
	pvList := &k8sapicore.PersistentVolumeList{}
	pvList.Items = make([]k8sapicore.PersistentVolume, 0)

	resources := make(k8sapicore.ResourceList)
	cpu := resource.Quantity{
		Format: resource.DecimalSI,
	}
	cpu.Set(int64(1))
	resources["cpu"] = cpu

	memory := resource.Quantity{
		Format: resource.BinarySI,
	}
	memory.Set(int64(1))
	resources["memory"] = memory

	driver := k8sapicore.PersistentVolumeSource{
		NFS: &k8sapicore.NFSVolumeSource{
			Server: "example.com",
		},
	}

	pvList.Items = append(pvList.Items, k8sapicore.PersistentVolume{
		ObjectMeta: k8smachinery.ObjectMeta{
			Name: "testpv",
		},
		Spec: k8sapicore.PersistentVolumeSpec{
			PersistentVolumeSource: driver,
			StorageClassName:       "testclass",
			Capacity:               resources,
		},
		Status: k8sapicore.PersistentVolumeStatus{
			Phase: k8sapicore.VolumePending,
		},
	})

	return pvList
}

// CreateTestNodeList create test node list
func CreateTestNodeList() *k8sapicore.NodeList {
	// Create node labels
	masterNodeLabels := make(map[string]string)
	masterNodeLabels["node-role.kubernetes.io/master"] = "true"

	masterNodeCapacity := make(k8sapicore.ResourceList)
	// Add CPU node usage
	cpuQuantity := resource.Quantity{
		Format: resource.DecimalSI,
	}
	cpuQuantity.Set(int64(2))
	masterNodeCapacity["cpu"] = cpuQuantity

	// Add node memory usage
	memoryQuantity := resource.Quantity{
		Format: resource.BinarySI,
	}
	memoryQuantity.Set(int64(2048))
	masterNodeCapacity["memory"] = memoryQuantity

	// Add pods
	podsQuantity := resource.Quantity{
		Format: resource.DecimalSI,
	}
	podsQuantity.Set(int64(10))
	masterNodeCapacity["pods"] = podsQuantity

	// Add resources that are available for scheduling
	allocatableResources := make(k8sapicore.ResourceList)

	allocatableMemoryQuantity := resource.Quantity{
		Format: resource.BinarySI,
	}
	allocatableMemoryQuantity.Set(int64(2048))
	masterNodeCapacity["memory"] = allocatableMemoryQuantity

	// Add pod count
	podList := &k8sapicore.PodList{}
	podList.Items = make([]k8sapicore.Pod, 0)
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "test-master",
		},
	})
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "test-master",
		},
	})
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "test-master",
		},
	})
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			NodeName: "not-this-node",
		},
	})

	namespaceList := make([]api.NamespaceResources, 0)
	namespaceList = append(namespaceList, api.NamespaceResources{
		PodList: podList,
	})

	// Init fake nodes
	nodes := make([]k8sapicore.Node, 0)
	nodes = append(nodes, k8sapicore.Node{
		ObjectMeta: k8smachinery.ObjectMeta{
			Name:   "test-master",
			Labels: masterNodeLabels,
		},
		Status: k8sapicore.NodeStatus{
			Capacity:    masterNodeCapacity,
			Allocatable: allocatableResources,
		},
	})

	return &k8sapicore.NodeList{
		Items: nodes,
	}
}

// CreateStorageClassList create storage class list
func CreateStorageClassList() *k8sapistorage.StorageClassList {
	storageClassList := &k8sapistorage.StorageClassList{}
	storageClassList.Items = make([]k8sapistorage.StorageClass, 0)
	storageClassList.Items = append(storageClassList.Items, k8sapistorage.StorageClass{
		ObjectMeta: k8smachinery.ObjectMeta{
			Name: "testclass",
		},
		Provisioner: "testprovisioner",
	})

	return storageClassList
}

// CreateTestNameSpaceList create test namespace list
func CreateTestNameSpaceList() []api.NamespaceResources {
	namespaces := make([]api.NamespaceResources, 0)

	namespaces = append(namespaces, api.NamespaceResources{
		NamespaceName: "testNamespace",
		PodList:       CreateTestPodList(),
		RouteList:     CreateTestRouteList(),
	})

	return namespaces
}

// CreateTestPodList test pod list
func CreateTestPodList() *k8sapicore.PodList {
	podList := &k8sapicore.PodList{}
	podList.Items = make([]k8sapicore.Pod, 0)
	timeStamp, _ := time.Parse(time.RFC1123Z, "Tue, 17 Nov 2009 21:34:58 +0100")
	podList.Items = append(podList.Items, k8sapicore.Pod{
		ObjectMeta: k8smachinery.ObjectMeta{
			Name:              "test-pod1",
			CreationTimestamp: k8smachinery.NewTime(timeStamp),
		},
	})

	podList.Items = append(podList.Items, k8sapicore.Pod{
		ObjectMeta: k8smachinery.ObjectMeta{
			Name:              "test-pod2",
			CreationTimestamp: k8smachinery.NewTime(timeStamp),
		},
	})

	return podList
}

// CreateTestRouteList create test route list
func CreateTestRouteList() *O7tapiroute.RouteList {
	routeList := &O7tapiroute.RouteList{}
	routeList.Items = make([]O7tapiroute.Route, 0)

	alternateBackends := make([]O7tapiroute.RouteTargetReference, 0)
	alternateBackends = append(alternateBackends, O7tapiroute.RouteTargetReference{
		Kind: "testkind",
		Name: "testname",
	})

	to := O7tapiroute.RouteTargetReference{
		Kind: "testkindTo",
		Name: "testTo",
	}

	tls := &O7tapiroute.TLSConfig{
		Termination: O7tapiroute.TLSTerminationEdge,
	}

	routeList.Items = append(routeList.Items, O7tapiroute.Route{
		ObjectMeta: k8smachinery.ObjectMeta{
			Name: "route1",
		},
		Spec: O7tapiroute.RouteSpec{
			AlternateBackends: alternateBackends,
			Host:              "testhost",
			Path:              "testpath",
			To:                to,
			TLS:               tls,
			WildcardPolicy:    O7tapiroute.WildcardPolicyNone,
		},
	})

	return routeList
}

// CreateTestPodResourceList create test resources
func CreateTestPodResourceList() *k8sapicore.PodList {
	resources := make(k8sapicore.ResourceList)
	cpu := resource.Quantity{
		Format: resource.DecimalSI,
	}
	cpu.Set(int64(1))
	resources["cpu"] = cpu

	memory := resource.Quantity{
		Format: resource.BinarySI,
	}
	memory.Set(int64(1))
	resources["memory"] = memory

	containers := make([]k8sapicore.Container, 0)
	containers = append(containers, k8sapicore.Container{
		Resources: k8sapicore.ResourceRequirements{
			Requests: resources,
		},
	})
	containers = append(containers, k8sapicore.Container{
		Resources: k8sapicore.ResourceRequirements{
			Requests: resources,
		},
	})

	podList := &k8sapicore.PodList{}
	podList.Items = make([]k8sapicore.Pod, 0)
	podList.Items = append(podList.Items, k8sapicore.Pod{
		Spec: k8sapicore.PodSpec{
			Containers: containers,
		},
	})

	return podList
}
