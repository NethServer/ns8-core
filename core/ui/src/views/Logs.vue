<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("logs.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <NsInlineNotification
          class="landscape-warning"
          kind="warning"
          :title="$t('common.use_landscape_mode')"
          :description="$t('common.use_landscape_mode_description')"
        />
      </div>
    </div>
    <cv-tile :light="true">
      <div class="bx--row">
        <div class="bx--col-md-8">
          <cv-form @submit.prevent="searchLogs">
            <div>
              <cv-date-picker
                kind="single"
                :cal-options="calOptions"
                :date-label="$t('logs.start_date')"
                v-model="startDateString"
                class="interval-date"
                :invalid-message="error.timeInterval"
              >
              </cv-date-picker>
              <cv-time-picker
                :label="$t('logs.start_time_label')"
                :time.sync="startTime"
                ampm="24"
                :pattern="timePattern"
                :placeholder="timePlaceholder"
                :form-item="true"
                class="interval-time"
                :invalid-message="error.timeInterval"
              >
              </cv-time-picker>
            </div>
            <div>
              <cv-date-picker
                kind="single"
                :cal-options="calOptions"
                :date-label="$t('logs.end_date')"
                v-model="endDateString"
                class="interval-date"
                :invalid-message="error.timeInterval"
              >
              </cv-date-picker>
              <cv-time-picker
                :label="$t('logs.end_time_label')"
                :time.sync="endTime"
                ampm="24"
                :pattern="timePattern"
                :placeholder="timePlaceholder"
                :form-item="true"
                class="interval-time"
                :invalid-message="error.timeInterval"
              >
              </cv-time-picker>
            </div>
            <cv-multi-select
              v-model="usersSelected"
              :options="users"
              :title="$t('logs.users') + ' (' + $t('common.optional') + ')'"
              :label="$t('logs.users_label')"
              :helper-text="usersHelperText"
            >
            </cv-multi-select>
            <NsInlineNotification
              v-if="error.auditUsers"
              kind="error"
              :title="$t('logs.cannot_retrieve_audit_users')"
              :description="error.auditUsers"
              :showCloseButton="false"
            />
            <cv-multi-select
              v-model="actionsSelected"
              :options="actions"
              :title="$t('logs.actions') + ' (' + $t('common.optional') + ')'"
              :label="$t('logs.actions_label')"
              :helper-text="actionsHelperText"
            >
            </cv-multi-select>
            <NsInlineNotification
              v-if="error.auditActions"
              kind="error"
              :title="$t('logs.cannot_retrieve_audit_actions')"
              :description="error.auditActions"
              :showCloseButton="false"
            />
            <cv-text-input
              :label="$t('logs.auditInfo') + ' (' + $t('common.optional') + ')'"
              v-model="auditInfo"
              :placeholder="$t('logs.audit_info_placeholder')"
            >
            </cv-text-input>
            <cv-number-input
              :label="$t('logs.max_results')"
              min="1"
              v-model="maxResults"
              class="logs-max-results"
            >
            </cv-number-input>
            <div class="buttons">
              <NsButton
                kind="primary"
                class="search-button"
                :icon="Search20"
                :loading="loading.auditLogs"
                :disabled="loading.auditLogs"
                >{{ $t("logs.search_logs") }}</NsButton
              >
              <NsButton
                kind="ghost"
                :icon="Reset20"
                @click="initFilters"
                type="button"
                >{{ $t("common.reset_filters") }}</NsButton
              >
            </div>
          </cv-form>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-md-8">
          <NsInlineNotification
            v-if="error.auditLogs"
            kind="error"
            :title="$t('logs.cannot_retrieve_audit_logs')"
            :description="error.auditLogs"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-md-8">
          <NsEmptyState
            v-if="!this.tableRows.length && !loading.auditLogs"
            :title="$t('logs.no_logs_found')"
          >
            <template #description>{{
              $t("logs.no_logs_found_description")
            }}</template>
          </NsEmptyState>
          <div v-else>
            <cv-data-table-skeleton
              v-if="loading.auditLogs"
              :columns="i18nTableColumns"
              :rows="10"
            ></cv-data-table-skeleton>
            <cv-data-table
              v-if="
                !loading.auditLogs &&
                !error.auditUsers &&
                !error.auditActions &&
                !error.auditLogs
              "
              :sortable="true"
              :columns="i18nTableColumns"
              @sort="sortTable"
              :pagination="pagination"
              @pagination="paginateTable"
              ref="table"
            >
              <template slot="data">
                <cv-data-table-row
                  v-for="(row, rowIndex) in tablePage"
                  :key="`${rowIndex}`"
                  :value="`${rowIndex}`"
                >
                  <cv-data-table-cell>{{ row.timestamp }}</cv-data-table-cell>
                  <cv-data-table-cell>{{ row.user }}</cv-data-table-cell>
                  <cv-data-table-cell>{{ row.action }}</cv-data-table-cell>
                  <cv-data-table-cell
                    :class="{ 'audit-info-collapsed': row.auditInfoCollapsed }"
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
                      @click="row.auditInfoCollapsed = !row.auditInfoCollapsed"
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
              </template></cv-data-table
            >
          </div>
        </div>
      </div>
    </cv-tile>
  </div>
