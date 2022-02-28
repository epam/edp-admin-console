package pipelinestage

import (
	"context"
	"fmt"
	"strconv"

	"edp-admin-console/k8s"
)

type ApplicationStage struct {
	Name       string `json:"name"`
	BranchName string `json:"branchName"`
	InputIs    string `json:"inputIs"`
	OutputIs   string `json:"outputIs"`
}

func BuildApplicationStages(ctx context.Context, namespacedClient *k8s.RuntimeNamespacedClient, inputIS, outputIS []string, applicationNames []string) ([]ApplicationStage, error) {
	if len(inputIS) != len(outputIS) || len(inputIS) != len(applicationNames) || len(outputIS) != len(applicationNames) {
		return nil, fmt.Errorf("inputIS, outputIS, applicationNames not the same size. InputIS size %v. OutputIS size %v. applicationNames size %v", len(inputIS), len(outputIS), len(applicationNames))
	}
	var applications []ApplicationStage
	for index := range inputIS {
		branchName := getBranchNameByISName(ctx, namespacedClient, inputIS[index])
		tmpApplication := ApplicationStage{
			Name:       applicationNames[index],
			InputIs:    inputIS[index],
			OutputIs:   outputIS[index],
			BranchName: branchName,
		}
		applications = append(applications, tmpApplication)
	}
	return applications, nil
}

type StageMainData struct {
	Description     string // can be found in Stage.Spec.Description by Stage Name
	TriggerType     string // can be found in Stage.Spec.TriggerType by Stage Name
	Order           string // can be found in Stage.Spec.Order by Stage Name
	JobProvisioning string // can be found in Stage.Spec.JobProvisioning by Stage Name
	QualityGates    []QualityGate
}

type QualityGate struct {
	Id               int             `json:"id"`
	QualityGateType  string          `json:"qualityGateType"`
	StepName         string          `json:"stepName"`
	CdStageId        *int            `json:"cdStageId"`
	CodebaseBranchId *int            `json:"branchId"`
	Autotest         *Codebase       `json:"autotest"`
	Branch           *CodebaseBranch `json:"codebaseBranch"`
}

type Codebase struct {
	Id                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	Language             string                 `json:"language"`
	BuildTool            string                 `json:"build_tool"`
	Framework            string                 `json:"framework"`
	Strategy             string                 `json:"strategy"`
	GitUrl               string                 `json:"git_url"`
	Type                 string                 `json:"type"`
	Status               string                 `json:"status"`
	TestReportFramework  string                 `json:"testReportFramework"`
	Description          string                 `json:"description"`
	GitServer            *string                `json:"gitServer"`
	GitProjectPath       *string                `json:"gitProjectPath"`
	JenkinsSlave         string                 `json:"jenkinsSlave"`
	JobProvisioning      string                 `json:"jobProvisioning"`
	DeploymentScript     string                 `json:"deploymentScript"`
	VersioningType       string                 `json:"versioningType"`
	StartVersioningFrom  *string                `json:"startFrom"`
	JiraServer           *string                `json:"jiraServer"`
	CommitMessagePattern string                 `json:"commitMessagePattern"`
	TicketNamePattern    string                 `json:"ticketNamePattern"`
	CiTool               string                 `json:"ciTool"`
	Perf                 *Perf                  `json:"perf"`
	DefaultBranch        string                 `json:"defaultBranch"`
	JiraIssueFields      map[string]interface{} `json:"jiraIssueFields"`
	EmptyProject         bool                   `json:"emptyProject"`
}

type CodebaseBranch struct {
	Id               int     `json:"id"`
	Name             string  `json:"branchName"`
	FromCommit       string  `json:"from_commit"`
	Status           string  `json:"status"`
	Version          *string `json:"version"`
	Build            *string `json:"build_number"`
	LastSuccessBuild *string `json:"last_success_build"`
	VCSLink          string  `json:"branchLink"`
	CICDLink         string  `json:"jenkinsLink"`
	AppName          string  `json:"appName"`
	Release          bool    `json:"release"`
}

type Perf struct {
	Name        string   `json:"name"`
	DataSources []string `json:"dataSources"`
}

