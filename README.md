# gang-scheduler
Gang schedulre with [k8s scheduler framework](https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/20180409-scheduling-framework.md)

Depends on [job management repo](https://github.com/vincent-pli/job-management) 

It implements a `Permit` plugin and depends on CRD [XJob](https://github.com/vincent-pli/job-management/blob/master/pkg/apis/job/v1alpha1/xjob_types.go) like this:
```
apiVersion: batch.xscheduler.vincent-pli.com/v1alpha1
kind: XJob
metadata:
  name: xjob-sample
spec:
  minAvailable: 1
  priorityClassName: mid-priority
  tasks:
    - name: task-1
      replicas: 1
      template:
        spec:
          containers:
            - name: sleep
              image: docker
              args: ["sleep", "600"]
```

If number of `Reserve`d pod reach the `minAvailable`, all pod will be `Bind` otherwise will be `Pending`.
