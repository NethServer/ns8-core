<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-tile class="log-search" :light="light">
    <button
      v-if="showCloseButton"
      class="bx--modal-close"
      type="button"
      @click="onClose"
      :aria-label="closeAriaLabel"
    >
      <Close16 class="bx--modal-close__icon" />
    </button>
    <cv-grid
      fullWidth
      :class="['no-padding', { 'mg-top-md': showCloseButton }]"
    >
      <template v-if="filtersShown">
        <cv-row>
          <cv-column :md="verticalLayout ? 8 : 4">
            <label class="bx--label">
              {{ $t("system_logs.context") }}
            </label>
            <cv-content-switcher
              class="mg-bottom-md"
              @selected="onContextSelected"
            >
              <cv-content-switcher-button
                owner-id="cluster"
                :selected="csbClusterSelected"
                >{{
                  $t("system_logs.context_cluster")
                }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="node"
                :selected="csbNodeSelected"
                >{{
                  $t("system_logs.context_node")
                }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="module"
                :selected="csbModuleSelected"
                >{{
                  $t("system_logs.context_module")
                }}</cv-content-switcher-button
              >
            </cv-content-switcher>
          </cv-column>
          <cv-column v-show="csbNodeSelected" :md="verticalLayout ? 8 : 4">
            <cv-combo-box
              v-model="internalSelectedNodeId"
              :label="$t('common.choose')"
              :title="$t('system_logs.context_node')"
              :invalid-message="error.selectedNode"
              :auto-filter="true"
              :auto-highlight="true"
              :options="nodes"
              class="mg-bottom-md"
              key="csbNode"
            >
            </cv-combo-box>
          </cv-column>
          <cv-column v-show="csbModuleSelected" :md="verticalLayout ? 8 : 4">
            <cv-combo-box
              v-model="internalSelectedAppId"
              :label="$t('common.choose')"
              :title="$t('system_logs.context_module')"
              :invalid-message="error.selectedApp"
              :auto-filter="true"
              :auto-highlight="true"
              :options="apps"
              :disabled="loadingApps"
              class="mg-bottom-md"
              key="csbApp"
            >
            </cv-combo-box>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column :md="verticalLayout ? 8 : 4">
            <cv-text-input
              :label="
                $t('system_logs.search_query') +
                ' (' +
                $t('common.optional') +
                ')'
              "
              v-model.trim="internalSearchQuery"
              :placeholder="$t('common.search_placeholder')"
              :helper-text="$t('common.case_sensitive')"
              @keypress.enter="onEnterKeyPress()"
              class="search-query mg-bottom-md"
            >
            </cv-text-input>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column :md="verticalLayout ? 8 : 4" :xlg="verticalLayout ? 8 : 4">
            <label class="bx--label">
              {{ $t("system_logs.mode") }}
            </label>
            <cv-content-switcher
              class="mg-bottom-md"
              @selected="onModeSelected"
            >
              <cv-content-switcher-button
                owner-id="dump"
                :selected="csbDumpModeSelected"
                >{{ $t("system_logs.dump") }}</cv-content-switcher-button
              >
              <cv-content-switcher-button
                owner-id="follow"
                :selected="csbFollowModeSelected"
                >{{ $t("system_logs.follow") }}</cv-content-switcher-button
              >
            </cv-content-switcher>
          </cv-column>
          <cv-column v-if="!internalFollowLogs" :md="verticalLayout ? 4 : 2">
            <cv-text-input
              :label="$t('system_logs.max_lines')"
              v-model="internalMaxLines"
              :invalid-message="error.maxLines"
              type="number"
              min="1"
              :max="MAX_LINES_LIMIT"
              @keypress.enter="onEnterKeyPress()"
              class="narrow mg-bottom-md"
            >
            </cv-text-input>
          </cv-column>
        </cv-row>
        <cv-row v-if="!internalFollowLogs">
          <cv-column :md="verticalLayout ? 8 : 4">
            <cv-date-picker
              kind="single"
              :cal-options="calOptions"
              :date-label="$t('system_logs.start_date')"
              v-model="internalStartDate"
              class="interval-date mg-bottom-md"
              :invalid-message="error.startDate"
            >
            </cv-date-picker>
            <NsTimePicker
              hideClearButton
              v-model="internalStartTime"
              :label="$t('system_logs.start_time_label')"
              class="interval-time mg-bottom-md"
            ></NsTimePicker>
          </cv-column>
          <cv-column :md="verticalLayout ? 8 : 4">
            <cv-date-picker
              kind="single"
              :cal-options="calOptions"
              :date-label="$t('system_logs.end_date')"
              v-model="internalEndDate"
              class="interval-date mg-bottom-md"
              :invalid-message="error.endDate"
            >
            </cv-date-picker>
            <NsTimePicker
              hideClearButton
              v-model="internalEndTime"
              :label="$t('system_logs.end_time_label')"
              class="interval-time mg-bottom-md"
            ></NsTimePicker>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column :md="verticalLayout ? 8 : 4">
            <cv-select
              :label="$t('system_logs.timezone')"
              class="mg-bottom-md"
              v-model="internalTimezone"
            >
              <cv-select-option disabled selected hidden>{{
                $t("system_logs.select_timezone")
              }}</cv-select-option>
              <cv-select-option value="local">
                {{ $t("system_logs.local") }}
              </cv-select-option>
              <cv-select-option value="utc">
                {{ $t("system_logs.utc") }}
              </cv-select-option>
            </cv-select>
          </cv-column>
        </cv-row>
        <cv-row v-if="lokiInstances.length > 1">
          <cv-column :md="verticalLayout ? 8 : 4">
            <NsInlineNotification
              kind="info"
              :title="$t('system_logs.loki_instance_title')"
              :description="$t('system_logs.loki_instance_description')"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row v-if="lokiInstances.length > 1">
          <cv-column :md="verticalLayout ? 8 : 4">
            <NsComboBox
              v-model="internalSelectedLokiId"
              :label="$t('common.choose')"
              :title="$t('system_logs.loki_instance')"
              :invalid-message="error.selectedLoki"
              :auto-filter="true"
              :auto-highlight="true"
              :options="lokiInstances"
              :disabled="loadingLoki"
              class="mg-bottom-md"
              tooltipAlignment="center"
              tooltipDirection="top"
              ><template slot="tooltip">
                {{ $t("system_logs.loki_tooltip") }}
              </template>
            </NsComboBox>
          </cv-column>
        </cv-row>
      </template>
    </cv-grid>
    <cv-grid fullWidth class="no-padding">
      <cv-row v-if="!filtersShown" class="mg-bottom-md">
        <cv-column>
          <div
            :class="[
              'paragraph-line-height',
              { 'pad-right-lg': showCloseButton },
            ]"
          >
            <span class="filter-collapsed">
              <strong>{{ $t("system_logs.context") }}</strong
              >:
              <span v-if="internalContext == 'cluster'">{{
                $t("system_logs.context_cluster")
              }}</span>
              <span v-else-if="internalContext == 'node'">
                {{ selectedNode ? selectedNode.label : "-" }}
              </span>
              <span v-else-if="internalContext == 'module'">
                {{ selectedApp ? selectedApp.label : "-" }}</span
              >
            </span>
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.search_query") }}</strong
              >: "<span>{{ internalSearchQuery }}</span
              >"</span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.start") }}</strong
              >:
              <span>{{
                internalStartDate + " " + internalStartTime
              }}</span></span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.end") }}</strong
              >:
              <span>{{ internalEndDate + " " + internalEndTime }}</span></span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.max_lines") }}</strong
              >: <span>{{ internalMaxLines }}</span></span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.follow_logs") }}</strong
              >:
              <span>{{
                internalFollowLogs
                  ? $t("common.enabled")
                  : $t("common.disabled")
              }}</span></span
            >
            <span class="filter-collapsed"
              ><strong>{{ $t("system_logs.timezone") }}</strong
              >: <span>{{ internalTimezone }}</span></span
            >
          </div>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-link @click="toggleFilters" class="toggle-filters">
            {{
              filtersShown
                ? $t("system_logs.collapse_filters")
                : $t("system_logs.expand_filters")
            }}
          </cv-link>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column class="logs-output-toolbar">
          <NsButton
            v-if="isFollowing"
            kind="primary"
            class="search-button mg-bottom-sm"
            :icon="Close20"
            :loading="loading.stopFollowing"
            :disabled="loading.stopFollowing"
            @click="logsStop"
            key="stop-following"
            >{{ $t("system_logs.stop_following") }}</NsButton
          >
          <NsButton
            v-else
            kind="primary"
            class="search-button mg-bottom-sm"
            :icon="Search20"
            :loading="loading.logs"
            :disabled="
              !isWebsocketConnected ||
              loading.logs ||
              loadingApps ||
              loadingLoki
            "
            @click="logsStart"
            key="search"
            >{{
              internalFollowLogs
                ? $t("system_logs.follow")
                : $t("system_logs.search")
            }}</NsButton
          >
          <template v-if="searchStarted">
            <NsButton
              kind="secondary"
              :icon="Erase20"
              @click="clearLogs"
              class="mg-right mg-bottom-sm"
              >{{ $t("system_logs.clear_logs") }}</NsButton
            >
            <div class="checkbox-filter">
              <cv-checkbox
                :label="$t('system_logs.wrap_text')"
                v-model="wrapText"
                value="checkWrapText"
                class="mg-right mg-bottom-sm"
              />
            </div>
            <div class="checkbox-filter">
              <cv-checkbox
                :label="$t('system_logs.scroll_to_bottom')"
                v-model="scrollToBottom"
                value="checkScrollToBottom"
                class="item mg-bottom-sm"
              />
            </div>
          </template>
        </cv-column>
      </cv-row>
      <cv-row v-if="searchStarted">
        <cv-column>
          <LogOutput
            :searchId="searchId"
            :scrollToBottom="scrollToBottom"
            :wrapText="wrapText"
            :loading="loading.logs"
            :verticalLayout="verticalLayout"
            :numSearches="numSearches"
            :noLogsFound="noLogsFound"
            :light="false"
            :outputLines="outputLines"
            :highlight="highlight"
            :key="'logOutput-' + searchId"
          />
        </cv-column>
      </cv-row>
    </cv-grid>
  </cv-tile>
