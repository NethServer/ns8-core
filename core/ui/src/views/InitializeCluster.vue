<template>
  <cv-grid class="welcome-grid">
    <cv-loading
      :active="isCreatingCluster || isJoiningCluster"
      :overlay="true"
    ></cv-loading>
    <cv-row>
      <cv-column>
        <div class="logo">
          <img
            :src="require('@/assets/logo.png')"
            :alt="this.$root.config.PRODUCT_NAME + ' logo'"
          />
        </div>
      </cv-column>
    </cv-row>
    <template v-if="q.page === 'welcome'">
      <cv-row>
        <cv-column class="welcome">
          <h2>
            {{
              $t("init.welcome", { product: this.$root.config.PRODUCT_NAME })
            }}
          </h2>
        </cv-column>
      </cv-row>
      <!-- create / join / restore cluster -->
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-grid fullWidth class="mg-top-lg">
              <cv-row>
                <cv-column :md="4">
                  <NsTile
                    kind="clickable"
                    :icon="EdgeCluster32"
                    @click="selectCreateCluster"
                    large
                  >
                    <h6 class="mg-bottom-sm">
                      {{ $t("init.create_cluster") }}
                    </h6>
                    <div class="tile-description">
                      {{ $t("init.create_cluster_description") }}
                    </div>
                  </NsTile>
                </cv-column>
                <cv-column :md="4">
                  <NsTile
                    kind="clickable"
                    :icon="Connect32"
                    @click="selectJoinCluster"
                    large
                  >
                    <h6 class="mg-bottom-sm">
                      {{ $t("init.join_cluster") }}
                    </h6>
                    <div class="tile-description">
                      {{ $t("init.join_cluster_description") }}
                    </div>
                  </NsTile>
                </cv-column>
              </cv-row>
              <cv-row>
                <cv-column class="horizontal-divider"> </cv-column>
              </cv-row>
              <cv-row>
                <cv-column>
                  <NsTile
                    kind="clickable"
                    :icon="Reset32"
                    @click="selectRestoreCluster"
                    large
                    class="restore-card"
                  >
                    <h6 class="mg-bottom-sm">
                      {{ $t("init.restore_cluster") }}
                    </h6>
                    <div class="tile-description">
                      {{ $t("init.restore_cluster_description") }}
                    </div>
                  </NsTile>
                </cv-column>
              </cv-row>
            </cv-grid>
          </cv-tile>
        </cv-column>
      </cv-row>
    </template>
    <template v-else-if="q.page === 'create'">
      <template v-if="isPasswordChangeNeeded">
        <cv-row>
          <!-- password change needed -->
          <cv-column class="welcome">
            <h2>{{ $t("init.create_cluster") }}</h2>
            <div class="title-description">
              {{ $t("init.create_cluster_description") }}
            </div>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <NsInlineNotification
              kind="info"
              :title="$t('init.change_admin_password')"
              :description="$t('init.change_admin_password_description')"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <cv-tile light>
              <cv-form @submit.prevent="changePassword">
                <cv-text-input
                  :label="$t('init.current_password')"
                  v-model="currentPassword"
                  :invalid-message="$t(error.currentPassword)"
                  type="password"
                  ref="currentPassword"
                >
                </cv-text-input>
                <NsPasswordInput
                  :newPasswordLabel="$t('init.new_password')"
                  :confirmPasswordLabel="$t('init.new_password_confirm')"
                  v-model="newPassword"
                  @passwordValidation="onPasswordValidation"
                  :newPasswordInvalidMessage="$t(error.newPassword)"
                  :confirmPasswordInvalidMessage="$t(error.confirmPassword)"
                  :passwordHideLabel="$t('password.hide_password')"
                  :passwordShowLabel="$t('password.show_password')"
                  :lengthLabel="$t('password.long_enough')"
                  :lowercaseLabel="$t('password.lowercase_letter')"
                  :uppercaseLabel="$t('password.uppercase_letter')"
                  :numberLabel="$t('password.number')"
                  :symbolLabel="$t('password.symbol')"
                  :equalLabel="$t('password.equal')"
                  :focus="focusPasswordField"
                />
                <cv-button-set class="footer-buttons">
                  <NsButton
                    type="button"
                    kind="secondary"
                    :icon="ChevronLeft20"
                    size="lg"
                    @click="goToWelcomePage"
                    >{{ $t("common.go_back") }}
                  </NsButton>
                  <NsButton
                    kind="primary"
                    :loading="isChangingPassword"
                    :disabled="isChangingPassword"
                    :icon="Password20"
                    size="lg"
                    >{{ $t("init.change_password") }}
                  </NsButton>
                </cv-button-set>
              </cv-form>
            </cv-tile>
          </cv-column>
        </cv-row>
      </template>
      <!-- admin password was changed -->
      <template v-else>
        <!-- create cluster form -->
        <cv-row>
          <cv-column class="welcome">
            <h2>{{ $t("init.create_cluster") }}</h2>
            <div class="title-description">
              {{ $t("init.create_cluster_description") }}
            </div>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <cv-tile light>
              <cv-form @submit.prevent="createCluster">
                <cv-text-input
                  :label="$t('init.vpn_endpoint_address')"
                  v-model.trim="vpnEndpointAddress"
                  :invalid-message="$t(error.vpnEndpointAddress)"
                  ref="vpnEndpointAddress"
                >
                </cv-text-input>
                <cv-text-input
                  :label="$t('init.vpn_endpoint_port')"
                  v-model.trim="vpnEndpointPort"
                  :invalid-message="$t(error.vpnEndpointPort)"
                  type="number"
                  class="narrow"
                  ref="vpnEndpointPort"
                >
                </cv-text-input>
                <cv-text-input
                  :label="$t('init.vpn_cidr')"
                  v-model.trim="vpnCidr"
                  :invalid-message="$t(error.vpnCidr)"
                  class="narrow"
                  ref="vpnCidr"
                >
                </cv-text-input>
                <cv-button-set class="footer-buttons">
                  <NsButton
                    type="button"
                    kind="secondary"
                    :icon="ChevronLeft20"
                    size="lg"
                    @click="goToWelcomePage"
                    >{{ $t("common.go_back") }}
                  </NsButton>
                  <NsButton
                    kind="primary"
                    :loading="isCreatingCluster"
                    :disabled="isCreatingCluster"
                    :icon="EdgeCluster20"
                    size="lg"
                    >{{ $t("init.create_cluster") }}
                  </NsButton>
                </cv-button-set>
              </cv-form>
            </cv-tile>
          </cv-column>
        </cv-row>
      </template>
    </template>
    <template v-else-if="q.page === 'join'">
      <!-- join cluster form -->
      <cv-row>
        <cv-column class="welcome">
          <h2>{{ $t("init.join_cluster") }}</h2>
          <div class="title-description">
            {{ $t("init.join_cluster_description") }}
          </div>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <cv-tile light>
            <cv-form @submit.prevent="joinCluster">
              <cv-text-area
                :label="$t('common.join_code')"
                v-model.trim="joinCode"
                :invalid-message="$t(error.joinCode)"
                :helper-text="
                  $t('init.join_code_helper_text') +
                  ' https://LEADER_NODE_IP/cluster-admin/#/nodes?isShownAddNodeModal=true'
                "
                class="join-code"
                ref="joinCode"
              >
              </cv-text-area>
              <cv-checkbox
                :label="$t('init.tls_verify')"
                v-model="tlsVerify"
                value="checkTlsVerify"
              />
              <cv-button-set class="footer-buttons">
                <NsButton
                  type="button"
                  kind="secondary"
                  :icon="ChevronLeft20"
                  size="lg"
                  @click="goToWelcomePage"
                  >{{ $t("common.go_back") }}
                </NsButton>
                <NsButton
                  kind="primary"
                  :loading="isJoiningCluster"
                  :disabled="isJoiningCluster"
                  :icon="Connect20"
                  size="lg"
                  >{{ $t("init.join_cluster") }}
                </NsButton>
              </cv-button-set>
            </cv-form>
          </cv-tile>
        </cv-column>
      </cv-row>
    </template>
    <template v-else-if="q.page === 'restore'">
      <template v-if="isPasswordChangeNeeded">
        <cv-row>
          <!-- password change needed -->
          <cv-column class="welcome">
            <h2>{{ $t("init.restore_cluster") }}</h2>
            <div class="title-description">
              {{ $t("init.restore_cluster_description") }}
            </div>
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <NsInlineNotification
              kind="info"
              :title="$t('init.change_admin_password')"
              :description="$t('init.change_admin_password_description')"
              :showCloseButton="false"
            />
          </cv-column>
        </cv-row>
        <cv-row>
          <cv-column>
            <cv-tile light>
              <cv-form @submit.prevent="changePassword">
                <cv-text-input
                  :label="$t('init.current_password')"
                  v-model="currentPassword"
                  :invalid-message="$t(error.currentPassword)"
                  type="password"
                  ref="currentPassword"
                >
                </cv-text-input>
                <NsPasswordInput
                  :newPasswordLabel="$t('init.new_password')"
                  :confirmPasswordLabel="$t('init.new_password_confirm')"
                  v-model="newPassword"
                  @passwordValidation="onPasswordValidation"
                  :newPasswordInvalidMessage="$t(error.newPassword)"
                  :confirmPasswordInvalidMessage="$t(error.confirmPassword)"
                  :passwordHideLabel="$t('password.hide_password')"
                  :passwordShowLabel="$t('password.show_password')"
                  :lengthLabel="$t('password.long_enough')"
                  :lowercaseLabel="$t('password.lowercase_letter')"
                  :uppercaseLabel="$t('password.uppercase_letter')"
                  :numberLabel="$t('password.number')"
                  :symbolLabel="$t('password.symbol')"
                  :equalLabel="$t('password.equal')"
                  :focus="focusPasswordField"
                />
                <cv-button-set class="footer-buttons">
                  <NsButton
                    type="button"
                    kind="secondary"
                    :icon="ChevronLeft20"
                    size="lg"
                    @click="goToWelcomePage"
                    >{{ $t("common.go_back") }}
                  </NsButton>
                  <NsButton
                    kind="primary"
                    :loading="isChangingPassword"
                    :disabled="isChangingPassword"
                    :icon="Password20"
                    size="lg"
                    >{{ $t("init.change_password") }}
                  </NsButton>
                </cv-button-set>
              </cv-form>
            </cv-tile>
          </cv-column>
        </cv-row>
      </template>
      <!-- admin password was changed -->
      <template v-else>
        <cv-row>
          <cv-column class="welcome">
            <h2>{{ $t("init.restore_cluster") }}</h2>
            <div class="title-description">
              {{ $t("init.restore_cluster_description") }}
            </div>
          </cv-column>
        </cv-row>
        <!-- restore cluster -->
        <template v-if="!restore.type">
          <!-- restore from file / url -->
          <cv-row>
            <cv-column>
              <cv-tile light>
                <cv-grid fullWidth class="mg-top-lg no-padding">
                  <cv-row>
                    <cv-column :md="4">
                      <NsTile
                        kind="clickable"
                        :icon="Link32"
                        @click="selectRestoreFromUrl"
                        large
                      >
                        <h6>
                          {{ $t("init.restore_from_remote_url") }}
                        </h6>
                      </NsTile>
                    </cv-column>
                    <cv-column :md="4">
                      <NsTile
                        kind="clickable"
                        :icon="Document32"
                        @click="selectRestoreFromFile"
                        large
                      >
                        <h6>
                          {{ $t("init.restore_from_cluster_backup_file") }}
                        </h6>
                      </NsTile>
                    </cv-column>
                  </cv-row>
                </cv-grid>
                <cv-button-set class="footer-buttons">
                  <NsButton
                    type="button"
                    kind="secondary"
                    :icon="ChevronLeft20"
                    size="lg"
                    @click="goToWelcomePage"
                    >{{ $t("common.go_back") }}
                  </NsButton>
                </cv-button-set>
              </cv-tile>
            </cv-column>
          </cv-row>
        </template>
        <template v-else-if="restore.type == 'url'">
          <!-- restore from url -->
          <cv-row>
            <cv-column>
              <cv-tile light>
                <cv-form @submit.prevent="retrieveClusterBackupFromUrl">
                  <NsTextInput
                    v-model.trim="restore.url"
                    :label="$t('init.remote_url')"
                    placeholder="https://myserver/core-backup.json.gz.gpg"
                    :invalid-message="error.restore.url"
                    :disabled="loading.retrieveClusterBackup"
                    class="mg-bottom-sm"
                    ref="url"
                  />
                  <cv-checkbox
                    :label="$t('init.tls_verify')"
                    v-model="restore.tlsVerify"
                    value="checkTlsVerify"
                  />
                  <!-- advanced options -->
                  <cv-accordion ref="accordion" class="mg-top-lg">
                    <cv-accordion-item
                      :open="
                        toggleAccordion[0] ||
                        restore.isBackupPasswordAccordionOpenAndDisabled
                      "
                      :disabled="
                        restore.isBackupPasswordAccordionOpenAndDisabled
                      "
                    >
                      <template slot="title">{{
                        $t("common.advanced")
                      }}</template>
                      <template slot="content">
                        <NsTextInput
                          type="password"
                          :password-hide-label="$t('password.hide_password')"
                          :password-show-label="$t('password.show_password')"
                          v-model.trim="restore.backupPassword"
                          :label="$t('init.backup_password')"
                          :placeholder="$t('init.use_admin_password')"
                          :invalid-message="error.restore.backup_password"
                          :disabled="loading.retrieveClusterBackup"
                          tooltipAlignment="center"
                          tooltipDirection="right"
                          ref="backup_password"
                          class="mg-top-md"
                        >
                          <template #tooltip>
                            <span
                              v-html="$t('init.backup_password_tooltip')"
                            ></span>
                          </template>
                        </NsTextInput>
                      </template>
                    </cv-accordion-item>
                  </cv-accordion>
                  <cv-button-set class="footer-buttons">
                    <NsButton
                      type="button"
                      kind="secondary"
                      :icon="ChevronLeft20"
                      size="lg"
                      @click="restore.type = ''"
                      >{{ $t("common.go_back") }}
                    </NsButton>
                    <NsButton
                      kind="primary"
                      :loading="loading.retrieveClusterBackup"
                      :disabled="loading.retrieveClusterBackup"
                      :icon="ChevronRight20"
                      size="lg"
                      >{{ $t("common.next") }}
                    </NsButton>
                  </cv-button-set>
                </cv-form>
              </cv-tile>
            </cv-column>
          </cv-row>
        </template>
        <template v-else-if="restore.type == 'file'">
          <!-- restore from file -->
          <cv-row>
            <cv-column>
              <cv-tile light>
                <cv-form @submit.prevent="retrieveClusterBackupFromFile">
                  <div>
                    {{ $t("init.upload_cluster_backup_file") }}
                  </div>
                  <cv-file-uploader
                    :drop-target-label="$t('common.upload_drop_target_text')"
                    accept=".gpg"
                    :clear-on-reselect="true"
                    :multiple="false"
                    :removable="false"
                    :remove-aria-label="$t('common.remove')"
                    v-model="restore.filesUploaded"
                    @change="onFileUpload"
                    :class="[
                      'file-uploader',
                      { 'validation-failed': error.restore.backup_file },
                    ]"
                    ref="backup_file"
                  >
                  </cv-file-uploader>
                  <!-- invalid message for file uploader -->
                  <div
                    v-if="error.restore.backup_file"
                    class="validation-failed-invalid-message"
                    v-html="error.restore.backup_file"
                  ></div>
                  <!-- advanced options -->
                  <cv-accordion ref="accordion" class="mg-top-lg">
                    <cv-accordion-item
                      :open="
                        toggleAccordion[0] ||
                        restore.isBackupPasswordAccordionOpenAndDisabled
                      "
                      :disabled="
                        restore.isBackupPasswordAccordionOpenAndDisabled
                      "
                    >
                      <template slot="title">{{
                        $t("common.advanced")
                      }}</template>
                      <template slot="content">
                        <NsTextInput
                          type="password"
                          :password-hide-label="$t('password.hide_password')"
                          :password-show-label="$t('password.show_password')"
                          v-model.trim="restore.backupPassword"
                          :label="$t('init.backup_password')"
                          :placeholder="$t('init.use_admin_password')"
                          :invalid-message="error.restore.backup_password"
                          :disabled="loading.retrieveClusterBackup"
                          tooltipAlignment="center"
                          tooltipDirection="right"
                          ref="backup_password"
                          class="mg-top-md"
                        >
                          <template #tooltip>
                            <span
                              v-html="$t('init.backup_password_tooltip')"
                            ></span>
                          </template>
                        </NsTextInput>
                      </template>
                    </cv-accordion-item>
                  </cv-accordion>
                  <NsInlineNotification
                    v-if="error.retrieveClusterBackup"
                    kind="error"
                    :title="$t('action.retrieve-cluster-backup')"
                    :description="error.retrieveClusterBackup"
                    :showCloseButton="false"
                  />
                  <cv-button-set class="footer-buttons">
                    <NsButton
                      type="button"
                      kind="secondary"
                      :icon="ChevronLeft20"
                      size="lg"
                      @click="restore.type = ''"
                      >{{ $t("common.go_back") }}
                    </NsButton>
                    <NsButton
                      kind="primary"
                      :loading="loading.retrieveClusterBackup"
                      :disabled="loading.retrieveClusterBackup"
                      :icon="ChevronRight20"
                      size="lg"
                      >{{ $t("common.next") }}
                    </NsButton>
                  </cv-button-set>
                </cv-form>
              </cv-tile>
            </cv-column>
          </cv-row>
        </template>
      </template>
    </template>
    <template v-else-if="q.page === 'redirect'">
      <cv-row>
        <cv-column class="welcome">
          <h2>{{ $t("init.redirect_cluster") }}</h2>
        </cv-column>
      </cv-row>
      <cv-row>
        <cv-column>
          <div class="title-description">
            {{ $t("init.redirect_cluster_description") }}
          </div>
        </cv-column>
      </cv-row>
      <cv-row class="mg-top-lg">
        <cv-column>
          <a
            :href="this.joinEndpoint + '/cluster-admin/'"
            class="external-link-button"
          >
            <NsButton kind="primary" :icon="Launch20">
              {{ $t("init.redirect_cluster_link") }}
            </NsButton>
          </a>
        </cv-column>
      </cv-row>
    </template>
  </cv-grid>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
  StorageService,
  NsPasswordInput,
  TaskService,
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";
import to from "await-to-js";
import NotificationService from "@/mixins/notification";

