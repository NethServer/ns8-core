<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("backup.title") }}</h2>
        </div>
      </div>
      <template v-if="loading.listBackupRepositories || loading.listBackups">
        <!-- repositories skeleton -->
        <div class="bx--row">
          <div class="bx--col">
            <cv-skeleton-text heading width="40%"></cv-skeleton-text>
          </div>
        </div>
        <div class="bx--row">
          <div
            v-for="index in 2"
            :key="index"
            class="bx--col-md-4 bx--col-max-4"
          >
            <cv-tile light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="8"
              ></cv-skeleton-text>
            </cv-tile>
          </div>
        </div>
        <!-- backups skeleton -->
        <div class="bx--row">
          <div class="bx--col">
            <cv-skeleton-text heading width="40%"></cv-skeleton-text>
          </div>
        </div>
        <div class="bx--row">
          <div
            v-for="index in 2"
            :key="index"
            class="bx--col-md-4 bx--col-max-4"
          >
            <cv-tile light>
              <cv-skeleton-text
                :paragraph="true"
                :line-count="8"
              ></cv-skeleton-text>
            </cv-tile>
          </div>
        </div>
      </template>
      <!-- empty state repositories -->
      <div v-else-if="!repositories.length" class="bx--row">
        <div class="bx--col">
          <NsEmptyState :title="$t('backup.no_backup_repository')">
            <template #pictogram>
              <HardDrive />
            </template>
            <template #description>
              <div>{{ $t("backup.empty_state_repository_description") }}</div>
              <NsButton
                kind="primary"
                :icon="Add20"
                @click="showAddRepoModal()"
                class="empty-state-button"
                >{{ $t("backup.add_repository") }}
              </NsButton>
            </template>
          </NsEmptyState>
        </div>
      </div>
      <template v-else>
        <!-- disabled backups warning -->
        <div class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              v-if="!loading.listBackups && disabledBackups.length"
              kind="warning"
              :title="$t('backup.disabled_backups')"
              :description="
                $tc(
                  'backup.disabled_backups_description',
                  disabledBackups.length,
                  {
                    num: disabledBackups.length,
                  }
                )
              "
              :showCloseButton="false"
            />
          </div>
        </div>
        <!-- unconfigured instances warning -->
        <div class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              v-if="!loading.listBackups && unconfiguredInstances.length"
              kind="warning"
              :title="$t('backup.app_instances_not_backed_up')"
              :description="
                $tc(
                  'backup.app_instances_not_backed_up_description',
                  unconfiguredInstances.length,
                  {
                    numInstances: unconfiguredInstances.length,
                  }
                )
              "
              :actionLabel="$t('backup.schedule_backup')"
              @action="showCreateBackupModal('notBackedUp')"
              :showCloseButton="false"
            />
          </div>
        </div>
        <!-- repositories -->
        <div class="bx--row">
          <div class="bx--col">
            <h4 class="mg-bottom-md">{{ $t("backup.repositories") }}</h4>
          </div>
        </div>
        <div class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              v-if="error.listBackupRepositories"
              kind="error"
              :title="$t('action.list-backup-repositories')"
              :description="error.listBackupRepositories"
              :showCloseButton="false"
            />
            <NsInlineNotification
              v-if="error.removeBackupRepository"
              kind="error"
              :title="$t('action.remove-backup-repository')"
              :description="error.removeBackupRepository"
              :showCloseButton="false"
            />
          </div>
        </div>
        <div class="bx--row">
          <div class="bx--col">
            <NsButton kind="secondary" :icon="Add20" @click="showAddRepoModal()"
              >{{ $t("backup.add_repository") }}
            </NsButton>
          </div>
        </div>
        <div class="bx--row">
          <div
            v-for="repo in repositories"
            :key="repo.id"
            class="bx--col-md-4 bx--col-max-4"
          >
            <NsInfoCard
              light
              :title="repo.name"
              :icon="DataBase32"
              :showOverflowMenu="true"
            >
              <template #menu>
                <cv-overflow-menu
                  :flip-menu="true"
                  tip-position="top"
                  tip-alignment="end"
                  class="top-right-overflow-menu"
                >
                  <cv-overflow-menu-item @click="showEditRepoModal(repo)">
                    {{ $t("common.edit") }}
                  </cv-overflow-menu-item>
                  <cv-overflow-menu-item
                    danger
                    @click="showDeleteRepoModal(repo)"
                  >
                    {{ $t("common.delete") }}
                  </cv-overflow-menu-item>
                </cv-overflow-menu>
              </template>
              <template #content>
                <div class="card-content">
                  <div class="row">
                    {{ $t("backup." + repo.provider) }}
                  </div>
                  <div class="row">
                    {{ repo.url }}
                  </div>
                  <div class="row actions">
                    <NsButton
                      kind="ghost"
                      :icon="ZoomIn20"
                      @click="showRepoDetailsModal(repo)"
                      >{{ $t("common.details") }}
                    </NsButton>
                  </div>
                </div>
              </template>
            </NsInfoCard>
          </div>
        </div>
        <!-- backups -->
        <div class="bx--row">
          <div class="bx--col">
            <h4 class="mg-bottom-md">{{ $t("backup.backup_scheduler") }}</h4>
          </div>
        </div>
        <div v-if="error.listBackups" class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              kind="error"
              :title="$t('action.list-backups')"
              :description="error.listBackups"
              :showCloseButton="false"
            />
          </div>
        </div>
        <div v-if="error.removeBackup" class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              kind="error"
              :title="$t('action.remove-backup')"
              :description="error.removeBackup"
              :showCloseButton="false"
            />
          </div>
        </div>
        <!-- empty state backups -->
        <div v-if="!backups.length" class="bx--row">
          <div class="bx--col">
            <NsEmptyState :title="$t('backup.no_backup_scheduled')">
              <template #description>
                <NsButton
                  kind="primary"
                  :icon="Add20"
                  @click="showCreateBackupModal('')"
                  class="empty-state-button-no-description"
                  >{{ $t("backup.schedule_backup") }}
                </NsButton>
              </template>
            </NsEmptyState>
          </div>
        </div>
        <template v-else>
          <div class="bx--row">
            <div class="bx--col">
              <NsButton
                kind="secondary"
                :icon="Add20"
                @click="showCreateBackupModal('')"
                >{{ $t("backup.schedule_backup") }}
              </NsButton>
            </div>
          </div>
          <div class="bx--row">
            <div
              v-for="backup in backups"
              :key="backup.id"
              class="bx--col-md-4 bx--col-max-4"
            >
              <NsInfoCard
                light
                :title="backup.name"
                :icon="Save32"
                :showOverflowMenu="true"
              >
                <template #menu>
                  <cv-overflow-menu
                    :flip-menu="true"
                    tip-position="top"
                    tip-alignment="end"
                    class="top-right-overflow-menu"
                  >
                    <cv-overflow-menu-item @click="runBackup(backup)">
                      {{ $t("backup.run_backup_now") }}
                    </cv-overflow-menu-item>
                    <cv-overflow-menu-item @click="showEditBackupModal(backup)">
                      {{ $t("common.edit") }}
                    </cv-overflow-menu-item>
                    <cv-overflow-menu-item
                      danger
                      @click="showDeleteBackupModal(backup)"
                    >
                      {{ $t("common.delete") }}
                    </cv-overflow-menu-item>
                  </cv-overflow-menu>
                </template>
                <template #content>
                  <div class="card-content">
                    <div class="row instance-to-repo">
                      <div class="backup-source">
                        <NsSvg :svg="Application20" class="icon" />
                        <span v-if="backup.instances.length == 1">
                          {{
                            backup.instances[0].ui_name
                              ? backup.instances[0].ui_name
                              : backup.instances[0].module_id
                          }}
                        </span>
                        <span v-else>
                          <!-- multiple instances -->
                          {{
                            backup.instances.length +
                            " " +
                            $tc("backup.instances", backup.instances.length)
                          }}
                        </span>
                      </div>
                      <div class="arrow-right">
                        <NsSvg :svg="ArrowRight20" />
                      </div>
                      <div class="backup-destination">
                        <NsSvg :svg="DataBase20" class="icon" />
                        <span :title="$t('backup.repository')">{{
                          backup.repoName
                        }}</span>
                      </div>
                    </div>
                    <!-- <div class="row icon-and-text"> ////
                      <NsSvg :svg="Time20" class="icon" />
                      <span :title="$t('backup.frequency')">{{
                        $t("backup." + backup.schedule)
                      }}</span>
                    </div> -->
                    <div class="row">
                      <cv-tag
                        v-if="backup.enabled"
                        kind="green"
                        :label="$t('common.enabled')"
                        :title="$t('backup.backup_enabled')"
                      ></cv-tag>
                      <cv-tag
                        v-else
                        kind="red"
                        :label="$t('common.disabled')"
                        :title="$t('backup.backup_disabled')"
                      ></cv-tag>
                    </div>
                    <div class="row actions">
                      <NsButton
                        kind="ghost"
                        :icon="ZoomIn20"
                        @click="showBackupDetailsModal(backup)"
                        >{{ $t("common.details") }}
                      </NsButton>
                    </div>
                  </div>
                </template>
              </NsInfoCard>
            </div>
          </div>
        </template>
      </template>
    </div>
    <!-- delete repository modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteRepoModal"
      :name="currentRepo.password"
      :title="$t('backup.delete_backup_repository')"
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('backup.delete_repository_confirm', {
          name: currentRepo.name,
        })
      "
      :typeToConfirm="
        $t('backup.type_repo_password_to_confirm', {
          password: currentRepo.password,
        })
      "
      @hide="hideDeleteRepoModal"
      @confirmDelete="deleteRepo(currentRepo)"
    >
      <template slot="explanation">
        <p
          class="mg-top-sm"
          v-html="$t('backup.delete_repo_explanation_1')"
        ></p>
        <p class="mg-top-sm">
          <strong>
            {{ $t("backup.delete_repo_explanation_2") }}
          </strong>
        </p>
        <p class="mg-top-md">
          {{ $t("backup.repository_password") }}
          <NsCodeSnippet
            :copyTooltip="$t('common.copy_to_clipboard')"
            :copy-feedback="$t('common.copied_to_clipboard')"
            :feedback-aria-label="$t('common.copied_to_clipboard')"
            :wrap-text="true"
            :moreText="$t('common.show_more')"
            :lessText="$t('common.show_less')"
            light
            hideExpandButton
            >{{ currentRepo.password }}</NsCodeSnippet
          >
        </p>
      </template>
    </NsDangerDeleteModal>
    <!-- add repository modal -->
    <AddRepositoryModal
      :isShown="q.isShownAddRepoModal"
      @hide="hideAddRepoModal"
      @repoCreated="listBackupRepositories"
    />
    <!-- create backup modal -->
    <CreateBackupModal
      :isShown="q.isShownCreateBackupModal"
      :repositories="repositories"
      :instanceSelection="instanceSelection"
      :instancesNotBackedUp="unconfiguredInstances"
      @hide="hideCreateBackupModal"
      @backupCreated="listBackupRepositories"
    />
    <!-- delete repository modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteBackupModal"
      :name="currentBackup.name"
      :title="$t('backup.delete_scheduled_backup')"
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('backup.delete_scheduled_backup_confirm', {
          name: currentBackup.name,
        })
      "
      @hide="hideDeleteBackupModal"
      @confirmDelete="deleteBackup(currentBackup)"
    >
      <template slot="explanation">
        <p
          v-if="currentBackup.instances.length"
          class="mg-top-sm"
          v-html="
            $tc(
              'backup.delete_backup_explanation_1',
              currentBackup.instances.length,
              {
                numInstances: currentBackup.instances.length,
                instanceName: currentBackup.instances[0].ui_name
                  ? currentBackup.instances[0].ui_name
                  : currentBackup.instances[0].module_id,
              }
            )
          "
        ></p>
        <p
          class="mg-top-sm"
          v-html="$t('backup.delete_backup_explanation_2')"
        ></p>
      </template>
    </NsDangerDeleteModal>
    <!-- repo details modal -->
    <RepoDetailsModal
      :isShown="isShownRepoDetailsModal"
      :repository="currentRepo"
      @hide="hideRepoDetailsModal"
    />
    <!-- edit repo modal -->
    <EditRepositoryModal
      :isShown="isShownEditRepoModal"
      :repository="currentRepo"
      @hide="hideEditRepoModal"
    />
    <!-- backup details modal -->
    <BackupDetailsModal
      :isShown="isShownBackupDetailsModal"
      :backup="currentBackup"
      @hide="hideBackupDetailsModal"
    />
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
} from "@nethserver/ns8-ui-lib";
import AddRepositoryModal from "@/components/backup/AddRepositoryModal";
import CreateBackupModal from "@/components/backup/CreateBackupModal";
import RepoDetailsModal from "@/components/backup/RepoDetailsModal";
import BackupDetailsModal from "@/components/backup/BackupDetailsModal";
import EditRepositoryModal from "@/components/backup/EditRepositoryModal";
import to from "await-to-js";

