<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <cv-grid fullWidth>
    <cv-row>
      <cv-column class="page-title">
        <h2>{{ $t("audit_trail.title") }}</h2>
      </cv-column>
    </cv-row>
    <cv-row class="landscape-warning">
      <cv-column>
        <NsInlineNotification
          kind="warning"
          :title="$t('common.use_landscape_mode')"
          :description="$t('common.use_landscape_mode_description')"
        />
      </cv-column>
    </cv-row>
    <cv-row>
      <cv-column>
        <cv-tile :light="true">
          <cv-grid fullWidth class="no-padding">
            <cv-row>
              <cv-column :lg="8">
                <cv-date-picker
                  kind="single"
                  :cal-options="calOptions"
                  :date-label="$t('audit_trail.start_date')"
                  v-model="startDateString"
                  class="interval-date mg-bottom-md"
                  :invalid-message="error.startDate"
                >
                </cv-date-picker>
                <NsTimePicker
                  hideClearButton
                  v-model="startTime"
                  :label="$t('audit_trail.start_time_label')"
                  class="interval-time mg-bottom-md"
                ></NsTimePicker>
              </cv-column>
              <cv-column :lg="8">
                <cv-date-picker
                  kind="single"
                  :cal-options="calOptions"
                  :date-label="$t('audit_trail.end_date')"
                  v-model="endDateString"
                  class="interval-date mg-bottom-md"
                  :invalid-message="error.endDate"
                >
                </cv-date-picker>
                <NsTimePicker
                  hideClearButton
                  v-model="endTime"
                  :label="$t('audit_trail.end_time_label')"
                  class="interval-time mg-bottom-md"
                ></NsTimePicker>
              </cv-column>
            </cv-row>
            <cv-row>
              <cv-column :lg="8">
                <cv-multi-select
                  v-model="usersSelected"
                  :options="users"
                  :title="
                    $t('audit_trail.users') + ' (' + $t('common.optional') + ')'
                  "
                  :label="$t('audit_trail.users_label')"
                  :helper-text="usersHelperText"
                  class="mg-bottom-md"
                >
                </cv-multi-select>
                <NsInlineNotification
                  v-if="error.auditUsers"
                  kind="error"
                  :title="$t('audit_trail.cannot_retrieve_audit_users')"
                  :description="error.auditUsers"
                  :showCloseButton="false"
                />
              </cv-column>
              <cv-column :lg="8">
                <cv-multi-select
                  v-model="actionsSelected"
                  :options="actions"
                  :title="
                    $t('audit_trail.actions') +
                    ' (' +
                    $t('common.optional') +
                    ')'
                  "
                  :label="$t('audit_trail.actions_label')"
                  :helper-text="actionsHelperText"
                  class="mg-bottom-md"
                >
                </cv-multi-select>
                <NsInlineNotification
                  v-if="error.auditActions"
                  kind="error"
                  :title="$t('audit_trail.cannot_retrieve_audit_actions')"
                  :description="error.auditActions"
                  :showCloseButton="false"
                />
              </cv-column>
            </cv-row>
            <cv-row>
              <cv-column :lg="8">
                <cv-text-input
                  :label="
                    $t('audit_trail.auditInfo') +
                    ' (' +
                    $t('common.optional') +
                    ')'
                  "
                  v-model="auditInfo"
                  :placeholder="$t('audit_trail.audit_info_placeholder')"
                  @keypress.enter="searchLogs()"
                  class="mg-bottom-md"
                >
                </cv-text-input>
              </cv-column>
              <cv-column :lg="8">
                <cv-text-input
                  :label="$t('audit_trail.max_results')"
                  v-model="maxResults"
                  :invalid-message="error.maxResults"
                  type="number"
                  min="1"
                  :max="MAX_RESULTS_LIMIT"
                  @keypress.enter="searchLogs()"
                  class="narrow mg-bottom-md"
                  ref="maxResults"
                >
                </cv-text-input>
              </cv-column>
            </cv-row>
            <cv-row class="search-button">
              <cv-column>
                <NsButton
                  kind="primary"
                  :icon="Search20"
                  :loading="loading.auditTrail"
                  :disabled="loading.auditTrail"
                  @click="searchLogs"
                  >{{ $t("audit_trail.search") }}</NsButton
                >
              </cv-column>
            </cv-row>
          </cv-grid>
          <cv-grid fullWidth class="no-padding">
            <cv-row>
              <cv-column>
                <NsDataTable
                  :allRows="logs"
                  :columns="i18nTableColumns"
                  :rawColumns="tableColumns"
                  :sortable="true"
                  :isSearchable="false"
                  :isLoading="loading.auditTrail"
                  :isErrorShown="!!error.auditTrail"
                  :errorTitle="$t('audit_trail.cannot_retrieve_audit_trail')"
                  :errorDescription="error.auditTrail"
                  :itemsPerPageLabel="$t('pagination.items_per_page')"
                  :rangeOfTotalItemsLabel="
                    $t('pagination.range_of_total_items')
                  "
                  :ofTotalPagesLabel="$t('pagination.of_total_pages')"
                  :backwardText="$t('pagination.previous_page')"
                  :forwardText="$t('pagination.next_page')"
                  :pageNumberLabel="$t('pagination.page_number')"
                  @updatePage="tablePage = $event"
                  ref="table"
                >
                  <template slot="empty-state">
                    <NsEmptyState
                      :title="$t('audit_trail.no_log_found')"
                      :animationData="GhostLottie"
                      animationTitle="ghost"
                      :loop="1"
                    >
                      <template #description>
                        <div>
                          {{ $t("audit_trail.no_log_found_description") }}
                        </div>
                      </template>
                    </NsEmptyState>
                  </template>
                  <template slot="data">
                    <cv-data-table-row
                      v-for="(row, rowIndex) in tablePage"
                      :key="`${rowIndex}`"
                      :value="`${rowIndex}`"
                    >
                      <cv-data-table-cell>{{
                        row.timestamp
                      }}</cv-data-table-cell>
                      <cv-data-table-cell>{{ row.user }}</cv-data-table-cell>
                      <cv-data-table-cell>{{ row.action }}</cv-data-table-cell>
                      <cv-data-table-cell
                        :class="{
                          'audit-info-collapsed': row.auditInfoCollapsed,
                        }"
                      >
                        <cv-icon-button
                          v-if="row.auditInfo"
                          kind="ghost"
                          size="sm"
                          :icon="
                            row.auditInfoCollapsed ? RowExpand20 : RowCollapse20
                          "
                          :label="
                            row.auditInfoCollapsed
                              ? $t('common.expand')
                              : $t('common.collapse')
                          "
                          @click="
                            row.auditInfoCollapsed = !row.auditInfoCollapsed
                          "
                          tip-position="left"
                          class="expand-audit-info-button"
                        />
                        <code>
                          <pre
                            v-if="!row.auditInfoCollapsed"
                            class="audit-info-pre"
                            >{{ tryParseJson(row.auditInfo) }}</pre
                          >
                          <span v-else :id="'audit-info-' + row.id">{{
                            tryParseJson(row.auditInfo)
                          }}</span>
                        </code>
                      </cv-data-table-cell>
                    </cv-data-table-row>
                  </template>
                </NsDataTable>
              </cv-column>
            </cv-row>
          </cv-grid>
        </cv-tile>
      </cv-column>
    </cv-row>
  </cv-grid>
