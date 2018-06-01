package mklog

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zhongxuqi/mklibs/common"
)

func TestNewLog(t *testing.T) {
	ml := New()
	if mlInstance, _ := ml.(*logger); mlInstance.level != LevelDebug || mlInstance.logID == "" {
		t.Fatalf("data error %+v", mlInstance)
	}

	req := httptest.NewRequest(http.MethodGet, "http://web.com", bytes.NewBufferString(""))

	// init with empty req
	ml = NewWithReq(req)
	if mlInstance, _ := ml.(*logger); mlInstance.level != LevelDebug || mlInstance.logID == "" || req.Header.Get(common.HttpLogID) != mlInstance.logID {
		t.Fatalf("data error %+v", mlInstance)
	}

	// init with header req
	req.Header.Set(common.HttpLogID, "test-log")
	ml = NewWithReq(req)
	if mlInstance, _ := ml.(*logger); mlInstance.level != LevelDebug || mlInstance.logID != "test-log" || req.Header.Get(common.HttpLogID) != mlInstance.logID {
		t.Fatalf("data error %+v", mlInstance)
	}

	ml = NewWithContext(context.TODO())
	if mlInstance, _ := ml.(*logger); mlInstance.level != LevelDebug || mlInstance.logID == "" {
		t.Fatalf("data error %+v", mlInstance)
	}
	ml1 := NewWithContext(ml.Context())
	if mlInstance, _ := ml.(*logger); mlInstance.level != LevelDebug || mlInstance.logID == "" {
		t.Fatalf("data error %+v", mlInstance)
	} else {
		if mlInstance1, _ := ml1.(*logger); mlInstance1.level != LevelDebug || mlInstance1.logID != mlInstance.logID {
			t.Fatalf("data error %+v %+v", mlInstance, mlInstance1)
		}
	}
}

func TestPrint(t *testing.T) {
	ml := New()
	ml.Debugf("test")
	ml.Infof("test")
	ml.Errorf("test")
	dist := bytes.NewBuffer(nil)
	ml.SetOutput(dist)

	// test print
	dist.Reset()
	ml.Debugf("test")
	if dist.Len() <= 0 {
		t.Fatalf("Debugf print error")
	}
	dist.Reset()
	ml.Infof("test")
	if dist.Len() <= 0 {
		t.Fatalf("Infof print error")
	}
	dist.Reset()
	ml.Errorf("test")
	if dist.Len() <= 0 {
		t.Fatalf("Errorf print error")
	}

	// test level print
	ml.SetLevel(LevelInfo)
	if mlInstance, _ := ml.(*logger); mlInstance.level != LevelInfo || mlInstance.logID == "" {
		t.Fatalf("data error %+v", mlInstance)
	}
	dist.Reset()
	ml.Debugf("test")
	if dist.Len() > 0 {
		t.Fatalf("Debugf print error")
	}
	dist.Reset()
	ml.Infof("test")
	if dist.Len() <= 0 {
		t.Fatalf("Infof print error")
	}
	dist.Reset()
	ml.Errorf("test")
	if dist.Len() <= 0 {
		t.Fatalf("Errorf print error")
	}
}
