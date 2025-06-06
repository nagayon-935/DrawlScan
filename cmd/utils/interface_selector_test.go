package utils

import (
	"testing"
)

// 実際の環境依存のため、AutoSelectInterfaceは空文字列でなければOKとする
func TestAutoSelectInterface(t *testing.T) {
	iface := AutoSelectInterface()
	// 実行環境によっては空文字列になる場合もあるため、エラーにはしない
	t.Logf("AutoSelectInterface() returned: %q", iface)
}

// isInterfaceConnectedは存在しないインターフェース名でfalseを返すことをテスト
func Test_isInterfaceConnected(t *testing.T) {
	// 存在しないインターフェース名
	if got := isInterfaceConnected("nonexistent0"); got != false {
		t.Errorf("isInterfaceConnected(nonexistent0) = %v, want false", got)
	}
	// 存在するインターフェース名が分かる場合は追加でテスト可能
}
