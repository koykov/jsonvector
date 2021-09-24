package jsonvector

import (
	"bytes"
	"testing"
)

type stageUnescape struct {
	origin, expect []byte
}

var (
	stagesUnescape = map[string]*stageUnescape{
		"fmt0": {
			origin: []byte(`Lorem \"ipsum\" dolor \"sit\" amet.`),
			expect: []byte(`Lorem "ipsum" dolor "sit" amet.`),
		},
		"fmt1": {
			origin: []byte(`\tchapter1\npage \\1\ncontents...\n\fpage \\2\rcontents...`),
			expect: []byte("	chapter1\npage \\1\ncontents...\n\fpage \\2\rcontents..."),
		},
		"fmt2": {
			origin: []byte(`unfinished \"example\`),
			expect: []byte("unfinished \"example\\"),
		},
		"unicodeChinese": {
			origin: []byte(`\u6751\u5909\u754C\u5E83\u5171\u6E08\u6975\u65AD\u77E5\u6B62\u904E\u8239\u8FD1\u7D00\u5DE7\u4EA4\u6D77\u8EE2\u9577\u3002`),
			expect: []byte(`村変界広共済極断知止過船近紀巧交海転長。`),
		},
		"unicodeArabic": {
			origin: []byte(`\u0643\u0644\u0651 \u0628\u0642\u0639\u0629 \u0623\u0645\u0644\u0627\u064B \u0627\u0646, \u0623\u0645\u0627 \u0645\u0627 \u064A\u0630\u0643\u0631 \u0646\u0647\u0627\u064A\u0629. \u064A\u0628\u0642 \u0660\u0668\u0660\u0664 \u0645\u0634\u0627\u0631\u0641 \u062A\u0643\u0627\u0644\u064A\u0641 \u062A\u0645, \u062F\u0648\u0644 \u0647\u0648 \u0625\u0639\u0627\u062F\u0629 \u0627\u0644\u064A\u0645\u064A\u0646\u064A \u0644\u064A\u062A\u0633\u0646\u0651\u0649. \u0639\u0631\u0636 \u0628\u062A\u0637\u0648\u064A\u0642 \u0644\u0647\u064A\u0645\u0646\u0629 \u0627\u0644\u0623\u0648\u0631\u0648\u0628\u064A \u0639\u0646. \u062C\u0647\u0629 \u0628\u0644 \u0627\u0644\u0625\u062B\u0646\u0627\u0646 \u0627\u0644\u0641\u0631\u0646\u0633\u064A. \u0623\u0645\u0627\u0645 \u0623\u062C\u0632\u0627\u0621 \u0627\u0644\u0639\u0627\u0644\u0645\u064A\u0629 \u0645\u0646 \u0642\u0627\u0645, \u0645\u064A\u0646\u0627\u0621 \u0648\u0627\u0639\u062A\u0644\u0627\u0621 \u0623\u064A \u0628\u062D\u0642.`),
			expect: []byte(`كلّ بقعة أملاً ان, أما ما يذكر نهاية. يبق ٠٨٠٤ مشارف تكاليف تم, دول هو إعادة اليميني ليتسنّى. عرض بتطويق لهيمنة الأوروبي عن. جهة بل الإثنان الفرنسي. أمام أجزاء العالمية من قام, ميناء واعتلاء أي بحق.`),
		},
		"unicodeGreek": {
			origin: []byte(`\u039B\u03BF\u03C1\u03B5\u03BC \u03B9\u03C0\u03C3\u03B8\u03BC \u03B4\u03BF\u03BB\u03BF\u03C1 \u03C3\u03B9\u03C4 \u03B1\u03BC\u03B5\u03C4, \u03B9\u03B4 c\u03BF\u03BD\u03B3\u03B8\u03B5 \u03B1cc\u03B8\u03C3\u03B1\u03BC v\u03B9\u03BE.`),
			expect: []byte(`Λορεμ ιπσθμ δολορ σιτ αμετ, ιδ cονγθε αccθσαμ vιξ.`),
		},
		"unicodeCyrillic": {
			origin: []byte(`\u041B\u043E\u0440\u0435\u043C \u0438\u043F\u0441\u0443\u043C \u0434\u043E\u043B\u043E\u0440 \u0441\u0438\u0442 \u0430\u043C\u0435\u0442, \u0442\u0435 \u0432\u043E\u0446\u0438\u0431\u0443\u0441 \u043D\u0443\u0441\u044F\u0443\u0430\u043C \u0442\u0438\u0431\u0438\u044F\u0443\u0435 \u0441\u0435\u0430, \u0446\u0443\u043C \u0446\u0443 \u0435\u0438\u0443\u0441 \u0435\u0438\u0440\u043C\u043E\u0434.`),
			expect: []byte(`Лорем ипсум долор сит амет, те воцибус нусяуам тибияуе сеа, цум цу еиус еирмод.`),
		},
		"unicodeSurrogate": {
			origin: []byte(`What is better: \uD834\uDD1E or \uD834\uDD22?`),
			expect: []byte("What is better: 𝄞 or 𝄢?"),
		},
		"complex": {
			origin: []byte(`You can see escaped surrogate characters below\n\u041D\u0438\u0436\u0435 \u0432\u044B \u0443\u0432\u0438\u0434\u0438\u0442\u0435 \u043F\u0440\u0438\u043C\u0435\u0440\u044B \u0437\u0430\u043A\u043E\u0434\u0438\u0440\u043E\u0432\u0430\u043D\u043D\u044B\u0445 \u0441\u0443\u0440\u0440\u043E\u0433\u0430\u0442\u043D\u044B\u0445 \u043F\u0430\u0440:\n\t\uD835\uDC9E - Mathematical script capital C\n\t\uD835\uDCAF - Mathematical script capital T\n\t\uD835\uDCAE - Mathematical script capital S\n\t\uD835\uDC9F - Mathematical script capital D\n\t\uD835\uDCB3 - Mathematical script capital X\n\t\uD834\uDD1E - Musical symbol G clef\n\t\uD834\uDD22 - Musical symbol F clef`),
			expect: []byte(`You can see escaped surrogate characters below
Ниже вы увидите примеры закодированных суррогатных пар:
	𝒞 - Mathematical script capital C
	𝒯 - Mathematical script capital T
	𝒮 - Mathematical script capital S
	𝒟 - Mathematical script capital D
	𝒳 - Mathematical script capital X
	𝄞 - Musical symbol G clef
	𝄢 - Musical symbol F clef`),
		},
	}
)

func getStageUnescape(key string) *stageUnescape {
	if st, ok := stagesUnescape[key]; ok {
		return st
	}
	return nil
}

func testUnescape(tb testing.TB, buf []byte) []byte {
	key := getTBName(tb)
	st := getStageUnescape(key)
	if st == nil {
		tb.Fatal("stage not found")
	}
	buf = append(buf[:0], st.origin...)
	buf = Unescape(buf)
	if !bytes.Equal(buf, st.expect) {
		tb.Error("unescape failed")
	}
	return buf
}

func benchUnescape(b *testing.B) {
	var buf []byte
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf = testUnescape(b, buf)
	}
}

func TestUnescape(t *testing.T) {
	t.Run("fmt0", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("fmt1", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("fmt2", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("unicodeChinese", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("unicodeArabic", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("unicodeGreek", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("unicodeCyrillic", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("unicodeSurrogate", func(t *testing.T) { testUnescape(t, nil) })
	t.Run("complex", func(t *testing.T) { testUnescape(t, nil) })
}

func BenchmarkUnescape(b *testing.B) {
	b.Run("fmt0", func(b *testing.B) { benchUnescape(b) })
	b.Run("fmt1", func(b *testing.B) { benchUnescape(b) })
	b.Run("fmt2", func(b *testing.B) { benchUnescape(b) })
	b.Run("unicodeChinese", func(b *testing.B) { benchUnescape(b) })
	b.Run("unicodeArabic", func(b *testing.B) { benchUnescape(b) })
	b.Run("unicodeGreek", func(b *testing.B) { benchUnescape(b) })
	b.Run("unicodeCyrillic", func(b *testing.B) { benchUnescape(b) })
	b.Run("unicodeSurrogate", func(b *testing.B) { benchUnescape(b) })
	b.Run("complex", func(b *testing.B) { benchUnescape(b) })
}
