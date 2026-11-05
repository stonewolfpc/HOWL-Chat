//go:build checkfa
package scanner

import "testing"

func TestCheckFa0e6cfc(t *testing.T) {
    s := New("H:/InvokeAI/models")
    idx, err := s.Scan()
    if err != nil { t.Fatal("Scan error:", err) }
    t.Logf("Models: %d, ControlNets: %d", len(idx.Models), len(idx.ControlNets))
    for _, e := range idx.Models {
        if "fa0e6cfc" <= e.Path && e.Path < "fa0e6cfcz" {
            t.Logf("  MODEL: %s | type=%s | family=%s | path=%s", e.Name, e.Type, e.Family, e.Path)
        }
    }
    for _, e := range idx.ControlNets {
        if "fa0e6cfc" <= e.Path && e.Path < "fa0e6cfcz" {
            t.Logf("  CN: %s | type=%s | family=%s | path=%s", e.Name, e.Type, e.Family, e.Path)
        }
    }
}
