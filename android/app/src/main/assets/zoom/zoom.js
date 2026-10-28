"use strict";
browser.runtime.sendNativeMessage("site-wrapper", { action: "zoom" }).then(response => {
  if (Number.isInteger(response?.percent) && response.percent >= 30 && response.percent <= 500) {
    const apply = () => { document.documentElement.style.zoom = `${response.percent}%`; };
    if (document.documentElement) apply();
    else document.addEventListener("DOMContentLoaded", apply, { once: true });
  }
}).catch(error => console.error("Site Wrapper zoom failed", error));
