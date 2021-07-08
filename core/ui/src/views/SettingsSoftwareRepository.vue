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
                      <div class="status-container">
                        <template v-if="row.status"
                          ><span class="enabled-badge"></span>
                          {{ $t("common.enabled") }}</template
                        >
                        <template v-else
                          ><span class="disabled-badge"></span
                          >{{ $t("common.disabled") }}</template
                        >
                      </div>
                    </cv-data-table-cell>
                    <cv-data-table-cell
                      ><div class="status-container">
                        <template v-if="row.testing"
                          ><span class="enabled-badge"></span>
                          {{ $t("common.enabled") }}</template
                        >
                        <template v-else
                          ><span class="disabled-badge"></span
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
                          @click="deleteRepository(row)"
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
import TaskService from "@/mixins/task";
import to from "await-to-js";
import UtilService from "@/mixins/util";
import IconService from "@/mixins/icon";
import QueryParamService from "@/mixins/queryParam";
import NsEmptyState from "@/components/NsEmptyState";
import NsButton from "@/components/NsButton";
import DataTableService from "../mixins/dataTable";

let nethserver = window.nethserver;

export default {
  name: "SettingsSoftwareRepository",
  components: { NsEmptyState, NsButton },
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
      loading: {
        repositories: true,
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
      console.log("beforeRouteEnter", to, from); ////
      nethserver.watchQueryData(vm);
      vm.queryParamsToData(vm, to.query);
    });
  },
  beforeRouteUpdate(to, from, next) {
    console.log("beforeRouteUpdate", to, from); ////
    this.queryParamsToData(this, to.query);
    next();
  },
  created() {
    // register to events
    this.$root.$on("validationFailed", this.validationFailed); ////

    // this.$root.$on("validationOk", this.validationOk);

    this.listRepositories();
  },
  // beforeDestroy() { ////
  //   // remove event listeners
  //   this.$root.$off("validationFailed");
  //   this.$root.$off("validationOk");
  // },
  methods: {
    validationFailed(validationErrors) {
      //// remove
      console.log("validationErrors", validationErrors); ////

      //// remove, it's always an array
      // if (Array.isArray(validationErrors)) {
      // list of errors
      for (const validationError of validationErrors) {
        const param = validationError.parameter;
        // set i18n error message
        this.error[param] = "settings_sw_repositories." + validationError.error;
      }
      // } else { ////
      //   // single error
      //   const param = validationErrors.parameter;
      //   // set i18n error message
      //   this.error[param] =
      //     "settings_sw_repositories." + validationErrors.error;
      // }
    },
    // validationOk(task) {
    //   //// remove
    //   console.log("validationOk, task", task); ////
    //   this.q.isShownCreateRepoModal = false;
    //   this.q.isShownEditRepoModal = false;
    // },
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
      const taskAction = "list-repositories";

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.listRepositoriesCompleted);

      const res = await to(
        this.createTask({
          action: taskAction,
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErroNotification(
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

      // //// remove
      // this.tableRows.push({
      //   name: "Afsdafds",
      //   url: "asdfsdlfjsdf",
      //   status: false,
      //   testing: false,
      // });

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
      console.log("createRepository"); ////

      if (!this.validateNewRepository()) {
        return;
      }

      ////
      console.log("create repo", this.q); ////

      const taskAction = "add-repository";

      // register to task validation
      this.$root.$on(
        taskAction + "-validation-ok",
        this.addRepositoryValidationOk
      );

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.addRepositoriesCompleted);

      const res = await to(
        this.createTask({
          action: taskAction,
          data: {
            name: this.q.newRepoName,
            url: this.q.newRepoUrl,
            status: this.q.isNewRepoEnabled,
            testing: this.q.isNewRepoTesting,
          },
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErroNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    addRepositoryValidationOk(task) {
      console.log("addRepositoryValidationOk", task); ////

      // hide modal after validation
      this.q.isShownCreateRepoModal = false;

      // unregister from event
      this.$root.$off("add-repository-validation-ok");
    },
    alterRepositoryValidationOk(task) {
      console.log("editRepositoryValidationOk", task); ////

      // hide modal after validation
      this.q.isShownEditRepoModal = false;

      // unregister from event
      this.$root.$off("alter-repository-validation-ok");
    },
    async editRepository() {
      console.log("editRepository"); ////

      const taskAction = "alter-repository"; ////

      // register to task validation
      this.$root.$on(
        taskAction + "-validation-ok",
        this.alterRepositoryValidationOk
      );

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.alterRepositoryCompleted);

      const res = await to(
        this.createTask({
          action: taskAction,
          data: {
            name: this.q.editRepoName,
            status: this.q.isEditRepoTesting,
            testing: this.q.isEditRepoEnabled,
          },
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErroNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    addRepositoriesCompleted() {
      console.log("alterRepositoryCompleted"); ////

      this.$root.$off("add-repository-completed");
      this.listRepositories();
    },
    alterRepositoryCompleted() {
      console.log("alterRepositoryCompleted"); ////

      this.$root.$off("alter-repository-completed");
      this.listRepositories();
    },
    async deleteRepository(repo) {
      console.log("deleteRepository", repo); ////

      const taskAction = "remove-repository"; ////

      // register to task completion
      this.$root.$on(taskAction + "-completed", this.removeRepositoryCompleted);

      const res = await to(
        this.createTask({
          action: taskAction,
          data: {
            name: repo.name,
          },
          extra: {
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        this.createTaskErroNotification(
          err,
          this.$t("task.cannot_create_task", { action: taskAction })
        );
        return;
      }
    },
    removeRepositoryCompleted() {
      console.log("removeRepositoryCompleted"); ////

      this.$root.$off("remove-repository-completed");
      this.listRepositories();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