</template>

<script>
import AuditService from "../mixins/audit";
import to from "await-to-js";
import {
  UtilService,
  IconService,
  DateTimeService,
  DataTableService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "Logs",
  mixins: [
    DateTimeService,
    DataTableService,
    IconService,
    AuditService,
    UtilService,
  ],
  pageTitle() {
    return this.$t("logs.title");
  },
  data() {
    return {
      maxResults: 0,
      startDateString: "",
      endDateString: "",
      startTime: "",
      endTime: "",
      timePattern: "([01]\\d|2[0-3]):?([0-5]\\d)",
      timePlaceholder: "hh:mm",
      calOptions: {
        dateFormat: "Y-m-d",
      },
      users: [],
      usersSelected: [],
      actions: [],
      actionsSelected: [],
      auditInfo: "",
      tableColumns: ["timestamp", "user", "action", "auditInfo"],
      tableRows: [],
      loading: {
        auditLogs: true,
      },
      error: {
        auditUsers: "",
        auditActions: "",
        auditLogs: "",
        timeInterval: "",
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
        return this.$t("logs." + column);
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
      this.error.auditLogs = "";
      this.error.timeInterval = "";

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
      this.maxResults = 100; //// save to local storage?
    },
    async getAuditData() {
      // get audit users
      const [errUsers, responseUsers] = await to(this.getAuditUsers());

      if (errUsers) {
        console.error("error retrieving audit users", errUsers);
        this.error.auditUsers = this.getErrorMessage(errUsers);
        this.loading.auditLogs = false;
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
        this.loading.auditLogs = false;
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
      // using Date constructor: new Date('1995-12-17T03:24:00')
      const startLocal = new Date(
        this.startDateString + "T" + this.startTime + ":00"
      );
      const endLocal = new Date(
        this.endDateString + "T" + this.endTime + ":59"
      );

      if (this.dateIsBefore(endLocal, startLocal)) {
        this.error.timeInterval = this.$t("logs.invalid_interval");
        return false;
      }

      return true;
    },
    async searchLogs() {
      this.error.auditLogs = "";
      this.error.timeInterval = "";
      const validationOk = this.validateSearchLogs();

      if (!validationOk) {
        return;
      }

      const auditParams = this.getAuditParams();
      this.loading.auditLogs = true;

      const [err, response] = await to(
        this.getAuditLogs(
          auditParams.users,
          auditParams.actions,
          auditParams.data,
          auditParams.from,
          auditParams.to,
          auditParams.limit
        )
      );
      this.loading.auditLogs = false;

      if (err) {
        console.error("error retrieving audit logs", err);
        this.error.auditLogs = this.getErrorMessage(err);
        return;
      }

      const audits = response.data.data.audits ? response.data.data.audits : [];

      let logs = audits.map((audit) => {
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

      this.tableRows = logs;
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

.buttons {
  margin-top: $spacing-07;
  margin-bottom: $spacing-09;
}

.logs-max-results {
  max-width: 12rem;
}
</style>
