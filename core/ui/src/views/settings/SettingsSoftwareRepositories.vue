<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-breadcrumb
            aria-label="breadcrumb"
            :no-trailing-slash="true"
            class="breadcrumb"
          >
            <cv-breadcrumb-item>
              <cv-link to="/settings">{{ $t("settings.title") }}</cv-link>
            </cv-breadcrumb-item>
            <cv-breadcrumb-item>
              <span>{{ $t("settings_sw_repositories.title") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16 subpage-title">
          <h3>{{ $t("settings_sw_repositories.title") }}</h3>
        </div>
      </div>
      <div class="bx--row landscape-warning">
        <div class="bx--col-lg-16">
          <NsInlineNotification
            kind="warning"
            :title="$t('common.use_landscape_mode')"
            :description="$t('common.use_landscape_mode_description')"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <!-- repository being deleted -->
          <NsInlineNotification
            v-if="repoToDelete"
            kind="warning"
            :title="
              $t('settings_sw_repositories.repository_is_going_to_be_deleted', {
                object: repoToDelete.name,
              })
            "
            :actionLabel="$t('common.cancel')"
            @action="cancelDeleteRepository()"
            :showCloseButton="false"
            :timer="DELETE_DELAY"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-tile :light="true">
            <div class="toolbar">
              <NsButton
                v-if="repositories.length"
                kind="secondary"
                :icon="Add20"
                @click="showCreateRepoModal()"
                :disabled="loading.repositories"
                >{{
                  $t("settings_sw_repositories.create_repository")
                }}</NsButton
              >
            </div>
            <NsDataTable
              :allRows="repositories"
              :columns="i18nTableColumns"
              :rawColumns="tableColumns"
              :sortable="true"
              :isSearchable="false"
              :isLoading="loading.repositories"
              :skeletonRows="5"
              :overflow-menu="true"
              :isErrorShown="!!error.listRepositories"
              :errorTitle="$t('action.list-repositories')"
              :errorDescription="error.listRepositories"
              :itemsPerPageLabel="$t('pagination.items_per_page')"
              :rangeOfTotalItemsLabel="$t('pagination.range_of_total_items')"
              :ofTotalPagesLabel="$t('pagination.of_total_pages')"
              :backwardText="$t('pagination.previous_page')"
              :forwardText="$t('pagination.next_page')"
              :pageNumberLabel="$t('pagination.page_number')"
              @updatePage="tablePage = $event"
              ref="table"
            >
              <template slot="empty-state">
                <NsEmptyState
                  :title="$t('settings_sw_repositories.no_sw_repositories')"
                >
                  <template #description>
                    <div>
                      {{
                        $t(
                          "settings_sw_repositories.no_sw_repositories_description"
                        )
                      }}
                    </div>
                    <NsButton
                      kind="primary"
                      :icon="Add20"
                      @click="showCreateRepoModal()"
                      class="empty-state-button"
                      >{{
                        $t("settings_sw_repositories.create_repository")
                      }}</NsButton
                    ></template
                  >
                </NsEmptyState>
              </template>
              <template slot="data">
                <cv-data-table-row
                  v-for="(row, rowIndex) in tablePage"
                  :key="`${rowIndex}`"
                  :value="`${rowIndex}`"
                >
                  <cv-data-table-cell>{{ row.name }}</cv-data-table-cell>
                  <cv-data-table-cell class="break-word">{{
                    row.url
                  }}</cv-data-table-cell>
                  <cv-data-table-cell>
                    <cv-tag
                      v-if="row.status"
                      kind="green"
                      :label="$t('common.enabled')"
                      class="no-margin"
                    ></cv-tag>
                    <cv-tag
                      v-else
                      kind="high-contrast"
                      :label="$t('common.disabled')"
                      class="no-margin"
                    ></cv-tag>
                    <!-- <div class="badge-container"> ////
                        <template v-if="row.status"
                          ><span class="green-badge left-badge"></span>
                          {{ $t("common.enabled") }}</template
                        >
                        <template v-else
                          ><span class="gray-badge left-badge"></span
                          >{{ $t("common.disabled") }}</template
                        >
                      </div> -->
                  </cv-data-table-cell>
                  <cv-data-table-cell>
                    <cv-tag
                      v-if="row.testing"
                      kind="green"
                      :label="$t('common.enabled')"
                      class="no-margin"
                    ></cv-tag>
                    <cv-tag
                      v-else
                      kind="high-contrast"
                      :label="$t('common.disabled')"
                      class="no-margin"
                    ></cv-tag>
                    <!-- <div class="badge-container"> ////
                        <template v-if="row.testing"
                          ><span class="green-badge left-badge"></span>
                          {{ $t("common.enabled") }}</template
                        >
                        <template v-else
                          ><span class="gray-badge left-badge"></span
                          >{{ $t("common.disabled") }}</template
                        >
                      </div> -->
                  </cv-data-table-cell>
                  <cv-data-table-cell class="table-overflow-menu-cell">
                    <cv-overflow-menu flip-menu class="table-overflow-menu">
                      <cv-overflow-menu-item @click="showEditRepoModal(row)">
                        <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
                      </cv-overflow-menu-item>
                      <NsMenuDivider />
                      <cv-overflow-menu-item
                        danger
                        @click="willDeleteRepository(row)"
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('common.delete')"
                        />
                      </cv-overflow-menu-item>
                    </cv-overflow-menu>
                  </cv-data-table-cell>
                </cv-data-table-row>
              </template>
            </NsDataTable>
          </cv-tile>
        </div>
      </div>
    </div>
    <!-- new repository modal -->
    <NsModal
      size="default"
      :visible="q.isShownCreateRepoModal"
      @modal-hidden="q.isShownCreateRepoModal = false"
      @primary-click="createRepository"
      :primary-button-disabled="loading.createRepository"
    >
      <template slot="title">{{
        $t("settings_sw_repositories.create_repository")
      }}</template>
      <template slot="content">
        <cv-form @submit.prevent="createRepository">
          <cv-text-input
            :label="$t('settings_sw_repositories.name')"
            v-model.trim="q.newRepoName"
            :invalid-message="$t(error.name)"
            ref="name"
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('settings_sw_repositories.url')"
            v-model.trim="q.newRepoUrl"
            :invalid-message="$t(error.url)"
            ref="url"
          >
          </cv-text-input>
          <cv-toggle
            :label="$t('settings_sw_repositories.status')"
            value="statusValue"
            :form-item="true"
            v-model="q.isNewRepoEnabled"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle>
          <cv-toggle
            :label="$t('settings_sw_repositories.testing')"
            value="testingValue"
            :form-item="true"
            v-model="q.isNewRepoTesting"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle>
        </cv-form>
        <NsInlineNotification
          v-if="error.addRepository"
          kind="error"
          :title="$t('action.add-repository')"
          :description="error.addRepository"
          :showCloseButton="false"
        />
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
      <template slot="primary-button">{{
        $t("settings_sw_repositories.create_repository")
      }}</template>
    </NsModal>
    <!-- edit repository modal -->
    <NsModal
      size="default"
      :visible="q.isShownEditRepoModal"
      @modal-hidden="q.isShownEditRepoModal = false"
      @primary-click="editRepository"
      :primary-button-disabled="loading.editRepository"
    >
      <template slot="title">{{
        $t("settings_sw_repositories.edit_repository")
      }}</template>
      <template slot="content">
        <cv-form @submit.prevent="editRepository">
          <cv-text-input
            :label="$t('settings_sw_repositories.name')"
            v-model.trim="q.editRepoName"
            :invalid-message="$t(error.name)"
            disabled
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('settings_sw_repositories.url')"
            v-model.trim="q.editRepoUrl"
            :invalid-message="$t(error.url)"
            disabled
          >
          </cv-text-input>
          <cv-toggle
            :label="$t('settings_sw_repositories.status')"
            value="statusValue"
            :form-item="true"
            v-model="q.isEditRepoEnabled"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle>
          <cv-toggle
            :label="$t('settings_sw_repositories.testing')"
            value="testingValue"
            :form-item="true"
            v-model="q.isEditRepoTesting"
          >
            <template slot="text-left">{{ $t("common.disabled") }}</template>
            <template slot="text-right">{{ $t("common.enabled") }}</template>
          </cv-toggle>
        </cv-form>
        <div v-if="error.alterRepository" class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              kind="error"
              :title="$t('action.alter-repository')"
              :description="error.alterRepository"
              :showCloseButton="false"
            />
          </div>
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
      <template slot="primary-button">{{
        $t("settings_sw_repositories.edit_repository")
      }}</template>
    </NsModal>
  </div>
</template>

<script>
import to from "await-to-js";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";

export default {
  name: "SettingsSoftwareRepositories",
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("settings_sw_repositories.title");
  },
  data() {
    return {
      q: {
        isShownCreateRepoModal: false,
        isShownEditRepoModal: false,
        newRepoName: "",
        newRepoUrl: "",
        isNewRepoTesting: false,
        isNewRepoEnabled: true,
        editRepoName: "",
        editRepoUrl: "",
        isEditRepoTesting: false,
        isEditRepoEnabled: true,
      },
      repositories: [],
      tablePage: [],
      tableColumns: ["name", "url", "status", "testing"],
      repoToDelete: null,
      loading: {
        repositories: true,
        createRepository: false,
        editRepository: false,
      },
      error: {
        name: "",
        url: "",
        listRepositories: "",
        addRepository: "",
        alterRepository: "",
      },
    };
  },
  computed: {
    i18nTableColumns() {
      return this.tableColumns.map((column) => {
        return this.$t("settings_sw_repositories." + column);
      });
    },
  },
  watch: {
    "q.isShownCreateRepoModal": function () {
      if (this.q.isShownCreateRepoModal) {
        setTimeout(() => {
          this.focusElement("name");
        }, 300);
      }
    },
  },
  beforeRouteEnter(to, from, next) {
    next((vm) => {
      vm.watchQueryData(vm);
      vm.queryParamsToDataForCore(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    this.queryParamsToDataForCore(this, to.query);
    next();
  },
  created() {
    this.listRepositories();
  },
  methods: {
    addRepositoryValidationFailed(validationErrors) {
      // enable "Create repository" button
      this.loading.createRepository = false;

      // enable "Edit repository" button
      this.loading.editRepository = false;

      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error[param] = "settings_sw_repositories." + validationError.error;

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    showCreateRepoModal() {
      this.clearErrors(this);
      this.q.isShownCreateRepoModal = true;
    },
    showEditRepoModal(repo) {
      this.clearErrors(this);
      this.q.editRepoName = repo.name;
      this.q.editRepoUrl = repo.url;
      this.q.isEditRepoTesting = repo.testing;
      this.q.isEditRepoEnabled = repo.status;
      this.q.isShownEditRepoModal = true;
    },
    async listRepositories() {
      this.loading.repositories = true;
      this.error.listRepositories = "";
      const taskAction = "list-repositories";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listRepositoriesCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.listRepositories = this.getErrorMessage(err);
        this.loading.repositories = false;
        return;
      }
    },
    listRepositoriesCompleted(taskContext, taskResult) {
      this.repositories = taskResult.output;
      this.loading.repositories = false;
    },
    validateNewRepository() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.q.newRepoName) {
        this.error.name = "common.required";
        isValidationOk = false;
      }

      if (!this.q.newRepoUrl) {
        this.error.url = "common.required";
        isValidationOk = false;
      }

      return isValidationOk;
    },
    async createRepository() {
      this.loading.createRepository = true;
      this.error.addRepository = "";

      if (!this.validateNewRepository()) {
        this.loading.createRepository = false;
        return;
      }

      const taskAction = "add-repository";
      const eventId = this.getUuid();

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.addRepositoryAborted
      );

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-ok-${eventId}`,
        this.addRepositoryValidationOk
      );
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.addRepositoryValidationFailed
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.addRepositoryCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.q.newRepoName,
            url: this.q.newRepoUrl,
            status: this.q.isNewRepoEnabled,
            testing: this.q.isNewRepoTesting,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.addRepository = this.getErrorMessage(err);
        this.loading.createRepository = false;
        return;
      }
    },
    addRepositoryAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.addRepository = this.$t("error.generic_error");
      this.loading.createRepository = false;

      // hide modal so that user can see error notification
      this.q.isShownCreateRepoModal = false;
    },
    addRepositoryValidationOk() {
      // hide modal after validation
      this.q.isShownCreateRepoModal = false;

      // enable "Create repository" button
      this.loading.createRepository = false;
    },
    alterRepositoryValidationOk() {
      // hide modal after validation
      this.q.isShownEditRepoModal = false;

      // enable "Edit repository" button
      this.loading.editRepository = false;
    },
    async editRepository() {
      this.loading.editRepository = true;
      this.error.alterRepository = "";
      const taskAction = "alter-repository";

      // register to task validation
      this.$root.$off(taskAction + "-validation-ok");
      this.$root.$once(
        taskAction + "-validation-ok",
        this.alterRepositoryValidationOk
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.alterRepositoryCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.q.editRepoName,
            status: this.q.isEditRepoEnabled,
            testing: this.q.isEditRepoTesting,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.alterRepository = this.getErrorMessage(err);
        this.loading.editRepository = false;
        return;
      }
    },
    addRepositoryCompleted() {
      this.q.newRepoName = "";
      this.q.newRepoUrl = "";
      this.q.isNewRepoTesting = false;
      this.q.isNewRepoEnabled = true;
      this.listRepositories();
    },
    alterRepositoryAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.alterRepository = this.$t("error.generic_error");
      this.loading.editRepository = false;

      // hide modal so that user can see error notification
      this.q.isShownEditRepoModal = false;
    },
    alterRepositoryCompleted() {
      this.listRepositories();
    },
    willDeleteRepository(repo) {
      const timeout = setTimeout(() => {
        this.deleteRepository(repo);
        this.repoToDelete = null;
      }, this.DELETE_DELAY);

      repo.timeout = timeout;
      this.repoToDelete = repo;

      // remove repo from table
      this.repositories = this.repositories.filter((r) => {
        return r.name != repo.name;
      });
    },
    async deleteRepository(repo) {
      const taskAction = "remove-repository";

      // register to task completion (using $on instead of $once for multiple revertable deletions)
      this.$root.$on(taskAction + "-completed", this.removeRepositoryCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: repo.name,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    removeRepositoryCompleted() {
      this.listRepositories();
    },
    cancelDeleteRepository() {
      clearTimeout(this.repoToDelete.timeout);
      this.repoToDelete = null;
      this.listRepositories();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.break-word {
  word-wrap: break-word;
  max-width: 20vw;
}
</style>
