<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <mark class="log-line" :class="levelClass"
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
      return this.splitByRanges(this.beforeTag, 0);
    },
    tagParts() {
      return this.splitByRanges(this.tagText, 1);
    },
    restParts() {
      return this.splitByRanges(this.afterTag, 2);
    },
  },
  methods: {
    // splits a segment so the search term keeps its own clickable mark, even
    // when the whole containing line is also colorized by level
    splitByRanges(text, segment) {
      const flat = this.ranges ? this.ranges[segment] : null;

      if (!flat || !flat.length) {
        return [{ text, isMatch: false }];
      }
      const parts = [];
      let cursor = 0;

      for (let i = 0; i < flat.length; i += 2) {
        if (flat[i] > cursor) {
          parts.push({ text: text.slice(cursor, flat[i]), isMatch: false });
        }
        parts.push({ text: text.slice(flat[i], flat[i + 1]), isMatch: true });
        cursor = flat[i + 1];
      }
      if (cursor < text.length) {
        parts.push({ text: text.slice(cursor), isMatch: false });
      }
      return parts.length ? parts : [{ text, isMatch: false }];
    },
  },
};
</script>
