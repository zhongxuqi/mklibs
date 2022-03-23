package mkerr

import "testing"

func TestError(t *testing.T) {
	err := NewError(1, "2", "3")
	if err.ErrNo() != 1 || err.Error() != "2" {
		t.Fatalf("err data error %+v", err)
	}
	var _ error = err
	t.Fatalf("%+v", err.ErrDetail())
}
