# CLI 테스트 옵션 비교

## 현재 상황
- ✅ Unit 테스트 (67개 통과)
- ✅ Expect 스크립트 (4개 시나리오)
- ❌ 골든 파일 테스트 없음
- ❌ 시각적 회귀 테스트 없음

## 추천 도구 (2025년 최신)

### 1. teatest (★★★★★ 추천)

**설치:**
```bash
go get github.com/charmbracelet/x/exp/teatest@latest
```

**장점:**
- 공식 BubbleTea 테스트 라이브러리
- 골든 파일 자동 관리
- Model 상태 검증 가능
- 빠른 회귀 테스트

**적용 예:**
```go
// internal/ui/model_integration_test.go
package ui_test

import (
    "testing"
    "github.com/charmbracelet/x/exp/teatest"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/bunnyholes/pokerhole/client/internal/ui"
)

func TestFullGameFlow(t *testing.T) {
    m := ui.NewModel(nil, false)
    tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 40))

    // 게임 시작
    tm.Send(ui.SwitchToOfflineModeMsg{})
    tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

    // All-in
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

    // 쇼다운까지 대기
    teatest.WaitFor(t, tm.Output(),
        func(bts []byte) bool {
            return bytes.Contains(bts, []byte("SHOWDOWN"))
        },
        teatest.WithDuration(5*time.Second),
    )

    // 골든 파일과 비교
    teatest.RequireEqualOutput(t, tm.FinalOutput(t))
}
```

**실행:**
```bash
# 최초 실행 (골든 파일 생성)
go test ./internal/ui -run TestFullGameFlow -update

# 회귀 테스트
go test ./internal/ui -run TestFullGameFlow

# 골든 파일 업데이트
go test ./internal/ui -run TestFullGameFlow -update
```

**골든 파일 위치:**
```
internal/ui/testdata/TestFullGameFlow.golden
```

---

### 2. VHS - 데모 + 회귀 테스트 (★★★★☆)

**설치:**
```bash
brew install vhs
# 또는
go install github.com/charmbracelet/vhs@latest
```

**적용 예 (poker_demo.tape):**
```tape
Output poker_demo.gif
Output poker_test.txt

Set FontSize 24
Set Width 1200
Set Height 800
Set Theme "Catppuccin Mocha"

Type "./poker-client"
Enter
Sleep 1s

# 게임 시작
Enter
Sleep 1s

# All-in 테스트
Type "a"
Sleep 2s

Screenshot poker_allin.png

# 새 게임
Type "n"
Sleep 1s

# 종료
Ctrl+C
Sleep 1s

# 테스트: ASCII 출력 확인
Assert "SHOWDOWN"
Assert "칩"
```

**실행:**
```bash
# GIF + ASCII 출력 생성
vhs poker_demo.tape

# ASCII만 (회귀 테스트)
vhs poker_demo.tape --output poker_test.txt

# CI에서 실행
vhs poker_demo.tape --quiet
```

**CI 통합 (.github/workflows/cli-test.yml):**
```yaml
name: CLI Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install VHS
        run: |
          go install github.com/charmbracelet/vhs@latest
          sudo apt-get install -y ttyd ffmpeg

      - name: Build
        run: go build -o poker-client cmd/poker-client/main.go

      - name: Run VHS tests
        run: vhs poker_demo.tape --output test.txt

      - name: Check output
        run: |
          grep "SHOWDOWN" test.txt
          grep "칩" test.txt
```

**장점:**
- README용 GIF 자동 생성
- 데모 자동화
- ASCII 출력으로 회귀 테스트
- CI/CD 통합

---

### 3. catwalk (★★★☆☆)

**설치:**
```bash
go get github.com/knz/catwalk@latest
```

**적용 예:**
```go
// internal/ui/model_catwalk_test.go
package ui_test

import (
    "testing"
    "github.com/knz/catwalk"
    "github.com/bunnyholes/pokerhole/client/internal/ui"
)

func TestPokerGameFlow(t *testing.T) {
    catwalk.RunTest(t, "testdata/poker_allin",
        func() tea.Model {
            return ui.NewModel(nil, false)
        })
}
```

**testdata/poker_allin:**
```
# 게임 시작
enter
----
Output contains "라운드: PRE_FLOP"

# All-in
type a
----
Output contains "All-in"
Model.offlineGame.pot > 0

# 쇼다운
wait 2s
----
Output contains "SHOWDOWN"
```

**실행:**
```bash
# 테스트 실행
go test ./internal/ui -run TestPokerGameFlow

# 출력 업데이트
go test ./internal/ui -run TestPokerGameFlow -rewrite
```

---

## 비교표

| 도구 | 설치 | 골든 파일 | 상태 검증 | CI/CD | 데모 | 난이도 |
|-----|------|---------|---------|-------|-----|-------|
| **teatest** | ⭐️⭐️⭐️⭐️⭐️ | ✅ | ✅ | ✅ | ❌ | 쉬움 |
| **VHS** | ⭐️⭐️⭐️⭐️ | ✅ (ASCII) | ❌ | ✅ | ✅ (GIF) | 보통 |
| **catwalk** | ⭐️⭐️⭐️ | ✅ | ✅ | ✅ | ❌ | 보통 |
| **Expect** | ⭐️⭐️⭐️⭐️ | ❌ | ❌ | ✅ | ❌ | 어려움 |

## 추천 조합

### 옵션 A: 빠른 적용 (teatest만)
```bash
go get github.com/charmbracelet/x/exp/teatest@latest
```
- 골든 파일 테스트 추가
- 회귀 테스트 자동화
- 가장 쉬움

### 옵션 B: 완벽한 테스트 (teatest + VHS)
```bash
go get github.com/charmbracelet/x/exp/teatest@latest
brew install vhs
```
- teatest로 회귀 테스트
- VHS로 데모 GIF 자동 생성
- CI/CD에서 둘 다 실행

### 옵션 C: 현재 유지 (Expect + 수동)
```bash
# 현재 방식 유지
./test_scenarios.exp 1
./debug_helper.sh stats
```
- 추가 설치 없음
- 수동 테스트 계속

## 다음 단계

### 1. teatest 적용 (추천)
```bash
# 1. 설치
go get github.com/charmbracelet/x/exp/teatest@latest

# 2. 테스트 작성
# internal/ui/model_teatest_test.go 생성

# 3. 골든 파일 생성
go test ./internal/ui -run TestFullGameFlow -update

# 4. CI에 추가
# .github/workflows/test.yml
```

### 2. VHS 데모 (선택)
```bash
# 1. 설치
brew install vhs

# 2. 데모 스크립트 작성
# poker_demo.tape

# 3. GIF 생성
vhs poker_demo.tape

# 4. README에 추가
```

## 결론

**지금 당장 적용할 것:**
1. ✅ **teatest** - 공식, 쉬움, 강력함
2. 🔄 **VHS** - 데모 필요하면 추가

**나중에 고려:**
- catwalk (teatest로 충분하면 불필요)

**현재 Expect는:**
- 유지해도 됨 (간단한 시나리오용)
- teatest로 대체 가능
