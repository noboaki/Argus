# Argus

Go로 구현한 경량 인프라 모니터링 시스템입니다. Prometheus에서 영감을 받아, Agent가 각 서버의 시스템 메트릭을 수집하고 중앙 서버(Argus Server)로 실시간 전송하는 구조로 설계되었습니다.

<br>

## 아키텍처

```
┌─────────────────┐         gRPC Streaming         ┌──────────────────────┐
│   Argus Agent   │  ────────────────────────────► │   Argus Server       │
│                 │                                │                      │
│ - CPU 수집       │   MetricPayload (protobuf)     │ - Agent 등록/관리    │
│ - Memory 수집    │                                │ - 저장소 선택 가능    │
│ - Disk 수집      │                                │ - Agent 상태 추적    │
└─────────────────┘                                 └──────────┬──────────┘
                                                               │
                                                  ┌────────────▼────────────┐
                                                  │  Memory / InfluxDB / S3  │
                                                  │     (환경변수로 선택)     │
                                                  └─────────────────────────┘
```

<br>

## 프로젝트 구조

모노레포(Go Workspace)로 구성되어 있습니다.

```
argus/
├── go.work
├── go.work.sum
├── docker-compose.yaml
├── .dockerignore
├── .gitignore
│
├── proto/                           # 공유 protobuf 정의
│   ├── go.mod
│   ├── go.sum
│   ├── metrics.proto
│   ├── metrics.pb.go                # 자동 생성 (수정 금지)
│   └── metrics_grpc.pb.go           # 자동 생성 (수정 금지)
│
├── argus-server/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── config/
│   │   └── config.go                # 환경변수 중앙 관리
│   ├── cmd/server/
│   │   └── main.go
│   └── internal/
│       ├── grpc/
│       │   └── handler.go           # gRPC 스트림 핸들러
│       └── store/
│           ├── store.go             # MetricStore, AgentStore 인터페이스
│           ├── factory.go           # 백엔드 기반 저장소 선택
│           ├── memory_metric.go     # MetricStore 인메모리 구현
│           ├── memory_agent.go      # AgentStore 인메모리 구현
│           ├── influxdb.go          # MetricStore InfluxDB 구현
│           └── s3.go                # MetricStore S3/MinIO 구현
│
└── argus-agent/
    ├── go.mod
    ├── go.sum
    ├── Dockerfile
    ├── cmd/agent/
    │   └── main.go
    └── internal/
        ├── collector/               # 메트릭 수집
        │   ├── collector.go         # Collector 인터페이스 & Metrics 타입
        │   ├── cpu.go
        │   ├── memory.go
        │   └── disk.go
        └── sender/
            └── sender.go            # gRPC 전송
```

<br>

## 주요 기술

| 항목 | 선택 | 이유 |
|---|---|---|
| 통신 방식 | gRPC Client Streaming | 실시간 단방향 스트리밍, HTTP/2 기반 고성능 |
| 직렬화 | Protocol Buffers | JSON 대비 바이너리 인코딩으로 전송 효율 향상 |
| 메트릭 수집 | gopsutil v4 | 크로스플랫폼 지원, Prometheus/Datadog도 사용 |
| 데이터 저장 | Memory / InfluxDB / S3 | 환경변수로 선택 가능 |

<br>

## 수집 메트릭

- **CPU** 사용률 (%)
- **Memory** 사용률 (%)
- **Disk** 사용률 (%)

<br>

## 환경변수

### Argus Agent

| 변수 | 기본값 | 설명 |
|---|---|---|
| `ARGUS_SERVER_ADDR` | `localhost:50051` | Argus Server 주소 |
| `ARGUS_AGENT_ID` | 호스트명 | Agent 식별자. 없으면 호스트명 사용 |

### Argus Server

