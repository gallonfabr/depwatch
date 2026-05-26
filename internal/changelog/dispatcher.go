package changelog

// Dispatcher routes processed entries to named output channels
// based on a routing decision function. It is designed to sit at
// the end of a pipeline and fan-out results to multiple consumers.
type Dispatcher struct {
	router  *Router
	sinks   map[string][]Entry
}

// NewDispatcher creates a Dispatcher backed by the given Router.
// Panics if router is nil.
func NewDispatcher(r *Router) *Dispatcher {
	if r == nil {
		panic("dispatcher: router must not be nil")
	}
	return &Dispatcher{
		router: r,
		sinks:  make(map[string][]Entry),
	}
}

// Dispatch routes each entry through the router and accumulates
// results per channel name. Previous dispatch results are cleared
// on every call.
func (d *Dispatcher) Dispatch(entries []Entry) {
	d.sinks = make(map[string][]Entry)
	for _, e := range entries {
		channels := d.router.Route(e)
		for _, ch := range channels {
			d.sinks[ch] = append(d.sinks[ch], e)
		}
	}
}

// Channels returns the sorted list of channel names that received
// at least one entry in the last Dispatch call.
func (d *Dispatcher) Channels() []string {
	keys := make([]string, 0, len(d.sinks))
	for k := range d.sinks {
		keys = append(keys, k)
	}
	return keys
}

// EntriesFor returns the entries routed to the named channel.
// Returns nil if the channel received no entries.
func (d *Dispatcher) EntriesFor(channel string) []Entry {
	return d.sinks[channel]
}
