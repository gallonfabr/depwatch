package changelog

// Router dispatches entries to named channels based on label rules.
// Each rule maps a label to a channel name; entries matching multiple
// rules are delivered to all matching channels. Entries with no match
// are routed to the fallback channel (default: "default").
type Router struct {
	rules    []routeRule
	fallback string
}

type routeRule struct {
	label   string
	channel string
}

// RouterOption configures a Router.
type RouterOption func(*Router)

// WithRouteRule adds a rule that sends entries carrying label to channel.
func WithRouteRule(label, channel string) RouterOption {
	return func(r *Router) {
		r.rules = append(r.rules, routeRule{label: label, channel: channel})
	}
}

// WithRouterFallback sets the channel used when no rule matches.
func WithRouterFallback(channel string) RouterOption {
	return func(r *Router) {
		r.fallback = channel
	}
}

// NewRouter creates a Router with the supplied options.
func NewRouter(opts ...RouterOption) *Router {
	r := &Router{fallback: "default"}
	for _, o := range opts {
		o(r)
	}
	return r
}

// Route distributes entries into a map of channel name → entries.
// Every entry appears in at least one channel.
func (r *Router) Route(entries []Entry) map[string][]Entry {
	out := make(map[string][]Entry)
	for _, e := range entries {
		matched := false
		for _, rule := range r.rules {
			if containsLabel(e.Labels, rule.label) {
				out[rule.channel] = append(out[rule.channel], e)
				matched = true
			}
		}
		if !matched {
			out[r.fallback] = append(out[r.fallback], e)
		}
	}
	return out
}

func containsLabel(labels []string, target string) bool {
	for _, l := range labels {
		if l == target {
			return true
		}
	}
	return false
}
