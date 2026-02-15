package entity

type Profile struct {
	name           string
	displayName    string
	description    string
	prompts        PromptConfig
	sourceGroups   map[string][]string
	defaultSources []string
	active         bool
}

type PromptConfig struct {
	Summarizer string `yaml:"summarizer"`
	Suggester  string `yaml:"suggester"`
}

func NewProfile(name, displayName string) *Profile {
	return &Profile{
		name:           name,
		displayName:    displayName,
		sourceGroups:   make(map[string][]string),
		defaultSources: []string{},
		active:         false,
	}
}

func (p *Profile) Name() string                      { return p.name }
func (p *Profile) DisplayName() string               { return p.displayName }
func (p *Profile) Description() string               { return p.description }
func (p *Profile) Prompts() PromptConfig             { return p.prompts }
func (p *Profile) SourceGroups() map[string][]string { return p.sourceGroups }
func (p *Profile) DefaultSources() []string          { return p.defaultSources }
func (p *Profile) Active() bool                      { return p.active }

func (p *Profile) SetDescription(d string)                { p.description = d }
func (p *Profile) SetPrompts(pr PromptConfig)             { p.prompts = pr }
func (p *Profile) SetSourceGroups(sg map[string][]string) { p.sourceGroups = sg }
func (p *Profile) SetDefaultSources(ds []string)          { p.defaultSources = ds }
func (p *Profile) SetActive(a bool)                       { p.active = a }

func (p *Profile) ToDTO() *ProfileDTO {
	return &ProfileDTO{
		Name:           p.name,
		DisplayName:    p.displayName,
		Description:    p.description,
		Prompts:        p.prompts,
		SourceGroups:   p.sourceGroups,
		DefaultSources: p.defaultSources,
		Active:         p.active,
	}
}

type ProfileDTO struct {
	Name           string              `yaml:"name" json:"name"`
	DisplayName    string              `yaml:"display_name" json:"display_name"`
	Description    string              `yaml:"description" json:"description"`
	Prompts        PromptConfig        `yaml:"prompts" json:"prompts"`
	SourceGroups   map[string][]string `yaml:"source_groups" json:"source_groups"`
	DefaultSources []string            `yaml:"default_sources" json:"default_sources"`
	Active         bool                `yaml:"active" json:"active"`
}

func ProfileFromDTO(dto *ProfileDTO) *Profile {
	p := NewProfile(dto.Name, dto.DisplayName)
	p.description = dto.Description
	p.prompts = dto.Prompts
	p.sourceGroups = dto.SourceGroups
	p.defaultSources = dto.DefaultSources
	p.active = dto.Active
	if p.sourceGroups == nil {
		p.sourceGroups = make(map[string][]string)
	}
	return p
}
