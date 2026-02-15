package entity

type Source struct {
	id           string
	name         string
	description  string
	sourceType   string
	config       map[string]any
	fieldMapping map[string]string
	transforms   []Transform
	display      DisplayConfig
	enabled      bool
}

type Transform struct {
	Type   string `yaml:"type"`
	Field  string `yaml:"field"`
	Format string `yaml:"format,omitempty"`
	Value  string `yaml:"value,omitempty"`
}

type DisplayConfig struct {
	Icon     string `yaml:"icon"`
	Color    string `yaml:"color"`
	Priority int    `yaml:"priority"`
}

func NewSource(id, name, sourceType string) *Source {
	return &Source{
		id:           id,
		name:         name,
		sourceType:   sourceType,
		config:       make(map[string]any),
		fieldMapping: make(map[string]string),
		transforms:   []Transform{},
		enabled:      true,
	}
}

func (s *Source) ID() string                      { return s.id }
func (s *Source) Name() string                    { return s.name }
func (s *Source) Description() string             { return s.description }
func (s *Source) Type() string                    { return s.sourceType }
func (s *Source) Config() map[string]any          { return s.config }
func (s *Source) FieldMapping() map[string]string { return s.fieldMapping }
func (s *Source) Transforms() []Transform         { return s.transforms }
func (s *Source) Display() DisplayConfig          { return s.display }
func (s *Source) Enabled() bool                   { return s.enabled }

func (s *Source) SetDescription(d string)              { s.description = d }
func (s *Source) SetConfig(c map[string]any)           { s.config = c }
func (s *Source) SetFieldMapping(fm map[string]string) { s.fieldMapping = fm }
func (s *Source) SetTransforms(t []Transform)          { s.transforms = t }
func (s *Source) SetDisplay(d DisplayConfig)           { s.display = d }
func (s *Source) SetEnabled(e bool)                    { s.enabled = e }

func (s *Source) ToDTO() *SourceDTO {
	return &SourceDTO{
		ID:           s.id,
		Name:         s.name,
		Description:  s.description,
		Type:         s.sourceType,
		Config:       s.config,
		FieldMapping: s.fieldMapping,
		Transforms:   s.transforms,
		Display:      s.display,
		Enabled:      s.enabled,
	}
}

type SourceDTO struct {
	ID           string            `yaml:"id" json:"id"`
	Name         string            `yaml:"name" json:"name"`
	Description  string            `yaml:"description" json:"description"`
	Type         string            `yaml:"type" json:"type"`
	Config       map[string]any    `yaml:"config" json:"config"`
	FieldMapping map[string]string `yaml:"field_mapping" json:"field_mapping"`
	Transforms   []Transform       `yaml:"transforms" json:"transforms"`
	Display      DisplayConfig     `yaml:"display" json:"display"`
	Enabled      bool              `yaml:"enabled" json:"enabled"`
}

func SourceFromDTO(dto *SourceDTO) *Source {
	s := NewSource(dto.ID, dto.Name, dto.Type)
	s.description = dto.Description
	s.config = dto.Config
	s.fieldMapping = dto.FieldMapping
	s.transforms = dto.Transforms
	s.display = dto.Display
	s.enabled = dto.Enabled
	if s.config == nil {
		s.config = make(map[string]any)
	}
	if s.fieldMapping == nil {
		s.fieldMapping = make(map[string]string)
	}
	return s
}
