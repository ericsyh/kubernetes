/*
Copyright 2021 The Kubernetes Authors.

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

package auth

import (
	"testing"

	utilfeature "k8s.io/apiserver/pkg/util/feature"
	featuregatetesting "k8s.io/component-base/featuregate/testing"
	kubeapiservertesting "k8s.io/kubernetes/cmd/kube-apiserver/app/testing"
	"k8s.io/kubernetes/pkg/capabilities"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/test/integration/framework"
	podsecuritytest "k8s.io/pod-security-admission/test"
)

func TestPodSecurity(t *testing.T) {
	defer featuregatetesting.SetFeatureGateDuringTest(t, utilfeature.DefaultFeatureGate, features.PodSecurity, true)()
	server := kubeapiservertesting.StartTestServerOrDie(t, kubeapiservertesting.NewDefaultTestServerOptions(), []string{
		"--anonymous-auth=false",
		"--enable-admission-plugins=PodSecurity",
		"--allow-privileged=true",
		// TODO: "--admission-control-config-file=" + admissionConfigFile.Name(),
	}, framework.SharedEtcd())
	defer server.TearDownFn()

	// ensure the global is set to allow privileged containers
	capabilities.SetForTests(capabilities.Capabilities{AllowPrivileged: true})

	opts := podsecuritytest.Options{
		ClientConfig: server.ClientConfig,

		// TODO
		ExemptClient:         nil,
		ExemptNamespaces:     []string{},
		ExemptRuntimeClasses: []string{},
	}
	podsecuritytest.Run(t, opts)
}