</template>

<script>
import AuditService from "../mixins/audit";
import to from "await-to-js";
import {
  UtilService,
  IconService,
  DateTimeService,
  LottieService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "AuditTrail",
  mixins: [
    DateTimeService,
    IconService,
    AuditService,
    UtilService,
    LottieService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("audit_trail.title");
  },
  data() {
    return {
      MAX_RESULTS_LIMIT: 2000,
      maxResults: 0,
      startDateString: "",
      endDateString: "",
      startTime: "",
      endTime: "",
      calOptions: {
        dateFormat: "Y-m-d",
      },
      logs: [],
      users: [],
      usersSelected: [],
      actions: [],
      actionsSelected: [],
      auditInfo: "",
      tableColumns: ["timestamp", "user", "action", "auditInfo"],
      tablePage: [],
      loading: {
        auditTrail: true,
      },
      error: {
        auditUsers: "",
        auditActions: "",
        auditTrail: "",
        startDate: "",
        startTime: "",
        endDate: "",
        endTime: "",
        maxResults: "",
      },
    };
  },
  computed: {
    usersHelperText() {
      if (!this.usersSelected.length) {
        return "";
      } else {
        return (
          this.$t("common.selected") + ": " + this.usersSelected.join(", ")
        );
      }
    },
    actionsHelperText() {
      if (!this.actionsSelected.length) {
        return "";
      } else {
        return (
          this.$t("common.selected") + ": " + this.actionsSelected.join(", ")
        );
      }
    },
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("audit_trail." + column);
      });
    },
  },
  created() {
    this.initFilters();

    // call audit api
    this.getAuditData();
  },
  methods: {
    initFilters() {
      this.clearErrors(this);

      // set time interval
      this.startDateString = this.formatDate(
        this.subDays(new Date(), 1),
        "yyyy-MM-dd"
      );
      this.endDateString = this.formatDate(new Date(), "yyyy-MM-dd");

      this.startTime = "00:00";
      this.endTime = "23:59";
      this.usersSelected = [];
      this.actionsSelected = [];
      this.auditInfo = "";
      this.maxResults = "100";
    },
    async getAuditData() {
      // get audit users
      const [errUsers, responseUsers] = await to(this.getAuditUsers());

      if (errUsers) {
        console.error("error retrieving audit users", errUsers);
        this.error.auditUsers = this.getErrorMessage(errUsers);
        this.loading.auditTrail = false;
        return;
      }

      const users = responseUsers.data.data.users;

      this.users = users.map((u) => {
        return { label: u, value: u, name: u };
      });

      // get audit actions
      const [errActions, responseActions] = await to(this.getAuditActions());

      if (errActions) {
        console.error("error retrieving audit actions", errActions);
        this.error.auditActions = this.getErrorMessage(errActions);
        this.loading.auditTrail = false;
        return;
      }

      const actions = responseActions.data.data.actions;

      this.actions = actions.map((u) => {
        return { label: u, value: u, name: u };
      });

      this.searchLogs();
    },
    getAuditParams() {
      // using Date constructor: new Date('1995-12-17T03:24:00')
      const startLocal = new Date(
        this.startDateString + "T" + this.startTime + ":00"
      );
      const endLocal = new Date(
        this.endDateString + "T" + this.endTime + ":59"
      );

      const format = "yyyy-MM-dd'T'HH:mm:ssX";

      const startUtcString = this.formatInTimeZone(startLocal, format, "UTC");

      const endUtcString = this.formatInTimeZone(endLocal, format, "UTC");

      return {
        from: startUtcString,
        to: endUtcString,
        users: this.usersSelected,
        actions: this.actionsSelected,
        data: this.auditInfo,
        limit: this.maxResults,
      };
    },
    validateSearchLogs() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.startTime) {
        this.error.startTime = this.$t("common.required");
        isValidationOk = false;
      }

      if (!this.endTime) {
        this.error.endTime = this.$t("common.required");
        isValidationOk = false;
      }

      if (this.startTime && this.endTime) {
        // using Date constructor: new Date('1995-12-17T03:24:00')
        const startLocal = new Date(
          this.startDateString + "T" + this.startTime + ":00"
        );
        const endLocal = new Date(
          this.endDateString + "T" + this.endTime + ":59"
        );

        if (this.dateIsBefore(endLocal, startLocal)) {
          this.error.startDate = this.$t("error.invalid_time_interval");
          this.error.startTime = this.$t("error.invalid_time_interval");
          this.error.endDate = this.$t("error.invalid_time_interval");
          this.error.endTime = this.$t("error.invalid_time_interval");

          isValidationOk = false;
        }
      }

      if (this.maxResults < 1 || this.maxResults > this.MAX_RESULTS_LIMIT) {
        this.error.maxResults = this.$t("error.must_be_between", {
          min: 1,
          max: this.MAX_RESULTS_LIMIT,
        });
        isValidationOk = false;
      }

      return isValidationOk;
    },
    async searchLogs() {
      if (!this.validateSearchLogs()) {
        return;
      }

      const auditParams = this.getAuditParams();
      this.loading.auditTrail = true;

      const [err, response] = await to(
        this.getAuditTrail(
          auditParams.users,
          auditParams.actions,
          auditParams.data,
          auditParams.from,
          auditParams.to,
          auditParams.limit
        )
      );
      this.loading.auditTrail = false;

      if (err) {
        console.error("error retrieving audit trail", err);
        this.error.auditTrail = this.getErrorMessage(err);
        return;
      }

      const audits = response.data.data.audits ? response.data.data.audits : [];

      this.logs = audits.map((audit) => {
        const localTimestamp = this.formatDate(
          this.parseIsoDate(audit.timestamp),
          "yyyy-MM-dd HH:mm:ss"
        );

        return {
          timestamp: localTimestamp,
          user: audit.user,
          action: audit.action,
          auditInfo: audit.data ? audit.data : "",
          auditInfoCollapsed: true,
          id: audit.id,
        };
      });
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.interval-date {
  display: inline-flex;
  margin-right: $spacing-06;
}

.interval-time {
  display: inline-flex;
}

.audit-info-collapsed {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 15rem;
}

.audit-info-pre {
  white-space: pre-wrap;
}

.expand-audit-info-button {
  float: right;
}

.search-button {
  margin-top: $spacing-05;
  margin-bottom: $spacing-07;
}
</style>
