<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <NsInlineNotification
      v-if="highlightNotice"
      kind="info"
      :title="highlightNotice.title"
      :description="highlightNotice.description"
      :showCloseButton="false"
    />
    <div
      class="logs-output cv-code-snippet-multiline"
      :class="[
        { 'reduced-output-height': numSearches > 1 && !verticalLayout },
        `bx--snippet`,
        `bx--snippet--multi`,
        {
          'bx--snippet--disabled': disabled,
          'bx--snippet--wraptext': wrapText,
          'bx--snippet--light': isLight,
        },
      ]"
      data-code-snippet
    >
      <cv-skeleton-text
        v-if="loading"
        :paragraph="true"
        :line-count="4"
      ></cv-skeleton-text>
      <NsEmptyState
        v-else-if="noLogsFound"
        :title="$t('system_logs.no_log_found')"
        :animationData="GhostDarkBgLottie"
        animationTitle="ghost"
        :loop="1"
        class="margin-auto"
      >
        <template #description>
          <div>{{ $t("system_logs.try_changing_search_filters") }}</div>
        </template>
      </NsEmptyState>
      <div
        v-else
        class="bx--snippet-container"
        :ref="'logsContainer-' + searchId"
        :key="'logsContainer-' + searchId"
      >
        <pre><template v-for="(line, index) in outputLines"
          ><LogHighlightMark
            :key="'line-' + index"
            :text="line"
            :ranges="lineRanges[index]"
          />{{ index === outputLines.length - 1 ? "" : "\n" }}</template
        ></pre>
      </div>
    </div>
  </div>
</template>

<script>
import { carbonPrefixMixin, themeMixin } from "@carbon/vue/src/mixins";
import { UtilService, LottieService } from "@nethserver/ns8-ui-lib";
import LogHighlightMark from "./LogHighlightMark.vue";
import { createHighlightWorker } from "./logHighlightWorker";

// 2000 lines cost 1-30 ms for every realistic pattern, so this only ever trips
// on one that would not have returned at all
const HIGHLIGHT_TIMEOUT_MS = 400;

// render cost, not match cost: the browser parses about 100000 marks per second
const MAX_MARKS = 100000;

// coalesces the burst of appends that follow mode produces
const RECOMPUTE_DEBOUNCE_MS = 120;

