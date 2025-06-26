package gmp

//go:generate go tool oapi-codegen --config=config_client.yaml --templates=. schema.yaml
//go:generate go tool oapi-codegen --config=config_models.yaml --templates=. schema.yaml
//go:generate go tool oapi-codegen --config=config_server.yaml --templates=. schema.yaml
//go:generate go tool oapi-codegen --config=config_spec.yaml --templates=. schema.yaml
