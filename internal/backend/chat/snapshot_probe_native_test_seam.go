//go:build native_test_seam

package chat

// NativeTestSnapshotIDs returns only the IDs of the loaded memory snapshot.
// It exists for the isolated native test seam and is absent from production.
func (s *Service) NativeTestSnapshotIDs() []string {
	ids := make([]string, 0, len(s.memorySnapshot))
	for _, entry := range s.memorySnapshot {
		ids = append(ids, entry.ID)
	}
	return ids
}

// NativeTestSaveChatData exercises the chat-owned save path, including the
// generated-lore synchronization loop, from the isolated native seam. It is
// absent from production builds.
func (s *Service) NativeTestSaveChatData(chatName string) error {
	if s.memoryService != nil {
		s.memoryService.SetChatID(chatName)
	}
	return s.SaveChatData(chatName, true)
}
