package apis

import (
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// Versions lists all API versions, sorted from oldest to most recent.
var Versions = []string{
	v1alpha5.CurrentGroupVersion,
}
