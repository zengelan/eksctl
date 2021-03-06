package builder

import (
	"fmt"

	gfn "github.com/awslabs/goformation/cloudformation"
)

var servicePrincipalPartitionMappings = map[string]map[string]string{
	"aws": {
		"EC2":            "ec2.amazonaws.com",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	"aws-cn": {
		"EC2":            "ec2.amazonaws.com.cn",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
}

const servicePrincipalPartitionMapName = "ServicePrincipalPartitionMap"

func makeFnFindInMap(mapName string, args ...*gfn.Value) *gfn.Value {
	return gfn.MakeIntrinsic(gfn.FnFindInMap, append([]*gfn.Value{gfn.NewString(mapName)}, args...))
}

// MakeServiceRef returns a reference to an intrinsic map function that looks up the servicePrincipalName
// in servicePrincipalPartitionMappings
func MakeServiceRef(servicePrincipalName string) *gfn.Value {
	return makeFnFindInMap(servicePrincipalPartitionMapName, gfn.RefPartition, gfn.NewString(servicePrincipalName))
}

func makePolicyARNs(policyNames ...string) []*gfn.Value {
	policyARNs := make([]*gfn.Value, len(policyNames))
	for i, policy := range policyNames {
		policyARNs[i] = gfn.MakeFnSubString(fmt.Sprintf("arn:${%s}:iam::aws:policy/%s", gfn.Partition, policy))
	}
	return policyARNs
}

func addARNPartitionPrefix(s string) *gfn.Value {
	return gfn.MakeFnSubString(fmt.Sprintf("arn:${%s}:%s", gfn.Partition, s))
}