| 변수 | 기본값 | 설명 |
|---|---|---|
| `ARGUS_SERVER_PORT` | `50051` | gRPC 포트 |
| `ARGUS_STORE_BACKEND` | `memory` | 저장소 선택 (`memory` \| `influxdb` \| `s3`) |
| `INFLUXDB_URL` | - | InfluxDB 주소 |
| `INFLUXDB_TOKEN` | - | InfluxDB 인증 토큰 |
| `INFLUXDB_ORG` | - | InfluxDB Organization |
| `INFLUXDB_BUCKET` | - | InfluxDB Bucket |
| `AWS_BUCKET` | - | S3 / MinIO 버킷명 |
| `AWS_REGION` | `us-east-1` | S3 / MinIO 리전 |
| `AWS_ACCESS_KEY_ID` | - | S3 / MinIO Access Key |
| `AWS_SECRET_ACCESS_KEY` | - | S3 / MinIO Secret Key |
| `S3_ENDPOINT` | - | MinIO 사용 시 엔드포인트. 없으면 AWS S3 |

<br>

## 저장소 선택

`ARGUS_STORE_BACKEND` 환경변수로 저장소를 선택합니다.

```
ARGUS_STORE_BACKEND=memory    → 인메모리 (기본값, 재시작 시 데이터 유실)
ARGUS_STORE_BACKEND=influxdb  → InfluxDB (시계열 데이터에 최적화)
ARGUS_STORE_BACKEND=s3        → S3 / MinIO (장기 보관에 적합)
```

`AgentStore`(Agent 상태 관리)는 실시간 접근이 필요하므로 어떤 백엔드를 선택하든 항상 인메모리를 사용합니다.

```
MetricStore  →  ARGUS_STORE_BACKEND에 따라 선택
AgentStore   →  항상 인메모리
```

<br>

## 시작하기

### 요구사항

- Go 1.24+
- Docker & Docker Compose
- protoc (Protocol Buffers 컴파일러)
- protoc-gen-go, protoc-gen-go-grpc

### protoc 플러그인 설치

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 저장소 클론 및 proto 컴파일

```bash
git clone https://github.com/noboaki/argus.git
cd argus

protoc \
  --go_out=. --go_opt=module=argus \
  --go-grpc_out=. --go-grpc_opt=module=argus \
  proto/metrics.proto
```

### 실행 (인메모리, 기본값)

```bash
docker compose up --build
```

### 실행 (InfluxDB)

```bash
ARGUS_STORE_BACKEND=influxdb \
INFLUXDB_URL=http://influxdb:8086 \
INFLUXDB_TOKEN=mytoken \
INFLUXDB_ORG=argus \
INFLUXDB_BUCKET=metrics \
docker compose up --build
```

### 실행 (MinIO)

```bash
ARGUS_STORE_BACKEND=s3 \
S3_ENDPOINT=http://minio:9000 \
AWS_BUCKET=argus-metrics \
AWS_REGION=us-east-1 \
AWS_ACCESS_KEY_ID=minioadmin \
AWS_SECRET_ACCESS_KEY=minioadmin \
docker compose up --build
```

### 여러 Agent 실행

```bash
# 환경변수로 Agent ID 구분
ARGUS_AGENT_ID=node-prod-1 go run ./cmd/agent
ARGUS_AGENT_ID=node-prod-2 go run ./cmd/agent

# Kubernetes DaemonSet 환경에서는 Pod 이름 자동 주입
env:
  - name: ARGUS_AGENT_ID
    valueFrom:
      fieldRef:
        fieldPath: metadata.name
```

<br>

## 설계 시 고려한 점

### 1. MetricStore / AgentStore 인터페이스 분리

저장소를 단일 인터페이스로 관리하면 S3처럼 실시간 상태 관리에 부적합한 백엔드도 Agent 상태 메서드를 구현해야 하는 문제가 생깁니다. 메트릭 저장과 Agent 상태 관리를 별도 인터페이스로 분리하여 각 백엔드가 자신에게 맞는 역할만 구현하도록 했습니다.

```go
// 백엔드마다 구현 (Memory / InfluxDB / S3)
type MetricStore interface {
    Save(metric Metric) error
    GetByAgent(agentID string) []Metric
    GetLatestMetric(agentID string) (*Metric, error)
}

// 항상 인메모리
type AgentStore interface {
    RegisterAgent(info AgentInfo) error
    UnregisterAgent(agentID string) error
    UpdateLastSeen(agentID string) error
    GetAgents() []AgentInfo
    GetAgentById(agentID string) (*AgentInfo, error)
}
```