export default {
  name: "InitializeCluster",
  components: { NsPasswordInput },
  mixins: [
    UtilService,
    IconService,
    QueryParamService,
    StorageService,
    TaskService,
    NotificationService,
  ],
  pageTitle() {
    return this.$t("init.welcome", { product: this.$root.config.PRODUCT_NAME });
  },
  data() {
    return {
      q: {
        page: "welcome",
        endpoint: null,
      },
      isPasswordChangeNeeded: false,
      currentPassword: "",
      newPassword: "",
      passwordValidation: null,
      focusPasswordField: { element: "" },
      vpnEndpointAddress: "",
      vpnEndpointPort: "",
      vpnCidr: "",
      joinCode: "",
      tlsVerify: true,
      joinEndpoint: this.$route.query.endpoint
        ? "https://" + this.$route.query.endpoint
        : "",
      joinPort: "",
      joinToken: "",
      isCreatingCluster: false,
      isJoiningCluster: false,
      isChangingPassword: false,
      isMaster: true,
      restore: {
        type: "",
        url: "",
        tlsVerify: true,
        backupPassword: "",
        filesUploaded: [],
        base64FileUploaded: "",
        isBackupPasswordAccordionOpenAndDisabled: false,
        isBackupPasswordAccordionDisabled: false,
      },
      loading: {
        retrieveClusterBackup: false,
      },
      error: {
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
        vpnEndpointAddress: "",
        vpnEndpointPort: "",
        vpnCidr: "",
        joinCode: "",
        retrieveClusterBackup: "",
        restore: {
          url: "",
          backup_password: "",
          backup_file: "",
        },
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
    this.retrieveClusterStatus();
    this.getDefaults();
  },
  methods: {
    ...mapActions(["setClusterInitializedInStore"]),
    async getDefaults() {
      const taskAction = "get-defaults";

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.getDefaultsCompleted);

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
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
    },
    getDefaultsCompleted(taskContext, taskResult) {
      const defaults = taskResult.output;
      this.vpnEndpointAddress = defaults.vpn.host;
      this.vpnEndpointPort = defaults.vpn.port.toString();
      this.vpnCidr = defaults.vpn.network;
    },
    selectCreateCluster() {
      this.$router.push({ path: "/init", query: { page: "create" } });

      if (this.isPasswordChangeNeeded) {
        this.focusElement("currentPassword");
      } else {
        this.focusElement("vpnEndpointAddress");
      }
    },
    selectJoinCluster() {
      this.$router.push({ path: "/init", query: { page: "join" } });
      this.focusElement("joinCode");
    },
    selectRestoreCluster() {
      this.$router.push({ path: "/init", query: { page: "restore" } });

      if (this.isPasswordChangeNeeded) {
        this.focusElement("currentPassword");
      }
    },
    selectRestoreFromFile() {
      this.restore.type = "file";
    },
    selectRestoreFromUrl() {
      this.restore.type = "url";
    },
    goToWelcomePage() {
      this.$router.push({ path: "/init", query: { page: "welcome" } });
    },
    async retrieveClusterStatus() {
      const taskAction = "get-cluster-status";

      // register to task completion
      this.$root.$once(
        taskAction + "-completed",
        this.getClusterStatusCompleted
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
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
    },
    getClusterStatusCompleted(taskContext, taskResult) {
      const clusterStatus = taskResult.output;

      //// remove mock
      // clusterStatus.initialized = false; ////

      if (clusterStatus.initialized && this.isMaster) {
        // redirect to status page
        this.setClusterInitializedInStore(true);
        this.$root.$emit("clusterInitialized");
        this.$router.replace("/status");
        return;
      }
      this.isPasswordChangeNeeded = clusterStatus.default_password;
    },
    onPasswordValidation(passwordValidation) {
      this.passwordValidation = passwordValidation;
    },
    validatePasswordChange() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.currentPassword) {
        this.error.currentPassword = "common.required";

        if (isValidationOk) {
          this.focusElement("currentPassword");
          isValidationOk = false;
        }
      }

      if (!this.newPassword) {
        this.error.newPassword = "common.required";

        if (isValidationOk) {
          this.focusPasswordField = { element: "newPassword" };
          isValidationOk = false;
        }
      } else {
        if (this.currentPassword === this.newPassword) {
          if (!this.error.newPassword) {
            this.error.newPassword =
              "password.old_new_passwords_must_be_different";
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (
          !this.passwordValidation.isLengthOk ||
          !this.passwordValidation.isLowercaseOk ||
          !this.passwordValidation.isUppercaseOk ||
          !this.passwordValidation.isNumberOk ||
          !this.passwordValidation.isSymbolOk
        ) {
          if (!this.error.newPassword) {
            this.error.newPassword = "password.password_not_secure";
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "newPassword" };
            isValidationOk = false;
          }
        }

        if (!this.passwordValidation.isEqualOk) {
          if (!this.error.newPassword) {
            this.error.newPassword = "password.passwords_do_not_match";
          }

          if (!this.error.confirmPassword) {
            this.error.confirmPassword = "password.passwords_do_not_match";
          }

          if (isValidationOk) {
            this.focusPasswordField = { element: "confirmPassword" };
            isValidationOk = false;
          }
        }
      }
      return isValidationOk;
    },
    async changePassword() {
      this.isChangingPassword = true;

      if (!this.validatePasswordChange()) {
        this.isChangingPassword = false;
        return;
      }

      const loginInfo = this.getFromStorage("loginInfo");
      const taskAction = "change-user-password";

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.changeUserPasswordCompleted
      );

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.changeUserPasswordValidationFailed
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            user: loginInfo.username,
            current_password: this.currentPassword,
            new_password: this.newPassword,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
    },
    changeUserPasswordCompleted() {
      this.isPasswordChangeNeeded = false;
    },
    changeUserPasswordValidationFailed(validationErrors) {
      this.isChangingPassword = false;

      for (const validationError of validationErrors) {
        // set i18n error message
        this.error.currentPassword = "password." + validationError.error;
        this.focusElement("currentPassword");
      }
    },
    validateCreateCluster() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.vpnEndpointAddress) {
        this.error.vpnEndpointAddress = "common.required";

        if (isValidationOk) {
          this.focusElement("vpnEndpointAddress");
          isValidationOk = false;
        }
      }

      if (!this.vpnEndpointPort) {
        this.error.vpnEndpointPort = "common.required";

        if (isValidationOk) {
          this.focusElement("vpnEndpointPort");
          isValidationOk = false;
        }
      } else {
        const vpnEndpointPortNumber = Number(this.vpnEndpointPort);

        if (
          !(
            Number.isInteger(vpnEndpointPortNumber) && vpnEndpointPortNumber > 0
          )
        ) {
          this.error.vpnEndpointPort = "error.invalid_port_number";

          if (isValidationOk) {
            this.focusElement("vpnEndpointPort");
            isValidationOk = false;
          }
        }
      }

      if (!this.vpnCidr) {
        this.error.vpnCidr = "common.required";

        if (isValidationOk) {
          this.focusElement("vpnCidr");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    async createCluster() {
      if (!this.validateCreateCluster()) {
        return;
      }

      const taskAction = "create-cluster";

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.createClusterCompleted);

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.createClusterAborted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            network: this.vpnCidr,
            endpoint: this.vpnEndpointAddress + ":" + this.vpnEndpointPort,
            listen_port: parseInt(this.vpnEndpointPort),
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            toastTimeout: 0, // persistent notification
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
        }
        return;
      }
      this.isCreatingCluster = true;
    },
    createClusterCompleted() {
      this.setClusterInitializedInStore(true);
      this.$root.$emit("clusterInitialized");
      this.$router.replace("/status");
      this.isCreatingCluster = false;
    },
    createClusterAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.isCreatingCluster = false;
    },
    validateJoinCluster() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.joinCode) {
        this.error.joinCode = "common.required";

        if (isValidationOk) {
          this.focusElement("joinCode");
          isValidationOk = false;
        }
      } else {
        let decoded;

        try {
          decoded = atob(this.joinCode);
        } catch (DOMException) {
          this.error.joinCode = "init.invalid_join_code";

          if (isValidationOk) {
            this.focusElement("joinCode");
            isValidationOk = false;
          }
        }

        if (!decoded) {
          this.error.joinCode = "init.invalid_join_code";

          if (isValidationOk) {
            this.focusElement("joinCode");
            isValidationOk = false;
          }
        } else {
          let [endpoint, port, token] = decoded.split("|");

          if (!(endpoint && port && token)) {
            this.error.joinCode = "init.invalid_join_code";

            if (isValidationOk) {
              this.focusElement("joinCode");
              isValidationOk = false;
            }
          } else {
            this.joinEndpoint = endpoint;
            this.joinPort = port;
            this.joinToken = token;
          }
        }
      }

      return isValidationOk;
    },
    async joinCluster() {
      if (!this.validateJoinCluster()) {
        return;
      }

      const taskAction = "join-cluster";

      // register to task events
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.joinClusterCompleted);

      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.joinClusterValidationFailed
      );

      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.joinClusterAborted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            url: this.joinEndpoint,
            jwt: this.joinToken,
            listen_port: parseInt(this.joinPort),
            tls_verify: this.tlsVerify,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            toastTimeout: 0, // persistent notification
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response.status == 403) {
          this.isMaster = false;
        } else {
          // persistent error notification
          const notification = {
            title: this.$t("task.cannot_create_task", { action: taskAction }),
            description: this.getErrorMessage(err),
            type: "error",
            toastTimeout: 0,
          };
          this.createNotification(notification);
          return;
        }
      }
      this.isJoiningCluster = true;
    },
    joinClusterCompleted(taskContext, taskResult) {
      console.log("joinClusterCompleted"); ////
      console.log("taskContext", taskContext); ////
      console.log("taskResult", taskResult); ////

      this.isJoiningCluster = false;
      this.$router.push("/init?page=redirect");

      const nodeId = taskResult.output.nodeId;

      console.log("nodeId", nodeId); ////

      const notification = {
        title: this.$t("action.join-cluster"),
        description: this.$t("init.node_added_to_cluster", { nodeId }),
        type: "success",
      };
      this.createNotification(notification);
    },
    joinClusterAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.isJoiningCluster = false;
    },
    joinClusterValidationFailed(validationErrors) {
      console.error("validation failed", validationErrors);
      this.isJoiningCluster = false;
      this.error.joinCode = "init.invalid_join_code";
      this.focusElement("joinCode");
    },
    async onFileUpload(files) {
      this.error.restore.backup_file = "";
      this.restore.base64FileUploaded = "";

      if (files.length > 1) {
        // keep only last uploaded file

        this.$set(this.restore.filesUploaded, 0, files[files.length - 1]);
        this.restore.filesUploaded.splice(1);
      }

      const fileUploaded = this.restore.filesUploaded[0];

      if (!fileUploaded || fileUploaded.invalidMessageTitle) {
        // invalid file
        return;
      }
      const file = fileUploaded.file;
      const result = await this.toBase64(file).catch((e) => Error(e));

      if (result instanceof Error) {
        console.log("error converting file to base64:", result.message);
        return;
      } else {
        this.restore.base64FileUploaded = result.split(
          "data:application/pgp-encrypted;base64,"
        )[1];
        console.log("base64!", this.restore.base64FileUploaded); ////
      }
    },
    //// move to ui-lib
    toBase64(file) {
      return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = (error) => reject(error);
      });
    },
    validateRetrieveClusterBackupFromFile() {
      let isValidationOk = true;
      this.clearErrors();

      if (
        !this.restore.filesUploaded.length ||
        !this.restore.base64FileUploaded
      ) {
        this.error.restore.backup_file = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("backup_file");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    validateRetrieveClusterBackupFromUrl() {
      let isValidationOk = true;
      this.clearErrors();

      if (!this.restore.url) {
        this.error.restore.url = this.$t("common.required");

        if (isValidationOk) {
          this.focusElement("url");
          isValidationOk = false;
        }
      }
      return isValidationOk;
    },
    retrieveClusterBackupFromFile() {
      if (!this.validateRetrieveClusterBackupFromFile()) {
        return;
      }

      const inputData = {
        type: "file",
        password: this.restore.backupPassword,
        file: this.restore.base64FileUploaded,
      };
      this.retrieveClusterBackup(inputData);
    },
    retrieveClusterBackupFromUrl() {
      if (!this.validateRetrieveClusterBackupFromUrl()) {
        return;
      }

      const inputData = {
        type: "url",
        password: this.restore.backupPassword,
        url: this.restore.url,
        tls_verify: this.restore.tlsVerify,
      };
      this.retrieveClusterBackup(inputData);
    },
    async retrieveClusterBackup(inputData) {
      // validation already performed in retrieveClusterBackupFromFile() and retrieveClusterBackupFromUrl()
      this.error.retrieveClusterBackup = "";
      this.loading.retrieveClusterBackup = true;
      const taskAction = "retrieve-cluster-backup";

      // register to task validation
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.retrieveClusterBackupValidationFailed
      );

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.retrieveClusterBackupAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.retrieveClusterBackupCompleted
      );

      console.log("retrieveClusterBackup, inputData", inputData); ////

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: inputData,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            toastTimeout: 0, // persistent notification
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.retrieveClusterBackup = this.getErrorMessage(err);
        return;
      }
    },
    retrieveClusterBackupValidationFailed(validationErrors) {
      console.log("retrieveClusterBackupValidationFailed", validationErrors); ////

      this.loading.retrieveClusterBackup = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        if (param == "backup_password") {
          this.restore.isBackupPasswordAccordionOpenAndDisabled = true;
        }

        console.log("validationError", validationError); ////

        // set i18n error message
        this.error.restore[param] = this.$t("init." + validationError.error);

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    retrieveClusterBackupAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.retrieveClusterBackup = false;
    },
    retrieveClusterBackupCompleted(taskContext, taskResult) {
      console.log("retrieveClusterBackupCompleted", taskResult.output); ////

      this.loading.retrieveClusterBackup = false;
      this.restore.isBackupPasswordAccordionOpenAndDisabled = false;

      //// go to next step
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.welcome-grid {
  max-width: 70rem;
}

.logo {
  width: 4rem;
  height: 4rem;
  margin-top: $spacing-07;
  flex-shrink: 0;
}

.logo img {
  width: 100%;
  height: 100%;
}

.welcome {
  margin-top: 2rem;
  margin-bottom: 4rem;
}

.tile-description {
  color: $text-02;
}

.bx--form .bx--form-item {
  margin-bottom: $spacing-06;
}

.horizontal-divider {
  margin-top: 0rem;
  margin-bottom: 2rem;
  height: 1px;
  background-color: $ui-04;
  max-width: 40rem;
  margin-left: auto;
  margin-right: auto;
}

.restore-card {
  max-width: 20rem;
  margin-left: auto;
  margin-right: auto;
}

.file-uploader {
  margin-bottom: 0 !important;
}

.footer-buttons {
  display: flex;
  justify-content: flex-end;
  margin-top: $spacing-07;
}

@media (max-width: $breakpoint-medium) {
  .footer-buttons.bx--btn-set .bx--btn {
    max-width: 9rem;
  }
}
</style>

<style lang="scss">
@import "../styles/carbon-utils";

// global styles

.join-code textarea {
  min-height: 5rem;
}

// highlight file uploaded on validation error
.validation-failed .bx--file__selected-file {
  outline: 2px solid $danger-01;
  outline-offset: -2px;
}
</style>
