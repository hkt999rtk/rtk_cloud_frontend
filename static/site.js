(function () {
  "use strict";

  function initHeader() {
    const header = document.querySelector(".site-header");
    const toggle = header && header.querySelector("[data-menu-toggle]");
    const navigation = header && header.querySelector("#site-navigation");
    if (!header || !toggle || !navigation) return;

    header.classList.add("is-nav-enhanced");

    function setOpen(open, returnFocus) {
      header.classList.toggle("is-menu-open", open);
      toggle.setAttribute("aria-expanded", String(open));
      if (open) {
        const firstLink = navigation.querySelector("a");
        if (firstLink) firstLink.focus();
      } else if (returnFocus) {
        toggle.focus();
      }
    }

    toggle.addEventListener("click", function () {
      setOpen(!header.classList.contains("is-menu-open"), false);
    });

    navigation.addEventListener("click", function (event) {
      if (event.target.closest("a")) setOpen(false, false);
    });

    document.addEventListener("keydown", function (event) {
      if (event.key === "Escape" && header.classList.contains("is-menu-open")) {
        setOpen(false, true);
      }
    });

    document.addEventListener("click", function (event) {
      if (header.classList.contains("is-menu-open") && !header.contains(event.target)) {
        setOpen(false, false);
      }
    });
  }

  function initCoreShowcase() {
    document.querySelectorAll("[data-core-showcase]").forEach(function (component) {
      const tabs = Array.from(component.querySelectorAll('[role="tab"]'));
      const panels = Array.from(component.querySelectorAll('.core-panel'));
      const accordions = Array.from(component.querySelectorAll('[data-core-accordion]'));
      if (!tabs.length || tabs.length !== panels.length) return;
      component.classList.add("is-enhanced");

      function activate(panelID, moveFocus) {
        tabs.forEach(function (candidate) {
          const selected = candidate.getAttribute("aria-controls") === panelID;
          candidate.setAttribute("aria-selected", String(selected));
          candidate.tabIndex = selected ? 0 : -1;
          if (selected && moveFocus) candidate.focus();
        });
        panels.forEach(function (panel) { panel.hidden = panel.id !== panelID; });
        accordions.forEach(function (button) { button.setAttribute("aria-expanded", String(button.getAttribute("data-core-accordion") === panelID)); });
      }

      tabs.forEach(function (tab, index) {
        tab.addEventListener("click", function () { activate(tab.getAttribute("aria-controls"), false); });
        tab.addEventListener("keydown", function (event) {
          let nextIndex = index;
          if (event.key === "ArrowRight" || event.key === "ArrowDown") nextIndex = (index + 1) % tabs.length;
          else if (event.key === "ArrowLeft" || event.key === "ArrowUp") nextIndex = (index - 1 + tabs.length) % tabs.length;
          else if (event.key === "Home") nextIndex = 0;
          else if (event.key === "End") nextIndex = tabs.length - 1;
          else return;
          event.preventDefault();
          activate(tabs[nextIndex].getAttribute("aria-controls"), true);
        });
      });
      accordions.forEach(function (button) { button.addEventListener("click", function () { activate(button.getAttribute("data-core-accordion"), false); }); });
      const selected = tabs.find(function (tab) { return tab.getAttribute("aria-selected") === "true"; }) || tabs[0];
      activate(selected.getAttribute("aria-controls"), false);
    });
  }

  function init() {
    initHeader();
    initCoreShowcase();
  }

  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init, { once: true });
  else init();
})();
