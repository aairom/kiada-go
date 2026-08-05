# Kubernetes in Action, 2nd Edition — Architecture & Book Structure

## Book Structure Flowchart

```mermaid
flowchart TD
    BOOK["📚 Kubernetes in Action\n2nd Edition"]

    BOOK --> P1["Part I\nGetting Started with Kubernetes"]
    BOOK --> P2["Part II\nRunning Applications on Kubernetes"]
    BOOK --> P3["Part III\nManaging App Configuration & Storage"]
    BOOK --> P4["Part IV\nConnecting & Exposing Applications"]
    BOOK --> P5["Part V\nManaging Applications with Controllers"]

    P1 --> C1["Ch 1 — Introducing Kubernetes\n(concepts only)"]
    P1 --> C2["Ch 2 — Understanding Containers\nDocker · kiada v0.1 Node.js app\nDockerfile · pod.kiada.yaml"]
    P1 --> C3["Ch 3 — Deploying Your First App\nkubectl · kind cluster · Deployment\ncreate-deployment.sh"]
    P1 --> C4["Ch 4 — Navigating the K8s API\n(concepts only)"]

    P2 --> C5["Ch 5 — Running Apps with Pods\nkiada v0.2 · init containers\nsidecar containers · SSL proxy"]
    P2 --> C6["Ch 6 — Pod Lifecycle & Health\nkiada v0.3 · liveness probes\nstartup probes · SIGTERM handler"]
    P2 --> C7["Ch 7 — Namespaces & Labels\nkiada-suite: kiada + quiz + quote\nnode affinity · annotations"]

    P3 --> C8["Ch 8 — ConfigMaps & Secrets\nkiada v0.4 · env vars · downward API\nTLS secrets · image pull secrets"]
    P3 --> C9["Ch 9 — Volumes\nquiz-api v0.1 Go app · MongoDB\nemptyDir · hostPath · projected"]
    P3 --> C10["Ch 10 — PersistentVolumes\nPVC · StorageClass · VolumeSnapshot\nlocal PV · dynamic provisioning"]

    P4 --> C11["Ch 11 — Services\nkiada v0.5 · ClusterIP · NodePort\nLoadBalancer · headless · DNS"]
    P4 --> C12["Ch 12 — Ingress\nIngress NGINX · TLS · host rules\npath routing · affinity"]
    P4 --> C13["Ch 13 — Gateway API\nHTTPRoute · TLSRoute · TCPRoute\nGateway · ReferenceGrant · Istio"]

    P5 --> C14["Ch 14 — ReplicaSets\nrs.kiada.yaml · pod ownership\nlabel selectors · scaling"]
    P5 --> C15["Ch 15 — Deployments\nkiada v0.6/0.7/0.8 · RollingUpdate\nRecreate · canary · rollback"]
    P5 --> C16["Ch 16 — StatefulSets\nquiz-api v0.2 · MongoDB replSet\nordered/parallel · PVC retention"]
    P5 --> C17["Ch 17 — DaemonSets\nkiada v0.9 · node-agent · hostNetwork\nhostPort · tolerations"]
    P5 --> C18["Ch 18 — Jobs & CronJobs\nquiz-init · aggregate-responses\nindexed jobs · work queues"]

    classDef part fill:#3b82d4,color:#fff,stroke:#2563eb
    classDef chapter fill:#f7f8fa,stroke:#e5e7eb,color:#1f2328
    class P1,P2,P3,P4,P5 part
    class C1,C2,C3,C4,C5,C6,C7,C8,C9,C10,C11,C12,C13,C14,C15,C16,C17,C18 chapter
```

---

## Kiada Application Architecture (End-to-End)

The **kiada** demo application suite introduced progressively across the book chapters:

