function t(t){return/^(https?:|http?:|mailto:|tel:)/.test(t)}function r(t){if(!t&&"object"!=typeof t)throw new Error("error arguments","deepClone");const e=t.constructor===Array?[]:{};return Object.keys(t).forEach((o=>{t[o]&&"object"==typeof t[o]?e[o]=r(t[o]):e[o]=t[o]})),e}export{r as d,t as i};
//# sourceMappingURL=index.B6fdKp2E.js.map
