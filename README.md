# Argus

Go로 구현한 경량 인프라 모니터링 시스템입니다. Prometheus에서 영감을 받아, Agent가 각 서버의 시스템 메트릭을 수집하고 중앙 서버(Argus Server)로 실시간 전송하는 구조로 설계되었습니다.

<br>

## 아키텍처

```
┌─────────────────┐         gRPC Streaming            ┌──────────────────────┐
│   Argus Agent   │  ────────────────────────────►    │   Argus Server       │
│                 │                                   │                      │
│ - CPU 수집      │    MetricBatch (protobuf)         │ - Agent 등록/관리     │
│ - Memory 수집   │    {agent_id, hostname,           │ - 저장소 선택 가능    │
│ - Disk 수집     │      []Metric{name,value,labels}} │ - Agent 상태 추적     │
└─────────────────┘                                   └──────────┬────────────┘
                                                                 │
                                                    ┌────────────▼────────────┐
                                                    │  Memory / InfluxDB / S3 │
                                                    │     (환경변수로 선택)    │
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
    ├── config/
    │   └── config.go                # 환경변수 중앙 관리
    ├── domain/
    │   ├── metric.go                # ArgusMetric 타입 정의
    │   └── labels.go                # Labels 타입 정의
    ├── cmd/agent/
    │   └── main.go
    └── internal/
        ├── collector/               # 메트릭 수집
        │   ├── collector.go         # Collector 인터페이스
        │   ├── cpu.go
        │   ├── memory.go
        │   └── disk.go
        ├── processor/               # 메트릭 후처리 (라벨 추가 등)
        │   ├── processor.go         # Processor 인터페이스
        │   └── simple_processor.go  # 기본 라벨 추가 구현체
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

## 메트릭 구조

기존 고정 필드(cpu_usage, mem_usage, disk_usage) 방식에서 `name / value / labels` 구조로 변경하여 커스텀 메트릭을 자유롭게 추가할 수 있습니다.

```protobuf
message Metric {
  string name                = 1;
  double value               = 2;
  int64  timestamp           = 3;
  map<string, string> labels = 4;  // 커스텀 라벨
}

message MetricBatch {
  string   agent_id = 1;
  string   hostname = 2;
  repeated Metric metrics = 3;  // 배치 전송
}
```

Agent는 수집한 메트릭을 `MetricBatch`로 묶어서 한 번에 전송합니다.

<br>

## 수집 메트릭

| 이름 | 설명 |
|---|---|
| `cpu_usage` | CPU 사용률 (%) |
| `mem_usage` | 메모리 사용률 (%) |
| `disk_usage` | 디스크 사용률 (%) |

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

InfluxDB에서 메트릭은 `name`을 Field명으로, 커스텀 라벨은 Tag로 저장합니다.

```
measurement: metrics
tags:    agent_id, hostname, {커스텀 라벨}
fields:  {metric name} = {metric value}
```

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

## 시작하기

### 요구사항

- Go 1.24+
- Docker & Docker Compose
- protoc (Protocol Buffers 컴파일러)
- protoc-gen-go, protoc-gen-go-grpc

### 설치

```bash
# protoc Go 플러그인 설치
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 저장소 클론
git clone https://github.com/noboaki/argus.git
cd argus

# proto 컴파일
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

### 1. 메트릭 구조 - 고정 필드 → name/value/labels

기존에는 CPU/메모리/디스크가 proto의 고정 필드였습니다. 새로운 메트릭을 추가하려면 proto 파일과 핸들러를 수정해야 했습니다. `name / value / labels` 구조로 변경하여 proto와 핸들러 수정 없이 어떤 메트릭이든 추가할 수 있습니다.

### 2. Processor 레이어

수집(Collector)과 전송(Sender) 사이에 Processor 레이어를 두어 메트릭 후처리를 담당합니다. 현재는 커스텀 라벨을 추가하는 `SimpleProcessor`가 구현되어 있으며, 메트릭 필터링이나 변환 등의 기능을 추가할 때 기존 코드를 수정하지 않아도 됩니다.

```
Collector → Processor → Sender
  수집       후처리      전송
```

### 3. 배치 전송 (MetricBatch)

메트릭을 한 건씩 전송하면 네트워크 오버헤드가 크고, 저장소 Lock/Unlock이 메트릭 수만큼 반복됩니다. `MetricBatch`로 묶어서 한 번에 전송하고 저장소도 배치 단위로 한 번만 Lock합니다.

### 4. 인메모리 저장 구조 - 이중 맵

Agent ID와 메트릭 이름 두 가지로 인덱싱하여 특정 Agent의 특정 메트릭을 빠르게 조회합니다.

```go
// agentID → metricName → []Metric
map[string]map[string][]*proto.Metric
```

### 5. Timestamp 변환 시점

각 레이어가 자신에게 맞는 타입을 사용합니다. Go 코드에서는 `time.Time`, 전송 구간에서만 `int64` Unix timestamp를 사용합니다.

```
ArgusMetric.Timestamp (time.Time)
        │  sender.go에서 proto 변환 시
        ▼
proto.Metric.Timestamp (int64 Unix)  ← 전송
        │  InfluxDB 저장 시
        ▼
time.Unix(m.Timestamp, 0)  ← 복원
```

### 6. MetricStore / AgentStore 인터페이스 분리

S3처럼 실시간 상태 관리에 부적합한 백엔드도 Agent 상태 메서드를 구현해야 하는 문제를 피하기 위해 두 인터페이스로 분리했습니다. `AgentStore`는 실시간 접근이 필요하므로 어떤 백엔드를 선택하든 항상 인메모리를 사용합니다.

### 7. 환경변수 중앙 관리 (config.go)

`os.Getenv`가 코드 전체에 흩어지지 않도록 각 서비스마다 `config/config.go`에서 모든 환경변수를 한 곳에 정의하고 시작 시 유효성 검사를 수행합니다.

### 8. S3와 MinIO 동시 지원

AWS SDK를 그대로 사용하되 `S3_ENDPOINT` 환경변수 유무로 AWS S3와 MinIO를 자동 분기합니다. 로컬/개발 환경은 MinIO, 운영 환경은 AWS S3로 환경변수만 바꿔서 전환할 수 있습니다.

### 9. Agent ID 전략

호스트명만으로 Agent를 식별하면 컨테이너 환경에서 Pod 재시작 시 ID가 바뀌는 문제가 있습니다. `ARGUS_AGENT_ID` 환경변수를 우선 사용하고, 없으면 호스트명으로 폴백합니다.

### 10. Agent 상태 추적

gRPC 스트림의 생명주기를 Agent 상태와 연결했습니다. `defer`로 스트림이 정상 종료되든 비정상 종료되든 반드시 오프라인 처리가 되도록 보장합니다.

### 11. Agent 자동 재연결

네트워크 단절이나 서버 재시작 시 Agent가 자동으로 재연결을 시도합니다.

```go
func runWithRetry() {
    for {
        s, err := sender.New(cfg)
        if err != nil {
            log.Printf("연결 실패, 5초 후 재시도: %v", err)
            time.Sleep(5 * time.Second)
            continue
        }
        if err := run(s, cfg); err != nil {
            log.Printf("스트림 에러, 5초 후 재연결: %v", err)
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
