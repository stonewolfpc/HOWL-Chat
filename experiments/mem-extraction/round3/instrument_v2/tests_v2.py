"""Behavioural tests for instrument v2.

Run from this directory:

    python -m unittest tests_v2 -v

Each test states the expected behaviour, because "it ran" is not a result. The
tests use synthetic cases; they do not read or write the historical results.
"""

from __future__ import annotations

import copy
import os
import sys
import unittest
from pathlib import Path

sys.path.insert(0, os.path.dirname(__file__))

import atomicity          # noqa: E402
import evaluate           # noqa: E402
import reporting          # noqa: E402
import rubric_v2          # noqa: E402


def req(i: str, dec: str = "KEEP", *keywords: str) -> dict:
    return {"id": i, "decision": dec, "keys": [[k] for k in keywords], "points": ""}


def case(cid: str = "T1", required=None, forbidden=None, **kw) -> dict:
    c = {"id": cid, "category": "test", "length": "short", "conversation": "",
         "required": required or [], "forbidden": forbidden or []}
    c.update(kw)
    return c


def run(content: str | None, finish: str = "stop", **kw) -> dict:
    r = {"run_code": "r1", "label": "test", "draw": 1, "case": "T1",
         "kind": "corpus", "content": content, "finish_reason": finish}
    r.update(kw)
    return r


FOUR = [req("F1", "KEEP", "alpha"), req("F2", "KEEP", "bravo"),
        req("F3", "KEEP", "charlie"), req("F4", "KEEP", "delta")]

FOUR_ROWS = ("C1 | KEEP | alpha is a durable fact | durable\n"
             "C2 | KEEP | bravo is a durable fact | durable\n"
             "C3 | KEEP | charlie is a durable fact | durable\n"
             "C4 | KEEP | delta is a durable fact | durable")

FOUR_MERGED = ("C1 | KEEP | alpha, bravo, charlie and delta are all durable "
               "facts; alpha matters, bravo matters, charlie matters and "
               "delta matters | durable")


class TestAtomicity(unittest.TestCase):
    def test_four_facts_in_four_records_get_full_recall(self):
        s = rubric_v2.score_run(run(FOUR_ROWS), case(required=FOUR))
        self.assertEqual(s["status"], "VALID_KEEP")
        self.assertEqual(s["parts"]["recall"], 30.0)
        self.assertEqual(s["contract_parts"]["atomicity"], 25.0)
        self.assertTrue(s["production_usable"])

    def test_four_facts_merged_into_one_record_do_not_get_full_recall(self):
        s = rubric_v2.score_run(run(FOUR_MERGED), case(required=FOUR))
        self.assertEqual(s["parts"]["recall"], 15.0)
        verdicts = {a["verdict"] for a in s["atomicity"]}
        self.assertIn("MERGED_DISTINCT", verdicts)
        self.assertEqual(s["contract_parts"]["atomicity"], 0.0)
        self.assertLess(s["semantic"], 100.0)
        self.assertFalse(s["production_usable"])

    def test_declared_joint_fact_is_a_legitimate_compound(self):
        c = case(required=[req("J1", "KEEP", "halden"), req("J2", "KEEP", "cartographer")],
                 compound_ok=[["J1", "J2"]])
        out = "C1 | KEEP | Halden trained Sable as a cartographer | durable"
        s = rubric_v2.score_run(run(out), c)
        self.assertEqual({a["verdict"] for a in s["atomicity"]}, {"LEGITIMATE_COMPOUND"})
        self.assertEqual(s["parts"]["recall"], 30.0)

    def test_off_topic_padding_is_flagged_not_full_credit(self):
        c = case(required=[req("N1", "KEEP", "alpha")], noise_markers=["weather"])
        out = "C1 | KEEP | alpha matters and the weather was grey | durable"
        s = rubric_v2.score_run(run(out), c)
        self.assertEqual(s["atomicity"][0]["verdict"], "UNRELATED")
        self.assertEqual(s["parts"]["recall"], 15.0)


