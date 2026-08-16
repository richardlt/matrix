package animate

import "testing"

func fill(a animation, distinct int) {
	// a real animation reuses a small palette across the frame
	for i := 0; i+2 < len(a); i += 3 {
		v := byte((i / 3) % distinct)
		a[i], a[i+1], a[i+2] = v*7, v*13, v*29
	}
}

func BenchmarkReadFrameRealistic(b *testing.B) {
	const w, h = 16, 9
	anim := make(animation, w*h*3*10)
	fill(anim, 8) // 8 distinct colours, typical of pixel art
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = anim.readFrame(w, h, i%10)
	}
}

func BenchmarkReadFrameWorstCase(b *testing.B) {
	const w, h = 16, 9
	anim := make(animation, w*h*3*10)
	fill(anim, 144) // every pixel a different colour
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = anim.readFrame(w, h, i%10)
	}
}