func StageViewByCRName(ctx context.Context, namespacedClient *k8s.RuntimeNamespacedClient, stageCRName string) (*StageMainData, error) {
	var particleStageView StageMainData
	stage, err := namespacedClient.GetCDStage(ctx, stageCRName)
	if err != nil {
		return nil, err
	}
	particleStageView.Description = stage.Spec.Description
	particleStageView.TriggerType = stage.Spec.TriggerType
	particleStageView.Order = strconv.Itoa(stage.Spec.Order)
	particleStageView.JobProvisioning = stage.Spec.JobProvisioning

	var qualityGates []QualityGate
	for _, v := range stage.Spec.QualityGates {

		autotest := &Codebase{}
		if v.AutotestName != nil {
			autotest = getCodebaseByCRName(ctx, *v.AutotestName, namespacedClient)
		} else {
			autotest = nil
		}

		branch := &CodebaseBranch{
			Name: emptyIfNil(v.BranchName),
		}

		//database compatibility
		cdStageID := 0
		tmpQG := QualityGate{
			QualityGateType: v.QualityGateType,
			StepName:        v.StepName,
			Autotest:        autotest,
			Branch:          branch,
			CdStageId:       &cdStageID,
		}
		qualityGates = append(qualityGates, tmpQG)
	}

	particleStageView.QualityGates = qualityGates

	return &particleStageView, err
}

func getCodebaseByCRName(ctx context.Context, crName string, namespacedClient *k8s.RuntimeNamespacedClient) *Codebase {
	codebase, err := namespacedClient.GetCodebase(ctx, crName)
	if err != nil {
		return nil
	}
	var autoTest Codebase
	autoTest.Name = codebase.Name
	autoTest.Language = codebase.Spec.Lang
	autoTest.BuildTool = codebase.Spec.BuildTool
	autoTest.Framework = emptyIfNil(codebase.Spec.Framework)
	autoTest.Strategy = string(codebase.Spec.Strategy)
	autoTest.GitUrl = emptyIfNil(codebase.Spec.GitUrlPath)
	autoTest.Type = codebase.Spec.Type
	autoTest.Status = codebase.Status.Status
	autoTest.TestReportFramework = emptyIfNil(codebase.Spec.TestReportFramework)
	autoTest.Description = emptyIfNil(codebase.Spec.Description)
	autoTest.GitServer = &codebase.Spec.GitServer
	autoTest.JenkinsSlave = emptyIfNil(codebase.Spec.JenkinsSlave)
	autoTest.JobProvisioning = emptyIfNil(codebase.Spec.JobProvisioning)
	autoTest.DeploymentScript = codebase.Spec.DeploymentScript
	autoTest.VersioningType = string(codebase.Spec.Versioning.Type)
	autoTest.StartVersioningFrom = codebase.Spec.Versioning.StartFrom
	autoTest.JiraServer = codebase.Spec.JiraServer
	autoTest.CommitMessagePattern = emptyIfNil(codebase.Spec.CommitMessagePattern)
	autoTest.TicketNamePattern = emptyIfNil(codebase.Spec.TicketNamePattern)
	autoTest.CiTool = codebase.Spec.CiTool
	if codebase.Spec.Perf != nil {
		autoTest.Perf.Name = codebase.Spec.Perf.Name
		autoTest.Perf.DataSources = codebase.Spec.Perf.DataSources
	}
	autoTest.DefaultBranch = codebase.Spec.DefaultBranch
	autoTest.EmptyProject = codebase.Spec.EmptyProject
	return &autoTest
}

func emptyIfNil(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func getBranchNameByISName(ctx context.Context, namespacedClient *k8s.RuntimeNamespacedClient, crName string) string {
	stream, err := namespacedClient.GetCBBranch(ctx, crName)
	if err != nil {
		return ""
	}
	return stream.Spec.BranchName
}

func CdPipelineAppNamesByCRName(ctx context.Context, namespacedClient *k8s.RuntimeNamespacedClient, cdPipeCRName string) ([]string, error) {
	cdPipeline, err := namespacedClient.GetCDPipeline(ctx, cdPipeCRName)
	if err != nil {
		return nil, err
	}
	return cdPipeline.Spec.ApplicationsToPromote, nil
}
