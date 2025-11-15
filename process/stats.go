package process

import "time"

type ProcessStatus int

const (
	Running ProcessStatus = iota
	Completed
	Failed
)

func (s ProcessStatus) String() string {
	return [...]string{"running", "completed", "failed"}[s]
}

type GlobalInformation struct {
	Start, End     time.Time
	ProcessesStats map[string]*ProcessInformation
}

type ProcessInformation struct {
	Start, End  time.Time
	Status      ProcessStatus
	ProcessName string
	Err         error

	Passed bool
}

func (p *ProcessInformation) String() {
	
}

func newGlobalInformation() *GlobalInformation {
	return &GlobalInformation{
		Start:          time.Now(),
		ProcessesStats: make(map[string]*ProcessInformation),
	}
}

func (info *GlobalInformation) start(processName string) {
	if info == nil {
		return
	}
	info.ProcessesStats[processName] = &ProcessInformation{
		ProcessName: processName,
		Status:      Running,
		Start:       time.Now(),
	}
}

func (info *GlobalInformation) end(processName string, err error) {
	if info == nil {
		return
	}
	info.get(processName).End = time.Now()

	if err != nil {
		info.get(processName).Err = err
		info.get(processName).Status = Failed
	} else {
		info.get(processName).Status = Completed
	}
}

func (info *GlobalInformation) get(processName string) *ProcessInformation {
	return info.ProcessesStats[processName]
}

func (info *GlobalInformation) setError(processName string, err error) {
	if info == nil {
		return
	}
	info.ProcessesStats[processName].Err = err
}
