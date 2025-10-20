package namer

import (
	"fmt"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

type NamerInterface interface {
	Name() string
	Base(baseName string) *NamerData
	Prefix(namePrefix string) *NamerData
	Suffix(nameSuffix string) *NamerData
	Generate()
}

type NamerData struct {
	baseName string
	prefix   string
	suffix   string
	name     string
}

func (n *NamerData) Generate() {
	if len(n.prefix) > 0 {
		n.name = n.prefix + "-" + n.baseName
	} else {
		n.name = n.baseName
	}
	if len(n.suffix) > 0 {
		n.name = n.name + "-" + n.suffix
	}
}

//func NewNamer() *NamerData {
//	var namer NamerData
//	return &namer
//}

func (n *NamerData) Base(baseName string) *NamerData {
	n.baseName = baseName
	return n
}

func (n *NamerData) Prefix(namePrefix string) *NamerData {
	n.prefix = namePrefix
	return n
}

func (n *NamerData) Suffix(nameSuffix string) *NamerData {
	n.suffix = nameSuffix
	return n
}

func (n *NamerData) SetName(override string) {
	n.name = override
}

func (n *NamerData) Name() string {
	return n.name
}

func CrToSS(crName string) string {
	return crName + "-ss"
}

func CrToSSOrdinal(crName string, ordinal int) string {
	return CrToSS(crName) + "-" + strconv.Itoa(ordinal)
}

func SSToCr(ssName string) string {
	return strings.TrimSuffix(ssName, "-ss")
}

func OrdinalFromPodName(crName string, podName string) (int, error) {
	ordinalString := strings.SplitAfter(podName, CrToSS(crName)+"-")
	if len(ordinalString) > 1 {
		fields := strings.Fields(ordinalString[1])
		if len(fields) > 0 {
			if val, err := strconv.Atoi(fields[0]); err == nil {
				return val, nil
			} else {
				return -1, fmt.Errorf("ordinal int not found in subtring field %s", fields[1])
			}
		} else {
			return -1, fmt.Errorf("ordinal int not found in subtring %s", ordinalString[1])
		}
	}
	return -1, fmt.Errorf("ordinal not found in %s", podName)
}

// This function returns whether a pod belongs to a statefulset
// and if so, also returns the pod's number attached to its name
// It depends on the naming convention between a statefulset and its pods:
// The pods always take the form <ssName>-<n>, where <ssName> is the owning
// statefulset name and <n> is a non-negative numerical number.
// for example if a statefulset's name is "ex-aao-ss", it's pod name should
// be "ex-aao-ss-0", or "ex-aao-ss-1" and so on so forth.
// Note: the statefulset and pod must be in the same namespace.
func PodBelongsToStatefulset(pod *types.NamespacedName, ssName *types.NamespacedName) (error, bool, int) {

	//first check that their namespaces must match
	if ssName.Namespace != pod.Namespace {
		err := fmt.Errorf("the pod and statefulset are not in the same namespace, pod ns: %v, ss ns: %v", pod.Namespace, ssName.Namespace)
		return err, false, -1
	}
	//next pod's name should be prefixed with statefulset's
	if !strings.HasPrefix(pod.Name, ssName.Name) {
		err := fmt.Errorf("the pod's name doesn't have the statefulset's name as its prefix, pod name: %v, ss name: %v", pod.Name, ssName.Name)
		return err, false, -1
	}
	//next try to extract the pod name's number part
	podSerial := pod.Name[len(ssName.Name)+1:]
	//convert to int
	i, err := strconv.Atoi(podSerial)
	if err != nil || i < 0 {
		err := fmt.Errorf("failed to convert pod number, pod name: %v, ss name: %v", pod.Name, ssName.Name)
		return err, false, -1
	}

	return nil, true, i
}
