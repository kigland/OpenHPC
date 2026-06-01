package inspect

import "testing"

func TestParseNvidiaSmiCSV(t *testing.T) {
	data := []byte("GPU-aaaaaaaa-bbbb, 1234, /usr/bin/python, 2048\nGPU-cccc, 5678, \"python,worker\", [N/A]\n")

	got, err := ParseNvidiaSmiCSV(data)
	if err != nil {
		t.Fatalf("ParseNvidiaSmiCSV() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(ParseNvidiaSmiCSV()) = %d, want 2", len(got))
	}
	if got[0].GPUUUID != "GPU-aaaaaaaa-bbbb" || got[0].PID != 1234 || got[0].ProcessName != "/usr/bin/python" || got[0].UsedMemoryMiB != 2048 {
		t.Fatalf("first process = %+v", got[0])
	}
	if got[1].ProcessName != "python,worker" || got[1].UsedMemoryMiB != -1 {
		t.Fatalf("second process = %+v", got[1])
	}
}

func TestParseNvidiaSmiCSVEmpty(t *testing.T) {
	got, err := ParseNvidiaSmiCSV([]byte("\n"))
	if err != nil {
		t.Fatalf("ParseNvidiaSmiCSV() error = %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("len(ParseNvidiaSmiCSV()) = %d, want 0", len(got))
	}
}
