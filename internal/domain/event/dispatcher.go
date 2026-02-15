package event

type SimpleDispatcher struct {
	handlers map[string][]EventHandler
}

func NewDispatcher() *SimpleDispatcher {
	return &SimpleDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

func (d *SimpleDispatcher) Subscribe(eventType string, handler EventHandler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

func (d *SimpleDispatcher) Dispatch(event Event) {
	if handlers, ok := d.handlers[event.Type()]; ok {
		for _, h := range handlers {
			go h(event)
		}
	}
}
