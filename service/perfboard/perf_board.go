package perfboard

import (
	ctx "context"
	"edp-admin-console/models/query"
	"edp-admin-console/repository/perfboard"
	"edp-admin-console/service/logger"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	coreV1Client "k8s.io/client-go/kubernetes/typed/core/v1"
	"strconv"
)

var log = logger.GetLogger()

type PerfBoard struct {
	PerfRepo perfboard.IPerfServer
	CoreClient coreV1Client.CoreV1Interface
}

func (s PerfBoard) GetPerfServers() ([]*query.PerfServer, error) {
	log.Info("start retrieving PERF servers from DB")
	ps, err := s.PerfRepo.GetPerfServers()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get PERF server from DB")
	}
	return ps, nil
}

func (s PerfBoard) GetPerfServerName(id int) (*query.PerfServer, error) {
	log.Info("start retrieving PERF server from DB", zap.Int("id", id))
	ps, err := s.PerfRepo.GetPerfServerName(id)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get PERF server from DB")
	}
	return ps, nil
}

func (s PerfBoard) GetCodebaseDataSources(codebaseId int) ([]string, error) {
	log.Info("start retrieving PERF data sources", zap.Int("codebase id", codebaseId))
	ds, err := s.PerfRepo.GetCodebaseDataSources(codebaseId)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get PERF data sources from DB")
	}
	return ds, nil
}


func (s PerfBoard) IsPerfEnabled(namespace string) (bool, error) {
	cm,err:=s.CoreClient.ConfigMaps(namespace).Get(ctx.TODO(),"edp-config", metav1.GetOptions{} )
	if err != nil {
		return false, err
	}
	r := cm.Data["perf_integration_enabled"]
	if len(r) == 0 {
		return false, fmt.Errorf("there is no key perf_integration_enabled in cm edp-config" )
	}
	pe, err := strconv.ParseBool(r)
	if err != nil {
		return false, err
	}
	return pe, nil
}