```mermaid
flowchart LR
    subgraph Client["Client Layer"]
        Browser["🌐 Web Browser / curl"]
    end

    subgraph Ingress["Ingress / Gateway Layer (Ch 12-13)"]
        ING["Ingress\nkiada.example.com\napi.example.com"]
    end

    subgraph KiadaNS["Namespace: kiada"]
        subgraph KiadaDeploy["Deployment: kiada (Ch 15)\nReplicas: 3"]
            KP1["Pod: kiada\n📦 kiada-go:1.0\n📦 envoy (TLS proxy)"]
            KP2["Pod: kiada\n📦 kiada-go:1.0\n📦 envoy (TLS proxy)"]
            KP3["Pod: kiada\n📦 kiada-go:1.0\n📦 envoy (TLS proxy)"]
        end
        KiadaSvc["Service: kiada\nClusterIP :80/:443"]

        subgraph QuizDeploy["Deployment: quiz (Ch 15)\nReplicas: 1"]
            QP1["Pod: quiz\n📦 quiz-api (Go)\n📦 mongo:7"]
        end
        QuizSvc["Service: quiz\nClusterIP :80"]

        subgraph QuoteDeploy["Deployment: quote (Ch 15)\nReplicas: 3"]
            QTP1["Pod: quote\n📦 quote-writer\n📦 nginx"]
            QTP2["Pod: quote\n📦 quote-writer\n📦 nginx"]
            QTP3["Pod: quote\n📦 quote-writer\n📦 nginx"]
        end
        QuoteSvc["Service: quote\nClusterIP :80"]

        CM["ConfigMap\nkiada-ssl-config\n(envoy.yaml)"]
        SEC["Secret\nkiada-tls\n(TLS cert/key)"]
    end

    Browser -->|"HTTPS :443"| ING
    ING -->|"/"| KiadaSvc
    ING -->|"/quote"| QuoteSvc
    ING -->|"/questions"| QuizSvc
    KiadaSvc --> KP1
    KiadaSvc --> KP2
    KiadaSvc --> KP3
    KP1 -->|"http://quote/quote"| QuoteSvc
    KP1 -->|"http://quiz/questions/random"| QuizSvc
    QuizSvc --> QP1
    QuoteSvc --> QTP1
    QuoteSvc --> QTP2
    QuoteSvc --> QTP3
    CM -->|"mounted"| KP1
    SEC -->|"mounted"| KP1
```

---

## Go App (kiada-go) Component Diagram

```mermaid
flowchart TD
    subgraph kiada-go["kiada-go Application"]
        main["main.go\nFlag parsing · startup logging\nHTTP server init"]
        server["server.go\nHTTP handlers\nGET / · GET /healthz/ready\nGET /info · GET /quote (proxy)\nGET /quiz (proxy)"]
        handlers["handlers.go\nJSON response helpers\nerror wrappers\nenv-var injection"]
        main --> server
        main --> handlers
        server --> handlers
    end

    subgraph EnvVars["Environment Variables"]
        E1["LISTEN_PORT (default :8080)"]
        E2["POD_NAME, POD_IP"]
        E3["NODE_NAME, NODE_IP"]
        E4["QUOTE_URL"]
        E5["QUIZ_URL"]
        E6["INITIAL_STATUS_MESSAGE"]
    end

    subgraph Upstream["Upstream Services"]
        QUOTE["quote service\nhttp://quote/quote"]
        QUIZ["quiz service\nhttp://quiz/questions/random"]
    end

    EnvVars --> main
    server -->|"HTTP proxy GET"| QUOTE
    server -->|"HTTP proxy GET"| QUIZ
```

---

## Kubernetes Resource Relationships (Ch 14-17)

```mermaid
flowchart TD
    Deploy["Deployment\nkiada"] -->|"owns"| RS["ReplicaSet\nkiada-XXXX"]
    RS -->|"creates/owns"| P1["Pod\nkiada-XXXX-yyy"]
    RS -->|"creates/owns"| P2["Pod\nkiada-XXXX-zzz"]
    RS -->|"creates/owns"| P3["Pod\nkiada-XXXX-www"]
    P1 & P2 & P3 -->|"selected by"| SVC["Service\nkiada"]
    SVC -->|"routes to"| EP["Endpoints"]
    ING["Ingress / HTTPRoute"] -->|"backend"| SVC
    CM["ConfigMap"] -->|"volume/env"| P1
    SEC["Secret"] -->|"volume/env"| P1
    PVC["PersistentVolumeClaim"] -->|"bound"| PV["PersistentVolume"]
    PV -->|"mounted"| P1
```

---

## Deployment Progression (kiada versions across chapters)

```mermaid
timeline
    title kiada Application Evolution
    Chapter 2  : kiada v0.1 — Simple HTTP server (Node.js), serve hostname/IP
    Chapter 5  : kiada v0.2 — Add sidecar SSL proxy (envoy), init containers
    Chapter 6  : kiada v0.3 — Add liveness & startup probes, SIGTERM handler
    Chapter 8  : kiada v0.4 — ConfigMap env-vars, downward API, pod/node metadata
    Chapter 11 : kiada v0.5 — Add /proxy/quote, /proxy/quiz, readiness probe
    Chapter 15 : kiada v0.6/0.7/0.8 — Deployments, rolling updates, canary
    Chapter 17 : kiada v0.9 — DaemonSet with hostPort, per-node agent
```
