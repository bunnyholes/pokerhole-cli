# Golden Test 능력 분석 보고서

## 실험 날짜
2025-10-06

## 목적
골든 테스트가 **텍스트가 같아도 색상(opacity) 값이 다르면** 감지할 수 있는지 검증

---

## 핵심 발견

### ✅ Golden Test가 **가능한 것**

1. **색상 변화 감지** ✅
   - ANSI escape code를 통해 RGB 색상 변화를 정확히 감지
   - 예시: `colorGold` (#D4AF37) vs `colorGoldDim` (#B8960F)
   - ANSI 코드: `\x1b[38;2;211;175;55m` vs `\x1b[38;2;184;150;15m`

2. **스타일 변화 감지** ✅
   - Bold, Italic, Underline 등의 스타일 변화 감지
   - 예시: `\x1b[1;38;2;211;175;55m` (bold) vs `\x1b[38;2;211;175;55m` (normal)

3. **레이아웃 변화 감지** ✅
   - 공백, 줄바꿈, 정렬 등 모든 문자 단위 변화 감지

4. **완전한 회귀 테스트** ✅
   - 의도치 않은 UI 변경사항을 자동으로 감지

### ❌ Golden Test의 **한계**

1. **환경 의존성** ⚠️
   - **문제**: Lipgloss는 비터미널 환경(테스트)에서 자동으로 색상을 끔
   - **결과**: 기본 테스트 환경에서는 색상 코드가 제거됨
   - **해결책**: `lipgloss.SetColorProfile(termenv.TrueColor)` 명시 필요

2. **시각적 검증 불가** ❌
   - Golden test는 문자열만 비교, 실제 터미널 렌더링 결과는 검증 불가
   - 예시: 같은 ANSI 코드라도 터미널마다 다르게 보일 수 있음

3. **시간 의존적 코드** ❌
   - 애니메이션 타이밍, 랜덤값 등은 테스트 불가능
   - 해결책: State-based 디자인 (Tick counter 사용)

---

## 실험 과정

### 1단계: 초기 테스트 (실패)
```go
// model_test.go
SubtitleOpacity: 0.5  // colorGold
```
- **결과**: 테스트 통과 ✅
- **Golden 파일**: 평문만 저장 (ANSI 코드 없음)

### 2단계: Opacity 변경 (예상: 실패, 실제: 통과!)
```go
SubtitleOpacity: 0.2  // colorGoldDim으로 변경
```
- **결과**: 테스트 통과 ✅ (예상과 다름!)
- **원인**: Lipgloss가 색상을 끄고 있었음

### 3단계: 검증 - RenderIntroSubtitle 직접 호출
```go
output05 := intro.RenderIntroSubtitle(0.5)
output02 := intro.RenderIntroSubtitle(0.2)
fmt.Println("Are they equal?", output05 == output02)
```
- **결과**: `Are they equal? false` ✅
- **확인**: ANSI 코드가 실제로 다름
  - 0.5: `1b5b33383b323b3231313b3137353b35356d` 
  - 0.2: `1b5b33383b323b3138343b3135303b31356d`

### 4단계: TrueColor 강제 활성화 (성공!)
```go
lipgloss.SetColorProfile(termenv.TrueColor)
SubtitleOpacity: 0.2  // colorGoldDim
```
- **결과**: 테스트 실패 ❌ (예상대로!)
- **에러 메시지**:
```
expected: \x1b[38;2;211;175;55mT E X A S   H O L D ' E M\x1b[0m
got:      \x1b[38;2;184;150;15mT E X A S   H O L D ' E M\x1b[0m
```

---

## 최종 결론

### Golden Test의 올바른 사용법

#### ✅ DO - 이렇게 사용하세요

1. **Color Profile 명시**
```go
func TestIntroView_Golden(t *testing.T) {
    lipgloss.SetColorProfile(termenv.TrueColor)  // 필수!
    // ... 테스트 코드
}
```

2. **State-based 애니메이션**
```go
state := State{
    CharsRevealed:   5,      // 시간 대신 상태 사용
    SubtitleOpacity: 0.5,
    Tick:            10,
}
```

3. **Pure View Functions**
```go
// ✅ Good: 순수 함수
func (v *View) Render(state State) string {
    return RenderLogo(state.CharsRevealed)
}

// ❌ Bad: side effect 있음
func (v *View) Render() string {
    time.Sleep(100 * time.Millisecond)
    return RenderLogo(time.Now())
}
```

#### ❌ DON'T - 이런 것은 피하세요

1. 시간 의존적 코드
2. 랜덤값 사용
3. 외부 API 호출
4. 파일 시스템 접근

---

## 권장사항

### 프로젝트에 적용할 패턴

1. **모든 View 테스트에서 ColorProfile 설정**
```go
func TestMain(m *testing.M) {
    lipgloss.SetColorProfile(termenv.TrueColor)
    os.Exit(m.Run())
}
```

2. **Opacity 단계별 Golden 테스트 추가** (선택적)
```go
opacityTests := []struct {
    name    string
    opacity float64
}{
    {"fade_start", 0.1},
    {"fade_mid", 0.5},
    {"fade_end", 1.0},
}
```

3. **기존 PhaseSubtitle 테스트 수정**
```go
func TestIntroView_PhaseSubtitle_Golden(t *testing.T) {
    lipgloss.SetColorProfile(termenv.TrueColor)  // 추가!
    // ... 기존 코드
}
```

---

## 최종 답변

> **"텍스트가 똑같더라도 내부에 색상값도 검증하는게 골든 테스트의 목적아닌가?"**

**YES! 맞습니다!** ✅

Golden test는:
- ✅ 텍스트뿐만 아니라 **색상 변화도 감지**합니다
- ✅ ANSI escape code를 통해 **정확한 RGB 값**을 비교합니다
- ⚠️ **단, `lipgloss.SetColorProfile(termenv.TrueColor)` 설정이 필수**입니다

설정 없이는 Lipgloss가 자동으로 색상을 끄기 때문에, 색상 테스트가 무력화됩니다.

---

## 다음 단계

1. ✅ 모든 테스트에 ColorProfile 설정 추가
2. ✅ color_test.go 정리 (opacity 0.5로 복구)
3. ✅ PhaseSubtitle 테스트에 ColorProfile 추가
4. ⏭️ Menu Scene으로 진행 (ROADMAP Phase 3)