export default {
  name: "Backup",
  components: {
    AddRepositoryModal,
    CreateBackupModal,
    RepoDetailsModal,
    BackupDetailsModal,
    EditRepositoryModal,
  },
  mixins: [TaskService, UtilService, IconService, QueryParamService],
  pageTitle() {
    return this.$t("backup.title");
  },
  data() {
    return {
      q: {
        isShownAddRepoModal: false,
        isShownCreateBackupModal: false,
      },
      isShownDeleteRepoModal: false,
      isShownRepoDetailsModal: false,
      isShownEditRepoModal: false,
      isShownEditBackupModal: false,
      isShownDeleteBackupModal: false,
      isShownBackupDetailsModal: false,
      repositories: [],
      backups: [],
      unconfiguredInstances: [],
      instanceSelection: "",
      currentRepo: {
        name: "",
        password: "",
      },
      currentBackup: {
        name: "",
        instances: [],
      },
      loading: {
        listBackupRepositories: true,
        listBackups: true,
      },
      error: {
        listBackupRepositories: "",
        listBackups: "",
        removeBackupRepository: "",
        removeBackup: "",
        runBackup: "",
      },
    };
  },
  computed: {
    disabledBackups() {
      return this.backups.filter((b) => !b.enabled);
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
    this.listBackupRepositories();
  },
  methods: {
    async listBackupRepositories() {
      this.loading.listBackupRepositories = true;
      this.error.listBackupRepositories = "";
      const taskAction = "list-backup-repositories";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.listBackupRepositoriesCompleted
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
        this.error.listBackupRepositories = this.getErrorMessage(err);
        return;
      }
    },
    listBackupRepositoriesCompleted(taskContext, taskResult) {
      console.log("listBackupRepositoriesCompleted", taskResult.output); ////

      this.repositories = taskResult.output;
      this.loading.listBackupRepositories = false;
      this.listBackups();
    },
    async listBackups() {
      this.loading.listBackups = true;
      this.error.listBackups = "";
      const taskAction = "list-backups";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.listBackupsCompleted);

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
        this.error.listBackups = this.getErrorMessage(err);
        return;
      }
    },
    listBackupsCompleted(taskContext, taskResult) {
      this.unconfiguredInstances = taskResult.output.unconfigured_instances;
      let backups = taskResult.output.backups;

      console.log("listBackupsCompleted", backups, this.unconfiguredInstances); ////

      // repository name

      for (const backup of backups) {
        const repo = this.repositories.find((r) => r.id == backup.repository);

        if (repo) {
          backup.repoName = repo.name;
        }
      }
      this.backups = backups;

      this.loading.listBackups = false;
    },
    showAddRepoModal() {
      this.q.isShownAddRepoModal = true;
    },
    hideAddRepoModal() {
      this.q.isShownAddRepoModal = false;
    },
    showDeleteRepoModal(repo) {
      this.currentRepo = repo;
      this.isShownDeleteRepoModal = true;
    },
    hideDeleteRepoModal() {
      this.isShownDeleteRepoModal = false;
    },
    showRepoDetailsModal(repo) {
      this.currentRepo = repo;
      this.isShownRepoDetailsModal = true;
    },
    hideRepoDetailsModal() {
      this.isShownRepoDetailsModal = false;
    },
    showEditRepoModal(repo) {
      this.currentRepo = repo;
      this.isShownEditRepoModal = true;
    },
    hideEditRepoModal() {
      this.isShownEditRepoModal = false;
    },
    showCreateBackupModal(instanceSelection = "") {
      this.instanceSelection = instanceSelection;
      this.q.isShownCreateBackupModal = true;
    },
    hideCreateBackupModal() {
      this.q.isShownCreateBackupModal = false;
    },
    showEditBackupModal(backup) {
      console.log("showEditBackupModal", backup); ////

      this.currentBackup = backup;
      this.isShownEditBackupModal = true;
    },
    showDeleteBackupModal(backup) {
      this.currentBackup = backup;
      this.isShownDeleteBackupModal = true;
    },
    hideDeleteBackupModal() {
      this.isShownDeleteBackupModal = false;
    },
    showBackupDetailsModal(backup) {
      this.currentBackup = backup;
      this.isShownBackupDetailsModal = true;
    },
    hideBackupDetailsModal() {
      this.isShownBackupDetailsModal = false;
    },
    async deleteRepo(repo) {
      this.error.removeBackupRepository = "";
      const taskAction = "remove-backup-repository";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.removeBackupRepositoryCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            id: repo.id,
            password: repo.password,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeBackupRepository = this.getErrorMessage(err);
        return;
      }

      this.isShownDeleteRepoModal = false;
    },
    removeBackupRepositoryCompleted() {
      // reload backup configuration
      this.listBackupRepositories();
    },
    async deleteBackup(backup) {
      this.error.removeBackup = "";
      const taskAction = "remove-backup";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.removeBackupCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            id: backup.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.removeBackup = this.getErrorMessage(err);
        return;
      }

      this.isShownDeleteBackupModal = false;
    },
    removeBackupCompleted() {
      // reload backup configuration
      this.listBackupRepositories();
    },
    async runBackup(backup) {
      this.error.runBackup = "";
      const taskAction = "run-backup";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.runBackupCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            id: backup.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("common.processing"),
            backupName: backup.name,
            completion: {
              i18nString: "backup.backup_completed_successfully",
              extraTextParams: ["backupName"],
            },
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.runBackup = this.getErrorMessage(err);
        return;
      }
    },
    runBackupCompleted() {
      this.listBackups();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.card-content {
  // let flex children use all horizontal space
  flex-grow: 1;
}

.card-content .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.card-content .row:last-child {
  margin-bottom: 0;
}

.instance-to-repo {
  display: flex;
  align-items: center;
  justify-content: center;
}

.arrow-right {
  margin-left: $spacing-05;
  margin-right: $spacing-05;
  flex-grow: 1;
  flex-basis: 0;
}

.instance-to-repo .icon {
  margin-right: $spacing-02;
}

.backup-source {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  flex-grow: 6;
  flex-basis: 0;
}

.backup-destination {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  flex-grow: 6;
  flex-basis: 0;
}
</style>
