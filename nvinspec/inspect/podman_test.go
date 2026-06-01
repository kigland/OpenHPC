package inspect

import "testing"

func TestParsePodmanPSJSON(t *testing.T) {
	data := []byte(`[
		{
			"Id": "abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789",
			"Names": ["train"],
			"Image": "localhost/train:latest"
		},
		{
			"ID": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			"Names": "infer",
			"ImageName": "localhost/infer:latest"
		}
	]`)

	got, err := ParsePodmanPSJSON(data)
	if err != nil {
		t.Fatalf("ParsePodmanPSJSON() error = %v", err)
	}
	if len(got.containers) != 2 {
		t.Fatalf("len(ParsePodmanPSJSON()) = %d, want 2", len(got.containers))
	}

	info, ok := got.find("abcdef012345")
	if !ok {
		t.Fatal("find short id failed")
	}
	if info.Name != "train" || info.Image != "localhost/train:latest" {
		t.Fatalf("first container = %+v", info)
	}
}

func TestParsePodmanTopHPIDs(t *testing.T) {
	data := []byte("HPID PID COMMAND\n1234 1 python train.py\n5678 2 /bin/sh\n")

	got := ParsePodmanTopHPIDs(data)
	if len(got) != 2 || got[0] != 1234 || got[1] != 5678 {
		t.Fatalf("ParsePodmanTopHPIDs() = %+v", got)
	}
}
