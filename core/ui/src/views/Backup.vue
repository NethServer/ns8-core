<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <cv-grid fullWidth>
      <cv-row>
        <cv-column class="page-title">
          <h2>
            {{ $t("backup.title") }}
            <cv-interactive-tooltip
              alignment="start"
              direction="right"
              class="info"
            >
              <template slot="trigger">
                <Information16 />
              </template>
              <template slot="content">
                <div class="mg-bottom-sm">
                  {{ $t("backup.backup_page_tooltip_1") }}
                </div>
                <div class="mg-bottom-sm">
                  {{ $t("backup.backup_page_tooltip_2") }}
                </div>
                <div class="mg-bottom-sm">
                  {{
                    $t("backup.backup_page_tooltip_8", {
                      productName: $root.config.PRODUCT_NAME,
                    })
                  }}
                </div>
                <i18n path="backup.backup_page_tooltip_manual" tag="p">
                  <template v-slot:manualLink>
                    <cv-link
                      href="https://docs.nethserver.org/projects/ns8/en/latest/backup.html"
                      target="_blank"
                      rel="noreferrer"
                      class="inline"
                    >
                      {{ $t("backup.go_to_backup_manual_page") }}
                    </cv-link>
                  </template>
                </i18n>
              </template>
            </cv-interactive-tooltip>
          </h2>
        </cv-column>
      </cv-row>
      <!-- cluster configuration -->
      <cv-row>
        <cv-column>
          <h4 class="mg-bottom-md">
            {{ $t("backup.cluster_configuration") }}
            <cv-interactive-tooltip
              alignment="start"
              direction="right"
              class="info"
            >
              <template slot="trigger">
                <Information16 />
              </template>
              <template slot="content">
                {{ $t("backup.cluster_configuration_tooltip") }}
              </template>
            </cv-interactive-tooltip>
          </h4>
        </cv-column>
      </cv-row>
      <cv-row class="mg-bottom-lg">
        <cv-column
          v-if="!isSetClusterBackupPassword && !loading.listBackupRepositories"
        >
          <cv-tile kind="standard" :light="true">
            <NsEmptyState :title="$t('backup.cluster_backup_password_warning')">
              <template #pictogram>
                <ExclamationMarkPictogram />
              </template>
              <template #description>
                <div class="mg-bottom-lg">
                  <p>
                    {{
                      $t("backup.cluster_backup_password_warning_description_1")
                    }}
                  </p>
                  <p>
                    {{
                      $t("backup.cluster_backup_password_warning_description_2")
                    }}
                  </p>
                </div>
                <div>
                  <NsButton
                    kind="primary"
                    :icon="Password20"
                    @click="showBackupPasswordModal()"
                    >{{ $t("backup.set_cluster_backup_password") }}
                  </NsButton>
                </div>
              </template>
            </NsEmptyState>
          </cv-tile>
        </cv-column>
        <cv-column
          v-if="isSetClusterBackupPassword && !loading.listBackupRepositories"
        >
          <NsButton
            kind="secondary"
            :icon="Download20"
            @click="downloadClusterConfigurationBackup()"
            :disabled="loading.downloadClusterBackup"
            :loading="loading.downloadClusterBackup"
            >{{ $t("backup.download_cluster_configuration_backup") }}
          </NsButton>
          <NsButton
            kind="tertiary"
            :icon="Password20"
            @click="showBackupPasswordModal()"
            class="mg-left-md"
            :disabled="loading.downloadClusterBackup"
            >{{ $t("backup.change_cluster_backup_password") }}
          </NsButton>
        </cv-column>
      </cv-row>
      <cv-row v-if="isSetClusterBackupPassword">
        <cv-column>
          <h4 class="mg-bottom-md">
            {{ $t("backup.apps") }}
            <cv-interactive-tooltip
              alignment="start"
              direction="right"
              class="info"
            >
              <template slot="trigger">
                <Information16 />
              </template>
              <template slot="content">
                {{ $t("backup.apps_configuration_tooltip") }}
              </template>
            </cv-interactive-tooltip>
          </h4>
        </cv-column>
      </cv-row>
      <template v-if="loading.listBackupRepositories || loading.listBackups">
        <!-- repositories skeleton -->
        <cv-row>
          <cv-column>
            <cv-skeleton-text heading width="40%"></cv-skeleton-text>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <!-- skeleton card grid -->
            <div
              class="card-grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 3xl:grid-cols-4"
            >
              <cv-tile v-for="index in 2" :key="index" light>
                <cv-skeleton-text
                  :paragraph="true"
                  :line-count="8"
                ></cv-skeleton-text>
              </cv-tile>
            </div>
          </cv-column>
        </cv-row>
        <!-- backups skeleton -->
        <cv-row>
          <cv-column>
            <cv-skeleton-text heading width="40%"></cv-skeleton-text>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <!-- skeleton card grid -->
            <div
              class="card-grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 3xl:grid-cols-4"
            >
              <cv-tile v-for="index in 2" :key="index" light>
                <cv-skeleton-text
                  :paragraph="true"
                  :line-count="8"
                ></cv-skeleton-text>
              </cv-tile>
            </div>
          </cv-column>
        </cv-row>
      </template>
      <!-- empty state repositories -->
      <cv-row v-else-if="isSetClusterBackupPassword && !repositories.length">
        <cv-column>
          <cv-tile kind="standard" :light="true">
            <NsEmptyState :title="$t('backup.no_backup_repository')">
              <template #pictogram>
                <HardDrivePictogram />
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
          </cv-tile>
        </cv-column>
      </cv-row>
      <template v-else-if="isSetClusterBackupPassword">
        <!-- errored backups -->
        <cv-row v-if="!loading.listBackups && erroredBackups.length">
          <cv-column>
            <NsInlineNotification
              kind="error"
              :title="$t('backup.some_backups_failed')"
              :description="$t('backup.errored_backups_description')"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <!-- disabled backups warning -->
        <cv-row v-if="!loading.listBackups && disabledBackups.length">
          <cv-column>
            <NsInlineNotification
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
          </cv-column>
        </cv-row>
        <!-- unconfigured instances warning -->
        <cv-row v-if="!loading.listBackups && unconfiguredInstances.length">
          <cv-column>
            <NsInlineNotification
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
              @action="showCreateOrEditBackupModal('notBackedUp')"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <!-- repositories -->
        <cv-row>
          <cv-column>
            <h6 class="mg-bottom-md">
              {{ $t("backup.repositories") }}
              <cv-interactive-tooltip
                alignment="start"
                direction="right"
                class="info"
              >
                <template slot="trigger">
                  <Information16 />
                </template>
                <template slot="content">
                  {{ $t("backup.repositories_tooltip") }}
                </template>
              </cv-interactive-tooltip>
            </h6>
          </cv-column>
        </cv-row>
        <cv-row v-if="error.listBackupRepositories">
          <cv-column>
            <NsInlineNotification
              kind="error"
              :title="$t('action.list-backup-repositories')"
              :description="error.listBackupRepositories"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row v-if="error.removeBackupRepository">
          <cv-column>
            <NsInlineNotification
              kind="error"
              :title="$t('action.remove-backup-repository')"
              :description="error.removeBackupRepository"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row class="toolbar">
          <cv-column>
            <NsButton kind="secondary" :icon="Add20" @click="showAddRepoModal()"
              >{{ $t("backup.add_repository") }}
            </NsButton>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <!-- card grid -->
            <div
              class="card-grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 3xl:grid-cols-4"
            >
              <NsInfoCard
                v-for="repo in repositories"
                :key="repo.id"
                light
                :title="repo.name"
                :icon="DataBase32"
              >
                <template #menu>
                  <cv-overflow-menu
                    :flip-menu="true"
                    tip-position="top"
                    tip-alignment="end"
                    class="top-right-overflow-menu"
                  >
                    <cv-overflow-menu-item @click="showEditRepoModal(repo)">
                      <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
                    </cv-overflow-menu-item>
                    <NsMenuDivider />
                    <cv-overflow-menu-item
                      danger
                      @click="showDeleteRepoModal(repo)"
                    >
                      <NsMenuItem
                        :icon="TrashCan20"
                        :label="$t('common.delete')"
                      />
                    </cv-overflow-menu-item>
                  </cv-overflow-menu>
                </template>
                <template #content>
                  <div class="card-content-backup">
                    <div class="row">
                      {{ $t("backup." + repo.provider) }}
                    </div>
                    <div class="row">
                      {{ repo.url }}
                    </div>
                    <div class="row actions">
                      <NsButton
                        kind="ghost"
                        :icon="ArrowRight20"
                        @click="showRepoDetailsModal(repo)"
                        >{{ $t("common.see_details") }}
                      </NsButton>
                    </div>
                  </div>
                </template>
              </NsInfoCard>
            </div>
          </cv-column>
        </cv-row>
        <!-- backups -->
        <cv-row>
          <cv-column>
            <h6 class="mg-bottom-md">{{ $t("backup.schedules") }}</h6>
          </cv-column>
        </cv-row>
        <cv-row v-if="error.listBackups">
          <cv-column>
            <NsInlineNotification
              kind="error"
              :title="$t('action.list-backups')"
              :description="error.listBackups"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row v-if="error.removeBackup">
          <cv-column>
            <NsInlineNotification
              kind="error"
              :title="$t('action.remove-backup')"
              :description="error.removeBackup"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row v-if="error.runBackup">
          <cv-column>
            <NsInlineNotification
              kind="error"
              :title="$t('action.run-backup')"
              :description="error.runBackup"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row v-if="error.alterBackup">
          <cv-column>
            <NsInlineNotification
              kind="error"
              :title="$t('action.alter-backup')"
              :description="error.alterBackup"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <!-- empty state backups -->
        <cv-row v-if="!backups.length">
          <cv-column>
            <cv-tile light>
              <NsEmptyState :title="$t('backup.no_backup_scheduled')">
                <template #description>
                  <NsButton
                    kind="primary"
                    :icon="Add20"
                    @click="showCreateOrEditBackupModal('')"
                    class="empty-state-button-no-description"
                    >{{ $t("backup.schedule_backup") }}
                  </NsButton>
                </template>
              </NsEmptyState>
            </cv-tile>
          </cv-column>
        </cv-row>
        <template v-else>
          <cv-row class="toolbar">
            <cv-column>
              <NsButton
                kind="secondary"
                :icon="Add20"
                @click="showCreateOrEditBackupModal('')"
                >{{ $t("backup.schedule_backup") }}
              </NsButton>
            </cv-column>
          </cv-row>
          <cv-row>
            <cv-column>
              <!-- card grid -->
              <div
                class="card-grid grid-cols-1 md:grid-cols-2 2xl:grid-cols-3 3xl:grid-cols-4"
              >
                <NsInfoCard
                  v-for="backup in backups"
                  :key="backup.id"
                  light
                  :title="backup.name"
                  :icon="Save32"
                >
                  <template #menu>
                    <cv-overflow-menu
                      :flip-menu="true"
                      tip-position="top"
                      tip-alignment="end"
                      class="top-right-overflow-menu"
                    >
                      <cv-overflow-menu-item
                        @click="runBackup(backup)"
                        :disabled="!backup.enabled"
                      >
                        <NsMenuItem
                          :icon="Save20"
                          :label="$t('backup.run_backup_now')"
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        @click="toggleBackupStatus(backup)"
                        :disabled="loading.alterBackup"
                      >
                        <NsMenuItem
                          :icon="Power20"
                          :label="
                            backup.enabled
                              ? $t('common.disable')
                              : $t('common.enable')
                          "
                        />
                      </cv-overflow-menu-item>
                      <cv-overflow-menu-item
                        @click="showEditBackupModal(backup)"
                      >
                        <NsMenuItem :icon="Edit20" :label="$t('common.edit')" />
                      </cv-overflow-menu-item>
                      <NsMenuDivider />
                      <cv-overflow-menu-item
                        danger
                        @click="showDeleteBackupModal(backup)"
                      >
                        <NsMenuItem
                          :icon="TrashCan20"
                          :label="$t('common.delete')"
                        />
                      </cv-overflow-menu-item>
                    </cv-overflow-menu>
                  </template>
                  <template #content>
                    <div class="card-content-backup">
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
                      <div class="row">
                        <cv-tag
                          v-if="backup.enabled"
                          kind="green"
                          :label="$t('common.enabled')"
                          :title="$t('backup.backup_enabled')"
                        ></cv-tag>
                        <cv-tag
                          v-else
                          kind="high-contrast"
                          :label="$t('common.disabled')"
                          :title="$t('backup.backup_disabled')"
                        ></cv-tag>
                      </div>
                      <div
                        v-if="backup.errorInstances.length"
                        class="row icon-and-text"
                      >
                        <NsSvg :svg="ErrorFilled16" class="icon ns-error" />
                        <span>
                          {{
                            $tc(
                              "backup.backup_of_instances_failed_c",
                              backup.errorInstances.length,
                              {
                                instance: backup.errorInstances[0].ui_name
                                  ? backup.errorInstances[0].ui_name
                                  : backup.errorInstances[0].module_id,
                                num: backup.errorInstances.length,
                              }
                            )
                          }}
                        </span>
                      </div>
                      <div class="row actions">
                        <NsButton
                          kind="ghost"
                          :icon="ArrowRight20"
                          @click="showBackupDetailsModal(backup)"
                          >{{ $t("common.see_details") }}
                        </NsButton>
                      </div>
                    </div>
                  </template>
                </NsInfoCard>
              </div>
            </cv-column>
          </cv-row>
        </template>
        <cv-row>
          <cv-column>
            <h6 class="mg-bottom-md">
              {{ $t("backup.restore") }}
              <cv-interactive-tooltip
                alignment="start"
                direction="right"
                class="info"
              >
                <template slot="trigger">
                  <Information16 />
                </template>
                <template slot="content">
                  {{ $t("backup.restore_tooltip") }}
                </template>
              </cv-interactive-tooltip>
            </h6>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <NsButton
              kind="secondary"
              :icon="Reset20"
              @click="showRestoreModal()"
              :disabled="!repositories.length"
              class="mg-bottom-xlg"
              >{{ $t("backup.restore_app") }}
            </NsButton>
          </cv-column>
        </cv-row>
      </template>
    </cv-grid>
    <!-- delete repository modal -->
    <NsDangerDeleteModal
      :isShown="isShownDeleteRepoModal"
      :name="currentRepo.name"
      :title="$t('backup.delete_backup_repository')"
      :warning="$t('common.please_read_carefully')"
      :description="
        $t('backup.delete_repository_confirm', {
          name: currentRepo.name,
        })
      "
      :typeToConfirm="$t('common.type_to_confirm', { name: currentRepo.name })"
      @hide="hideDeleteRepoModal"
      @confirmDelete="deleteRepo(currentRepo)"
    >
      <template slot="explanation">
        <p class="mg-top-sm">{{ $t("backup.delete_repo_explanation_1") }}</p>
      </template>
    </NsDangerDeleteModal>
    <!-- add repository modal -->
    <AddRepositoryModal
      :isShown="q.isShownAddRepoModal"
      :repositories="repositories"
      @hide="hideAddRepoModal"
      @repoCreated="listBackupRepositories"
    />
    <!-- create backup modal -->
    <CreateOrEditBackupModal
      :isShown="isShownCreateOrEditBackupModal"
      :isEditing="isEditingBackup"
      :repositories="repositories"
      :instanceSelection="instanceSelection"
      :instancesNotBackedUp="unconfiguredInstances"
      :backup="currentBackup"
      :backups="backups"
      @hide="hideCreateOrEditBackupModal"
      @backupCreated="onBackupCreated"
      @backupAltered="onBackupAltered"
    />
    <!-- delete backup modal -->
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
      :typeToConfirm="
        $t('common.type_to_confirm', { name: currentBackup.name })
      "
      @hide="hideDeleteBackupModal"
      @confirmDelete="deleteBackup(currentBackup)"
    >
      <template slot="explanation">
        <p v-if="currentBackup.instances.length" class="mg-top-sm">
          {{
            $tc(
              "backup.delete_backup_explanation_1",
              currentBackup.instances.length,
              {
                numInstances: currentBackup.instances.length,
                instanceName: currentBackup.instances[0].ui_name
                  ? currentBackup.instances[0].ui_name
                  : currentBackup.instances[0].module_id,
              }
            )
          }}
        </p>
        <p class="mg-top-sm">{{ $t("backup.delete_backup_explanation_2") }}</p>
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
      @repoAltered="listBackupRepositories"
    />
    <!-- backup details modal -->
    <BackupDetailsModal
      :isShown="isShownBackupDetailsModal"
      :backup="currentBackup"
      @hide="hideBackupDetailsModal"
    />
    <!-- restore modal -->
    <RestoreSingleInstanceModal
      :isShown="isShownRestoreModal"
      @hide="hideRestoreModal"
    />
    <!-- download backup modal -->
    <BackupPasswordModal
      :isShown="isShownBackupPasswordModal"
      @hide="hideBackupPasswordModal"
      @password-set.once="reloadBackupRepositories"
    />
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  TaskService,
  IconService,
  DateTimeService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import AddRepositoryModal from "@/components/backup/AddRepositoryModal";
