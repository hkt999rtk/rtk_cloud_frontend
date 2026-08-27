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

  function initFeatureTabs() {
    document.querySelectorAll("[data-feature-tabs]").forEach(function (component) {
      const tabs = Array.from(component.querySelectorAll('[role="tab"]'));
      const panels = Array.from(component.querySelectorAll('[role="tabpanel"]'));
      if (!tabs.length || tabs.length !== panels.length) return;

      component.classList.add("is-enhanced");

      function activate(tab, moveFocus) {
        tabs.forEach(function (candidate) {
          const selected = candidate === tab;
          candidate.setAttribute("aria-selected", String(selected));
          candidate.tabIndex = selected ? 0 : -1;
          const panel = component.querySelector("#" + candidate.getAttribute("aria-controls"));
          if (panel) panel.hidden = !selected;
        });
        if (moveFocus) tab.focus();
      }

      tabs.forEach(function (tab, index) {
        tab.addEventListener("click", function () { activate(tab, false); });
        tab.addEventListener("keydown", function (event) {
          let nextIndex = index;
          if (event.key === "ArrowRight" || event.key === "ArrowDown") nextIndex = (index + 1) % tabs.length;
          else if (event.key === "ArrowLeft" || event.key === "ArrowUp") nextIndex = (index - 1 + tabs.length) % tabs.length;
          else if (event.key === "Home") nextIndex = 0;
          else if (event.key === "End") nextIndex = tabs.length - 1;
          else return;
          event.preventDefault();
          activate(tabs[nextIndex], true);
        });
      });

      activate(tabs.find(function (tab) { return tab.getAttribute("aria-selected") === "true"; }) || tabs[0], false);
    });
  }

  function init() {
    initHeader();
    initFeatureTabs();
  }

  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init, { once: true });
  else init();
})();
