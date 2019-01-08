// User: SubLuLu
// Date: 2019-01-08 14:53
package grapheme_splitter

const (
	NotBreak = iota
	BreakStart
	Break
	BreakLastRegional
	BreakPenultimateRegional
)

// Split breaks the given string into an array of grapheme cluster strings
func Split(str string) []string {
	var res []string
	if len(str) == 0 {
		return res
	}
	
	var index, brk int
	runes := []rune(str)
	l := len(runes)
	
	for {
		brk = nextBreak(str, index)
		if brk < l {
			res = append(res, string(runes[index:brk]))
			index = brk
		} else {
			break
		}
	}
	
	if index < l {
		res = append(res, string(runes[index:]))
	}
	
	return res
}

// Counter returns the number of grapheme clusters there are in the given string
func Counter(str string) int {
	if len(str) == 0 {
		return 0
	}
	
	var count, index, brk int
	runes := []rune(str)
	
	l := len(runes)
	
	for {
		brk = nextBreak(str, index)
		if brk < l {
			index = brk
			count++
		} else {
			break
		}
	}
	
	if index < l {
		count++
	}
	return count
}

func nextBreak(str string, index int) int {
	runes := []rune(str)
	
	l := len(runes)
	
	if index >= l-1 {
		return l
	}
	
	prev := graphemeBreakProperty(runes[index])
	
	var mid = make([]int, 0)
	
	for i := index + 1; i < l; i++ {
		// check for already processed low surrogates
		if isSurrogate(runes[i], runes[i-1]) {
			continue
		}
		
		next := graphemeBreakProperty(runes[i])
		
		if shouldBreak(prev, mid, next) != 0 {
			return i
		}
		
		mid = append(mid, next)
	}
	
	return l
}

func shouldBreak(start int, mid []int, end int) int {
	all := append([]int{start}, append(mid, end)...)
	
	prev, next := all[len(all)-2], end
	
	// Lookahead termintor for:
	// GB10. (E_Base | EBG) Extend* ?	E_Modifier
	eModifierIndex := lastIndex(all, EModifier)
	
	if eModifierIndex > 1 &&
		every(all[1:eModifierIndex], func(c int) bool {
			return c == Extend
		}) && firstIndex([]int{
		Extend,
		EBase,
		EBaseGaz}, start) == -1 {
		return Break
	}
	
	// Lookahead termintor for:
	// GB12. ^ (RI RI)* RI	?	RI
	// GB13. [^RI] (RI RI)* RI	?	RI
	var rIIndex = lastIndex(all, RegionalIndicator)
	
	if rIIndex > 0 &&
		every(all[1:rIIndex], func(c int) bool {
			return c == RegionalIndicator
		}) && firstIndex([]int{
		Prepend,
		RegionalIndicator}, prev) == -1 {
		if len(filter(all, RegionalIndicator))%2 == 1 {
			return BreakLastRegional
		} else {
			return BreakPenultimateRegional
		}
	}
	
	// GB3. CR X LF
	if prev == CR && next == LF {
		return NotBreak;
	} else if prev == Control ||
		prev == CR ||
		prev == LF { // GB4. (Control|CR|LF) ÷
		if next == EModifier &&
			every(mid, func(c int) bool {
				return c == Extend
			}) {
			return Break
		} else {
			return BreakStart
		}
	} else if next == Control ||
		next == CR ||
		next == LF { // GB5. ÷ (Control|CR|LF)
		return BreakStart
	} else if prev == L &&
		(next == L ||
			next == V ||
			next == LV ||
			next == LVT) { // GB6. L X (L|V|LV|LVT)
		return NotBreak
	} else if prev == LV ||
		prev == V &&
			(next == V || next == T) { // GB7. (LV|V) X (V|T)
		return NotBreak
	} else if (prev == LVT || prev == T) &&
		next == T { // GB8. (LVT|T) X (T)
		return NotBreak
	} else if next == Extend || next == ZWJ { // GB9. X (Extend|ZWJ)
		return NotBreak
	} else if next == SpacingMark { // GB9a. X SpacingMark
		return NotBreak
	} else if prev == Prepend { // GB9b. Prepend X
		return NotBreak
	}
	
	// GB10. (E_Base | EBG) Extend* ?	E_Modifier
	var previousNonExtendIndex int
	if firstIndex(all, Extend) != -1 {
		previousNonExtendIndex = lastIndex(all, Extend) - 1
	} else {
		previousNonExtendIndex = len(all) - 2
	}
	
	if firstIndex([]int{
		EBase,
		EBaseGaz},
		all[previousNonExtendIndex]) != -1 &&
		previousNonExtendIndex <= len(all)-3 &&
		every(all[previousNonExtendIndex+1:len(all)-2], func(c int) bool {
			return c == Extend
		}) &&
		next == EModifier {
		return NotBreak
	}
	
	// GB11. ZWJ ? (Glue_After_Zwj | EBG)
	if prev == ZWJ &&
		firstIndex([]int{
			GlueAfterZwj,
			EBaseGaz}, next) != -1 {
		return NotBreak
	}
	
	// GB12. ^ (RI RI)* RI ? RI
	// GB13. [^RI] (RI RI)* RI ? RI
	if firstIndex(mid, RegionalIndicator) != -1 {
		return Break
	}
	
	if prev == RegionalIndicator &&
		next == RegionalIndicator {
		return NotBreak
	}
	
	// GB999. Any ? Any
	return BreakStart
}

func lastIndex(src []int, ele int) int {
	l := len(src)
	if l == 0 {
		return -1
	}
	var index = -1
	
	for i := 0; i <= int(l/2); i++ {
		li := l - i - 1
		
		if src[li] == ele {
			index = li
			break
		}
		
		if src[i] == ele {
			index = i
		}
	}
	
	return index
}

func firstIndex(src []int, ele int) int {
	l := len(src)
	if l == 0 {
		return -1
	}
	var index = -1
	
	for i := 0; i <= int(l/2); i++ {
		li := l - i - 1
		
		if src[i] == ele {
			index = i
			break
		}
		
		if src[li] == ele {
			index = li
		}
	}
	
	return index
}

func every(src []int, f func(int) bool) bool {
	for _, v := range src {
		if !f(v) {
			return false
		}
	}
	return true
}

func filter(arr []int, ele int) []int {
	result := make([]int, 0)
	for _, v := range arr {
		if v == ele {
			result = append(result, v)
		}
	}
	return result
}

func isSurrogate(cur, next rune) bool {
	return 0xd800 <= cur && cur <= 0xdbff &&
		0xdc00 <= next && next <= 0xdfff
}
