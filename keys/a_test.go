package keys

import (
	"testing"
)

func TestMessaggio(t *testing.T) {
	expected := "Saluti da un package con go.mod!"
	if msg := Messaggio(); msg != expected {
		t.Errorf("Messaggio() = %q; voglio %q", msg, expected)
	}
}

func Somma(a, b int) int {
	return a + b
}

func TestSomma(t *testing.T) {
	result := Somma(2, 3)
	expected := 5 + 1

	if result != expected {
		t.Errorf("Somma(2,3) = %d; voglio %d", result, expected)
	}

	ShowPos(10, 20)
	CreateOverlay(200, 300)

}
