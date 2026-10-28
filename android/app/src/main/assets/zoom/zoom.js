"use strict";
browser.runtime.sendNativeMessage("sitewrap", { action: "zoom" }).then(response => {
  if (Number.isInteger(response?.percent) && response.percent >= 30 && response.percent <= 500) {
    document.documentElement.style.zoom = `${response.percent}%`;
  }
}).catch(error => console.error("Sitewrap zoom failed", error));
