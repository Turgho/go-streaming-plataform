# go-streaming-plataform

```text
streaming-platform/
├── services/
│   ├── user/
│   ├── upload/
│   ├── transcode/
│   └── notification/
├── proto/
│   ├── user/
│   ├── upload/
│   └── transcode/
├── gateway/
├── k8s/
│   ├── user/
│   │   ├── deployment.yaml
│   │   ├── service.yaml
│   │   └── configmap.yaml
│   ├── upload/
│   ├── transcode/
│   └── ingress.yaml
├── .github/
│   └── workflows/
│       ├── test.yml
│       ├── build.yml
│       └── deploy.yml
└── docker-compose.yml  ← pra desenvolvimento local
```