import CreateOrEditBackupModal from "@/components/backup/CreateOrEditBackupModal";
import RepoDetailsModal from "@/components/backup/RepoDetailsModal";
import BackupDetailsModal from "@/components/backup/BackupDetailsModal";
import EditRepositoryModal from "@/components/backup/EditRepositoryModal";
import RestoreSingleInstanceModal from "@/components/backup/RestoreSingleInstanceModal";
import BackupPasswordModal from "@/components/backup/BackupPasswordModal";
import to from "await-to-js";
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  name: "Backup",
  components: {
    AddRepositoryModal,
    CreateOrEditBackupModal,
    RepoDetailsModal,
    BackupDetailsModal,
    EditRepositoryModal,
    RestoreSingleInstanceModal,
    BackupPasswordModal,
    Information16,
  },
  mixins: [
    TaskService,
    UtilService,
    IconService,
    QueryParamService,
    DateTimeService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("backup.title");
  },
  data() {
    return {
      q: {
        isShownAddRepoModal: false,
      },
      isShownCreateOrEditBackupModal: false,
      isShownDeleteRepoModal: false,
      isShownRepoDetailsModal: false,
      isShownEditRepoModal: false,
      isShownDeleteBackupModal: false,
      isShownBackupDetailsModal: false,
      isShownRestoreModal: false,
      isShownBackupPasswordModal: false,
      repositories: [],
      backups: [],
      unconfiguredInstances: [],
      instanceSelection: "",
      isEditingBackup: false,
      isSetClusterBackupPassword: false,
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
        alterBackup: false,
        downloadClusterBackup: false,
        passwordClusterBackup: false,
      },
      error: {
        listBackupRepositories: "",
        listBackups: "",
        removeBackupRepository: "",
        removeBackup: "",
        runBackup: "",
        alterBackup: "",
        downloadClusterBackup: "",
        passwordClusterBackup: "",
      },
    };
  },
  computed: {
    disabledBackups() {
      return this.backups.filter((b) => !b.enabled);
    },
    erroredBackups() {
      return this.backups.filter((b) => b.errorInstances.length);
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
      let repositories = taskResult.output.repositories.sort(
        this.sortByProperty("name")
      );
      this.repositories = repositories;
      this.loading.listBackupRepositories = false;
      this.isSetClusterBackupPassword = taskResult.output.password_exists;
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
      backups.sort(this.sortByProperty("name"));

      for (const backup of backups) {
        // repository name
        const repo = this.repositories.find((r) => r.id == backup.repository);
        if (repo) {
          backup.repoName = repo.name;
        }

        // schedule
        backup.scheduleExpression = backup.schedule;
        backup.schedule = backup.schedule_hint;

        // error instances
        backup.errorInstances = backup.instances.filter(
          (i) => i.status && i.status.success == false
        );
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
    showCreateOrEditBackupModal(instanceSelection = "") {
      this.instanceSelection = instanceSelection;
      this.isEditingBackup = false;
      this.isShownCreateOrEditBackupModal = true;
    },
    hideCreateOrEditBackupModal() {
      this.isShownCreateOrEditBackupModal = false;
    },
    showEditBackupModal(backup) {
      this.currentBackup = backup;
      this.isEditingBackup = true;
      this.isShownCreateOrEditBackupModal = true;
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
    showRestoreModal() {
      this.isShownRestoreModal = true;
    },
    hideRestoreModal() {
      this.isShownRestoreModal = false;
    },
    showBackupPasswordModal() {
      this.isShownBackupPasswordModal = true;
    },
    hideBackupPasswordModal() {
      this.isShownBackupPasswordModal = false;
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
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.runBackupCompleted);

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.runBackupAborted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            id: backup.id,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$t("backup.running_backup_name", {
              backupName: backup.name,
            }),
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
    runBackupAborted() {
      this.listBackups();
    },
    onBackupCreated(runBackupOnFinish, createdBackup) {
      this.listBackupRepositories();

      if (runBackupOnFinish) {
        // run created backup now
        this.runBackup(createdBackup);
      }
    },
    onBackupAltered(runBackupOnFinish, createdBackup) {
      this.listBackupRepositories();

      if (runBackupOnFinish) {
        // run created backup now
        this.runBackup(createdBackup);
      }
    },
    async toggleBackupStatus(backup) {
      this.error.alterBackup = "";
      this.loading.alterBackup = true;
      const taskAction = "alter-backup";
      const taskDescription = backup.enabled
        ? this.$t("backup.disabling_backup", { backupName: backup.name })
        : this.$t("backup.enabling_backup", { backupName: backup.name });

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.alterBackupCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            id: backup.id,
            name: backup.name,
            schedule: backup.scheduleExpression,
            schedule_hint: backup.schedule,
            retention: backup.retention,
            instances: backup.instances.map((i) => i.module_id),
            enabled: !backup.enabled,
          },
          extra: {
            title: this.$t("action." + taskAction),
            description: taskDescription,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.alterBackup = this.getErrorMessage(err);
        return;
      }
    },
    alterBackupCompleted() {
      this.loading.alterBackup = false;
      this.listBackups();
    },
    async downloadClusterConfigurationBackup() {
      this.loading.downloadClusterBackup = true;
      const taskAction = "download-cluster-backup";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.downloadClusterBackupAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.downloadClusterBackupCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: { password: this.clusterPassword },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.downloadClusterBackup = this.getErrorMessage(err);
        return;
      }
    },
    downloadClusterBackupAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.downloadClusterBackup = false;
    },
    downloadClusterBackupCompleted(taskContext, taskResult) {
      this.loading.downloadClusterBackup = false;
      const downloadUrl = `${window.location.protocol}//${window.location.hostname}/cluster-admin/backup/${taskResult.output.path}`;

      const fileName =
        "cluster-configuration-backup " +
        this.formatDate(new Date(), "yyyy-MM-dd HH.mm") +
        ".json.gz.gpg";

      this.axios({
        url: downloadUrl,
        method: "GET",
        responseType: "blob",
      }).then((response) => {
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute("download", fileName);
        document.body.appendChild(link);
        link.click();
      });
    },
    reloadBackupRepositories() {
      this.isSetClusterBackupPassword = true;
      this.listBackupRepositories();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.card-content-backup {
  // let flex children use all horizontal space
  flex-grow: 1;
}

.card-content-backup .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.card-content-backup .row:last-child {
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

.inline {
  display: inline;
}
</style>
