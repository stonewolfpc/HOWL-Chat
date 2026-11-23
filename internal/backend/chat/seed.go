package chat

// currentSeed reports the seed recorded for this chat's generations.
// llama treats -1 as "pick a random seed", so a stored -1 stays meaningful;
// any other value is the seed the user pinned.
func (s *Service) currentSeed() int {
	// The seed authority is the sampler settings the host applies to each
	// generation. This service no longer holds inference options of its own, so
	// the legacy chat.json copy records the default rather than a stale mirror.
	return -1
}
