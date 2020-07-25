package jsonvector

import (
	"bytes"
	"testing"
)

var (
	esc0   = []byte(`Lorem \"ipsum\" dolor \"sit\" amet.`)
	unesc0 = []byte(`Lorem "ipsum" dolor "sit" amet.`)
	esc1   = []byte(`\tchapter1\npage \\1\ncontents...\n\fpage \\2\rcontents...`)
	unesc1 = []byte("	chapter1\npage \\1\ncontents...\n\fpage \\2\rcontents...")
	esc2   = []byte(`unfinished \"example\`)
	unesc2 = []byte("unfinished \"example\\")

	escChin   = []byte(`\u6751\u5909\u754C\u5E83\u5171\u6E08\u6975\u65AD\u77E5\u6B62\u904E\u8239\u8FD1\u7D00\u5DE7\u4EA4\u6D77\u8EE2\u9577\u3002`)
	unescChin = []byte(`村変界広共済極断知止過船近紀巧交海転長。`)
	escArab   = []byte(`\u0643\u0644\u0651 \u0628\u0642\u0639\u0629 \u0623\u0645\u0644\u0627\u064B \u0627\u0646, \u0623\u0645\u0627 \u0645\u0627 \u064A\u0630\u0643\u0631 \u0646\u0647\u0627\u064A\u0629. \u064A\u0628\u0642 \u0660\u0668\u0660\u0664 \u0645\u0634\u0627\u0631\u0641 \u062A\u0643\u0627\u0644\u064A\u0641 \u062A\u0645, \u062F\u0648\u0644 \u0647\u0648 \u0625\u0639\u0627\u062F\u0629 \u0627\u0644\u064A\u0645\u064A\u0646\u064A \u0644\u064A\u062A\u0633\u0646\u0651\u0649. \u0639\u0631\u0636 \u0628\u062A\u0637\u0648\u064A\u0642 \u0644\u0647\u064A\u0645\u0646\u0629 \u0627\u0644\u0623\u0648\u0631\u0648\u0628\u064A \u0639\u0646. \u062C\u0647\u0629 \u0628\u0644 \u0627\u0644\u0625\u062B\u0646\u0627\u0646 \u0627\u0644\u0641\u0631\u0646\u0633\u064A. \u0623\u0645\u0627\u0645 \u0623\u062C\u0632\u0627\u0621 \u0627\u0644\u0639\u0627\u0644\u0645\u064A\u0629 \u0645\u0646 \u0642\u0627\u0645, \u0645\u064A\u0646\u0627\u0621 \u0648\u0627\u0639\u062A\u0644\u0627\u0621 \u0623\u064A \u0628\u062D\u0642.`)
	unescArab = []byte(`كلّ بقعة أملاً ان, أما ما يذكر نهاية. يبق ٠٨٠٤ مشارف تكاليف تم, دول هو إعادة اليميني ليتسنّى. عرض بتطويق لهيمنة الأوروبي عن. جهة بل الإثنان الفرنسي. أمام أجزاء العالمية من قام, ميناء واعتلاء أي بحق.`)
	escGrk    = []byte(`\u039B\u03BF\u03C1\u03B5\u03BC \u03B9\u03C0\u03C3\u03B8\u03BC \u03B4\u03BF\u03BB\u03BF\u03C1 \u03C3\u03B9\u03C4 \u03B1\u03BC\u03B5\u03C4, \u03B9\u03B4 c\u03BF\u03BD\u03B3\u03B8\u03B5 \u03B1cc\u03B8\u03C3\u03B1\u03BC v\u03B9\u03BE.`)
	unescGrk  = []byte(`Λορεμ ιπσθμ δολορ σιτ αμετ, ιδ cονγθε αccθσαμ vιξ.`)
	escCyr    = []byte(`\u041B\u043E\u0440\u0435\u043C \u0438\u043F\u0441\u0443\u043C \u0434\u043E\u043B\u043E\u0440 \u0441\u0438\u0442 \u0430\u043C\u0435\u0442, \u0442\u0435 \u0432\u043E\u0446\u0438\u0431\u0443\u0441 \u043D\u0443\u0441\u044F\u0443\u0430\u043C \u0442\u0438\u0431\u0438\u044F\u0443\u0435 \u0441\u0435\u0430, \u0446\u0443\u043C \u0446\u0443 \u0435\u0438\u0443\u0441 \u0435\u0438\u0440\u043C\u043E\u0434.`)
	unescCyr  = []byte(`Лорем ипсум долор сит амет, те воцибус нусяуам тибияуе сеа, цум цу еиус еирмод.`)
)

func testUnescape(t testing.TB, key string, src, dst []byte) {
	buf = append(buf[:0], src...)
	buf = unescape(buf)
	if !bytes.Equal(buf, dst) {
		t.Error(key, "assertion failed")
	}
}

func TestUnescape0(t *testing.T) {
	testUnescape(t, "unescape 0", esc0, unesc0)
}

func TestUnescape1(t *testing.T) {
	testUnescape(t, "unescape 1", esc1, unesc1)
}

func TestUnescape2(t *testing.T) {
	testUnescape(t, "unescape 2", esc2, unesc2)
}

func BenchmarkUnescape0(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testUnescape(b, "unescape 0", esc0, unesc0)
	}
}

func BenchmarkUnescape1(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testUnescape(b, "unescape 1", esc0, unesc0)
	}
}

func BenchmarkUnescape2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		testUnescape(b, "unescape 2", esc0, unesc0)
	}
}
