package model

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type CTXKey string

const (
	AuthHeader      = "Authorization"
	ReqIDKey        = CTXKey("request-id")
	TokenKey        = CTXKey("aspire_token")
	UserKey         = CTXKey("aspire_user")
	ResponseCodeKey = CTXKey("response-code")

	ComponentLabel       = "component"
	ComponentDataspace   = "dataspace"
	ComponentWorkspace   = "workspace"
	ComponentBatchJob    = "batch-job"
	ComponentInternalJob = "internal-job"

	JobIDLabel   = "job-id"
	JobTypeLabel = "job-type"

	ServiceLabel        = "service"
	UserIDLabel         = "user"
	DeploymentNameLabel = "name"
	MilestoneService    = "milestone"
	ComputeService      = "compute"

	FQDNAnnotationKey          = "fqdn"
	BatchFilePathAnnotationKey = "batch-file-path"
	StandartISOFormat          = "2006-01-02T15:04:05-0700"
)

var (
	ErrInternalServerFail error = errors.New("internal server error")
	// string verification regex
	alphanumericRegEx = regexp.MustCompile("[^a-zA-Z0-9_]+")
)

func (key CTXKey) String() string {
	return string(key)
}

func GetReqIDFromRequest(r *http.Request) (string, error) {
	reqId := r.Header.Get(ReqIDKey.String())
	if reqId == "" {
		return "", errors.New("failed to get request ID from the incoming request")
	}
	return reqId, nil
}

func GetReqIDFromCTX(ctx context.Context) (string, error) {
	reqId := ctx.Value(ReqIDKey)
	if reqId == nil {
		return "", errors.New("failed to get request ID from the context")
	}
	return reqId.(string), nil
}

func GetCompositeKey(key1, key2 string) string {
	return fmt.Sprintf("%s+%s", strings.ToLower(key1), strings.ToLower(key2))
}

func CleanString(input string) string {
	target := strings.ToLower(input)
	target = strings.ReplaceAll(target, ".", " ")
	target = strings.ReplaceAll(target, "-", " ")
	target = strings.Join(strings.Fields(target), " ") // we remove the consecutive blanks
	target = strings.ReplaceAll(target, " ", "_")
	target = alphanumericRegEx.ReplaceAllString(target, "")
	return target
}
