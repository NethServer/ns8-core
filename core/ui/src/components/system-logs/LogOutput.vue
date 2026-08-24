<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <NsInlineNotification
      v-if="highlightUnavailable"
      kind="info"
      :title="$t('system_logs.highlight_unavailable')"
      :description="$t('system_logs.highlight_unavailable_description')"
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

// A pattern RE2 accepts can still be exponential for the JavaScript engine, so
// the worker gets a deadline and is killed past it. Matching 2000 lines costs
// 1-30 ms for every realistic pattern, so this leaves a wide margin and only
// ever trips on a pattern that would not have returned at all.
const HIGHLIGHT_TIMEOUT_MS = 400;

// safety net on render cost rather than match cost: past this the marks alone
// would take longer to build than they are worth
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
    // identifies the newest request, so a reply for a superseded one is dropped
    // instead of being applied to a different set of lines
    this.highlightRequestId = 0;
    this.highlightWorker = null;
    this.highlightTimer = null;
    this.debounceTimer = null;
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
    // outputLines is trimmed from the front in follow mode, so index bookkeeping
    // across updates would silently misalign: recompute the whole buffer, which
    // the debounce above keeps cheap
    computeHighlight() {
      this.stopHighlightWorker();
      const requestId = ++this.highlightRequestId;

      if (!this.highlight || !this.outputLines.length) {
        this.lineRanges = [];
        this.highlightUnavailable = false;
        return;
      }
      const pattern =
        this.highlight instanceof RegExp
          ? { source: this.highlight.source, flags: this.highlight.flags }
          : { source: escapeRegExp(this.highlight), flags: "gi" };

      const worker = createHighlightWorker();

      if (!worker) {
        this.giveUpHighlight();
        return;
      }
      this.highlightWorker = worker;

      worker.onmessage = (event) => {
        if (event.data.id !== this.highlightRequestId) {
          return;
        }
        this.stopHighlightWorker();

        if (event.data.error) {
          this.giveUpHighlight();
          return;
        }
        this.lineRanges = event.data.ranges;
        this.highlightUnavailable = false;
      };
      worker.onerror = () => {
        if (requestId === this.highlightRequestId) {
          this.stopHighlightWorker();
          this.giveUpHighlight();
        }
      };
      this.highlightTimer = setTimeout(() => {
        this.stopHighlightWorker();
        this.giveUpHighlight();
      }, HIGHLIGHT_TIMEOUT_MS);

      worker.postMessage({
        id: requestId,
        source: pattern.source,
        flags: pattern.flags,
        lines: this.outputLines,
        maxMarks: MAX_MARKS,
      });
    },
    // the lines themselves are still the ones the backend matched, so they stay
    // on screen: only the highlighting is dropped
    giveUpHighlight() {
      this.lineRanges = [];
      this.highlightUnavailable = true;
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
  // every line is wrapped in a <mark> (for level/tag detection), so drop the
  // native browser <mark> look on that wrapper: the UA sets color too, and
  // black-on-black would make the whole line invisible. The level rules below
  // have the same specificity and come later, so they still win. Nested
  // mark.log-search-match keeps its default look, still legible on black.
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
