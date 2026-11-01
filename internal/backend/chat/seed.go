package chat

// currentSeed reports the seed recorded for this chat's generations.
// llama treats -1 as "pick a random seed", so a stored -1 stays meaningful;
// any other value is the seed the user pinned.
func (s *Service) currentSeed() int {
	if s.inferenceOpts == nil {
		return -1
	}
	return s.inferenceOpts.Seed
}
