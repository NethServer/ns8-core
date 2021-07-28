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
              <span>{{ $t("settings.sw_repositories") }}</span>
            </cv-breadcrumb-item>
          </cv-breadcrumb>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16 page-subtitle">
          <h3>{{ $t("settings.sw_repositories") }}</h3>
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <!-- repository being deleted -->
          <NsInlineNotification
            v-if="repoToDelete"
            kind="warning"
            :title="
              $t('settings_sw_repositories.repository_deleted', {
                repo: repoToDelete.name,
              })
            "
            :actionLabel="$t('common.undo')"
            @action="cancelDeleteRepository()"
            :showCloseButton="false"
          />
        </div>
      </div>
      <div class="bx--row">
        <div class="bx--col-lg-16">
          <cv-tile :light="true" class="content-tile">
            <div v-if="!tableRows.length && !loading.repositories">
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
            </div>
            <div v-else>
              <div class="toolbar">
                <NsButton
                  kind="secondary"
                  :icon="Add20"
                  @click="showCreateRepoModal()"
                  :disabled="loading.repositories"
                  >{{
                    $t("settings_sw_repositories.create_repository")
                  }}</NsButton
                >
              </div>
              <cv-data-table-skeleton
                v-if="loading.repositories"
                :columns="i18nTableColumns"
                :rows="5"
              ></cv-data-table-skeleton>
              <cv-data-table
                v-if="!loading.repositories"
                :sortable="true"
                :columns="i18nTableColumns"
                @sort="sortTable"
                :pagination="pagination"
                @pagination="paginateTable"
                :overflow-menu="true"
                ref="table"
              >
                <template slot="data">
                  <cv-data-table-row
                    v-for="(row, rowIndex) in tablePage"
                    :key="`${rowIndex}`"
                    :value="`${rowIndex}`"
                  >
                    <cv-data-table-cell>{{ row.name }}</cv-data-table-cell>
                    <cv-data-table-cell>{{ row.url }}</cv-data-table-cell>
                    <cv-data-table-cell>
                      <div class="badge-container">
                        <template v-if="row.status"
                          ><span class="green-badge left-badge"></span>
                          {{ $t("common.enabled") }}</template
                        >
                        <template v-else
                          ><span class="gray-badge left-badge"></span
                          >{{ $t("common.disabled") }}</template
                        >
                      </div>
                    </cv-data-table-cell>
                    <cv-data-table-cell
                      ><div class="badge-container">
                        <template v-if="row.testing"
                          ><span class="green-badge left-badge"></span>
                          {{ $t("common.enabled") }}</template
                        >
                        <template v-else
                          ><span class="gray-badge left-badge"></span
                          >{{ $t("common.disabled") }}</template
                        >
                      </div></cv-data-table-cell
                    >
                    <cv-data-table-cell>
                      <cv-overflow-menu flip-menu class="table-overflow-menu">
                        <cv-overflow-menu-item
                          @click="showEditRepoModal(row)"
                          >{{ $t("common.edit") }}</cv-overflow-menu-item
                        >
                        <cv-overflow-menu-item
                          danger
                          @click="willDeleteRepository(row)"
                          >{{ $t("common.delete") }}</cv-overflow-menu-item
                        >
                      </cv-overflow-menu>
                    </cv-data-table-cell>
                  </cv-data-table-row>
                </template></cv-data-table
              >
            </div>
          </cv-tile>
        </div>
      </div>
    </div>
    <!-- new repository modal -->
    <cv-modal
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
          >
          </cv-text-input>
          <cv-text-input
            :label="$t('settings_sw_repositories.url')"
            v-model.trim="q.newRepoUrl"
            :invalid-message="$t(error.url)"
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
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
      <template slot="primary-button">{{
        $t("settings_sw_repositories.create_repository")
      }}</template>
    </cv-modal>
    <!-- edit repository modal -->
    <cv-modal
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
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
      <template slot="primary-button">{{
        $t("settings_sw_repositories.edit_repository")
      }}</template>
    </cv-modal>
  </div>
</template>

<script>
import to from "await-to-js";
import DataTableService from "../mixins/dataTable";
import NsInlineNotification from "../components/NsInlineNotification.vue";
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "andrelib"; ////

export default {
  name: "SettingsSoftwareRepository",
  components: { NsInlineNotification },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    DataTableService,
  ],
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
      tableColumns: ["name", "url", "status", "testing"],
      tableRows: [],
      repoToDelete: null,
      deleteRepoDelay: 5000, // you have 5 seconds to undo repository deletion
      loading: {
        repositories: true,
        createRepository: false,
        editRepository: false,
      },
      error: {
        name: "",
        url: "",
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

      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error[param] = "settings_sw_repositories." + validationError.error;
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
      const taskAction = "list-repositories";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.listRepositoriesCompleted);

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
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    listRepositoriesCompleted(taskContext, taskResult) {
      // unregister from event
      this.$root.$off("list-repositories-completed");

      this.tableRows = taskResult.output;
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

      if (!this.validateNewRepository()) {
        this.loading.createRepository = false;
        return;
      }

      const taskAction = "add-repository";

      // register to task validation
      this.$root.$on(
        taskAction + "-validation-ok",
        this.addRepositoryValidationOk
      );
      this.$root.$on(
        taskAction + "-validation-failed",
        this.addRepositoryValidationFailed
      );

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.addRepositoriesCompleted);

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
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    addRepositoryValidationOk() {
      // hide modal after validation
      this.q.isShownCreateRepoModal = false;

      // enable "Create repository" button
      this.loading.createRepository = false;

      // unregister from event
      this.$root.$off("add-repository-validation-ok");
    },
    alterRepositoryValidationOk() {
      // hide modal after validation
      this.q.isShownEditRepoModal = false;

      // enable "Edit repository" button
      this.loading.editRepository = false;

      // unregister from event
      this.$root.$off("alter-repository-validation-ok");
    },
    async editRepository() {
      this.loading.editRepository = true;
      const taskAction = "alter-repository";

      // register to task validation
      this.$root.$on(
        taskAction + "-validation-ok",
        this.alterRepositoryValidationOk
      );

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.alterRepositoryCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: this.q.editRepoName,
            status: this.q.isEditRepoTesting,
            testing: this.q.isEditRepoEnabled,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    addRepositoriesCompleted() {
      this.$root.$off("add-repository-completed");
      this.listRepositories();
    },
    alterRepositoryCompleted() {
      this.$root.$off("alter-repository-completed");
      this.listRepositories();
    },
    willDeleteRepository(repo) {
      const timeout = setTimeout(() => {
        this.deleteRepository(repo);
        this.repoToDelete = null;
      }, this.deleteRepoDelay);

      repo.timeout = timeout;
      this.repoToDelete = repo;

      // remove repo from table
      this.tableRows = this.tableRows.filter((r) => {
        return r.name != repo.name;
      });
    },
    async deleteRepository(repo) {
      const taskAction = "remove-repository";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.removeRepositoryCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            name: repo.name,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErrorNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    removeRepositoryCompleted() {
      this.$root.$off("remove-repository-completed");
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
@import "../styles/carbon-utils";
</style>
