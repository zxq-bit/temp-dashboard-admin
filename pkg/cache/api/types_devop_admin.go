package api

import "time"

// Workspace represents the isolated space for your work.
type Workspace struct {
	ID          string `bson:"_id" json:"-"`
	Name        string `bson:"name" json:"name"`
	Alias       string `bson:"alias" json:"alias"`
	Description string `bson:"-" json:"description"`
	Owner       string `bson:"owner" json:"owner"`
	// SCM *cyclonev1.SCMConfig `bson:"-" json:"scm"`
	Tenant            string          `bson:"tenant" json:"-"`
	CycloneProject    string          `bson:"cycloneProject" json:"-"`
	PipelineCount     int             `bson:"-" json:"pipelineCount"`
	Cargo             *Cargo          `bson:"cargo" json:"cargo"`
	CacheDependencies map[string]bool `bson:"cacheDependencies" json:"cacheDependencies"`
	// WorkerQuota *cyclonev1.WorkerQuota `bson:"-" json:"workerQuota"`
	CreationTime   string `bson:"creationTime" json:"creationTime"`
	LastUpdateTime string `bson:"lastUpdateTime" json:"lastUpdateTime"`
}

// Cargo represents the config of Cargo.
type Cargo struct {
	Name    string `bson:"name" json:"name"`
	Project string `bson:"project" json:"project"`
	Host    string `bson:"host" json:"host"`
}

// ListMeta represents metadata that list resources must have.
type DaListMeta struct {
	Total int `json:"total"`
}

type WorkspaceList struct {
	Metadata DaListMeta  `json:"metadata"`
	Items    []Workspace `json:"items"`
}

// ErrorResponse represents response of error.
type DaErrorResponse struct {
	Message string `json:"message,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Details string `json:"details,omitempty"`
}

// come from cyclone

type Pipeline struct {
	ID          string `bson:"_id,omitempty" json:"id,omitempty" description:"id of the pipeline"`
	Name        string `bson:"name,omitempty" json:"name,omitempty" description:"name of the pipeline，unique in one project"`
	Alias       string `bson:"alias,omitempty" json:"alias,omitempty" description:"alias of the pipeline"`
	Description string `bson:"description,omitempty" json:"description,omitempty" description:"description of the pipeline"`
	Owner       string `bson:"owner,omitempty" json:"owner,omitempty" description:"owner of the pipeline"`
	ProjectID   string `bson:"projectID,omitempty" json:"projectID,omitempty" description:"id of the project which the pipeline belongs to"`
	// Build *Build `bson:"build,omitempty" json:"build,omitempty" description:"build spec of the pipeline"`
	// AutoTrigger *AutoTrigger `bson:"autoTrigger,omitempty" json:"autoTrigger,omitempty" description:"auto trigger strategy of the pipeline"`
	CreationTime         time.Time        `bson:"creationTime,omitempty" json:"creationTime,omitempty" description:"creation time of the pipeline"`
	LastUpdateTime       time.Time        `bson:"lastUpdateTime,omitempty" json:"lastUpdateTime,omitempty" description:"last update time of the pipeline"`
	RecentRecords        []PipelineRecord `bson:"-" json:"recentRecords,omitempty" description:"recent records of the pipeline"`
	RecentSuccessRecords []PipelineRecord `bson:"-" json:"recentSuccessRecords,omitempty" description:"recent success records of the pipeline"`
	RecentFailedRecords  []PipelineRecord `bson:"-" json:"recentFailedRecords,omitempty" description:"recent failed records of the pipeline"`
}

type PipelineRecord struct {
	ID            string                 `bson:"_id,omitempty" json:"id,omitempty" description:"id of the pipeline record"`
	Name          string                 `bson:"name,omitempty" json:"name,omitempty" description:"name of the pipeline record"`
	PipelineID    string                 `bson:"pipelineID,omitempty" json:"pipelineID,omitempty" description:"id of the related pipeline which the pipeline record belongs to"`
	Trigger       string                 `bson:"trigger,omitempty" json:"trigger,omitempty" description:"trigger of the pipeline record"`
	PerformParams *PipelinePerformParams `bson:"performParams,omitempty" json:"performParams,omitempty" description:"perform params of the pipeline record"`
	// StageStatus *StageStatus `bson:"stageStatus,omitempty" json:"stageStatus,omitempty" description:"status of each pipeline stage"`
	Status       string    `bson:"status,omitempty" json:"status,omitempty" description:"status of the pipeline record"`
	ErrorMessage string    `bson:"errorMessage,omitempty" json:"errorMessage,omitempty" description:"error message for the pipeline failure"`
	StartTime    time.Time `bson:"startTime,omitempty" json:"startTime,omitempty" description:"start time of the pipeline record"`
	EndTime      time.Time `bson:"endTime,omitempty" json:"endTime,omitempty" description:"end time of the pipeline record"`
}

// PipelinePerformParams the params to perform the pipeline.
type PipelinePerformParams struct {
	Ref             string   `bson:"ref,omitempty" json:"ref,omitempty" description:"reference of git repo, support branch, tag"`
	Name            string   `bson:"name,omitempty" json:"name,omitempty" description:"name of this running of pipeline"`
	Description     string   `bson:"description,omitempty" json:"description,omitempty" description:"description of this running of pipeline"`
	CreateSCMTag    bool     `bson:"createScmTag,omitempty" json:"createScmTag,omitempty" description:"whether create tag in SCM"`
	CacheDependency bool     `bson:"cacheDependency,omitempty" json:"cacheDependency,omitempty" description:"whether use dependency cache to speedup"`
	Stages          []string `bson:"stages,omitempty" json:"stages,omitempty" description:"stages to be executed"`
}

type PipelineList struct {
	Metadata DaListMeta `json:"metadata"`
	Items    []Pipeline `json:"items"`
}
