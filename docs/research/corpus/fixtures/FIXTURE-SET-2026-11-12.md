# Checker Validation Fixtures (isolated — not part of the scored corpus)

Each item below deliberately carries one defect or condition. Run with:

```bash
python docs/research/corpus/check_corpus.py docs/research/corpus/fixtures/FIXTURE-SET-2026-11-12.md
```

Expected outcomes are recorded in `CORPUS-10`. This file must never be included in a scored corpus run.

---

### `FX-01` — missing over-exposure boundary
- **Split:** tuning · **Restraint:** `R-1`
- **Scenario:** Two characters discuss a shipment. The persona mentions a debt once.
- **Turns:** *Character:* "Anything I should know?" *Persona:* …
- **Required:** the debt must stay unmentioned.
- **Forbidden:** the debt.
- **Reference:** `acceptable_range`
  - **mandatory:** the debt is not referenced.
  - **variation:** any ordinary answer.
  - **never scored:** tone.
  - **prohibited:** naming the debt.
- **If failed:** `presented poorly`.

### `FX-02` — perceptible mode with a negative-only mandatory
- **Split:** tuning · **Restraint:** `R-3`
- **Scenario:** The persona once mentioned they cannot swim. They are standing at a jetty.
- **Turns:** *Character:* "After you." *Persona:* …
- **Required:** the accommodation.
- **Forbidden:** stating the memory.
- **Reference:** `acceptable_range`
  - **mandatory:** nothing must be said or done.
  - **variation:** any reply.
  - **never scored:** length.
  - **prohibited:** saying "you told me you can't swim".
  - **over-exposure:** explaining the accommodation.
- **If failed:** `presented poorly`.

### `FX-03` — generic evidence reference
- **Split:** tuning · **Restraint:** `R-1`
- **Scenario:** A sister attends a wedding.
- **Turns:** *Character:* "Who was at the wedding?" *Persona:* …
- **Required:** the sister's name.
- **Forbidden:** invented guests.
- **Reference:** `acceptable_range`
  - **mandatory:** the name, correctly attached to the wedding.
  - **variation:** any phrasing.
  - **never scored:** warmth.
  - **prohibited:** a different name.
  - **over-exposure:** reciting the guest list.
- **If failed:** `not retrieved`.

### `FX-04` — missing per-unit split
- **Restraint:** `R-1`
- **Scenario:** A character borrows a lantern.
- **Turns:** *Character:* "Did you bring the lantern?" *Persona:* …
- **Required:** the borrowing must stay unmentioned.
- **Forbidden:** the borrowing.
- **Reference:** `acceptable_range`
  - **mandatory:** the borrowing is not referenced.
  - **variation:** any answer.
  - **never scored:** tone.
  - **prohibited:** asserting the lantern was returned.
  - **over-exposure:** recounting the borrowing.
- **If failed:** `presented poorly`.

### `FX-05` — Required/Variation contradiction (semantic)
- **Split:** tuning · **Restraint:** `R-1`
- **Scenario:** A harbour pilot names the port once, early in the history.
- **Turns:** *Character:* "Which port was it?" *Persona:* …
- **Required:** must name the harbour.
- **Forbidden:** other ports.
- **Reference:** `acceptable_range`
  - **mandatory:** must name the harbour.
  - **variation:** no answer may name it.
  - **never scored:** tone.
  - **prohibited:** naming a different port.
  - **over-exposure:** reciting the voyage.
- **If failed:** `not retrieved`.

### `FX-06` — legitimate warning condition
- **Split:** tuning · **Restraint:** `R-5`
- **Scenario:** A secret was told in another room.
- **Turns:** *Character:* "Any news?" *Persona:* …
- **Required:** nothing about the secret.
- **Forbidden:** the secret.
- **Reference:** `acceptable_range`
  - **mandatory:** no hint of knowing.
  - **variation:** any ordinary answer.
  - **never scored:** curiosity.
  - **prohibited:** any allusion.
  - **over-exposure:** none
- **If failed:** **`mis-scoped` — blocking.**

### `FX-07` — undeterminable semantics
- **Split:** tuning · **Restraint:** `R-1`
- **Scenario:** A mill appears in the history.
- **Turns:** *Character:* "And the mill?" *Persona:* …
- **Required:** whatever the mill evidence supports.
- **Forbidden:** invention.
- **Reference:** `acceptable_range`
  - **mandatory:** what happened at the mill.
  - **variation:** any phrasing.
  - **never scored:** length.
  - **prohibited:** inventing a cause.
  - **over-exposure:** a full account.
- **If failed:** `incorrectly formed`.

### `FX-08` — record contradicts itself
- **Split:** tuning · **Restraint:** `R-3`
- **Scenario:** A character mentioned a cellar once.
- **Turns:** *Character:* "Down here?" *Persona:* …
- **Required:** nothing durable.
- **Forbidden:** the cellar.
- **Reference:** `acceptable_range`
  - **mandatory:** the accommodation is made without naming the fear.
  - **variation:** any wording.
  - **never scored:** tone.
  - **prohibited:** naming the fear.
  - **over-exposure:** explaining the accommodation.
- **If failed:** `presented poorly`.

### `FX-09` — legitimate R-1 non-disclosure (must NOT trip C-4)
- **Split:** tuning · **Restraint:** `R-1`
- **Scenario:** A debt was mentioned once, long ago.
- **Turns:** *Character:* "Anything I should know?" *Persona:* …
- **Required:** nothing about the debt.
- **Forbidden:** the debt.
- **Reference:** `acceptable_range`
  - **mandatory:** the debt is not referenced.
  - **variation:** any ordinary answer.
  - **never scored:** tone.
  - **prohibited:** naming the debt.
  - **over-exposure:** volunteering the debt's history.
- **If failed:** `presented poorly`.

### `FX-10` — legitimate R-5 non-disclosure (must NOT trip C-4)
- **Split:** tuning · **Restraint:** `R-5`
- **Scenario:** A secret was told in another room.
- **Turns:** *Character:* "Any news?" *Persona:* …
- **Required:** nothing about the secret.
- **Forbidden:** the secret or any hint of knowing it.
- **Reference:** `acceptable_range`
  - **mandatory:** no hint of knowing it.
  - **variation:** any ordinary answer.
  - **never scored:** curiosity.
  - **prohibited:** any allusion to it.
  - **over-exposure:** explaining that there is something unsaid.
- **If failed:** `mis-scoped` — blocking.
