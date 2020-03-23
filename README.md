# gang-scheduler
Gang schedulre with [k8s scheduler framework](https://github.com/kubernetes/enhancements/blob/master/keps/sig-scheduling/20180409-scheduling-framework.md)

It implements a `Permit` plugin and depends on another CRD like this:
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
