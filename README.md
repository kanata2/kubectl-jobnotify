kubectl-jobnotify
---

This plugin make you notify k8s's job completion through some tools (currently only Slack).  
However this is **just a toy**(i.e. not production-ready) for studying Kubernetes :)

### Setup

    $ go get git@github.com:kanata2/kubectl-jobnotify
    $ make build
    $ mv kubectl-jobnotify /path/to/bin
    $ export SLACK_WEBHOOK_URL=https://hooks.slack.com/services/xxxxxxx/xxxxxxx/xxxxxxx


### Usage

    $ kubectl jobnotify --job <your_job_name>

More info: `kubectl jobnotify help`

### Example

1. make job yaml and apply
```
$ cat <<YAML > job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: sleeper
spec:
  completions: 5
  parallelism: 3
  template:
    spec:
      containers:
      - name: sleep
        image: alpine
        command: ["sh", "-c"]
        args:
        - |
          sleep 3
      restartPolicy: Never
YAML

$ kubectl apply -f job.yaml
job.batch/sleeper created
```

2. `$ kubectl jobnotify --job sleeper`

3. notify to you
![result](https://user-images.githubusercontent.com/7460883/63013081-66434f00-bec6-11e9-8b9f-06231312f175.png)