class TestOutputValidity(unittest.TestCase):
    def test_valid_none_scores_full_on_a_restraint_case(self):
        c = case(required=[], forbidden=["rain"])
        s = rubric_v2.score_run(run("NONE"), c)
        self.assertEqual(s["status"], "VALID_NONE")
        self.assertEqual(s["semantic"], 100.0)
        self.assertTrue(s["production_usable"])

    def test_valid_reject_only_is_a_valid_output(self):
        c = case(required=[], forbidden=["rain"])
        s = rubric_v2.score_run(run("C1 | REJECT | rain on the shutters | transient"),
                                c)
        self.assertEqual(s["status"], "VALID_REJECT")
        self.assertEqual(s["semantic"], 100.0)

    def test_valid_hold_is_recognised(self):
        c = case(required=[req("H1", "HOLD", "birthday", "frostfell")])
        s = rubric_v2.score_run(
            run("C1 | HOLD | birthday is 3 frostfell, maybe canon | uncertain"), c)
        self.assertEqual(s["status"], "VALID_HOLD")
        self.assertEqual(s["parts"]["uncertainty"], 5.0)
        self.assertEqual(s["parts"]["recall"], 30.0)

    def test_malformed_header_only_gets_no_restraint_credit(self):
        c = case(required=[], forbidden=["rain"])
        s = rubric_v2.score_run(run("C1 | DECISION | MEMORY | REASON"), c)
        self.assertEqual(s["status"], "MALFORMED")
        self.assertEqual(s["semantic"], 0.0)
        self.assertEqual(s["parts"]["restraint"], 0.0)
        self.assertFalse(s["production_usable"])

    def test_correct_fact_in_invalid_format_keeps_semantic_but_not_contract(self):
        c = case(required=[req("A1", "KEEP", "nurse")])
        out = "C1 | REALITY CHECK | The user is an ER nurse | because dosages"
        s = rubric_v2.score_run(run(out), c)
        self.assertEqual(s["status"], "UNPARSEABLE")
        self.assertEqual(s["parts"]["recall"], 30.0)
        self.assertLess(s["contract"], 50.0)
        self.assertFalse(s["production_usable"])

    def test_silence_on_a_required_case_earns_no_semantic_credit(self):
        c = case(required=FOUR)
        for out in ("NONE", ""):
            s = rubric_v2.score_run(run(out), c)
            self.assertEqual(s["semantic"], 0.0, out)
            self.assertEqual(s["parts"]["recall"], 0.0, out)

    def test_missing_output(self):
        s = rubric_v2.score_run(run(""), case(required=FOUR))
        self.assertEqual(s["status"], "MISSING")
        self.assertFalse(s["quality_eligible"])

    def test_truncated_output_is_flagged_and_not_quality_eligible(self):
        c = case(required=FOUR)
        out = ("C1 | KEEP | alpha is durable | durable\n"
               "C2 | KEEP | bravo is dura")
        s = rubric_v2.score_run(run(out, finish="length"), c)
        self.assertTrue(s["truncated"])
        self.assertFalse(s["quality_eligible"])
        self.assertEqual(s["contract_parts"]["termination"], 0.0)
        self.assertIsNotNone(s["semantic"])

    def test_truncated_with_no_rows_is_incomplete(self):
        s = rubric_v2.score_run(run("", finish="length"), case(required=FOUR))
        self.assertEqual(s["status"], "INCOMPLETE")

    def test_execution_failure_is_not_scored(self):
        s = rubric_v2.score_run({"run_code": "x", "label": "test", "draw": 1,
                                 "case": "T1", "kind": "corpus",
                                 "error": "connection refused"}, case(required=FOUR))
        self.assertEqual(s["status"], "EXECUTION_FAILURE")
        self.assertIsNone(s["semantic"])
        self.assertEqual([e["code"] for e in s["events"]], ["EXECUTION_FAILURE"])


