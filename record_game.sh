#!/bin/bash

# 게임 녹화 스크립트
# ttyrec나 asciinema가 설치되어 있어야 합니다

echo "=== 포커 게임 녹화 도구 ==="
echo ""
echo "사용 가능한 녹화 도구:"
echo "1. asciinema (추천)"
echo "2. script (기본 제공)"
echo ""

read -p "선택 (1/2): " choice

case $choice in
    1)
        if command -v asciinema &> /dev/null; then
            echo "asciinema로 녹화를 시작합니다..."
            asciinema rec poker_game_$(date +%Y%m%d_%H%M%S).cast -c "./poker-client"
        else
            echo "asciinema가 설치되지 않았습니다."
            echo "설치: brew install asciinema"
        fi
        ;;
    2)
        output_file="poker_game_$(date +%Y%m%d_%H%M%S).log"
        echo "script로 녹화를 시작합니다... (저장: $output_file)"
        script -q "$output_file" ./poker-client
        echo ""
        echo "녹화 완료: $output_file"
        ;;
    *)
        echo "잘못된 선택입니다."
        ;;
esac
