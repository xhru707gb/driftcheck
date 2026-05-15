package tfstate

// ApplyWatchlist returns a new State containing only resources matched by the
// watchlist, with attributes filtered to those being watched.
// If the watchlist is empty, the original state is returned unchanged.
func ApplyWatchlist(s *State, wl *Watchlist) *State {
	if s == nil {
		return nil
	}
	if wl == nil || len(wl.Entries()) == 0 {
		return s
	}

	out := NewState()
	for _, key := range s.Keys() {
		if !wl.Matches(key) {
			continue
		}
		res, ok := s.Get(key)
		if !ok {
			continue
		}

		watched := wl.WatchedAttributes(key)
		if len(watched) == 0 {
			// watch all attributes
			out.Add(res)
			continue
		}

		// build a filtered copy of the resource
		watchSet := toWatchSet(watched)
		filteredAttrs := make(map[string]interface{}, len(watched))
		for _, attr := range watched {
			if v, exists := res.Attributes[attr]; exists {
				filteredAttrs[attr] = v
			}
		}
		_ = watchSet

		filtered := Resource{
			Type:       res.Type,
			Name:       res.Name,
			ID:         res.ID,
			Attributes: filteredAttrs,
		}
		out.Add(filtered)
	}
	return out
}

func toWatchSet(attrs []string) map[string]struct{} {
	s := make(map[string]struct{}, len(attrs))
	for _, a := range attrs {
		s[a] = struct{}{}
	}
	return s
}
