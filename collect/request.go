package collect

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

type Task struct {
	Name        string
	Url         string
	Cookie      string
	WaitTime    time.Duration
	Reload      bool
	MaxDepth    int
	Visited     map[string]bool
	VisitedLock sync.Mutex
	Fetcher     Fetcher
	Rule        RuleTree
}

type Context struct {
	Body []byte
	Req  *Request
}

type Request struct {
	unique   string
	Task     *Task
	Url      string
	Method   string
	Depth    int
	Priority int
	RuleName string
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{}
}

func (r *Request) Check() error {
	if r.Depth > r.Task.MaxDepth {
		return errors.New("Max depth limit reached")
	}
	return nil
}

// 请求的唯一识别码
func (r *Request) Unique() string {
	block := md5.Sum([]byte(r.Url + r.Method))
	return hex.EncodeToString(block[:])
}
