# craftd
Craft Demo Project

- Unit tests
- Determine the time a pod has been in completed state (not start time)
- How would you operationlize it on a prod cluster
- how would you decrease the load on the Kube API server

## States

### Running
```
Status:       Running
Start Time:   Thu, 19 Aug 2021 07:54:05 -0700
Containers:
  nginx:
    State:          Running
      Started:      Thu, 19 Aug 2021 07:54:10 -0700
```

### Terminating
```
Status:        Terminating (lasts <invalid>)
Start Time:    Thu, 19 Aug 2021 07:48:46 -0700
Containers:
  main:
    State:          Running
      Started:      Thu, 19 Aug 2021 07:48:47 -0700
```

### Completed
```
Start Time:   Thu, 19 Aug 2021 12:43:59 -0700
Status:       Succeeded
Containers:
  main:
    State:          Terminated
      Reason:       Completed
      Exit Code:    0
      Started:      Thu, 19 Aug 2021 12:44:00 -0700
      Finished:     Thu, 19 Aug 2021 12:44:00 -0700
    Ready:          False
Conditions:
  Type              Status
  Initialized       True
  Ready             False
  ContainersReady   False
  PodScheduled      True
```

### Failed
#### CrashLoopBackOff - Single Container Pod
```
Status:       Running
Start Time:   Wed, 18 Aug 2021 20:37:07 -0700
Containers:
  main:
    State:          Waiting
      Reason:       CrashLoopBackOff
    Last State:     Terminated
      Reason:       Error
      Exit Code:    1
      Started:      Thu, 19 Aug 2021 08:19:11 -0700
      Finished:     Thu, 19 Aug 2021 08:19:11 -0700
```
#### CrashLoopBackOff - Multi Containers
```
Start Time:   Thu, 19 Aug 2021 07:55:23 -0700
Status:       Running
Container:
  main:
    State:          Waiting
      Reason:       CrashLoopBackOff
    Last State:     Terminated
      Reason:       Error
      Exit Code:    1
      Started:      Thu, 19 Aug 2021 08:11:35 -0700
      Finished:     Thu, 19 Aug 2021 08:11:35 -0700
  web:
    State:          Running
      Started:      Thu, 19 Aug 2021 07:55:25 -0700
```

#### ImagePullBackOff
```
Start Time:   Thu, 19 Aug 2021 14:25:37 -0700
Status:       Pending
Containers:
  main:
    State:          Waiting
      Reason:       ImagePullBackOff
Conditions:
  Type              Status
  Initialized       True
  Ready             False
  ContainersReady   False
  PodScheduled      True
```
