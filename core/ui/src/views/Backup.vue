<template>
  <div>
    <div class="bx--grid bx--grid--full-width">
      <div class="bx--row">
        <div class="bx--col-lg-16 page-title">
          <h2>{{ $t("backup.title") }}</h2>
        </div>
      </div>
      <div
        v-if="loading.listBackupRepositories || loading.listBackups"
        class="bx--row"
      >
        <div v-for="index in 4" :key="index" class="bx--col-md-4 bx--col-max-8">
          <cv-tile light>
            <cv-skeleton-text
              :paragraph="true"
              :line-count="8"
            ></cv-skeleton-text>
          </cv-tile>
        </div>
      </div>
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
                >{{ $t("backup.add_backup_repository") }}
              </NsButton>
            </template>
          </NsEmptyState>
        </div>
      </div>
      <template v-else>
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
          </div>
        </div>
        <div class="bx--row">
          <div class="bx--col">
            <NsButton kind="secondary" :icon="Add20" @click="showAddRepoModal()"
              >{{ $t("backup.add_backup_repository") }}
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
                  <cv-overflow-menu-item @click="showEditRepoNameModal(repo)">
                    {{ $t("backup.edit_repository_name") }}
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
                  <div class="row icon-and-text">
                    <NsSvg :svg="Link20" class="icon" />
                    <span>{{ repo.url }}</span>
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
            <h4 class="mg-bottom-md">{{ $t("backup.backups") }}</h4>
          </div>
        </div>
        <div class="bx--row">
          <div class="bx--col">
            <NsInlineNotification
              v-if="error.listBackups"
              kind="error"
              :title="$t('action.list-backups')"
              :description="error.listBackups"
              :showCloseButton="false"
            />
          </div>
        </div>
        <!-- empty state backups -->
        <div v-if="!backups.length" class="bx--row">
          <div class="bx--col">
            <NsEmptyState :title="$t('backup.no_backup_configured')">
              <template #description>
                <NsButton
                  kind="primary"
                  :icon="Add20"
                  @click="showCreateBackupModal()"
                  class="empty-state-button-no-description"
                  >{{ $t("backup.create_backup") }}
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
                @click="showCreateBackupModal()"
                >{{ $t("backup.create_backup") }}
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
                    <div class="row icon-and-text backup-repo">
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
                          backup.instances.length + " " + $t("backup.instances")
                        }}
                      </span>
                    </div>
                    <div class="row icon-and-text backup-repo">
                      <NsSvg :svg="ArrowDown32" class="icon" />
                    </div>
                    <div class="row icon-and-text">
                      <NsSvg :svg="DataBase20" class="icon" />
                      <span>{{ backup.repoName }}</span>
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
    <AddBackupRepositoryModal
      :isShown="q.isShownAddRepoModal"
      @hide="hideAddRepoModal"
      @repoCreated="listBackupRepositories"
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
import AddBackupRepositoryModal from "@/components/AddBackupRepositoryModal";

export default {
  name: "Backup",
  components: { AddBackupRepositoryModal },
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
      isShownEditRepoNameModal: false,
      isShownEditBackupModal: false,
      isShownDeleteBackupModal: false,
      isShownBackupDetailsModal: false,
      repositories: [],
      backups: [],
      currentRepo: {
        name: "",
        password: "",
      },
      loading: {
        listBackupRepositories: true,
        listBackups: true,
      },
      error: {
        listBackupRepositories: "",
        listBackups: "",
      },
    };
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
    listBackupRepositories() {
      console.log("listBackupRepositories"); ////

      //// remove mock
      // this.repositories = [];
      this.repositories = [
        {
          id: "48ce000a-79b7-5fe6-8558-177fd70c27b4",
          name: "BackBlaze repo1",
          provider: "backblaze",
          url: "b2:backupex1",
          password:
            "d59a90ec7ad2b2967257a7a308c82c96ac006efd138254bc1e58c8ea07c18400",
          parameters: {
            b2_account_id: "xxxxxxxxxxxxxx",
            b2_account_key: "yyyyyyyyyyyyyyyyyyyyyy",
          },
        },
        {
          id: "98ce000a-79b7-5fe6-8558-177fd70c27b4",
          name: "S3 repo",
          provider: "s3",
          url: "s3:backupex1",
          password:
            "d59a90ec7ad2b2967257a7a308c82c96ac006efd138254bc1e58c8ea07c18400",
          parameters: {
            s3_account_id: "xxxxxxxxxxxxxx",
            s3_account_key: "yyyyyyyyyyyyyyyyyyyyyy",
          },
        },
      ];

      this.loading.listBackupRepositories = false;
      this.listBackups();
    },
    listBackups() {
      console.log("listBackups"); ////

      //// remove mock
      let backups = [];
      // let backups = [
      //   {
      //     id: 1,
      //     name: "Dokuwiki backup",
      //     repository: "48ce000a-79b7-5fe6-8558-177fd70c27b4",
      //     schedule: "daily",
      //     retention: "7d",
      //     instances: [
      //       {
      //         module_id: "dokuwiki1",
      //         ui_name: "",
      //         repository_path: "dokuwiki1@2f72561e-89b2-4cdc-b4e4-425ca23bbec9",
      //         status: null,
      //       },
      //     ],
      //   },
      //   {
      //     id: 2,
      //     name: "Nextcloud backup",
      //     repository: "98ce000a-79b7-5fe6-8558-177fd70c27b4",
      //     schedule: "daily",
      //     retention: "7d",
      //     instances: [
      //       {
      //         module_id: "nextcloud1",
      //         ui_name: "",
      //         repository_path:
      //           "nextcloud1@2f72561e-89b2-4cdc-b4e4-425ca23bbec9",
      //         status: null,
      //       },
      //       {
      //         module_id: "nextcloud2",
      //         ui_name: "",
      //         repository_path:
      //           "nextcloud2@2f72561e-89b2-4cdc-b4e4-425ca23bbec9",
      //         status: null,
      //       },
      //     ],
      //   },
      // ];

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
      console.log("showDeleteRepoModal", repo); ////

      this.currentRepo = repo;
      this.isShownDeleteRepoModal = true;
    },
    hideDeleteRepoModal() {
      this.isShownDeleteRepoModal = false;
    },
    showRepoDetailsModal(repo) {
      console.log("showRepoDetailsModal", repo); ////

      this.isShownRepoDetailsModal = true;
    },
    showEditRepoNameModal(repo) {
      console.log("showEditRepoNameModal", repo); ////

      this.isShownEditRepoNameModal = true;
    },
    showCreateBackupModal() {
      console.log("showCreateBackupModal"); ////

      this.q.isShownCreateBackupModal = true;
    },
    showEditBackupModal(backup) {
      console.log("showEditRepoModal", backup); ////

      this.isShownEditBackupModal = true;
    },
    showDeleteBackupModal(backup) {
      console.log("showDeleteBackupModal", backup); ////

      this.isShownDeleteBackupModal = true;
    },
    showBackupDetailsModal(backup) {
      console.log("showBackupDetailsModal", backup); ////

      this.isShownBackupDetailsModal = true;
    },
    deleteRepo(repo) {
      console.log("deleteRepo", repo); ////
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.card-content .row {
  margin-bottom: $spacing-05;
  text-align: center;
}

.card-content .row.backup-repo {
  margin-bottom: $spacing-02;
}

.card-content .row:last-child {
  margin-bottom: 0;
}
</style>
