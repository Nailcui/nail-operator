nail-operator
---

> operator-example

[Kubernetes API 与 Operator，不为人知的开发者战争【知乎-阿里云云栖】](https://zhuanlan.zhihu.com/p/54633203)
[operator-sdk [github]](https://github.com/operator-framework/operator-sdk)
[operator-doc](https://sdk.operatorframework.io/)
[一个 CI/CD operator的Demo](https://segmentfault.com/a/1190000023931203)
[一个实践（有deployment代码）](https://blog.csdn.net/yunqiinsight/article/details/104929329)


## step by step
### 一、安装 operator-sdk，这里采用brew的方式安装:
```bash
# 中间碰到了一些问题，通过升级brew解决掉了
brew install operator-sdk
```
### 二、创建基于go的operator
#### 1、准备环境

- go
- operator-sdk
- kubectl
- docker
#### 2、开始
创建工作目录
```bash
mkdir nail-operator
cd nail-operator
```
初始化
```bash
operator-sdk init --domain nailcui.github.io --repo github.com/Nailcui/nail-operator
```
添加crd api
```bash
operator-sdk create api --group test --version v1alpha1 --kind Nail --controller
```
控制器镜像打包并推送到仓库
```bash
# 这一步失败了，提示有错误，连上vpn后手动通过docker build、push解决了，应该直接开vpn也行的
make docker-build docker-push IMG=naildocker/nail-operator:v0.0.1
```
创建crd资源
```bash
# build crd && 在k8s中创建crd
$ make install

# 之后查看资源，ok
$ kubectl get crd                                                                                 cuidingyu@cuideMacBook-Pro
NAME                                                  CREATED AT
nails.test.nailcui.github.io                          2021-04-17T18:27:41Z
```
创建一个自定义的
```bash
$ kubectl apply -f config/samples/test_v1alpha1_nail.yaml

# 验证
kubectl get nail                                                                                 cuidingyu@cuideMacBook-Pro
NAME          AGE
nail-sample   9s
```
### 
### 三、开发指南
简单demo中，我们仅仅修改 `./controllers/nail_controller.go`中的如下代码即可
```bash
func (r *NailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("nail", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}
```
#### 
#### 1、事件监听
上面的代码段中，函数入参：`ctx`、`req`，主要的是rep，定义如下：
```go
type Request = reconcile.Request

type Request struct {
	// NamespacedName is the name and namespace of the object to reconcile.
	types.NamespacedName
}

type NamespacedName struct {
	Namespace string
	Name      string
}
```
其实就是当有变化的时候，会将namespace、name通知到这里；先修改为：
```go
func (r *NailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("nail", req.NamespacedName)
	fmt.Printf("namespace: %s name: %s\n", req.Namespace, req.Name)

	nail := &testv1alpha1.Nail{}
	err := r.Get(ctx, req.NamespacedName, nail)
	if err != nil {
		fmt.Printf("get error: %s\n", err)
	} else {
		fmt.Printf("new replicas: %d\n", nail.Spec.Replicas)
	}
	return ctrl.Result{}, nil
}
```
直接本地执行（本地需要有kubeconfig）修改 manifest并apply后，便可以看到日志


#### 2、执行管理操作
查询事件对应资源的要求，然后进行处理
比如我们要进行一个deployment
```go
deployment := appv1.Deployment{
	// ...
}
r.Create(&deployment)
```


#### 3、后续操作
再看下需修改的代码的返回值
```go
type Result struct {
	// Requeue tells the Controller to requeue the reconcile key.  Defaults to false.
	Requeue bool

	// RequeueAfter if greater than 0, tells the Controller to requeue the reconcile key after the Duration.
	// Implies that Requeue is true, there is no need to set Requeue to true at the same time as RequeueAfter.
	RequeueAfter time.Duration
}

```

---

