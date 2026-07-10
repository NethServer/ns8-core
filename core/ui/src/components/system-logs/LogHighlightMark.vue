<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <mark :class="levelClass"
    ><span v-if="beforeTag" class="log-timestamp"
      ><template v-for="(part, i) in beforeParts"
        ><mark v-if="part.isMatch" :key="'b' + i" class="log-search-match">{{
          part.text
        }}</mark
        ><template v-else>{{ part.text }}</template></template
      ></span
    ><span v-if="tagText" class="log-process-tag"
      ><template v-for="(part, i) in tagParts"
        ><mark v-if="part.isMatch" :key="'t' + i" class="log-search-match">{{
          part.text
        }}</mark
        ><template v-else>{{ part.text }}</template></template
      ></span
    ><template v-for="(part, i) in restParts"
      ><mark v-if="part.isMatch" :key="'r' + i" class="log-search-match">{{
        part.text
      }}</mark
      ><template v-else>{{ part.text }}</template></template
    ></mark
  >
</template>

<script>
// "contains" check: a chunk can be a single keyword or a whole matching line
const LOG_LEVEL_CLASSIFIERS = [
  {
    class: "log-level-error",
    pattern: /\b(?:error|err|fatal|crit|critical)\b/i,
  },
  { class: "log-level-warn", pattern: /\b(?:warn|warning)\b/i },
  { class: "log-level-info", pattern: /\b(?:info|information|notice)\b/i },
  { class: "log-level-debug", pattern: /\b(?:debug|trace)\b/i },
];

// leading timestamp followed by the "[node_id:module_id:syslog_id]" tag
// added by the backend's LogQL line_format, e.g.
// "2026-07-05T12:03:17+02:00 [1:traefik1:traefik] ..."
const PROCESS_TAG_PATTERN = /^(\S+\s+)(\[[^\]]*\])/;

function escapeRegExp(str) {
  return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

export default {
  name: "LogHighlightMark",
  props: {
    text: {
      type: String,
      required: true,
    },
    searchTerm: {
      type: String,
      default: "",
    },
  },
  computed: {
    levelClass() {
      const level = LOG_LEVEL_CLASSIFIERS.find((l) =>
        l.pattern.test(this.text)
      );
      return level ? level.class : "";
    },
    tagMatch() {
      return PROCESS_TAG_PATTERN.exec(this.text);
    },
    beforeTag() {
      return this.tagMatch ? this.tagMatch[1] : "";
    },
    tagText() {
      return this.tagMatch ? this.tagMatch[2] : "";
    },
    afterTag() {
      return this.tagMatch
        ? this.text.slice(this.tagMatch[0].length)
        : this.text;
    },
    beforeParts() {
      return this.splitBySearchTerm(this.beforeTag);
    },
    tagParts() {
      return this.splitBySearchTerm(this.tagText);
    },
    restParts() {
      return this.splitBySearchTerm(this.afterTag);
    },
  },
  methods: {
    // splits a piece of text so the search term keeps its own clickable
    // mark, even when the whole containing line is also colorized by level
    splitBySearchTerm(text) {
      if (!this.searchTerm) {
        return [{ text, isMatch: false }];
      }
      const re = new RegExp(`(${escapeRegExp(this.searchTerm)})`, "gi");
      return text
        .split(re)
        .filter((part) => part !== "")
        .map((part) => ({
          text: part,
          isMatch: part.toLowerCase() === this.searchTerm.toLowerCase(),
        }));
    },
  },
};
</script>
