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

  function initManualNavigation() {
    const article = document.querySelector("[data-manual-article]");
    const toc = document.querySelector("[data-manual-toc]");
    if (!article || !toc) return;
    const headings = Array.from(article.querySelectorAll("h2"));
    if (headings.length < 2) return;
    const list = toc.querySelector("ol");
    headings.forEach(function (heading, index) {
      if (!heading.id) {
        let id = "manual-section-" + (index + 1);
        while (document.getElementById(id)) id += "-section";
        heading.id = id;
      }
      const item = document.createElement("li");
      const link = document.createElement("a");
      link.href = "#" + encodeURIComponent(heading.id);
      link.textContent = heading.textContent;
      link.addEventListener("click", function () {
        heading.tabIndex = -1;
        heading.focus({ preventScroll: true });
      });
      item.appendChild(link);
      list.appendChild(item);
    });
    toc.hidden = false;
  }

  function init() {
    initHeader();
    initCoreShowcase();
    initManualNavigation();
  }

  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init, { once: true });
  else init();
})();
