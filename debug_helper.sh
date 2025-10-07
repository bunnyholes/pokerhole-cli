#!/bin/bash

# 포커 게임 디버그 헬퍼 스크립트

LOGS_DIR="./logs"
LATEST_LOG=""

# 최신 로그 파일 찾기
find_latest_log() {
    LATEST_LOG=$(ls -t $LOGS_DIR/*.log 2>/dev/null | head -1)
}

# 실시간 로그 모니터링
watch_logs() {
    find_latest_log
    if [ -z "$LATEST_LOG" ]; then
        echo "로그 파일이 없습니다. 게임을 먼저 실행하세요."
        exit 1
    fi

    echo "=== 로그 모니터링: $LATEST_LOG ==="
    tail -f "$LATEST_LOG" | grep --line-buffered -E "(Player action|AI action|Winner|pot|chips|Round)"
}

# 게임 결과 통계
show_stats() {
    find_latest_log
    if [ -z "$LATEST_LOG" ]; then
        echo "로그 파일이 없습니다."
        exit 1
    fi

    echo "=== 게임 통계 분석 ==="
    echo ""
    echo "총 게임 수:"
    grep -c "Game restart" "$LATEST_LOG"
    echo ""
    echo "플레이어 승리:"
    grep "Winner" "$LATEST_LOG" | grep -c "Player"
    echo ""
    echo "AI 승리:"
    grep "Winner" "$LATEST_LOG" | grep -c "AI"
    echo ""
    echo "올인 횟수:"
    grep -c "Player action: All-in" "$LATEST_LOG"
    echo ""
    echo "폴드 횟수:"
    grep -c "Player action: Fold" "$LATEST_LOG"
}

# 마지막 게임 상세 정보
show_last_game() {
    find_latest_log
    if [ -z "$LATEST_LOG" ]; then
        echo "로그 파일이 없습니다."
        exit 1
    fi

    echo "=== 마지막 게임 상세 정보 ==="
    tail -100 "$LATEST_LOG" | grep -E "(Round|Player action|AI action|pot|chips|Winner|Hand)"
}

# 올인 버그 검증
verify_allin_bug() {
    find_latest_log
    if [ -z "$LATEST_LOG" ]; then
        echo "로그 파일이 없습니다."
        exit 1
    fi

    echo "=== 올인 버그 검증 ==="
    echo ""
    echo "올인 후 currentBet 체크:"
    grep -A 3 "All-in" "$LATEST_LOG" | grep "currentBet"
    echo ""
    echo "올인 후 pot 변화:"
    grep -A 2 "All-in" "$LATEST_LOG" | grep "pot"
}

# 메뉴
case "$1" in
    watch)
        watch_logs
        ;;
    stats)
        show_stats
        ;;
    last)
        show_last_game
        ;;
    verify)
        verify_allin_bug
        ;;
    *)
        echo "포커 게임 디버그 헬퍼"
        echo ""
        echo "사용법: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  watch   - 실시간 로그 모니터링"
        echo "  stats   - 게임 통계 분석"
        echo "  last    - 마지막 게임 상세 정보"
        echo "  verify  - 올인 버그 검증"
        echo ""
        echo "예제:"
        echo "  $0 watch     # 실시간 로그 보기"
        echo "  $0 stats     # 승률 통계"
        ;;
esac
