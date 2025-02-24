package types

type GetJobById struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type AddPluginByName struct {
	Name string `uri:"name" binding:"required"`
}
