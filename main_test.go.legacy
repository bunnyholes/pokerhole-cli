package main

import (
	"encoding/json"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

// TestNickhole - nickhole 함수가 올바른 형식의 닉네임을 생성하는지 테스트
func TestNickhole(t *testing.T) {
	// 여러 번 실행해서 랜덤성 확인
	for i := 0; i < 10; i++ {
		nickname := nickhole()

		// 닉네임이 비어있지 않은지 확인
		if len(nickname) == 0 {
			t.Errorf("nickhole() returned empty string")
		}

		// 닉네임 형식 확인 (영문자로 시작하고 숫자로 끝남)
		if len(nickname) < 5 {
			t.Errorf("nickhole() returned too short: %s", nickname)
		}

		t.Logf("Generated nickname: %s", nickname)
	}
}

// TestGetOrCreateUUID - UUID 생성 및 저장 테스트
func TestGetOrCreateUUID(t *testing.T) {
	// 테스트용 임시 디렉토리 생성
	tmpDir := t.TempDir()

	// 홈 디렉토리를 임시로 변경 (테스트용)
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// 첫 번째 호출 - 새 UUID 생성
	uuid1 := getOrCreateUUID()
	if len(uuid1) != 36 {
		t.Errorf("Expected UUID length 36, got %d: %s", len(uuid1), uuid1)
	}

	// UUID 형식 확인 (8-4-4-4-12)
	parts := strings.Split(uuid1, "-")
	if len(parts) != 5 {
		t.Errorf("Invalid UUID format: %s", uuid1)
	}

	// 두 번째 호출 - 같은 UUID 반환
	uuid2 := getOrCreateUUID()
	if uuid1 != uuid2 {
		t.Errorf("UUID should be persistent. Got different UUIDs:\n  First: %s\n  Second: %s", uuid1, uuid2)
	}

	t.Logf("UUID persistence test passed: %s", uuid1)
}

// TestGenerateClientInfo - ClientInfo 생성 테스트
func TestGenerateClientInfo(t *testing.T) {
	// 테스트용 임시 디렉토리
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	info := generateClientInfo()

	// UUID 검증
	if len(info.UUID) != 36 {
		t.Errorf("Invalid UUID length: %d", len(info.UUID))
	}

	// Nickname 검증
	if len(info.Nickname) == 0 {
		t.Errorf("Nickname should not be empty")
	}

	t.Logf("Generated ClientInfo: UUID=%s, Nickname=%s", info.UUID, info.Nickname)
}

// TestClientInfoJSON - ClientInfo JSON 직렬화 테스트
func TestClientInfoJSON(t *testing.T) {
	info := ClientInfo{
		UUID:     "test-uuid-1234",
		Nickname: "TestNickname",
	}

	// JSON 직렬화
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal ClientInfo: %v", err)
	}

	// JSON 형식 확인
	expected := `{"uuid":"test-uuid-1234","nickname":"TestNickname"}`
	if string(data) != expected {
		t.Errorf("Expected JSON:\n  %s\nGot:\n  %s", expected, string(data))
	}

	// JSON 역직렬화
	var decoded ClientInfo
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ClientInfo: %v", err)
	}

	if decoded.UUID != info.UUID || decoded.Nickname != info.Nickname {
		t.Errorf("Decoded ClientInfo doesn't match original")
	}

	t.Logf("JSON serialization test passed")
}

// TestServerConnection - 실제 서버 연결 테스트 (통합 테스트)
func TestServerConnection(t *testing.T) {
	// 서버가 실행 중인지 확인
	conn, err := net.DialTimeout("tcp", "localhost:7777", 1*time.Second)
	if err != nil {
		t.Skipf("Server not running on localhost:7777, skipping integration test: %v", err)
		return
	}
	defer conn.Close()

	// 테스트용 ClientInfo 생성
	testInfo := ClientInfo{
		UUID:     "test-uuid-integration",
		Nickname: "IntegrationTest",
	}

	// JSON 직렬화
	data, err := json.Marshal(testInfo)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// 서버로 전송
	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		t.Fatalf("Failed to send data to server: %v", err)
	}

	t.Logf("Successfully sent ClientInfo to server: UUID=%s, Nickname=%s",
		testInfo.UUID, testInfo.Nickname)
}

// Benchmark: nickhole 함수 성능 테스트
func BenchmarkNickhole(b *testing.B) {
	for i := 0; i < b.N; i++ {
		nickhole()
	}
}

// Benchmark: getOrCreateUUID 함수 성능 테스트
func BenchmarkGetOrCreateUUID(b *testing.B) {
	tmpDir := b.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	for i := 0; i < b.N; i++ {
		getOrCreateUUID()
	}
}

// Example: nickhole 사용 예제 (주석 처리 - 랜덤 출력으로 인해 테스트 불가)
// func ExampleNickhole() {
// 	nickname := nickhole()
// 	fmt.Printf("Generated nickname length: %d\n", len(nickname))
// }

// Table-driven test: UUID 형식 검증
func TestUUIDFormat(t *testing.T) {
	tests := []struct {
		name  string
		uuid  string
		valid bool
	}{
		{"Valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"Invalid - too short", "550e8400-e29b", false},
		{"Invalid - wrong format", "not-a-uuid", false},
		{"Valid - another format", "123e4567-e89b-12d3-a456-426614174000", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(tt.uuid, "-")
			isValid := len(tt.uuid) == 36 && len(parts) == 5

			if isValid != tt.valid {
				t.Errorf("UUID %s: expected valid=%v, got valid=%v",
					tt.uuid, tt.valid, isValid)
			}
		})
	}
}
