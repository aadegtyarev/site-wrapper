"use strict";

browser.runtime.sendMessage({ action: "mode" }).then(({ mode, language }) => {
  if (mode !== "window" || document.getElementById("site-wrapper-settings-entry")) return;
  const button = document.createElement("button");
  button.id = "site-wrapper-settings-entry";
  button.type = "button";
  button.textContent = "⚙";
  const labels = { en: "Open Site Wrapper settings", de: "Site Wrapper-Einstellungen öffnen", kk: "Site Wrapper баптауларын ашу", it: "Apri le impostazioni di Site Wrapper" };
  button.setAttribute("aria-label", labels[language] || labels.en);
  button.title = (labels[language] || labels.en) + " (Ctrl+Alt+S)";
  Object.assign(button.style, {
    position: "fixed", right: "12px", bottom: "12px", width: "40px", height: "40px",
    zIndex: "2147483647", border: "1px solid #b6c4d4", borderRadius: "12px",
    background: "#fff", color: "#18354f", cursor: "pointer", boxShadow: "0 4px 18px #001b3433"
  });
  button.addEventListener("click", () => { void browser.runtime.sendMessage({ action: "open_settings" }); });
  document.documentElement.append(button);
}).catch(error => console.error("Site Wrapper settings entry failed", error));
