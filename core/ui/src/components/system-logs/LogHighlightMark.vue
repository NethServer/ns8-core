<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <mark class="log-line" :class="levelClass" v-html="html"></mark>
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

// The log text is remote input, and it goes through v-html below. This is a
// text context, no attribute is ever built from it, so & and < are the two
// characters that can end it -- but they must be escaped without exception.
function escapeHtml(text) {
  return text.replace(/&/g, "&amp;").replace(/</g, "&lt;");
}

export default {
  name: "LogHighlightMark",
  props: {
    text: {
      type: String,
      required: true,
    },
    // match offsets computed off the main thread, one flat [start, end, ...]
    // list per segment: [timestamp, tag, message]. Null until they arrive, or
    // for good when the pattern turned out to be too expensive to run.
    ranges: {
      type: Array,
      default: null,
    },
  },
  computed: {
    levelClass() {
      const level = LOG_LEVEL_CLASSIFIERS.find((l) =>
        l.pattern.test(this.text)
      );
      return level ? level.class : "";
    },
    // One string per line instead of a vnode per fragment: a search matching
    // every other character produces tens of thousands of fragments, and the
    // template compiler spends ten times longer diffing them than the browser
    // spends parsing the equivalent markup.
    html() {
      const tag = PROCESS_TAG_PATTERN.exec(this.text);

      if (!tag) {
        return this.segmentHtml(this.text, 2);
      }
      return (
        '<span class="log-timestamp">' +
        this.segmentHtml(tag[1], 0) +
        '</span><span class="log-process-tag">' +
        this.segmentHtml(tag[2], 1) +
        "</span>" +
        this.segmentHtml(this.text.slice(tag[0].length), 2)
      );
    },
  },
  methods: {
    // wraps the matched offsets of one segment, so the search term keeps its
    // own mark even when the whole containing line is colorized by level
    segmentHtml(text, segment) {
      const flat = this.ranges ? this.ranges[segment] : null;

      if (!flat || !flat.length) {
        return escapeHtml(text);
      }
      const parts = [];
      let cursor = 0;

      for (let i = 0; i < flat.length; i += 2) {
        if (flat[i] > cursor) {
          parts.push(escapeHtml(text.slice(cursor, flat[i])));
        }
        parts.push(
          '<mark class="log-search-match">',
          escapeHtml(text.slice(flat[i], flat[i + 1])),
          "</mark>"
        );
        cursor = flat[i + 1];
      }
      if (cursor < text.length) {
        parts.push(escapeHtml(text.slice(cursor)));
      }
      return parts.join("");
    },
  },
};
</script>
