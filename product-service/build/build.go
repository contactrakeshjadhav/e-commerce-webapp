package build

import (
	_ "embed"
	"encoding/json"
	"time"

	commons "github.com/contactrakeshjadhav/e-commerce-webapp/product-service/commons/pkg/core/model"
	"github.com/contactrakeshjadhav/e-commerce-webapp/product-service/core/model"
)

//go:embed build.json
var buildInfo []byte

// LoadBuildInformation will marshall the buils json info into our BuildInfo model
// returns an error on fail
func LoadBuildInformation() (model.BuildInfo, error) {
	var info model.BuildInfo
	if err := json.Unmarshal(buildInfo, &info); err != nil {
		return model.BuildInfo{}, err
	}
	if info.BuildDate == "" {
		info.BuildDate = time.Now().Format(commons.StandartISOFormat)
	}
	return info, nil
}
