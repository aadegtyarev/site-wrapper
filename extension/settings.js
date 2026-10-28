import { text, validateDraft } from "./i18n.mjs";

const hostName = "io.github.aadegtyarev.sitewrapper";
const $ = id => document.getElementById(id);
let language = "en";
let firstSetup = false;
let currentPassword = "";
let launchURL = "";

function localize() {
  document.documentElement.lang = language;
  for (const element of document.querySelectorAll("[data-i18n]")) element.textContent = text(language, element.dataset.i18n);
  document.title = text(language, "title");
}

function notice(key) {
  $("notice").textContent = text(language, key);
  $("notice").hidden = false;
}

function clearErrors() {
  $("notice").hidden = true;
  for (const id of ["site-url", "zoom", "new-password", "confirm-password"]) $(id).removeAttribute("aria-invalid");
  for (const id of ["url-error", "zoom-error", "new-password-error", "confirm-password-error"]) $(id).textContent = "";
}

function showFieldError(field, key) {
  const map = { url: ["site-url", "url-error"], zoom_percent: ["zoom", "zoom-error"], new_password: ["new-password", "new-password-error"], confirm_password: ["confirm-password", "confirm-password-error"] };
  const [input, output] = map[field];
  $(input).setAttribute("aria-invalid", "true");
  $(output).textContent = text(language, key);
}

function fill(settings) {
  $("site-url").value = settings.url;
  $("zoom").value = settings.zoom_percent;
  $("mode").value = settings.mode;
  $("language").value = settings.language;
  language = settings.language;
  localize();
  $("unlock-form").hidden = true;
  $("settings-form").hidden = false;
}

async function host(message) {
  try {
    return await browser.runtime.sendNativeMessage(hostName, message);
  } catch (error) {
    console.error("Site Wrapper native helper unavailable", error);
    notice("bridgeFailed");
    return null;
  }
}

async function start() {
  try {
    const state = await browser.runtime.sendMessage({ action: "state" });
    firstSetup = state.needs_setup;
    launchURL = state.launch_url || "";
    language = state.language || "en";
    $("app-version").textContent = state.app_version;
    localize();
    if (firstSetup) {
      $("setup-badge").hidden = false;
      fill({ url: "", zoom_percent: 100, mode: "window", language: "en" });
    } else {
      $("unlock-form").hidden = false;
    }
  } catch (error) {
    console.error("Site Wrapper settings startup failed", error);
    notice("startupFailed");
  }
}

$("unlock-form").addEventListener("submit", async event => {
  event.preventDefault();
  const password = $("unlock-password").value;
  const response = await host({ action: "unlock", password });
  if (!response) return;
  if (!response.ok) {
    notice(response.error_code === "incorrect_password" ? "incorrectPassword" : "saveFailed");
    return;
  }
  currentPassword = password;
  $("unlock-password").value = "";
  clearErrors();
  fill(response.settings);
});

$("language").addEventListener("change", () => {
  language = $("language").value;
  localize();
});

$("settings-form").addEventListener("submit", async event => {
  event.preventDefault();
  clearErrors();
  const draft = { url: $("site-url").value.trim(), zoom_percent: Number($("zoom").value), mode: $("mode").value, language: $("language").value };
  const newPassword = $("new-password").value;
  const errors = validateDraft(draft, newPassword, $("confirm-password").value, firstSetup);
  if (Object.keys(errors).length) {
    for (const [field, key] of Object.entries(errors)) showFieldError(field, key);
    return;
  }
  $("save").disabled = true;
  const response = await host({ action: "save", password: currentPassword, new_password: newPassword, settings: draft });
  $("save").disabled = false;
  if (!response) return;
  if (!response.ok) {
    notice(response.error_code === "incorrect_password" ? "incorrectPassword" : "saveFailed");
    return;
  }
  try {
    await browser.runtime.sendMessage({ action: "saved", url: launchURL || response.settings.url, mode: response.settings.mode, zoom_percent: response.settings.zoom_percent, language: response.settings.language });
    const current = await browser.windows.getCurrent();
    await browser.windows.remove(current.id);
  } catch (error) {
    console.error("Site Wrapper could not open site", error);
    notice("startupFailed");
  }
});

$("cancel").addEventListener("click", async () => {
  const current = await browser.windows.getCurrent();
  await browser.windows.remove(current.id);
});
void start();