### 2. 환경변수 중앙 관리 (config.go)

`os.Getenv`가 코드 전체에 흩어지면 어떤 환경변수가 있는지 파악하기 어렵습니다. `config/config.go`에 모든 환경변수를 한 곳에서 정의하고, 시작 시 유효성 검사를 수행합니다.

```go
func (c *Config) validate() error {
    if c.StoreBackend == "influxdb" && c.InfluxDBURL == "" {
        return fmt.Errorf("influxdb 사용 시 INFLUXDB_URL 필수")
    }
    if c.StoreBackend == "s3" && c.S3Bucket == "" {
        return fmt.Errorf("s3 사용 시 AWS_BUCKET 필수")
    }
    return nil
}
```

### 3. S3와 MinIO 동시 지원

AWS SDK를 그대로 사용하되 `S3_ENDPOINT` 환경변수 유무로 AWS S3와 MinIO를 자동 분기합니다. 로컬/개발 환경은 MinIO, 운영 환경은 AWS S3로 환경변수만 바꿔서 전환할 수 있습니다.

```go
if endpoint != "" {
    o.BaseEndpoint = aws.String(endpoint)
    o.UsePathStyle = true  // MinIO는 path style 필수
}
```

### 4. Agent ID 전략

호스트명만으로 Agent를 식별하면 컨테이너 환경에서 Pod 재시작 시 ID가 바뀌는 문제가 있습니다. `ARGUS_AGENT_ID` 환경변수를 우선 사용하고, 없으면 호스트명으로 폴백합니다.

```go
func resolveAgentID(hostname string) string {
    if id := os.Getenv("ARGUS_AGENT_ID"); id != "" {
        return id
    }
    return hostname
}
```

### 5. Agent 상태 추적

gRPC 스트림의 생명주기를 Agent 상태와 연결했습니다. `defer`로 스트림이 정상 종료되든 비정상 종료되든 반드시 오프라인 처리가 되도록 보장합니다.

```go
defer func() {
    if agentID != "" {
        h.agentStore.UnregisterAgent(agentID)
    }
}()
```

### 6. Collector 인터페이스로 확장성 확보

메트릭 수집기를 인터페이스로 추상화하여 새로운 메트릭(네트워크 I/O, Docker 컨테이너별 메트릭 등)을 추가할 때 기존 코드를 수정하지 않아도 됩니다.

```go
type Collector interface {
    Collect() (float64, error)
    Name() string
}
```

### 7. Agent 자동 재연결

네트워크 단절이나 서버 재시작 시 Agent가 자동으로 재연결을 시도합니다.

```go
func runWithRetry(serverAddr string) {
    for {
        if err := run(s); err != nil {
            log.Printf("스트림 에러, 재연결: %v", err)
            time.Sleep(5 * time.Second)
        }
    }
}
```

<br>

## 현재 한계 및 로드맵

### 현재 한계

- 메트릭 조회 REST API 미구현
- TLS 미적용 (평문 통신)
- 메트릭 보존 건수 제한 미구현 (인메모리 모드에서 장기 운영 시 메모리 증가)
- 대시보드 UI 미구현

### 로드맵

```
1단계 (단기)
├── 메트릭 보존 건수 제한 (Agent당 최대 N건)
└── 메트릭 조회 REST API

2단계 (중기)
├── 수평 확장 (Traefik + Distributed 모드)
├── 대시보드 UI
└── 알림 기능 (임계치 초과 시)

3단계 (장기)
├── TLS 적용
├── 컨테이너별 메트릭 수집 (Docker API 연동)
└── K8s DaemonSet 배포 지원 (kubelet API 연동)
```

<br>

## 참고

- [gopsutil v4](https://github.com/shirou/gopsutil) - 시스템 메트릭 수집
- [gRPC-go](https://github.com/grpc/grpc-go) - gRPC Go 구현체
- [InfluxDB Go Client](https://github.com/influxdata/influxdb-client-go) - InfluxDB 연동
- [AWS SDK Go v2](https://github.com/aws/aws-sdk-go-v2) - S3/MinIO 연동
- [OpenTelemetry OpAMP](https://opentelemetry.io/docs/collector/management/) - Agent 관리 프로토콜 참고