function escapeRegExp(str) {
  return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

export default {
  name: "LogOutput",
  props: {
    searchId: {
      type: String,
      required: true,
    },
    outputLines: {
      type: Array,
      required: true,
    },
    scrollToBottom: {
      type: Boolean,
      default: true,
    },
    numSearches: {
      type: Number,
      default: 1,
    },
    highlight: {
      type: [String, RegExp],
      default: "",
    },
    // the query is a valid RE2 pattern the browser engine cannot compile
    highlightUnsupported: Boolean,
    verticalLayout: Boolean,
    loading: Boolean,
    wrapText: Boolean,
    light: Boolean,
    disabled: Boolean,
    noLogsFound: Boolean,
  },
  mixins: [carbonPrefixMixin, themeMixin, UtilService, LottieService],
  components: { LogHighlightMark },
  data() {
    return {
      // aligned with outputLines by index; empty means "render unhighlighted"
      lineRanges: [],
      highlightUnavailable: false,
    };
  },
  computed: {
    highlightNotice() {
      // next to "no log found" the notice reads as a failed search
      if (this.loading || this.noLogsFound || !this.outputLines.length) {
        return null;
      }

      if (this.highlightUnsupported) {
        return {
          title: this.$t("system_logs.highlight_unsupported"),
          description: this.$t("system_logs.highlight_unsupported_description"),
        };
      }

      if (this.highlightUnavailable) {
        return {
          title: this.$t("system_logs.highlight_unavailable"),
          description: this.$t("system_logs.highlight_unavailable_description"),
        };
      }
      return null;
    },
  },
  watch: {
    searchId: function () {
      this.$root.$on(`logsStart-${this.searchId}`, this.logsUpdated);
    },
    outputLines: function () {
      this.scheduleHighlight();
    },
    highlight: function () {
      this.scheduleHighlight();
    },
  },
  created() {
    this.$root.$on(`logsUpdated-${this.searchId}`, this.logsUpdated);
    // a reply for a superseded request must not land on a different set of lines
    this.highlightRequestId = 0;
    this.highlightWorker = null;
    this.highlightTimer = null;
    this.debounceTimer = null;
    // follow mode would otherwise start a doomed worker on every batch of lines
    this.abandonedPattern = "";
  },
  beforeDestroy() {
    this.$root.$off(`logsUpdated-${this.searchId}`);
    clearTimeout(this.debounceTimer);
    this.stopHighlightWorker();
  },
  methods: {
    scheduleHighlight() {
      clearTimeout(this.debounceTimer);
      this.debounceTimer = setTimeout(
        this.computeHighlight,
        RECOMPUTE_DEBOUNCE_MS
      );
    },
    stopHighlightWorker() {
      clearTimeout(this.highlightTimer);
      this.highlightTimer = null;

      if (this.highlightWorker) {
        // terminating is the only way to stop a regex already backtracking
        this.highlightWorker.terminate();
        this.highlightWorker = null;
      }
    },
    // outputLines is trimmed from the front in follow mode, so incremental
    // bookkeeping would misalign: recompute it all, the debounce keeps it cheap
    computeHighlight() {
      this.stopHighlightWorker();
      const requestId = ++this.highlightRequestId;

      if (!this.highlight || !this.outputLines.length) {
        this.lineRanges = [];
        this.highlightUnavailable = false;
        // a new search clears the buffer: the next result set may be smaller
        this.abandonedPattern = "";
        return;
      }
      const pattern =
        this.highlight instanceof RegExp
          ? { source: this.highlight.source, flags: this.highlight.flags }
          : { source: escapeRegExp(this.highlight), flags: "gi" };
      const patternKey = pattern.source + "\u0000" + pattern.flags;

      if (patternKey === this.abandonedPattern) {
        this.giveUpHighlight(patternKey);
        return;
      }
      const worker = createHighlightWorker();

      if (!worker) {
        this.giveUpHighlight(patternKey);
        return;
      }
      this.highlightWorker = worker;

      worker.onmessage = (event) => {
        if (event.data.id !== this.highlightRequestId) {
          return;
        }
        this.stopHighlightWorker();

        if (event.data.error) {
          this.giveUpHighlight(patternKey);
          return;
        }
        // Vue 2 would install a reactive accessor on each of the hundreds of
        // thousands of offsets; freezing makes observe() bail out
        this.lineRanges = Object.freeze(event.data.ranges);
        this.highlightUnavailable = false;
      };
      worker.onerror = () => {
        if (requestId === this.highlightRequestId) {
          this.stopHighlightWorker();
          this.giveUpHighlight(patternKey);
        }
      };
      this.highlightTimer = setTimeout(() => {
        this.stopHighlightWorker();
        this.giveUpHighlight(patternKey);
      }, HIGHLIGHT_TIMEOUT_MS);

      worker.postMessage({
        id: requestId,
        source: pattern.source,
        flags: pattern.flags,
        lines: this.outputLines,
        maxMarks: MAX_MARKS,
      });
    },
    // the lines stay, only the marking goes. A substring is escaped before it is
    // compiled, so it cannot be the expensive case the notice describes.
    giveUpHighlight(patternKey) {
      this.abandonedPattern = patternKey;
      this.lineRanges = [];
      this.highlightUnavailable = this.highlight instanceof RegExp;
    },
    logsUpdated() {
      // new log lines are displayed

      setTimeout(() => {
        if (this.scrollToBottom) {
          this.scrollToBottomOfContainer();
        }
      }, 50);
    },
    scrollToBottomOfContainer() {
      const el = this.$refs[`logsContainer-${this.searchId}`];

      if (el) {
        el.scrollTop = el.scrollHeight;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.logs-output.bx--snippet--multi {
  max-width: none;
  min-height: 4rem;
  max-height: 35rem;
  background-color: #161616 !important;
  color: #f4f4f4 !important;
}

.logs-output.reduced-output-height.bx--snippet--multi {
  max-height: 20rem;
}

.margin-auto {
  margin: auto;
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

.bx--snippet--multi .bx--snippet-container pre {
  overflow-x: visible;
}

// remove fade effect on right border of code snippet
.system-logs .logs-output.bx--snippet--multi .bx--snippet-container pre::after {
  background-image: none;
}

// show scrollbar, and grow so it lands on the edge, not on the longest line
.system-logs .logs-output.bx--snippet--multi .bx--snippet-container {
  overflow-y: auto !important;
  flex: 1 1 auto;
}

.logs-output {
  // the whole line is a <mark>: drop the native look, black-on-black would hide
  // it. The level rules below tie on specificity and come later, so they win.
  mark.log-line {
    background: transparent;
    color: inherit;
  }

  // log level colorization on the dark log viewer background (search-match
  // highlighting keeps its default native <mark> look, still legible on black)
  mark.log-level-error {
    background: transparent;
    color: $inverse-support-01;
    font-weight: 600;
  }

  mark.log-level-warn {
    background: transparent;
    color: $inverse-support-03;
    font-weight: 600;
  }

  mark.log-level-info {
    background: transparent;
    color: $inverse-support-04;
    font-weight: 600;
  }

  mark.log-level-debug {
    background: transparent;
    color: $ui-02;
    font-weight: 600;
  }

  // leading timestamp and "[node:module:syslog_id]" tag, set off from the
  // rest of the line
  .log-timestamp,
  .log-process-tag {
    color: $ui-02;
    font-weight: 700;
  }
}
</style>
