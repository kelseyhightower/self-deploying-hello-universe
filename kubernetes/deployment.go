package kubernetes

var deploymentConfig = `apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  labels:
    app: hello-universe
  name: hello-universe
spec:
  replicas: 3 
  template:
    metadata:
      labels:
        app: hello-universe
      annotations:
        pod.alpha.kubernetes.io/init-containers: '[
          {
            "name": "install",
            "image": "gcr.io/hightowerlabs/alpine",
            "command": [
              "wget", "-O", "/opt/bin/hello-universe",
              "https://storage.googleapis.com/hightowerlabs/hello-universe"
            ],
            "volumeMounts": [
              {
                "name": "bin",
                "mountPath": "/opt/bin"
              }
            ]
          },
          {
            "name": "configure",
            "image": "gcr.io/hightowerlabs/alpine",
            "command": ["chmod", "+x", "/opt/bin/hello-universe"],
            "volumeMounts": [
              {
                "name": "bin",
                "mountPath": "/opt/bin"
              }
            ]
          }
        ]'
    spec:
      containers:
        - name: hello-universe
          image: "gcr.io/hightowerlabs/alpine"
          command:
            - "/opt/bin/hello-universe"
          args:
            - "-http=0.0.0.0:443"
          volumeMounts:
            - name: bin
              mountPath: /opt/bin
            - name: tls
              mountPath: /etc/hello-universe/
      volumes:
        - name: "tls"
          secret:
            secretName: "hello-universe"
        - name: bin
          emptyDir: {}
`