</template>

<script>
import { mapState } from "vuex";
import LogOutput from "./LogOutput.vue";
import {
  UtilService,
  IconService,
  DateTimeService,
} from "@nethserver/ns8-ui-lib";
import Close16 from "@carbon/icons-vue/es/close/16";

export default {
  name: "LogSearch",
  props: {
    searchId: {
      type: String,
      required: true,
    },
    nodes: {
      type: Array,
      required: true,
    },
    apps: {
      type: Array,
      required: true,
    },
    lokiInstances: {
      type: Array,
      required: true,
    },
    numSearches: {
      type: Number,
      default: 1,
    },
    searchQuery: {
      type: String,
      default: "",
    },
    context: {
      type: String,
      default: "cluster",
      validator: (val) => ["cluster", "node", "module"].includes(val),
    },
    selectedNodeId: {
      type: String,
      default: "",
    },
    selectedAppId: {
      type: String,
      default: "",
    },
    selectedLokiId: {
      type: String,
      default: "",
    },
    maxLines: {
      type: String,
      default: "500",
    },
    startDate: {
      type: String,
      default: "",
    },
    startTime: {
      type: String,
      default: "",
    },
    endDate: {
      type: String,
      default: "",
    },
    endTime: {
      type: String,
      default: "",
    },
    timezone: {
      type: String,
      default: "",
    },
    startSearchCommand: {
      type: Number,
      default: 0,
    },
    mainSearch: Boolean,
    followLogs: Boolean,
    verticalLayout: Boolean,
    loadingApps: Boolean,
    loadingLoki: Boolean,
    closeAriaLabel: { type: String, default: "Close modal" },
    light: Boolean,
  },
  components: { LogOutput, Close16 },
  mixins: [UtilService, IconService, DateTimeService],
  data() {
    return {
      MAX_LINES_LIMIT: 2000,
      filtersShown: true,
      internalContext: "",
      internalSearchQuery: "",
      calOptions: {
        dateFormat: "Y-m-d",
      },
      internalStartDate: "",
      internalEndDate: "",
      internalStartTime: "",
      internalEndTime: "",
      internalMaxLines: "",
      internalTimezone: "",
      wrapText: true,
      internalFollowLogs: false,
      outputLines: [],
      isFollowing: false,
      internalSelectedNodeId: "",
      internalSelectedAppId: "",
      internalSelectedLokiId: "",
      scrollToBottom: true,
      searchStarted: false,
      noLogsFound: false,
      highlight: "",
      loading: {
        logs: false,
        stopFollowing: false,
      },
      error: {
        startDate: "",
        startTime: "",
        endDate: "",
        endTime: "",
        maxLines: "",
        selectedNode: "",
        selectedApp: "",
        selectedLoki: "",
      },
    };
  },
  computed: {
    ...mapState(["isWebsocketConnected"]),
    csbClusterSelected() {
      return this.internalContext === "cluster";
    },
    csbNodeSelected() {
      return this.internalContext === "node";
    },
    csbModuleSelected() {
      return this.internalContext === "module";
    },
    selectedApp() {
      return this.apps.find((app) => app.value == this.internalSelectedAppId);
    },
    selectedLoki() {
      return this.lokiInstances.find(
        (instance) => instance.value == this.internalSelectedLokiId
      );
    },
    selectedNode() {
      return this.nodes.find(
        (node) => node.value == this.internalSelectedNodeId
      );
    },
    showCloseButton() {
      return this.numSearches > 1;
    },
    csbDumpModeSelected() {
      return !this.internalFollowLogs;
    },
    csbFollowModeSelected() {
      return this.internalFollowLogs;
    },
  },
  watch: {
    searchQuery: function () {
      if (this.mainSearch) {
        this.internalSearchQuery = this.searchQuery;
      }
    },
    internalSearchQuery: function () {
      if (this.mainSearch) {
        this.$emit("updateSearchQuery", this.internalSearchQuery);
      }
    },
    timezone: function () {
      if (this.mainSearch) {
        this.internalTimezone = this.timezone;
      }
    },
    internalTimezone: function () {
      if (this.mainSearch) {
        this.$emit("updateTimezone", this.internalTimezone);
      }
    },
    context: function () {
      if (this.mainSearch) {
        this.internalContext = this.context;
      }
    },
    internalContext: function () {
      if (this.mainSearch) {
        this.$emit("updateContext", this.internalContext);
      }
    },
    selectedNodeId: function () {
      if (this.mainSearch && this.nodes.length) {
        this.$nextTick(() => {
          this.internalSelectedNodeId = this.selectedNodeId;
        });
      }
    },
    internalSelectedNodeId: function () {
      if (this.mainSearch) {
        this.$emit("updateSelectedNodeId", this.internalSelectedNodeId);
      }
    },
    nodes: function () {
      if (this.mainSearch && this.nodes.length) {
        this.$nextTick(() => {
          this.internalSelectedNodeId = this.selectedNodeId;
        });
      }
    },
    selectedAppId: function () {
      if (this.mainSearch && this.apps.length) {
        this.$nextTick(() => {
          this.internalSelectedAppId = this.selectedAppId;
        });
      }
    },
    internalSelectedAppId: function () {
      if (this.mainSearch) {
        this.$emit("updateSelectedAppId", this.internalSelectedAppId);
      }
    },
    apps: function () {
      if (this.mainSearch && this.apps.length) {
        this.$nextTick(() => {
          this.internalSelectedAppId = this.selectedAppId;
        });
      }
    },
    lokiInstances: function () {
      if (this.mainSearch && this.lokiInstances.length) {
        this.$nextTick(() => {
          this.internalSelectedLokiId = this.selectedLokiId;
        });
      }
    },
    selectedLokiId: function () {
      if (this.mainSearch && this.lokiInstances.length) {
        this.$nextTick(() => {
          this.internalSelectedLokiId = this.selectedLokiId;
        });
      }
    },
    internalSelectedLokiId: function () {
      if (this.mainSearch && this.lokiInstances.length) {
        this.$emit("updateSelectedLokiId", this.internalSelectedLokiId);
      }
    },
    followLogs: function () {
      if (this.mainSearch) {
        this.internalFollowLogs = this.followLogs;
      }
    },
    internalFollowLogs: function () {
      if (this.mainSearch) {
        this.$emit("updateFollowLogs", this.internalFollowLogs);
      }
    },
    maxLines: function () {
      if (this.mainSearch) {
        this.internalMaxLines = this.maxLines;
      }
    },
    internalMaxLines: function () {
      if (this.mainSearch) {
        this.$emit("updateMaxLines", this.internalMaxLines);
      }
    },
    startDate: function () {
      if (this.mainSearch) {
        this.internalStartDate = this.startDate;
      }
    },
    internalStartDate: function () {
      if (this.mainSearch) {
        this.$emit("updateStartDate", this.internalStartDate);
      }
    },
    endDate: function () {
      if (this.mainSearch) {
        this.internalEndDate = this.endDate;
      }
    },
    internalEndDate: function () {
      if (this.mainSearch) {
        this.$emit("updateEndDate", this.internalEndDate);
      }
    },
    startTime: function () {
      if (this.mainSearch) {
        this.internalStartTime = this.startTime;
      }
    },
    internalStartTime: function () {
      if (this.mainSearch) {
        this.$emit("updateStartTime", this.internalStartTime);
      }
    },
    endTime: function () {
      if (this.mainSearch) {
        this.internalEndTime = this.endTime;
      }
    },
    internalEndTime: function () {
      if (this.mainSearch) {
        this.$emit("updateEndTime", this.internalEndTime);
      }
    },
    startSearchCommand: function () {
      if (this.mainSearch && this.startSearchCommand > 0) {
        this.logsStart();
      }
    },
  },
  created() {
    this.initFilters();

    // register event listeners
    this.$root.$on(
      `collapseSystemLogsFilters-${this.searchId}`,
      this.collapseFilters
    );
    this.$root.$on("logSearchClosed", this.onLogSearchClosed);
  },
  beforeDestroy() {
    // remove event listeners
    this.$root.$off(`logsStart-${this.searchId}`);
    this.$root.$off(`logsStop-${this.searchId}`);
    this.$root.$off(`collapseSystemLogsFilters-${this.searchId}`);
    this.$root.$off("logSearchClosed");

    if (this.pid) {
      this.logsStop();
    }
  },
  methods: {
    toggleFilters() {
      this.filtersShown = !this.filtersShown;
    },
    collapseFilters() {
      this.filtersShown = false;
    },
    onContextSelected(value) {
      this.internalContext = value;
    },
    initFilters() {
      this.clearErrors(this);
      // set time interval
      this.internalStartDate = this.formatDate(
        this.subDays(new Date(), 1),
        "yyyy-MM-dd"
      );
      this.internalEndDate = this.formatDate(new Date(), "yyyy-MM-dd");
      this.internalStartTime = "00:00";
      this.internalEndTime = "23:59";
      this.internalContext = "cluster";
      this.internalSelectedAppId = "";
      this.internalSelectedLokiId = this.selectedLokiId;
      this.internalSelectedNodeId = "";
      this.internalSearchQuery = "";
      this.internalTimezone = "local";
      this.internalMaxLines = "500";
      this.internalFollowLogs = false;
    },
    validateSearchLogs() {
      this.clearErrors();
      let isValidationOk = true;

      if (this.internalContext == "node" && !this.internalSelectedNodeId) {
        this.error.selectedNode = this.$t("common.required");
        isValidationOk = false;
      }

      if (this.internalContext == "module" && !this.internalSelectedAppId) {
        this.error.selectedApp = this.$t("common.required");
        isValidationOk = false;
      }
      if (!this.internalSelectedLokiId) {
        this.error.selectedLoki = this.$t("common.required");
        isValidationOk = false;
      }

      if (!this.internalFollowLogs) {
        if (!this.internalStartTime) {
          this.error.startTime = this.$t("common.required");
          isValidationOk = false;
        }
        if (!this.internalEndTime) {
          this.error.endTime = this.$t("common.required");
          isValidationOk = false;
        }
        if (this.internalStartTime && this.internalEndTime) {
          // using Date constructor: new Date('1995-12-17T03:24:00')
          const startLocal = new Date(
            this.internalStartDate + "T" + this.internalStartTime + ":00"
          );
          const endLocal = new Date(
            this.internalEndDate + "T" + this.internalEndTime + ":59"
          );
          if (this.dateIsBefore(endLocal, startLocal)) {
            this.error.startDate = this.$t("error.invalid_time_interval");
            this.error.startTime = this.$t("error.invalid_time_interval");
            this.error.endDate = this.$t("error.invalid_time_interval");
            this.error.endTime = this.$t("error.invalid_time_interval");
            isValidationOk = false;
          }
        }

        if (
          this.internalMaxLines < 1 ||
          this.internalMaxLines > this.MAX_LINES_LIMIT
        ) {
          this.error.maxLines = this.$t("error.must_be_between", {
            min: 1,
            max: this.MAX_LINES_LIMIT,
          });
          isValidationOk = false;
        }
      }

      return isValidationOk;
    },
    logsStart() {
      if (!this.validateSearchLogs()) {
        return;
      }
      this.loading.logs = true;
      this.noLogsFound = false;
      this.searchStarted = true;
      const mode = this.internalFollowLogs ? "tail" : "dump";

      // using Date constructor: new Date('1995-12-17T03:24:00')
      const startLocal = new Date(
        this.internalStartDate + "T" + this.internalStartTime + ":00"
      );
      const endLocal = new Date(
        this.internalEndDate + "T" + this.internalEndTime + ":59"
      );
      const format = "yyyy-MM-dd'T'HH:mm:ssX";
      const startUtcString = this.formatInTimeZone(startLocal, format, "UTC");
      const endUtcString = this.formatInTimeZone(endLocal, format, "UTC");
      this.highlight = this.internalSearchQuery;
      this.clearLogs();
      let entityName;
      let timezone = "UTC";

      switch (this.internalContext) {
        case "cluster":
          entityName = "";
          break;
        case "node":
          entityName = this.internalSelectedNodeId;
          break;
        case "module":
          entityName = this.internalSelectedAppId;
          break;
      }

      if (this.internalFollowLogs) {
        this.$root.$on(`logsStart-${this.searchId}`, this.onLogsStartFollow);
      } else {
        this.$root.$once(`logsStart-${this.searchId}`, this.onLogsStartDump);
      }

      if (this.internalTimezone == "local") {
        timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
      }

      const logsStartObj = {
        action: "logs-start",
        payload: {
          mode: mode,
          lines: this.internalMaxLines,
          filter: this.internalSearchQuery,
          from: startUtcString,
          to: endUtcString,
          entity: this.internalContext,
          entity_name: String(entityName),
          id: this.searchId,
          timezone: timezone,
          instance: this.internalSelectedLokiId,
        },
      };

      this.$socket.sendObj(logsStartObj);
    },
    logsStop() {
      this.loading.stopFollowing = true;
      this.$root.$off(`logsStart-${this.searchId}`);
      this.$root.$once(`logsStop-${this.searchId}`, this.onLogsStop);

      const logsStopObj = {
        action: "logs-stop",
        payload: {
          pid: this.pid,
          id: this.searchId,
        },
      };
      this.$socket.sendObj(logsStopObj);
    },
    onLogsStop() {
      this.loading.stopFollowing = false;
      this.isFollowing = false;
      this.pid = "";
    },
    onLogsStartDump(payload) {
      this.loading.logs = false;

      if (!payload.message.length) {
        // show empty state in LogsOutput
        this.noLogsFound = true;
      } else {
        this.outputLines = [...this.outputLines, ...payload.message];
      }

      // signal LogOutput
      this.$root.$emit(`logsUpdated-${this.searchId}`);
    },
    onLogsStartFollow(payload) {
      if (this.loading.logs) {
        this.loading.logs = false;
        this.isFollowing = true;
        this.pid = payload.pid;
      }

      if (payload.message.length) {
        this.outputLines = [...this.outputLines, ...payload.message];
      }

      // limit output buffer
      if (this.outputLines.length > this.MAX_LINES_LIMIT) {
        const numLinesToDelete = this.outputLines.length - this.MAX_LINES_LIMIT;
        this.outputLines.splice(0, numLinesToDelete);
      }

      // signal LogOutput
      this.$root.$emit(`logsUpdated-${this.searchId}`);
    },
    clearLogs() {
      this.outputLines = [];
    },
    onEnterKeyPress() {
      if (!this.isFollowing) {
        this.logsStart();
      }
    },
    onClose() {
      this.$emit("close", this.searchId);
      this.$root.$emit("logSearchClosed");
    },
    onLogSearchClosed() {
      this.$nextTick(() => {
        if (this.mainSearch) {
          // update query params
          this.$emit("updateSearchQuery", this.internalSearchQuery);
          this.$emit("updateContext", this.internalContext);
          this.$emit("updateSelectedNodeId", this.internalSelectedNodeId);
          this.$emit("updateSelectedAppId", this.internalSelectedAppId);
          this.$emit("updateSelectedLokiId", this.internalSelectedLokiId);
          this.$emit("updateFollowLogs", this.internalFollowLogs);
          this.$emit("updateMaxLines", this.internalMaxLines);
          this.$emit("updateStartDate", this.internalStartDate);
          this.$emit("updateEndDate", this.internalEndDate);
          this.$emit("updateStartTime", this.internalStartTime);
          this.$emit("updateEndTime", this.internalEndTime);
          this.$emit("updateTimezone", this.internalTimezone);
        }
      });
    },
    onModeSelected(value) {
      if (value == "dump") {
        this.internalFollowLogs = false;
      } else if (value == "follow") {
        this.internalFollowLogs = true;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.log-search {
  // needed for close button positon
  position: relative;
}

.interval-date {
  display: inline-flex;
  margin-right: $spacing-06;
}

.interval-time {
  display: inline-flex;
}

.checkbox-filter {
  display: flex;
  align-items: center;
}

.toggle-filters {
  margin-bottom: $spacing-06;
}

.filter-collapsed {
  margin-right: $spacing-05;
}

.pad-right-lg {
  padding-right: $spacing-07;
}

.logs-output-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;

  .mg-right {
    margin-right: $spacing-06;
  }
}
</style>

<style lang="scss">
@import "../../styles/carbon-utils";

// global styles

@media (min-width: $breakpoint-medium) {
  .system-logs .search-query .bx--text-input,
  .system-logs .search-query .bx--text-input__field-wrapper {
    max-width: none;
  }
}
</style>
