package controllers

import (
	"reflect"
	"testing"

	"github.com/RHsyseng/operator-utils/pkg/resource/compare"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHexShaHashOfMap(t *testing.T) {

	nilOne := HexShaHashOfMap(nil)
	nilTwo := HexShaHashOfMap(nil)

	if nilOne != nilTwo {
		t.Errorf("HexShaHashOfMap(nil) = %v, want %v", nilOne, nilTwo)
	}

	props := map[string]string{"a": "a", "b": "b"}

	propsOriginal := HexShaHashOfMap(props)

	// modify
	props["c"] = "c"

	propsModified := HexShaHashOfMap(props)

	if propsOriginal == propsModified {
		t.Errorf("HexShaHashOfMap(props mod) = %v, want %v", propsOriginal, propsModified)
	}

	// revert
	delete(props, "c")

	if propsOriginal != HexShaHashOfMap(props) {
		t.Errorf("HexShaHashOfMap(props) with revert = %v, want %v", propsOriginal, HexShaHashOfMap(props))
	}

	// modify further
	delete(props, "a")

	if propsOriginal == HexShaHashOfMap(props) {
		t.Errorf("HexShaHashOfMap(props) with just a = %v, want %v", propsOriginal, HexShaHashOfMap(props))
	}

}

func TestMapComparatorForStatefulSet(t *testing.T) {

	ss := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1beta1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "ss",
			GenerateName:               "",
			Namespace:                  "a",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "1",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          &metav1.Time{},
			DeletionGracePeriodSeconds: new(int64),
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            []metav1.OwnerReference{},
			Finalizers:                 []string{},
			ClusterName:                "",
			ManagedFields:              []metav1.ManagedFieldsEntry{},
		},
		Spec:   appsv1.StatefulSetSpec{},
		Status: appsv1.StatefulSetStatus{},
	}

	ssMod := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1beta1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "ss",
			GenerateName:               "",
			Namespace:                  "a",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "1",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          &metav1.Time{},
			DeletionGracePeriodSeconds: new(int64),
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            []metav1.OwnerReference{},
			Finalizers:                 []string{},
			ClusterName:                "",
			ManagedFields:              []metav1.ManagedFieldsEntry{},
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:             new(int32),
			Selector:             &metav1.LabelSelector{},
			Template:             v1.PodTemplateSpec{},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{},
			ServiceName:          "ssMod",
			PodManagementPolicy:  "",
			UpdateStrategy:       appsv1.StatefulSetUpdateStrategy{},
			RevisionHistoryLimit: new(int32),
			MinReadySeconds:      0,
		},
		Status: appsv1.StatefulSetStatus{},
	}

	ss0 := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{Kind: "StatefulSet", APIVersion: "apps/v1beta1"},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "ss0",
			GenerateName:               "",
			Namespace:                  "a",
			SelfLink:                   "",
			UID:                        "",
			ResourceVersion:            "1",
			Generation:                 0,
			CreationTimestamp:          metav1.Time{},
			DeletionTimestamp:          &metav1.Time{},
			DeletionGracePeriodSeconds: new(int64),
			Labels:                     nil,
			Annotations:                nil,
			OwnerReferences:            []metav1.OwnerReference{},
			Finalizers:                 []string{},
			ClusterName:                "",
			ManagedFields:              []metav1.ManagedFieldsEntry{},
		},
		Spec:   appsv1.StatefulSetSpec{},
		Status: appsv1.StatefulSetStatus{},
	}

	requestedResources = append(requestedResources, ss0)

	requestedResources = append(requestedResources, ssMod)

	deployed = make(map[reflect.Type][]client.Object)
	var deployedSets []client.Object
	deployedSets = append(deployedSets, ss)

	ssType := reflect.ValueOf(ss).Elem().Type()
	deployed[ssType] = deployedSets

	requested := compare.NewMapBuilder().Add(requestedResources...).ResourceMap()
	comparator := compare.NewMapComparator()

	comparator.Comparator.SetDefaultComparator(semanticEquals)
	deltas := comparator.Compare(deployed, requested)

	for resourceType, delta := range deltas {
		t.Log("", "instances of ", resourceType, "Will create ", len(delta.Added), "update ", len(delta.Updated), "and delete", len(delta.Removed))

		for index := range delta.Added {
			resourceToAdd := delta.Added[index]
			t.Log("", "instances of ", resourceType, "Will add", resourceToAdd)
		}

		for index := range delta.Updated {
			resourceToUpdate := delta.Updated[index]
			t.Log("", "instances of ", resourceType, "Will update", resourceToUpdate)

		}

		for index := range delta.Removed {
			resourceToRemove := delta.Removed[index]
			t.Log("", "instances of ", resourceType, "Will remove", resourceToRemove)

		}
	}
	if len(deltas[ssType].Added) != 1 {
		t.Errorf("expect new addition to appear!")

	}
	if (len(deltas[ssType].Updated)) != 1 {
		t.Errorf("not good!, expect difference on ss to be respected as an update")
	}
}

func semanticEquals(a client.Object, b client.Object) bool {
	return equality.Semantic.DeepEqual(a, b)
}
