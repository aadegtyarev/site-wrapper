import { text } from "./i18n.mjs";

let language = "en";
try {
  const state = await browser.runtime.sendMessage({ action: "mode" });
  language = state.language || "en";
} catch (error) {
  console.error("Sitewrap could not determine error language", error);
}
document.documentElement.lang = language;
for (const element of document.querySelectorAll("[data-i18n]")) element.textContent = text(language, element.dataset.i18n);
document.getElementById("error-id").textContent = new URL(location.href).searchParams.get("id") || "-";
