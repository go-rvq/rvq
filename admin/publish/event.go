package publish

import (
	"github.com/qor5/admin/v3/presets"
	"gorm.io/gorm"
)

const (
	EventPublish            = "publish_EventPublish"
	EventRepublish          = "publish_EventRepublish"
	EventPublishOrRepublish = "publish_EventRepublishOrRepublish"
	EventUnpublish          = "publish_EventUnpublish"

	EventDuplicateVersion      = "publish_EventDuplicateVersion"
	eventSelectVersion         = "publish_eventSelectVersion"
	eventSchedulePublishDialog = "publish_eventSchedulePublishDialog"
	eventSchedulePublish       = "publish_eventSchedulePublish"

	eventRenameVersionDialog = "publish_eventRenameVersionDialog"
	eventRenameVersion       = "publish_eventRenameVersion"
	eventDeleteVersionDialog = "publish_eventDeleteVersionDialog"
	eventDeleteVersion       = "publish_eventDeleteVersion"

	ActivityPublish            = "Publish"
	ActivityRepublish          = "Republish"
	ActivityPublishOrRepublish = "PublishOrRepublish"
	ActivityUnPublish          = "UnPublish"

	ParamScriptAfterPublish = "publish_param_script_after_publish"
)

func registerEventFuncsForResource(mb *presets.ModelBuilder, publisher *Builder) {
	mb.RegisterEventHandler(EventPublish, publishAction(mb, publisher, ActivityPublish))
	mb.RegisterEventHandler(EventRepublish, publishAction(mb, publisher, ActivityRepublish))
	mb.RegisterEventHandler(EventPublishOrRepublish, publishAction(mb, publisher, ActivityPublishOrRepublish))
	mb.RegisterEventHandler(EventUnpublish, unpublishAction(mb, publisher, ActivityUnPublish))

	mb.RegisterEventHandler(EventDuplicateVersion, duplicateVersionAction(mb, publisher))
	mb.RegisterEventHandler(eventSelectVersion, selectVersion(mb))
	mb.RegisterEventHandler(eventSchedulePublishDialog, schedulePublishDialog(publisher.db, mb))
	mb.RegisterEventHandler(eventSchedulePublish, schedulePublish(publisher.db, mb))
}

func registerEventFuncsForVersion(mb *presets.ModelBuilder, pm *presets.ModelBuilder, db *gorm.DB) {
	mb.RegisterEventHandler(eventRenameVersionDialog, renameVersionDialog(mb))
	mb.RegisterEventHandler(eventRenameVersion, renameVersion(mb))
	mb.RegisterEventHandler(eventDeleteVersionDialog, deleteVersionDialog(mb))
	mb.RegisterEventHandler(eventDeleteVersion, deleteVersion(mb, pm, db))
}
