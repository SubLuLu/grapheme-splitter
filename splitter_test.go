// User: SubLuLu
// Date: 2019-01-08 15:36
package grapheme_splitter

import (
	"testing"
)

var (
	English = data{
		source: "Love me, love my dog",
		length: 20,
		graphemes: []string{
			"L", "o", "v", "e", " ",
			"m", "e", ",", " ",
			"l", "o", "v", "e", " ",
			"m", "y", " ",
			"d", "o", "g",
		},
	}
	Chinese = data{
		source: "学而时习之，不亦乐乎",
		length: 10,
		graphemes: []string{
			"学", "而", "时", "习", "之", "，",
			"不", "亦", "乐", "乎",
		},
	}
	Emoji = data{
		source: "👨‍👨‍👧‍👦👨‍👩‍👦💏🙍",
		length: 4,
		graphemes: []string{
			"👨‍👨‍👧‍👦",
			"👨‍👩‍👦",
			"💏",
			"🙍",
		},
	}
)

type data struct {
	source    string
	length    int
	graphemes []string
}

func TestSplit(t *testing.T) {
	empty := Split("")
	if !isEqual(empty, []string{}) {
		t.Errorf(`expects "" but result is %v`, empty)
	}
	
	english := Split(English.source)
	if !isEqual(English.graphemes, english) {
		t.Errorf("expects %v but result is %v", English.graphemes, english)
	}
	
	chinese := Split(Chinese.source)
	if !isEqual(Chinese.graphemes, chinese) {
		t.Errorf("expects %v but result is %v", Chinese.graphemes, chinese)
	}
	
	emoji := Split(Emoji.source)
	if !isEqual(Emoji.graphemes, emoji) {
		t.Errorf("expects %v but result is %v", Emoji.graphemes, emoji)
	}
}

func TestCounter(t *testing.T) {
	empty := Counter("")
	if empty != 0 {
		t.Errorf("expects 0 but result is %v", empty)
	}
	
	english := Counter(English.source)
	if english != English.length {
		t.Errorf("expects %d but result is %d", English.length, english)
	}
	
	chinese := Counter(Chinese.source)
	if chinese != Chinese.length {
		t.Errorf("expects %d but result is %d", Chinese.length, chinese)
	}
	
	emoji := Counter(Emoji.source)
	if emoji != Emoji.length {
		t.Errorf("expects %d but result is %d", Emoji.length, emoji)
	}
}

func BenchmarkSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Split(English.source)
		
		Split(Chinese.source)
		
		Split(Emoji.source)
	}
}

func BenchmarkCounter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Counter(English.source)
		
		Counter(Chinese.source)
		
		Counter(Emoji.source)
	}
}

func isEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 && len(b) == 0 {
		return true
	}
	for i, v := range b {
		if a[i] != v {
			return false
		}
	}
	return true
}
