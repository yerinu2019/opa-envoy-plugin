package decisionlog

import (
	"context"

	"github.com/pkg/errors"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/plugins"
	"github.com/open-policy-agent/opa/plugins/logs"
	"github.com/open-policy-agent/opa/server"
	"github.com/open-policy-agent/opa/storage"
	"github.com/open-policy-agent/opa/topdown"
	"github.com/open-policy-agent/opa/util"
	"github.com/yerinu2019/opa-envoy-plugin/envoyauth"
)

type internalError struct {
	Message string `json:"message"`
}

func (e *internalError) Error() string {
	return e.Message
}

func getAstValue(str string) (ast.Value, error) {
	var input interface{}
	if err := util.UnmarshalJSON([]byte(str), &input); err != nil {
		return nil, errors.Wrapf(err, "parameter contains malformed input document")
	}
	return ast.InterfaceToValue(input)
}

// LogDecision - Logs a decision log event
func LogDecision(ctx context.Context, manager *plugins.Manager, info *server.Info, result *envoyauth.EvalResult, err error) error {
	plugin := logs.Lookup(manager)
	if plugin == nil {
		return nil
	}

	//mask := "[{\"/input/attributes/request/http/body\": [true]}, {\"/input/parsed_body\": [true]}]"
	//info.InputAST, err = getAstValue(mask)

	info.Revision = result.Revision

	bundles := map[string]server.BundleInfo{}
	for name, rev := range result.Revisions {
		bundles[name] = server.BundleInfo{Revision: rev}
	}
	info.Bundles = bundles

	info.DecisionID = result.DecisionID
	info.Metrics = result.Metrics
	info.Txn = result.Txn

	if err != nil {
		switch err.(type) {
		case *storage.Error, *ast.Error, ast.Errors:
			break
		case *topdown.Error:
			if topdown.IsCancel(err) {
				err = &topdown.Error{
					Code:    topdown.CancelErr,
					Message: "context deadline reached during query execution",
				}
			}
		default:
			// Wrap errors that may not serialize to JSON well (e.g., fmt.Errorf, etc.)
			err = &internalError{Message: err.Error()}
		}
		info.Error = err
	} else {
		var x interface{}
		if result != nil {
			x = result.Decision
		}
		info.Results = &x
	}

	return plugin.Log(ctx, info)
}
