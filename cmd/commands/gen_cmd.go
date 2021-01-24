package commands

type GeneratorCommand struct {
	Server *ServerGenerateCmd `command:"server"`
	Client *ClientGenerateCmd `command:"client"`
	Spec   *SpecGenerateCmd   `command:"spec"`
}
