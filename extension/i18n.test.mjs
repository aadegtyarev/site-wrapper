import assert from "node:assert/strict";
import { test } from "node:test";
import { messages, text, validateDraft } from "./i18n.mjs";

test("all four languages contain the same settings and error keys", () => {
  const base = Object.keys(messages.en).sort();
  for (const language of ["en", "de", "kk", "it"]) {
    assert.deepEqual(Object.keys(messages[language]).sort(), base, language);
    for (const value of Object.values(messages[language])) assert.ok(value.trim());
  }
  assert.equal(text("unknown", "save"), messages.en.save);
});

test("settings validation rejects unsafe URL, zoom, and incomplete password", () => {
  const valid = { url: "https://example.org", zoom_percent: 100, mode: "window", language: "en" };
  assert.deepEqual(validateDraft(valid, "long-test-password", "long-test-password", true), {});
  assert.equal(validateDraft({ ...valid, url: "javascript:alert(1)" }, "long-test-password", "long-test-password", true).url, "invalidUrl");
  assert.equal(validateDraft({ ...valid, url: "https://user:pass@example.org" }, "long-test-password", "long-test-password", true).url, "invalidUrl");
  assert.equal(validateDraft({ ...valid, zoom_percent: 501 }, "long-test-password", "long-test-password", true).zoom_percent, "invalidZoom");
  assert.equal(validateDraft(valid, "", "", true).new_password, "passwordRequired");
  assert.equal(validateDraft(valid, "long-test-password", "different-password", true).confirm_password, "passwordMismatch");
});
