package chat

import "testing"

func TestLoadChatDoesNotBackfillWorldIDsFromNames(t *testing.T) {
	pm, err := NewPersistenceManagerWithDir(t.TempDir())
	if err != nil {
		t.Fatalf("NewPersistenceManagerWithDir: %v", err)
	}
	chatData, err := pm.NewChat("Legacy Names")
	if err != nil {
		t.Fatalf("NewChat: %v", err)
	}
	chatData.Meta.WorldName = "Ashfall"
	chatData.Meta.ScenarioName = "The Siege"
	if err := pm.SaveChat(chatData); err != nil {
		t.Fatalf("SaveChat: %v", err)
	}

	loaded, err := pm.LoadChat(chatData.Meta.ID)
	if err != nil {
		t.Fatalf("LoadChat: %v", err)
	}
	if loaded.Meta.WorldID != "" {
		t.Errorf("WorldID = %q, want unresolved empty", loaded.Meta.WorldID)
	}
	if loaded.Meta.ScenarioID != "" {
		t.Errorf("ScenarioID = %q, want unresolved empty", loaded.Meta.ScenarioID)
	}
	if loaded.Meta.WorldName != "Ashfall" || loaded.Meta.ScenarioName != "The Siege" {
		t.Errorf("display names changed: %q / %q", loaded.Meta.WorldName, loaded.Meta.ScenarioName)
	}
}