class TestAttributionAndOwnership(unittest.TestCase):
    def test_incorrect_attribution_is_an_event_and_zeroes_the_dimension(self):
        c = case(required=[req("A1", "KEEP", "nurse")],
                 attribution_forbidden=[{"tokens": ["sable", "nurse"],
                                         "why": "real user fact given to the character"}])
        out = "C1 | KEEP | Sable is a nurse | durable"
        s = rubric_v2.score_run(run(out), c)
        self.assertEqual(s["parts"]["attribution"], 0.0)
        self.assertIn("ATTRIBUTION_ERROR", [e["code"] for e in s["events"]])

    def _king_case(self):
        return case(required=[], ownership={
            "scenario": "court",
            "characters": {"jim": ["Jim"], "lisa": ["Lisa"], "bob": ["Bob"]},
            "claims": [
                {"id": "K1", "topic": "king", "topic_hints": ["king"],
                 "holders": ["jim"], "scope": "character_belief",
                 "polarity": "negative", "polarity_hints": ["criminal"],
                 "exclusive": True, "must_not_promote": True},
                {"id": "K2", "topic": "king", "topic_hints": ["king"],
                 "holders": ["lisa"], "scope": "character_belief",
                 "polarity": "positive", "polarity_hints": ["wonderful"],
                 "exclusive": True, "must_not_promote": True},
            ]})

    def test_character_belief_promoted_to_world_truth_is_critical(self):
        s = rubric_v2.score_run(
            run("C1 | KEEP | The king is a criminal | durable"), self._king_case())
        codes = [e["code"] for e in s["events"]]
        self.assertIn("PROMOTION", codes)
        self.assertEqual(s["parts"]["ownership"], 0.0)
        self.assertFalse(s["production_usable"])

    def test_two_contradictory_beliefs_collapsed_into_one_claim(self):
        out = "C1 | KEEP | The king is a criminal and wonderful | durable"
        s = rubric_v2.score_run(run(out), self._king_case())
        self.assertIn("CONTRADICTION_COLLAPSE", [e["code"] for e in s["events"]])

    def test_belief_kept_with_holder_is_not_a_violation(self):
        out = "C1 | KEEP | Jim believes the king is a criminal | durable"
        s = rubric_v2.score_run(run(out), self._king_case())
        self.assertEqual([f["kind"] for f in s["ownership"]["findings"]
                          if f["kind"] in ("PROMOTION", "INHERITANCE")], [])

    def test_private_fact_shared_with_another_character(self):
        c = case(required=[], ownership={
            "characters": {"jim": ["Jim"], "lisa": ["Lisa"], "bob": ["Bob"]},
            "claims": [{"id": "D1", "topic": "vault", "topic_hints": ["password"],
                        "holders": ["jim", "lisa"], "scope": "participant_disclosure",
                        "exclusive": True, "must_not_promote": True}]})
        out = "C1 | KEEP | Bob knows the password is hunter2 | durable"
        s = rubric_v2.score_run(run(out), c)
        self.assertIn("INHERITANCE", [e["code"] for e in s["events"]])

    def test_scenario_information_does_not_migrate(self):
        c = case(required=[], ownership={
            "scenario": "court", "characters": {},
            "scenario_markers": {"court": ["court", "herald"],
                                 "harbour": ["tide", "dock"]},
            "claims": []})
        out = "C1 | KEEP | The tide gate jams every winter | durable"
        s = rubric_v2.score_run(run(out), c)
        self.assertIn("SCENARIO_LEAK", [e["code"] for e in s["events"]])

    def test_unscored_ownership_is_reported_as_unscored(self):
        s = rubric_v2.score_run(run(FOUR_ROWS), case(required=FOUR))
        self.assertFalse(s["ownership"]["scored"])


