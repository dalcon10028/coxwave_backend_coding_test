# Coxwave Backend Coding Test

> Coupon Issuance System

## Stack

- Go
- SQLite 
- connectrpc

## Problem Solving Approach

### 동시성 제어 전략

- UPDATE 쿼리를 통한 낙관적 락(Optimistic Locking) 구현:
  ```sql
  UPDATE campaign 
  SET remaining_coupons = remaining_coupons - 1
  WHERE id = ? AND remaining_coupons > 0 AND start_at <= ?
  ```
- Row Lock으로 동시성 보장
- 별도의 락 메커니즘 없이 DB 레벨에서 동시성 제어

### 테스트

- 2000개의 동시 요청으로 테스트 (1000개 성공 예상)
- 고루틴을 사용한 동시 요청 시뮬레이션
- sync.Mutex를 통한 카운터 및 에러 맵 동기화

## Getting Started

### Install Dependencies
```bash
go mod download
```

### Run Server
```bash
go run cmd/server/main.go
```

### Run Client

```bash
go run cmd/client/main.go
```