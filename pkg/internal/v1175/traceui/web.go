// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package traceui

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"

	_ "net/http/pprof" // Required to use pprof

	"github.com/xyctruth/profiler/pkg/internal/v1175/trace"
)

type TraceUI struct {
	data   []byte
	loader struct {
		once sync.Once
		res  trace.ParseResult
		err  error
	}
	ranges   []Range
	Handlers map[string]http.HandlerFunc
	gsInit   sync.Once
	gs       map[uint64]*trace.GDesc
	mmuCache struct {
		m    map[trace.UtilFlags]*mmuCacheEntry
		lock sync.Mutex
	}
}

func NewTraceUI(data []byte) *TraceUI {
	traceUI := &TraceUI{
		data: data,
	}
	traceUI.mmuCache.m = make(map[trace.UtilFlags]*mmuCacheEntry)

	res, err := traceUI.parseTrace()
	if err != nil {
		dief("%v\n", err)
	}
	traceUI.ranges = traceUI.splitTrace(res)
	handlers := make(map[string]http.HandlerFunc)
	handlers["/"] = traceUI.httpMain
	handlers["/mmu"] = httpMMU
	handlers["/mmuPlot"] = traceUI.httpMMUPlot
	handlers["/mmuDetails"] = traceUI.httpMMUDetails
	handlers["/usertasks"] = traceUI.httpUserTasks
	handlers["/usertask"] = traceUI.httpUserTask
	handlers["/userregions"] = traceUI.httpUserRegions
	handlers["/userregion"] = traceUI.httpUserRegion
	handlers["/trace"] = traceUI.httpTrace
	handlers["/jsontrace"] = traceUI.httpJsonTrace
	handlers["/jsontrace"] = traceUI.httpJsonTrace
	handlers["/trace_viewer_html"] = httpTraceViewerHTML
	handlers["/webcomponents.min.js"] = webcomponentsJS
	handlers["/io"] = serveSVGProfile(traceUI.pprofByGoroutine(computePprofIO))
	handlers["/block"] = serveSVGProfile(traceUI.pprofByGoroutine(computePprofBlock))
	handlers["/syscall"] = serveSVGProfile(traceUI.pprofByGoroutine(computePprofSyscall))
	handlers["/sched"] = serveSVGProfile(traceUI.pprofByGoroutine(computePprofSched))
	handlers["/regionio"] = serveSVGProfile(traceUI.pprofByRegion(computePprofIO))
	handlers["/regionblock"] = serveSVGProfile(traceUI.pprofByRegion(computePprofBlock))
	handlers["/regionsyscall"] = serveSVGProfile(traceUI.pprofByRegion(computePprofSyscall))
	handlers["/regionsched"] = serveSVGProfile(traceUI.pprofByRegion(computePprofSched))
	handlers["/goroutines"] = traceUI.httpGoroutines
	handlers["/goroutine"] = traceUI.httpGoroutine

	traceUI.Handlers = handlers
	return traceUI
}

// parseEvents is a compatibility wrapper that returns only
// the Events part of trace.ParseResult returned by parseTrace.
func (traceUI *TraceUI) parseEvents() ([]*trace.Event, error) {
	res, err := traceUI.parseTrace()
	if err != nil {
		return nil, err
	}
	return res.Events, err
}

func (traceUI *TraceUI) parseTrace() (trace.ParseResult, error) {
	traceUI.loader.once.Do(func() {
		buf := bytes.NewBuffer(traceUI.data)
		// Parse and symbolize.
		res, err := trace.Parse(bufio.NewReader(buf), "")
		if err != nil {
			traceUI.loader.err = fmt.Errorf("failed to parse trace: %v", err)
			return
		}
		traceUI.loader.res = res
	})
	return traceUI.loader.res, traceUI.loader.err
}

// httpMain serves the starting page.
func (traceUI *TraceUI) httpMain(w http.ResponseWriter, r *http.Request) {
	if err := templMain.Execute(w, traceUI.ranges); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var templMain = template.Must(template.New("").Parse(`
<html>
<body>
{{if $}}
	{{range $e := $}}
		<a href="{{$e.URL}}">View trace ({{$e.Name}})</a><br>
	{{end}}
	<br>
{{else}}
	<a href="trace">View trace</a><br>
{{end}}
<a href="goroutines">Goroutine analysis</a><br>
<a href="usertasks">User-defined tasks</a><br>
<a href="userregions">User-defined regions</a><br>
<a href="mmu">Minimum mutator utilization</a><br>

<!-- <a href="io">Network blocking profile</a> -->
(<a href="io?raw=1" download="io.profile">Network blocking profile⬇</a>)<br> 
<!-- <a href="block">Synchronization blocking profile</a> -->
(<a href="block?raw=1" download="block.profile">Synchronization blocking profile⬇</a>)<br> 
<!-- <a href="syscall">Syscall blocking profile</a> --> 
(<a href="syscall?raw=1" download="syscall.profile">Syscall blocking profile⬇</a>)<br> 
<!-- <a href="sched">Scheduler latency profile</a> --> 
(<a href="sche?raw=1" download="sched.profile">Scheduler latency profile⬇</a>)<br>

</body>
</html>
`))

func dief(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
