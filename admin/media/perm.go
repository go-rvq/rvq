package media

import "net/http"

// The media library is a model of the admin (media_libraries, scope
// presets): what the fields of pictures and files do with it in a form is
// asked of it, as of any model —
//
//	uploading a file      :…:media_libraries:@create
//	deleting one          :…:media_libraries:<5>:@delete
//	describing one        :…:media_libraries:<5>:@edit
//
// — through its groups, or by its unique name (:media_libraries:…).

func (mb *Builder) uploadIsAllowed(r *http.Request) error {
	return mb.model.Permissioner().ReqCreator(r).IsAllowed()
}

func (mb *Builder) deleteIsAllowed(r *http.Request, obj interface{}) error {
	return mb.model.Permissioner().ReqObjectDeleter(r, obj).IsAllowed()
}

func (mb *Builder) updateDescIsAllowed(r *http.Request, obj interface{}) error {
	return mb.model.Permissioner().ReqObjectUpdater(r, obj).IsAllowed()
}