class TestIncidentAccounting(unittest.TestCase):
    def test_one_incident_reduces_one_dimension_and_stays_visible(self):
        c = case(required=[req("F1", "KEEP", "alpha")], forbidden=["alpha secret"])
        control = rubric_v2.score_run(run("C1 | KEEP | alpha fact | durable"), c)
        hurt = rubric_v2.score_run(
            run("C1 | KEEP | alpha secret fact | durable"), c)

        self.assertEqual(control["parts"]["fidelity"], 18.0)
        self.assertEqual(hurt["parts"]["fidelity"], 0.0)
        self.assertEqual(len(hurt["events"]), 1)
        ev = hurt["events"][0]
        self.assertEqual(ev["code"], "FORBIDDEN_CLAIM_STORED")
        self.assertEqual(ev["primary_dimension"], "fidelity")
        self.assertEqual(ev["secondary_impacts"], ["false_memory"])

        differing = [k for k in control["parts"]
                     if control["parts"][k] != hurt["parts"][k]]
        self.assertEqual(differing, ["fidelity"])

    def test_negated_forbidden_phrase_is_not_a_stored_falsehood(self):
        c = case(required=[req("N1", "KEEP", "toll")], forbidden=["triple"])
        out = ("C1 | KEEP | late salt doubles the toll, not triple | "
               "exact numeric rule")
        s = rubric_v2.score_run(run(out), c)
        self.assertEqual(s["parts"]["fidelity"], 18.0)
        self.assertNotIn("FORBIDDEN_CLAIM_STORED", [e["code"] for e in s["events"]])

    def test_one_incident_with_several_phrases_is_one_event(self):
        c = case(required=[req("F1", "KEEP", "password")],
                 forbidden=["hunter2", "password is hunter2"])
        out = "C1 | KEEP | The user's password is hunter2 | in-world dialogue"
        s = rubric_v2.score_run(run(out), c)
        events = [e for e in s["events"] if e["code"] == "FORBIDDEN_CLAIM_STORED"]
        self.assertEqual(len(events), 1)
        self.assertIn("hunter2", events[0]["detail"])

    def test_score_run_does_not_mutate_its_input(self):
        r = run(FOUR_ROWS)
        before = copy.deepcopy(r)
        rubric_v2.score_run(r, case(required=FOUR))
        self.assertEqual(r, before)


class TestReportingLanguage(unittest.TestCase):
    def test_single_draw_is_labelled_as_such(self):
        self.assertIn("single-draw", reporting.sample_qualifier(1))
        self.assertIn("underpowered", reporting.sample_qualifier(3))

    def test_thin_samples_cannot_claim_no_difference(self):
        v = reporting.compare(77.1, (64.6, 88.4), 3, 75.9, (64.3, 87.0), 1)
        self.assertEqual(v["claim"], "UNESTABLISHED")
        self.assertIn("not evidence that a difference does not exist", v["allowed"])
        self.assertIn("no measured", v["allowed"])

    def test_overlapping_intervals_are_not_equivalence(self):
        v = reporting.compare(70, (60, 80), 8, 68, (58, 78), 8)
        self.assertEqual(v["claim"], "NOT_ESTABLISHED")
        self.assertIn("equivalence is not established", v["allowed"])

    def test_scanner_finds_overclaiming_phrases(self):
        hits = reporting.scan("The result proves the model is safer and does not "
                              "beat the baseline.")
        codes = {h["code"] for h in hits}
        self.assertIn("PROOF", codes)
        self.assertIn("NEGATIVE_CLAIM", codes)
        for h in hits:
            self.assertIn(h["match"], h["sentence"])

    def test_isolation_and_timing_claims_are_bounded(self):
        self.assertIn("HTTP 501", reporting.ISOLATION_ALLOWED)
        self.assertIn("excludes model load", reporting.TIMING_ALLOWED)
        self.assertEqual(reporting.PREREGISTRATION["status"], "UNESTABLISHED")


class TestWriteGuards(unittest.TestCase):
    def test_refuses_to_write_into_historical_results(self):
        with self.assertRaises(SystemExit):
            evaluate.guard_out(Path("results/corpus_scores.json"), force=True)

    def test_refuses_to_overwrite_without_force(self):
        p = Path(__file__)
        with self.assertRaises(SystemExit):
            evaluate.guard_out(p, force=False)


if __name__ == "__main__":
    unittest.main(verbosity=2)
