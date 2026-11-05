package scanner

import (
	"fmt"
	"testing"
)

func TestLiveScanInvokeAIBundle(t *testing.T) {
	modelsDir := "H:/InvokeAI/models"
	s := New(modelsDir)
	idx, err := s.Scan()
	if err != nil {
		t.Fatal("Scan error:", err)
	}
	fmt.Printf("Scanned %s: %d models, %d VAEs, %d LoRAs, %d ControlNets\n",
		modelsDir, len(idx.Models), len(idx.VAEs), len(idx.LoRAs), len(idx.ControlNets))
	fmt.Println("\n=== UUID-NAMED MODELS ===")
	for _, m := range idx.Models {
		if len(m.Name) > 8 && m.Name[0] >= '0' && m.Name[0] <= '9' {
			fmt.Printf("  %s | family=%s | format=%s | size=%.1fMB\n", m.Name, m.Family, m.Format, m.SizeMB)
		}
	}
	fmt.Println("\n=== ALL MODELS ===")
	for _, m := range idx.Models {
		fmt.Printf("  %-40s family=%-10s format=%-10s size=%.1fMB\n", m.Name, m.Family, m.Format, m.SizeMB)
	}
	fmt.Println("\n=== CONTROL NETS ===")
	for _, c := range idx.ControlNets {
		fmt.Printf("  %-40s family=%-10s format=%-10s size=%.1fMB\n", c.Name, c.Family, c.Format, c.SizeMB)
	}
	fmt.Println("\n=== VAEs ===")
	for _, v := range idx.VAEs {
		fmt.Printf("  %-40s family=%-10s format=%-10s size=%.1fMB\n", v.Name, v.Family, v.Format, v.SizeMB)
	}
}
