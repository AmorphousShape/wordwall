package tests

import (
	"testing"

	"github.com/AmorphousShape/wordwall/pkg/wordwall"
)

func TestCensor(t *testing.T) {
	wordwall.SetCensoredWords([]string{"badword", "offensive"})

	tests := []struct {
		input     string
		expected  string
		hitCensor bool
	}{
		{"This is a b@dw0rd in a sentence.", "This is a ******* in a sentence.", true},
		{"Nothing to censor here.", "Nothing to censor here.", false},
		{"An 0ff3ns1ve word should be censored.", "An ********* word should be censored.", true},
	}

	for _, test := range tests {
		result, hitCensor, _, _ := wordwall.FilterString(test.input)
		if hitCensor != test.hitCensor || result != test.expected {
			t.Errorf("Expected '%s', got '%s' (hitCensor: %v)", test.expected, result, hitCensor)
		}
	}
}

func TestFilter(t *testing.T) {
	wordwall.SetFilteredWords([]string{"filterme", "ignorethis"})

	tests := []struct {
		input     string
		expected  string
		hitFilter bool
	}{
		{"This should f1lt3rme out.", "", true},
		{"1gn0r3 tH1s message.", "", true},
		{"Not going to be filtered.", "Not going to be filtered.", false},
	}

	for _, test := range tests {
		result, _, hitFilter, _ := wordwall.FilterString(test.input)
		if hitFilter != test.hitFilter || result != test.expected {
			t.Errorf("Expected '%s', got '%s' (hitFilter: %v)", test.expected, result, hitFilter)
		}
	}
}

func TestZeroTolerance(t *testing.T) {
	wordwall.SetZeroToleranceWords([]string{"zerotolerance", "bannedword"})

	tests := []struct {
		input            string
		expected         string
		hitZeroTolerance bool
	}{
		{"This contains a zero tolerance word.", "", true},
		{"A bannedword should trigger an instant ban.", "", true},
		{"This is safe and should not trigger anything.", "This is safe and should not trigger anything.", false},
	}

	for _, test := range tests {
		result, _, _, hitZeroTolerance := wordwall.FilterString(test.input)
		if hitZeroTolerance != test.hitZeroTolerance || result != test.expected {
			t.Errorf("Expected '%s', got '%s' (hitZeroTolerance: %v)", test.expected, result, hitZeroTolerance)
		}
	}
}

func TestLengthSorting(t *testing.T) {
	words := []string{"we", "like", "trains"}
	wordwall.SetCensoredWords(words)

	wordwall.FilterString("trains are cool we like them")

	expectedOutput := "****** are cool ** **** them"
	if result, _, _, _ := wordwall.FilterString("trains are cool we like them"); result != expectedOutput {
		t.Errorf("Expected '%s', got '%s'", expectedOutput, result)
	}
}

func TestNoise(t *testing.T) {
	words := []string{"ice"}

	wordwall.SetCensoredWords(words)
	result, hitCensor, _, _ := wordwall.FilterString("i like ivcven cream")
	if !hitCensor || result != "i like ****** cream" {
		t.Errorf("Expected 'i like ****** cream', got '%s' (hitCensor: %v)", result, hitCensor)
	}
}
