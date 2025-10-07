# CLI 테스트 가이드

CLI 포커 게임을 쉽게 테스트하는 방법입니다.

## 🎮 빠른 시작

### 1. 일반 플레이
```bash
./poker-client
```

### 2. 자동 시나리오 테스트
```bash
# 올인 버그 테스트 (5판 연속 올인)
./test_scenarios.exp 1

# 정상 플레이 (Call/Raise/Check)
./test_scenarios.exp 2

# 폴드 테스트
./test_scenarios.exp 3

# 10판 연속 랜덤 플레이
./test_scenarios.exp 4
```

### 3. 디버그 모니터링

**실시간 로그 보기:**
```bash
# 터미널 1: 게임 실행
./poker-client

# 터미널 2: 실시간 로그 모니터링
./debug_helper.sh watch
```

**게임 통계 확인:**
```bash
./debug_helper.sh stats
```

출력 예:
```
=== 게임 통계 분석 ===

총 게임 수: 15
플레이어 승리: 7
AI 승리: 8
올인 횟수: 12
폴드 횟수: 3
```

**마지막 게임 상세:**
```bash
./debug_helper.sh last
```

**올인 버그 검증:**
```bash
./debug_helper.sh verify
```

## 🐛 버그 확인 방법

### 올인 버그 테스트

**이전 (버그):**
- Player 올인 1000칩
- currentBet이 20으로 유지됨
- AI가 20만 Call
- **AI 부당 이득**

**현재 (수정됨):**
- Player 올인 1000칩
- currentBet이 1000으로 업데이트
- AI가 1000 Call or Fold
- **공정한 게임**

**확인 방법:**
```bash
# 1. 올인 테스트 실행
./test_scenarios.exp 1

# 2. 로그에서 currentBet 확인
./debug_helper.sh verify
```

### Ace-low Straight 버그 테스트

**버그 시나리오:**
- Player: A-2-3-4-5 (wheel)
- AI: 6-7-8-9-10
- **이전 버그**: Player 승 (잘못됨!)
- **수정 후**: AI 승 (올바름)

로그에서 확인:
```bash
grep "Straight" logs/*.log | tail -20
```

## 📊 승률 분석

100판 플레이 후 통계:
```bash
# 100판 자동 플레이 (백그라운드)
for i in {1..10}; do
    ./test_scenarios.exp 4 > /dev/null 2>&1
done

# 통계 확인
./debug_helper.sh stats
```

**정상 범위:**
- 플레이어 승률: 40-60%
- AI 승률: 40-60%

**비정상 (버그):**
- AI 승률: 80%+ ❌

## 🎯 게임 조작키

| 키 | 액션 |
|---|---|
| `c` | Call (콜) |
| `r` | Raise (레이즈 +50) |
| `a` | All-in (올인) |
| `k` | Check (체크) |
| `f` | Fold (폴드) |
| `F1` | Help (도움말) |
| `n` | New Game (쇼다운 후) |
| `ESC/q` | Quit (종료) |

## 🔍 로그 위치

```bash
ls -lh logs/
```

로그 파일명: `poker-client-<UUID>-<timestamp>.log`

## 💡 팁

**빠른 버그 재현:**
```bash
# 올인 5판 → 통계 확인
./test_scenarios.exp 1 && ./debug_helper.sh stats
```

**실시간 디버깅:**
```bash
# Split terminal
# Left: 게임 플레이
# Right: tail -f logs/poker-client-*.log
```

**특정 시나리오 반복:**
```bash
for i in {1..5}; do
    echo "=== Round $i ==="
    ./test_scenarios.exp 1
    sleep 1
done
```

## 🚨 알려진 이슈

- ~~칩이 0이어도 게임 계속됨~~ ✅ 수정됨
- ~~쇼다운에서 올인/폴드 가능~~ ✅ 수정됨
- ~~AI가 항상 이김 (올인 버그)~~ ✅ 수정됨
- ~~Ace-low straight 승자 뒤집힘~~ ✅ 수정됨

## 📝 테스트 체크리스트

- [ ] 올인 후 AI가 정확한 금액 콜하는가?
- [ ] 칩 0이면 게임 종료되는가?
- [ ] 쇼다운에서 입력이 막히는가?
- [ ] 승률이 대략 50:50인가?
- [ ] A-2-3-4-5 스트레이트가 6-7-8-9-10에게 지는가?

## 🔧 문제 해결

**expect 없음:**
```bash
brew install expect
```

**권한 오류:**
```bash
chmod +x test_scenarios.exp debug_helper.sh
```

**로그 너무 많음:**
```bash
rm logs/*.log
```
