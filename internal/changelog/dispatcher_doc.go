// Package changelog provides the dispatcher component which fans out
// processed changelog entries to named output channels using a Router.
//
// Usage:
//
//	router := changelog.NewRouter(
//		changelog.WithRouteRule("security", "security"),
//		changelog.WithRouterFallback("general"),
//	)
//	d := changelog.NewDispatcher(router)
//	d.Dispatch(entries)
//	for _, ch := range d.Channels() {
//		fmt.Println(ch, d.EntriesFor(ch))
//	}
package changelog
