"use strict";

const HOST = "io.github.aadegtyarev.sitewrapper";
let startup = null;
let booting = null;

async function host(message) {
  return browser.runtime.sendNativeMessage(HOST, message);
}

async function closeOtherWindows(keepId) {
  const windows = await browser.windows.getAll();
  for (const window of windows) {
    if (window.id !== keepId && window.type === "normal") {
      await browser.windows.remove(window.id);
    }
  }
}

async function openSite(url, mode, zoomPercent, language) {
  const window = await browser.windows.create({ url, type: "popup" });
  startup = { ...startup, needs_setup: false, launch_url: url, settings: { ...startup?.settings, mode, zoom_percent: zoomPercent, language } };
  if (mode === "kiosk") {
    await browser.windows.update(window.id, { state: "fullscreen" });
  }
  await closeOtherWindows(window.id);
  for (const tab of window.tabs || []) {
    if (tab.id !== undefined) {
      await browser.tabs.setZoom(tab.id, zoomPercent / 100);
    }
  }
}

async function openSettings() {
  const window = await browser.windows.create({ url: browser.runtime.getURL("settings.html"), type: "popup", width: 640, height: 780 });
  return window.id;
}

async function bootstrap() {
  if (booting) return booting;
  booting = (async () => {
    try {
      const result = await host({ action: "startup" });
      if (!result.ok) throw new Error(result.error_code || "startup_failed");
      startup = result;
      if (result.needs_setup) {
        const windowId = await openSettings();
        await closeOtherWindows(windowId);
      } else {
        await openSite(result.launch_url, result.settings.mode, result.settings.zoom_percent, result.settings.language);
      }
    } catch (error) {
      const errorId = crypto.randomUUID();
      console.error("Site Wrapper startup failed", errorId, error);
      await browser.windows.create({ url: browser.runtime.getURL("error.html?id=" + encodeURIComponent(errorId)), type: "popup", width: 560, height: 360 });
    }
  })();
  return booting;
}

browser.runtime.onStartup.addListener(() => { void bootstrap(); });
browser.runtime.onInstalled.addListener(() => { void bootstrap(); });
browser.commands.onCommand.addListener(command => {
  if (command === "open-settings") void openSettings();
});
browser.runtime.onMessage.addListener(async (message, sender) => {
  switch (message.action) {
    case "state":
      if (!startup) await bootstrap();
      return { needs_setup: startup?.needs_setup ?? true, launch_url: startup?.needs_setup ? startup.launch_url : "", language: startup?.settings?.language || "en", app_version: startup?.app_version || "dev" };
    case "open_settings":
      await openSettings();
      return { ok: true };
    case "mode":
      return { mode: startup?.settings?.mode || "window", language: startup?.settings?.language || "en" };
    case "saved":
      if (sender.url !== browser.runtime.getURL("settings.html")) return { ok: false, error_code: "forbidden" };
      await openSite(message.url, message.mode, message.zoom_percent, message.language);
      return { ok: true };
    default:
      return { ok: false, error_code: "unknown_action" };
  }
});
browser.tabs.onUpdated.addListener((tabId, change) => {
  if (change.status === "complete" && startup?.settings?.zoom_percent) {
    void browser.tabs.get(tabId).then(tab => {
      if (/^https?:\/\//.test(tab.url || "")) return browser.tabs.setZoom(tabId, startup.settings.zoom_percent / 100);
    }).catch(error => console.error("Site Wrapper zoom failed", error));
  }
});
