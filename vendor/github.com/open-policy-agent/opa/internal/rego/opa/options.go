package opa

import (
	"io"
	"time"

	"github.com/yerinu2019/opa/metrics"
	"github.com/yerinu2019/opa/topdown/cache"
)

// Result holds the evaluation result.
type Result struct {
	Result []byte
}

// EvalOpts define options for performing an evaluation.
type EvalOpts struct {
	Input                  *interface{}
	Metrics                metrics.Metrics
	Entrypoint             int32
	Time                   time.Time
	Seed                   io.Reader
	InterQueryBuiltinCache cache.InterQueryCache
}
