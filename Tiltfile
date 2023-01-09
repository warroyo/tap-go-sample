SOURCE_IMAGE = os.getenv("SOURCE_IMAGE", default='dev.registry.pivotal.io/warroyo/tap-go-sample-source')
LOCAL_PATH = os.getenv("LOCAL_PATH", default='.')
NAMESPACE = os.getenv("NAMESPACE", default='default')



k8s_custom_deploy(
    'tap-go-sample',
    apply_cmd="tanzu apps workload apply -f config/workload.yaml --live-update" +
               " --local-path " + LOCAL_PATH +
               " --source-image " + SOURCE_IMAGE +
               " --namespace " + NAMESPACE +
               " --yes >/dev/null" +
               " && kubectl get workload tap-go-sample --namespace " + NAMESPACE + " -o yaml",
    delete_cmd="tanzu apps workload delete -f config/workload.yaml --namespace " + NAMESPACE + " --yes",
    container_selector='workload',
    deps=['./build'],
    live_update=[
    sync('./build', '/tmp/tilt')  ,      
    run('cp -rf /tmp/tilt/* /layers/tanzu-buildpacks_go-build/targets/bin', trigger=['./build']),
  ]
)

# (Re)build locally when source code changes
local_resource('go-build',
  cmd='GOOS=linux GOARCH=amd64 go build -o ./build/ -buildmode pie .',
  deps=['./main.go','./pkg/'],
  ignore=['./build'],
  dir='.'
)

k8s_resource('tap-go-sample', port_forwards=["8080:8080"],
            extra_pod_selectors=[{'serving.knative.dev/service': 'tap-go-sample'}])
allow_k8s_contexts(['eks.eks-warroyo2.us-west-2.tap-iterate','arn:aws:eks:us-west-2:074754820263:cluster/tap-full'])