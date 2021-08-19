# craftd
Craft Demo Project

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
Pod:
  State: Succeeded
Containers:
  main:
    State: Terminated
      Reason: Completed

### Failed
#### Single Container Pod
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
#### Multi Containers
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
