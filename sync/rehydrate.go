package sync

import (
	"time"

	things "github.com/pdurlej/things-cloud-sdk"
)

func (s *Syncer) rehydrateLoggedChange(base baseChange, changeType, entityType, entityUUID, rawPayload string) (Change, error) {
	logged := LoggedChange{
		baseChange: base,
		changeType: changeType,
		entityType: entityType,
		entityUUID: entityUUID,
		payload:    rawPayload,
	}

	switch entityType {
	case "Task", "Project", "Heading":
		return s.rehydrateTaskChange(logged, base, changeType, entityUUID)
	case "Area":
		return s.rehydrateAreaChange(logged, base, changeType, entityUUID)
	case "Tag":
		return s.rehydrateTagChange(logged, base, changeType, entityUUID)
	case "ChecklistItem":
		return s.rehydrateChecklistChange(logged, base, changeType, entityUUID)
	default:
		if changeType == "UnknownChange" {
			return UnknownChange{baseChange: base, entityType: entityType, entityUUID: entityUUID, Details: rawPayload}, nil
		}
		return logged, nil
	}
}

func (s *Syncer) rehydrateTaskChange(fallback LoggedChange, base baseChange, changeType, entityUUID string) (Change, error) {
	task, err := s.getTask(entityUUID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return fallback, nil
	}

	tc := taskChange{baseChange: base, Task: task}
	switch changeType {
	case "TaskCreated":
		return TaskCreated{taskChange: tc}, nil
	case "TaskDeleted":
		return TaskDeleted{taskChange: tc}, nil
	case "TaskCompleted":
		return TaskCompleted{taskChange: tc}, nil
	case "TaskUncompleted":
		return TaskUncompleted{taskChange: tc}, nil
	case "TaskCanceled":
		return TaskCanceled{taskChange: tc}, nil
	case "TaskTitleChanged":
		return TaskTitleChanged{taskChange: tc}, nil
	case "TaskNoteChanged":
		return TaskNoteChanged{taskChange: tc}, nil
	case "TaskMovedToInbox":
		return TaskMovedToInbox{taskChange: tc, From: LocationUnknown}, nil
	case "TaskMovedToToday":
		return TaskMovedToToday{taskChange: tc, From: LocationUnknown}, nil
	case "TaskMovedToAnytime":
		return TaskMovedToAnytime{taskChange: tc, From: LocationUnknown}, nil
	case "TaskMovedToSomeday":
		return TaskMovedToSomeday{taskChange: tc, From: LocationUnknown}, nil
	case "TaskMovedToUpcoming":
		var scheduledFor time.Time
		if task.ScheduledDate != nil {
			scheduledFor = *task.ScheduledDate
		}
		return TaskMovedToUpcoming{taskChange: tc, From: LocationUnknown, ScheduledFor: scheduledFor}, nil
	case "TaskDeadlineChanged":
		return TaskDeadlineChanged{taskChange: tc}, nil
	case "TaskAssignedToProject":
		project, err := s.firstTaskRef(task.ParentTaskIDs)
		if err != nil {
			return nil, err
		}
		return TaskAssignedToProject{taskChange: tc, Project: project}, nil
	case "TaskAssignedToArea":
		area, err := s.firstAreaRef(task.AreaIDs)
		if err != nil {
			return nil, err
		}
		return TaskAssignedToArea{taskChange: tc, Area: area}, nil
	case "TaskTrashed":
		return TaskTrashed{taskChange: tc}, nil
	case "TaskRestored":
		return TaskRestored{taskChange: tc}, nil
	case "TaskTagsChanged":
		return TaskTagsChanged{taskChange: tc}, nil
	case "ProjectCreated":
		return ProjectCreated{projectChange: projectChange{baseChange: base, Project: task}}, nil
	case "ProjectDeleted":
		return ProjectDeleted{projectChange: projectChange{baseChange: base, Project: task}}, nil
	case "ProjectCompleted":
		return ProjectCompleted{projectChange: projectChange{baseChange: base, Project: task}}, nil
	case "ProjectTitleChanged":
		return ProjectTitleChanged{projectChange: projectChange{baseChange: base, Project: task}}, nil
	case "ProjectTrashed":
		return ProjectTrashed{projectChange: projectChange{baseChange: base, Project: task}}, nil
	case "ProjectRestored":
		return ProjectRestored{projectChange: projectChange{baseChange: base, Project: task}}, nil
	case "HeadingCreated":
		return HeadingCreated{headingChange: headingChange{baseChange: base, Heading: task}}, nil
	case "HeadingDeleted":
		return HeadingDeleted{headingChange: headingChange{baseChange: base, Heading: task}}, nil
	case "HeadingTitleChanged":
		return HeadingTitleChanged{headingChange: headingChange{baseChange: base, Heading: task}}, nil
	default:
		return fallback, nil
	}
}

func (s *Syncer) rehydrateAreaChange(fallback LoggedChange, base baseChange, changeType, entityUUID string) (Change, error) {
	area, err := s.getArea(entityUUID)
	if err != nil {
		return nil, err
	}
	if area == nil {
		return fallback, nil
	}
	ac := areaChange{baseChange: base, Area: area}
	switch changeType {
	case "AreaCreated":
		return AreaCreated{areaChange: ac}, nil
	case "AreaDeleted":
		return AreaDeleted{areaChange: ac}, nil
	case "AreaRenamed":
		return AreaRenamed{areaChange: ac}, nil
	default:
		return fallback, nil
	}
}

func (s *Syncer) rehydrateTagChange(fallback LoggedChange, base baseChange, changeType, entityUUID string) (Change, error) {
	tag, err := s.getTag(entityUUID)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return fallback, nil
	}
	tc := tagChange{baseChange: base, Tag: tag}
	switch changeType {
	case "TagCreated":
		return TagCreated{tagChange: tc}, nil
	case "TagDeleted":
		return TagDeleted{tagChange: tc}, nil
	case "TagRenamed":
		return TagRenamed{tagChange: tc}, nil
	case "TagShortcutChanged":
		return TagShortcutChanged{tagChange: tc}, nil
	default:
		return fallback, nil
	}
}

func (s *Syncer) rehydrateChecklistChange(fallback LoggedChange, base baseChange, changeType, entityUUID string) (Change, error) {
	item, err := s.getChecklistItem(entityUUID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return fallback, nil
	}
	task, err := s.firstTaskRef(item.TaskIDs)
	if err != nil {
		return nil, err
	}
	cc := checklistItemChange{baseChange: base, Item: item}
	switch changeType {
	case "ChecklistItemCreated":
		return ChecklistItemCreated{checklistItemChange: cc, Task: task}, nil
	case "ChecklistItemDeleted":
		return ChecklistItemDeleted{checklistItemChange: cc}, nil
	case "ChecklistItemCompleted":
		return ChecklistItemCompleted{checklistItemChange: cc, Task: task}, nil
	case "ChecklistItemUncompleted":
		return ChecklistItemUncompleted{checklistItemChange: cc, Task: task}, nil
	case "ChecklistItemTitleChanged":
		return ChecklistItemTitleChanged{checklistItemChange: cc}, nil
	default:
		return fallback, nil
	}
}

func (s *Syncer) firstTaskRef(uuids []string) (*things.Task, error) {
	if len(uuids) == 0 || uuids[0] == "" {
		return nil, nil
	}
	return s.getTask(uuids[0])
}

func (s *Syncer) firstAreaRef(uuids []string) (*things.Area, error) {
	if len(uuids) == 0 || uuids[0] == "" {
		return nil, nil
	}
	return s.getArea(uuids[0])
}
