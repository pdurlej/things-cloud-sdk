package main

import (
	"testing"
	"time"

	thingscloud "github.com/arthursoares/things-cloud-sdk"
	memory "github.com/arthursoares/things-cloud-sdk/state/memory"
)

func testStateForListFilters() *memory.State {
	state := memory.NewState()
	today := time.Now().UTC()
	tomorrow := today.Add(24 * time.Hour)

	state.Areas["area-1"] = &thingscloud.Area{UUID: "area-1", Title: "Work"}
	state.Tasks["project-1"] = &thingscloud.Task{
		UUID:  "project-1",
		Title: "Project Alpha",
		Type:  thingscloud.TaskTypeProject,
	}
	state.Tasks["inbox-1"] = &thingscloud.Task{
		UUID:     "inbox-1",
		Title:    "Inbox Task",
		Schedule: thingscloud.TaskScheduleInbox,
	}
	state.Tasks["today-1"] = &thingscloud.Task{
		UUID:          "today-1",
		Title:         "Today Task",
		Schedule:      thingscloud.TaskScheduleAnytime,
		ScheduledDate: &today,
	}
	state.Tasks["anytime-1"] = &thingscloud.Task{
		UUID:     "anytime-1",
		Title:    "Anytime Task",
		Schedule: thingscloud.TaskScheduleAnytime,
	}
	state.Tasks["someday-1"] = &thingscloud.Task{
		UUID:     "someday-1",
		Title:    "Someday Task",
		Schedule: thingscloud.TaskScheduleSomeday,
	}
	state.Tasks["upcoming-1"] = &thingscloud.Task{
		UUID:          "upcoming-1",
		Title:         "Upcoming Task",
		Note:          "needle in note",
		Schedule:      thingscloud.TaskScheduleSomeday,
		ScheduledDate: &tomorrow,
	}
	state.Tasks["project-task-1"] = &thingscloud.Task{
		UUID:          "project-task-1",
		Title:         "Project Task",
		Schedule:      thingscloud.TaskScheduleAnytime,
		ParentTaskIDs: []string{"project-1"},
	}
	state.Tasks["area-task-1"] = &thingscloud.Task{
		UUID:     "area-task-1",
		Title:    "Area Task",
		Schedule: thingscloud.TaskScheduleAnytime,
		AreaIDs:  []string{"area-1"},
	}
	state.Tasks["completed-1"] = &thingscloud.Task{
		UUID:     "completed-1",
		Title:    "Completed Task",
		Schedule: thingscloud.TaskScheduleAnytime,
		Status:   thingscloud.TaskStatusCompleted,
	}
	state.Tasks["trashed-1"] = &thingscloud.Task{
		UUID:     "trashed-1",
		Title:    "Trashed Task",
		Schedule: thingscloud.TaskScheduleAnytime,
		InTrash:  true,
	}
	return state
}

func outputUUIDs(tasks []TaskOutput) []string {
	uuids := make([]string, len(tasks))
	for i, task := range tasks {
		uuids[i] = task.UUID
	}
	return uuids
}

func requireUUIDs(t *testing.T, tasks []TaskOutput, want ...string) {
	t.Helper()
	got := outputUUIDs(tasks)
	if len(got) != len(want) {
		t.Fatalf("uuids = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("uuids = %v, want %v", got, want)
		}
	}
}

func TestListTasksLocationFilters(t *testing.T) {
	state := testStateForListFilters()

	requireUUIDs(t, listTasks(state, map[string]string{"today": "true"}), "today-1")
	requireUUIDs(t, listTasks(state, map[string]string{"inbox": "true"}), "inbox-1")
	requireUUIDs(t, listTasks(state, map[string]string{"someday": "true"}), "someday-1")
	requireUUIDs(t, listTasks(state, map[string]string{"upcoming": "true"}), "upcoming-1")

	anytime := outputUUIDs(listTasks(state, map[string]string{"anytime": "true"}))
	wantAnytime := []string{"anytime-1", "area-task-1", "project-task-1"}
	if len(anytime) != len(wantAnytime) {
		t.Fatalf("anytime = %v, want %v", anytime, wantAnytime)
	}
	for i := range wantAnytime {
		if anytime[i] != wantAnytime[i] {
			t.Fatalf("anytime = %v, want %v", anytime, wantAnytime)
		}
	}
}

func TestListTasksSearchAndContainerFilters(t *testing.T) {
	state := testStateForListFilters()

	requireUUIDs(t, listTasks(state, map[string]string{"search": "needle"}), "upcoming-1")
	requireUUIDs(t, listTasks(state, map[string]string{"search": "project"}), "project-task-1")
	requireUUIDs(t, listTasks(state, map[string]string{"area": "Work"}), "area-task-1")
	requireUUIDs(t, listTasks(state, map[string]string{"project": "Project Alpha"}), "project-task-1")
}
