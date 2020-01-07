package cron

import (
	"gopkg.in/robfig/cron.v2"
)

var (
	_cron    = cron.New()
	specJobs = make(map[string]cron.EntryID)
)

// SpecJob cron任务定义
type SpecJob interface {
	Id() string
	Spec() string
	Job()
}

// SpecJobErr 自定义任务添加错误
type SpecJobErr struct {
	SpecJob
	Err error
}

// AddSpecJob 添加SpecJob任务
// 如果forced 为true，即使已经有相同id的任务在调度执行，也会强制移除再添加给定的任务
// 如果forced 为false，相同id的任务已经被调度将不会被添加
func AddSpecJob(job SpecJob, forced bool) error {
	id := job.Id()
	entryId, exist := specJobs[id]
	if exist {
		// 非强制添加，直接跳过
		if !forced {
			return nil
		}
		_cron.Remove(entryId)
	}
	entryId, err := _cron.AddFunc(job.Spec(), job.Job)
	if err != nil {
		return err
	}
	specJobs[id] = entryId
	return nil
}

// AddSpecJobs 批量添加SpecJob任务
// 如果forced 为true，即使已经有相同id的任务在调度执行，也会强制移除再添加给定的任务
// 如果forced 为false，相同id的任务已经被调度将不会被添加
func AddSpecJobs(jobs []SpecJob, forced bool) (res []*SpecJobErr) {
	for _, sj := range jobs {
		if err := AddSpecJob(sj, forced); err != nil {
			res = append(res, &SpecJobErr{
				SpecJob: sj,
				Err:     err,
			})
			continue
		}
	}
	return
}

// RemoveSpecJob 移除定时任务
func RemoveSpecJob(job SpecJob) {
	id := job.Id()
	RemoveJobById(id)
}

// RemoveJobById 通过id移除定时任务
func RemoveJobById(id string) {
	if entryId, exist := specJobs[id]; exist {
		_cron.Remove(entryId)
		delete(specJobs, id)
	}
}

// StartCron 开启定时任务
func StartCron() {
	_cron.Start()
}

// RestartCron 重新启动定时任务
func RestartCron() {
	_cron.Stop()
	_cron.Start()
}

// StopCron 关闭定时任务
func StopCron() {
	for id, entryId := range specJobs {
		_cron.Remove(entryId)
		delete(specJobs, id)
	}
	_cron.Stop()
}
