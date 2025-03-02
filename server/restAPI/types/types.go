package types

type JobById struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type PluginByName struct {
	Name string `uri:"name" binding:"required"`
}

type BackendByName struct {
	Name string `uri:"name" binding:"required"`
}
