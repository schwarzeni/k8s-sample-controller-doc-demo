# Sample Controller for K8S

[参考项目: kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)

编程环境：
- Go v1.13.8
- MacOS 10.15.8
- Kubernetes v1.17.0

## Write a CRD

### Project Setup

注意，由于后期使用 code-generator 时，code-generator 不支持 go module 模式，所以项目虽然依然可以使用 go module ，但是目录需要建成 GOPATH 的样式，例如如下

```bash
go mod init github.com/schwarzeni/k8s-sample-controller-doc-demo
```

此时项目 k8s-sample-controller-doc-demo 父文件夹路径应该为 `github.com/schwarzeni/`

安装相关的依赖，对于 1.17 版本的 Kubernetes ，下载 0.17.0 版本的相关客户端，下载 code-generator 时可能会报错，这个不用管

```bash
go get -u k8s.io/code-generator@v0.17.0
go get -u k8s.io/client-go@v0.17.0
go get -u k8s.io/apimachinery@v0.17.0
go get -u k8s.io/api@v0.17.0
```

建立 `hack/tools.go` 并写入如下内容，便于之后使用 `go mod vendor`

```go
// +build tools

// This package imports things required by build scripts, to force `go mod` to see them as dependencies
package tools

import _ "k8s.io/code-generator"
```

之后执行命令：

```bash
go mod vendor
```

此时项目结构如下：

```txt
.
├── go.mod
├── go.sum
├── vendor
    └── ...
├── hack
    └── tools.go
```

文件 `vendor/k8s.io/code-generator/generate-groups.sh` 就是用来生成代码的脚本

---

### Code Generation

首先先定义相关的 K8S 资源，然后使用 K8S 提供的工具生成相关的访问代码。这里的代码全部参照 [kubernetes/sample-controller](https://github.com/kubernetes/sample-controller) 中的相关代码。

在根目录下新建文件夹 `pkg/apis/samplecontroller` 作为定义资源的目录，在此目录下新建四个文件并填入相应的内容，IDE 可能会报错，这个不用管

```txt
pkg/apis/samplecontroller
├── register.go
└── v1alpha1
    ├── doc.go
    ├── register.go
    └── types.go
```

register.go
```go
package samplecontroller

// GroupName is the group name used in this package
const GroupName = "samplecontroller.k8s.io"
```

v1alpha1/doc.go
```go
// +k8s:deepcopy-gen=package
// +groupName=samplecontroller.k8s.io

// Package v1alpha1 is the v1alpha1 version of the API.
package v1alpha1
```

v1alpha1/types.go
```go
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for a Foo resource
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooSpec   `json:"spec"`
	Status FooStatus `json:"status"`
}

// FooSpec is the spec for a Foo resource
type FooSpec struct {
	DeploymentName string `json:"deploymentName"`
	Replicas       *int32 `json:"replicas"`
}

// FooStatus is the status for a Foo resource
type FooStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type FooList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Foo `json:"items"`
}
```

v1alpha1/register.go
```go
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	samplecontroller "github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/apis/samplecontroller"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: samplecontroller.GroupName, Version: "v1alpha1"}

// Kind takes an unqualified kind and returns back a Group qualified GroupKind
func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder initializes a scheme builder
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	// AddToScheme is a global function that registers this API group & version to a scheme
	AddToScheme = SchemeBuilder.AddToScheme
)

// Adds the list of known types to Scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Foo{},
		&FooList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
```

再执行一次 `go mod vendor`

为自动生成的文件准备一个 header

hack/boilerplate.go.txt

```txt
/*
this is a header file
*/
```

执行如下命令生成代码：（我在 Goland 的终端中执行此命令会报错，但是在普通终端中执行就没问题了，目前还不太清楚原因）

```bash
bash vendor/k8s.io/code-generator/generate-groups.sh all \
github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/generated  github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/apis \
samplecontroller:v1alpha1 \
--go-header-file hack/boilerplate.go.txt \
--output-base /Users/nizhenyang/my-project/cloud/src
```

对于 `--output-base` ，假如你的项目绝对路径为 `/Users/nizhenyang/my-project/cloud/src/github.com/schwarzeni/k8s-sample-controller-doc-demo` ，将后半段 `github.com/schwarzeni/k8s-sample-controller-doc-demo` 视为 GOPATH，则这个参数的值就是前半段路径。

执行完代码之后，在 `pkg/` 下会出现新生成的 `generated` 文件夹， 大致的结构如下

```txt
pkg/generated
├── clientset
│   └── versioned
│       ├── clientset.go
│       ├── doc.go
│       ├── fake
│       ├── scheme
│       └── typed
├── informers
│   └── externalversions
│       ├── factory.go
│       ├── generic.go
│       ├── internalinterfaces
│       └── samplecontroller
└── listers
    └── samplecontroller
        └── v1alpha1
```



在`pkg/apis/samplecontroller/v1alpha1/` 下新生成了 `zz_generated.deepcopy.go` 文件。

生成代码任务完成，再执行一次 `go mod vendor`

---

### Testing

这里写一个测试程序来验证我们生产代码的可用性，首先，根据自定义资源的格式准备两份 yaml 配置文件：

crd.yaml

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: foos.samplecontroller.k8s.io
spec:
  group: samplecontroller.k8s.io
  version: v1alpha1
  names:
    kind: Foo
    plural: foos
  scope: Namespaced
```

example-foo.yaml

```yaml
apiVersion: samplecontroller.k8s.io/v1alpha1
kind: Foo
metadata:
  name: example-foo
spec:
  deploymentName: example-foo
  replicas: 1
```

执行命令 `kubectl apply -f crd.yaml` 创建资源定义

删除 vendor 目录，在项目根目录新建 `main.go`  ，内容如下

```go
package main

import (
	"log"
	"time"

	"github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/apis/samplecontroller/v1alpha1"
	"k8s.io/client-go/tools/cache"

	clientset "github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/generated/clientset/versioned"
	informers "github.com/schwarzeni/k8s-sample-controller-doc-demo/pkg/generated/informers/externalversions"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
  // 你的Kubernetes配置文件路径，一般为 ~/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", "config/config")
	if err != nil {
		panic(err)
	}
	exampleClient, err := clientset.NewForConfig(config)
	sharedInformers := informers.NewSharedInformerFactory(exampleClient, time.Second*2)
	informer := sharedInformers.Samplecontroller().V1alpha1().Foos().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(*v1alpha1.Foo)
			log.Printf("New Foo Added: %s", mObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(*v1alpha1.Foo)
			log.Printf("Delete Foo : %s", mObj.GetName())
		},
	})

	stopCh := make(chan struct{})
	defer close(stopCh)
	informer.Run(stopCh)
}
```

执行该程序

执行 `kubectl apply -f foo-example.yaml` ，发现程序输出了 “New Foo Added: example-foo”

执行 `kubectl delete -f foo-example.yaml` ，发现程序输出了 “Delete Foo : example-foo”

