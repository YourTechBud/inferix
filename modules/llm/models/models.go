package models

// Models stores the model configurations
type Models struct {
	models map[string]Config
}

// New initializes the models map with the given configurations
func New(modelConfigs []Config) *Models {
	models := make(map[string]Config)
	for _, model := range modelConfigs {
		if model.DefaultOptions == nil {
			model.DefaultOptions = DefaultModelOptions()
		}
		models[model.Name] = model
	}
	return &Models{models}
}
