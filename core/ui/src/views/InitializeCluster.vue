<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="initialize-cluster">
    <cv-grid class="welcome-grid">
      <cv-loading :active="isJoiningCluster" :overlay="true"></cv-loading>
      <template v-if="q.page === 'welcome'">
        <cv-row>
          <cv-column class="welcome">
            <h2>
              {{
                $t("init.welcome", { product: this.$root.config.PRODUCT_NAME })
              }}
            </h2>
            <WelcomeLogo />
          </cv-column>
        </cv-row>
        <!-- create / join / restore cluster -->
        <cv-row>
          <cv-column>
            <cv-tile light>
              <cv-grid fullWidth class="mg-top-xlg">
                <cv-row>
                  <cv-column>
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
                  <cv-column>
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
                  <cv-column>
                    <NsTile
                      kind="clickable"
                      :icon="Reset32"
                      @click="selectRestoreCluster"
                      large
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
      <template v-else-if="q.page === 'fqdn'">
        <template
          v-if="isPasswordChangeNeeded && $route.query.action === 'create'"
        >
          <cv-row>
            <!-- password change needed -->
            <cv-column class="welcome">
              <div>
                <h2>{{ $t("init.create_cluster") }}</h2>
                <div class="title-description">
                  {{ $t("init.create_cluster_description") }}
                </div>
              </div>
              <WelcomeLogo />
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
          <!-- fqdn form -->
          <SetFqdn
            :action="$route.query.action"
            @goBack="goToWelcomePage"
            @fqdnSet="onFqdnSet"
          />
        </template>
      </template>
      <template v-else-if="q.page === 'create'">
        <template>
          <!-- create cluster form -->
          <cv-row>
            <cv-column class="welcome">
              <div>
                <h2>{{ $t("init.create_cluster") }}</h2>
                <div class="title-description">
                  {{ $t("init.create_cluster_description") }}
                </div>
              </div>
              <WelcomeLogo />
            </cv-column>
          </cv-row>
          <cv-row v-if="!isCreatingCluster">
            <cv-column>
              <NsInlineNotification
                kind="info"
                :title="$t('init.vpn_network')"
                :description="$t('init.vpn_cidr_tooltip')"
                :showCloseButton="false"
              />
            </cv-column>
          </cv-row>
          <cv-row>
            <cv-column>
              <cv-tile light>
                <template v-if="!isCreatingCluster">
                  <NsInlineNotification
                    v-if="error.getDefaults"
                    kind="error"
                    :title="$t('action.get-defaults')"
                    :description="error.getFqdn"
                    :showCloseButton="false"
                  />
                  <cv-form v-else @submit.prevent="createCluster">
                    <cv-skeleton-text
                      v-if="loading.getDefaults"
                      :paragraph="true"
                      :line-count="3"
                      heading
                    ></cv-skeleton-text>
                    <template v-else>
                      <NsTextInput
                        :label="$t('init.vpn_cidr')"
                        v-model.trim="network"
                        :invalid-message="$t(error.network)"
                        :disabled="loading.getDefaults || isCreatingCluster"
                        tooltipAlignment="end"
                        tooltipDirection="right"
                        class="narrow"
                        ref="network"
                      />
                    </template>
                    <cv-button-set class="footer-buttons">
                      <NsButton
                        type="button"
                        kind="secondary"
                        :icon="ChevronLeft20"
                        size="lg"
                        @click="goToFqdnPage"
                        >{{ $t("common.go_back") }}
                      </NsButton>
                      <NsButton
                        kind="primary"
                        :loading="isCreatingCluster"
                        :disabled="loading.getDefaults || isCreatingCluster"
                        :icon="EdgeCluster20"
                        size="lg"
                        >{{ $t("init.create_cluster") }}
                      </NsButton>
                    </cv-button-set>
                  </cv-form>
                </template>
                <template v-else>
                  <NsInlineNotification
                    v-if="error.createCluster"
                    kind="error"
                    :title="$t('action.create-cluster')"
                    :description="error.createCluster"
                    :showCloseButton="false"
                  />
                  <NsEmptyState
                    :title="$t('init.creating_cluster')"
                    :animationData="GearsLottie"
                    animationTitle="gears"
                    :loop="true"
                  />
                  <NsProgressBar
                    :value="createClusterProgress"
                    :indeterminate="!createClusterProgress"
                    class="mg-bottom-xlg"
                  />
                </template>
              </cv-tile>
            </cv-column>
          </cv-row>
        </template>
      </template>
      <template v-else-if="q.page === 'join'">
        <!-- join cluster form -->
        <cv-row>
          <cv-column class="welcome">
            <div>
              <h2>{{ $t("init.join_cluster") }}</h2>
              <div class="title-description">
                {{ $t("init.join_cluster_description") }}
              </div>
            </div>
            <WelcomeLogo />
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
                    @click="goToFqdnPage"
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
              <div>
                <h2>{{ $t("init.restore_cluster") }}</h2>
                <div class="title-description">
                  {{ $t("init.restore_cluster_description") }}
                </div>
              </div>
              <WelcomeLogo />
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
              <div>
                <h2>{{ $t("init.restore_cluster") }}</h2>
                <div class="title-description">
                  {{ $t("init.restore_cluster_description") }}
                </div>
              </div>
              <WelcomeLogo />
            </cv-column>
          </cv-row>
          <!-- restore cluster -->
          <template v-if="restore.step == 'type'">
            <!-- restore from file / url -->
            <cv-row>
              <cv-column>
                <cv-tile light>
                  <cv-grid fullWidth class="mg-top-xlg no-padding">
                    <cv-row>
                      <cv-column :md="4">
                        <NsTile
                          kind="clickable"
                          :icon="Link32"
                          @click="goToRestoreFromUrl"
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
                          @click="goToRestoreFromFile"
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
          <template v-else-if="restore.step == 'url'">
            <!-- restore from url -->
            <cv-row>
              <cv-column>
                <cv-tile light>
                  <h5 class="mg-bottom-md">
                    {{ $t("init.restore_from_remote_url") }}
                  </h5>
                  <cv-form @submit.prevent="retrieveClusterBackupFromUrl">
                    <NsTextInput
                      v-model.trim="restore.url"
                      :label="$t('init.remote_url')"
                      placeholder="https://myserver/cluster-backup.json.gz.gpg"
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
                    <NsTextInput
                      type="password"
                      :password-hide-label="$t('password.hide_password')"
                      :password-show-label="$t('password.show_password')"
                      v-model="restore.backupPassword"
                      :label="$t('init.backup_password')"
                      :invalid-message="error.restore.backup_password"
                      :disabled="loading.retrieveClusterBackup"
                      tooltipAlignment="center"
                      tooltipDirection="right"
                      ref="backup_password"
                      class="mg-top-md"
                    >
                      <template #tooltip>
                        <span>{{ $t("init.backup_password_tooltip") }}</span>
                      </template>
                    </NsTextInput>
                    <cv-button-set class="footer-buttons">
                      <NsButton
                        type="button"
                        kind="secondary"
                        :icon="ChevronLeft20"
                        size="lg"
                        @click="goToRestoreType"
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
          <template v-else-if="restore.step == 'file'">
            <!-- restore from file -->
            <cv-row>
              <cv-column>
                <cv-tile light>
                  <h5 class="mg-bottom-md">
                    {{ $t("init.restore_from_cluster_backup_file") }}
                  </h5>
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
                    >
                      {{ error.restore.backup_file }}
                    </div>
                    <NsTextInput
                      type="password"
                      :password-hide-label="$t('password.hide_password')"
                      :password-show-label="$t('password.show_password')"
                      v-model="restore.backupPassword"
                      :label="$t('init.backup_password')"
                      :invalid-message="error.restore.backup_password"
                      :disabled="loading.retrieveClusterBackup"
                      tooltipAlignment="center"
                      tooltipDirection="right"
                      ref="backup_password"
                      class="mg-top-md"
                    >
                      <template #tooltip>
                        <span>{{ $t("init.backup_password_tooltip") }}</span>
                      </template>
                    </NsTextInput>
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
                        @click="goToRestoreType"
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
          <template v-else-if="restore.step == 'summary'">
            <!-- restore summary -->
            <cv-row>
              <cv-column>
                <cv-tile light>
                  <h5 class="mg-bottom-md">
                    {{ $t("init.restore_summary") }}
                  </h5>
                  <cv-form @submit.prevent="restoreCluster">
                    <div
                      v-if="restore.summary.cluster"
                      class="key-value-setting"
                    >
                      <span class="label">{{ $t("init.cluster_label") }}</span>
                      <span class="value">{{ restore.summary.cluster }}</span>
                    </div>
                    <div class="key-value-setting">
                      <span class="label">{{ $t("init.vpn_endpoint") }}</span>
                      <span class="value">{{ restore.summary.vpn }}</span>
                    </div>
                    <div class="key-value-setting">
                      <span class="label">{{
                        $t("init.external_domains")
                      }}</span>
                      <span class="value">{{ restore.summary.domains }}</span>
                    </div>
                    <div class="key-value-setting">
                      <span class="label">{{ $t("init.custom_routes") }}</span>
                      <span class="value">{{ restore.summary.routes }}</span>
                    </div>
                    <div class="key-value-setting">
                      <span class="label">{{
                        $t("init.custom_certificates")
                      }}</span>
                      <span class="value">{{
                        restore.summary.certificates
                      }}</span>
                    </div>
                    <div class="key-value-setting">
                      <span class="label">{{
                        $t("init.backup_repositories")
                      }}</span>
                      <span class="value">{{
                        restore.summary.backup_repositories
                      }}</span>
                    </div>
                    <NsInlineNotification
                      kind="info"
                      :title="$t('init.apps_restoration')"
                      :description="
                        restore.summary.single_node_cluster
                          ? $t('init.apps_restoration_single_node_description')
                          : $t('init.apps_restoration_multi_node_description')
                      "
                      :showCloseButton="false"
                      class="mg-top-lg"
                    />
                    <NsInlineNotification
                      v-if="error.restoreCluster"
                      kind="error"
                      :title="$t('action.restore-cluster')"
                      :description="error.restoreCluster"
                      :showCloseButton="false"
                    />
                    <cv-button-set class="footer-buttons">
                      <NsButton
                        type="button"
                        kind="secondary"
                        :icon="ChevronLeft20"
                        size="lg"
                        @click="goToRestoreFromFileOrUrl"
                        >{{ $t("common.go_back") }}
                      </NsButton>
                      <NsButton
                        kind="primary"
                        :loading="loading.restoreCluster"
                        :disabled="loading.restoreCluster"
                        :icon="Reset20"
                        size="lg"
                        >{{ $t("init.restore_cluster") }}
                      </NsButton>
                    </cv-button-set>
                  </cv-form>
                </cv-tile>
              </cv-column>
            </cv-row>
          </template>
          <template v-if="restore.step == 'restoringCluster'">
            <NsInlineNotification
              v-if="error.restoreCluster"
              kind="error"
              :title="$t('action.restore-cluster')"
              :description="error.restoreCluster"
              :showCloseButton="false"
            />
            <NsInlineNotification
              v-if="error.readBackupRepositories"
              kind="error"
              :title="$t('action.read-backup-repositories')"
              :description="error.readBackupRepositories"
              :showCloseButton="false"
            />
            <cv-row>
              <cv-column>
                <cv-tile light>
                  <NsEmptyState
                    :title="$t('init.restoring_cluster')"
                    :animationData="GearsLottie"
                    animationTitle="gears"
                    :loop="true"
                  >
                    <template #description>
                      <span v-if="restore.progress < 100">
                        {{ $t("init.traveling_back_in_time") }}
                      </span>
                      <span v-else-if="loading.readBackupRepositories">
                        {{ $t("init.reading_backup_repositories") }}
                      </span>
                      <span v-else>&nbsp;</span>
                    </template>
                  </NsEmptyState>
                  <NsProgressBar
                    :value="restore.progress"
                    :indeterminate="!restore.progress"
                    class="mg-bottom-md"
                  />
                </cv-tile>
              </cv-column>
            </cv-row>
          </template>
          <template v-if="restore.step == 'apps'">
            <!-- restore modules -->
            <cv-row>
              <cv-column>
                <NsInlineNotification
                  kind="info"
                  :title="$t('init.select_apps_to_restore')"
                  :description="$t('init.select_apps_to_restore_description')"
                  :showCloseButton="false"
                />
                <cv-tile light>
                  <cv-form @submit.prevent>
                    <RestoreMultipleInstancesSelector
                      :instances="restore.instances"
                      selection="all"
                      :loading="loading.retrieveClusterBackup"
                      @select="onSelectInstances"
                      :light="false"
                    />
                    <NsInlineNotification
                      v-if="error.restoreModules"
                      kind="error"
                      :title="$t('action.restore-modules')"
                      :description="error.restoreModules"
                      :showCloseButton="false"
                    />
                    <cv-button-set class="footer-buttons">
                      <NsButton
                        type="button"
                        kind="secondary"
                        :icon="ChevronRight20"
                        size="lg"
                        @click="showSkipRestoreAppsModal"
                        :disabled="loading.restoreModules"
                        >{{ $t("common.skip") }}
                      </NsButton>
                      <NsButton
                        kind="primary"
                        :loading="loading.restoreModules"
                        :disabled="
                          loading.restoreModules ||
                          !restore.selectedInstances.length
                        "
                        :icon="Reset20"
                        size="lg"
                        type="button"
                        @click.prevent="restoreModules"
                      >
                        {{
                          $tc(
                            "init.restore_num_apps",
                            restore.selectedInstances.length,
                            { num: restore.selectedInstances.length }
                          )
                        }}
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
            <WelcomeLogo />
          </cv-column>
        </cv-row>
        <cv-tile light>
          <cv-row>
            <cv-column>
              <div>
                {{ $t("init.redirect_cluster_description") }}
              </div>
            </cv-column>
          </cv-row>
          <cv-row class="mg-top-xlg">
            <cv-column>
              <a
                :href="this.joinEndpoint + '/cluster-admin/'"
                class="external-link-button"
              >
                <NsButton kind="primary" :icon="ArrowRight20">
                  {{ $t("init.redirect_cluster_link") }}
                </NsButton>
              </a>
            </cv-column>
          </cv-row>
        </cv-tile>
      </template>
    </cv-grid>
    <SkipRestoreAppsModal
      :isShown="restore.isShownSkipRestoreAppsModal"
      @hide="hideSkipRestoreAppsModal"
      @confirm="skipRestoreApps"
    />
  </div>
</template>

<script>
import {
  QueryParamService,
  UtilService,
  IconService,
  StorageService,
  NsPasswordInput,
  TaskService,
  LottieService,
  PageTitleService,
} from "@nethserver/ns8-ui-lib";
import { mapActions } from "vuex";
import to from "await-to-js";
import NotificationService from "@/mixins/notification";
import RestoreMultipleInstancesSelector from "@/components/backup/RestoreMultipleInstancesSelector";
import SkipRestoreAppsModal from "@/components/initialize-cluster/SkipRestoreAppsModal";
import WelcomeLogo from "@/components/initialize-cluster/WelcomeLogo";
import SetFqdn from "@/components/initialize-cluster/SetFqdn";

export default {
  name: "InitializeCluster",
  components: {
    NsPasswordInput,
    RestoreMultipleInstancesSelector,
    SkipRestoreAppsModal,
    WelcomeLogo,
    SetFqdn,
  },
  mixins: [
    UtilService,
    IconService,
    QueryParamService,
    StorageService,
    TaskService,
    NotificationService,
    LottieService,
    PageTitleService,
  ],
  pageTitle() {
    return this.$t("init.welcome", { product: this.$root.config.PRODUCT_NAME });
  },
  data() {
    return {
      q: {
        page: "welcome",
        endpoint: null,
        action: null,
      },
      isPasswordChangeNeeded: false,
      currentPassword: "",
      newPassword: "",
      passwordValidation: null,
      focusPasswordField: { element: "" },
      vpnEndpointAddress: "",
      vpnEndpointPort: "",
      network: "",
      joinCode: "",
      tlsVerify: true,
      joinEndpoint: this.$route.query.endpoint
        ? "https://" + this.$route.query.endpoint
        : "",
      joinToken: "",
      isCreatingCluster: false,
      isJoiningCluster: false,
      isChangingPassword: false,
      isMaster: true,
      createClusterProgress: 0,
      isOpenCreateClusterAccordion: false,
      restore: {
        step: "type",
        type: "",
        url: "",
        tlsVerify: true,
        backupPassword: "",
        filesUploaded: [],
        base64FileUploaded: "",
        summary: {},
        progress: 0,
        instances: [],
        selectedInstances: [],
        isShownSkipRestoreAppsModal: false,
      },
      loading: {
        getDefaults: true,
        retrieveClusterBackup: false,
        restoreCluster: false,
        readBackupRepositories: false,
        restoreModules: false,
      },
      error: {
        currentPassword: "",
        newPassword: "",
        confirmPassword: "",
        vpnEndpointAddress: "",
        vpnEndpointPort: "",
        network: "",
        joinCode: "",
        retrieveClusterBackup: "",
        restoreCluster: "",
        readBackupRepositories: "",
        restoreModules: "",
        createCluster: "",
        getDefaults: "",
        restore: {
          url: "",
          backup_password: "",
          backup_file: "",
        },
      },
    };
  },
  watch: {
    "q.page": function () {
      if (this.q.page === "create") {
        this.getDefaults();
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
    this.getClusterStatus();
    this.getDefaults();
  },
  methods: {
    ...mapActions(["setClusterInitializedInStore"]),
    async getDefaults() {
      this.loading.getDefaults = true;
      const taskAction = "get-defaults";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.getDefaultsAborted);

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(taskAction + "-completed", this.getDefaultsCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            toastTimeout: 0, // persistent notification
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response && err.response.status == 403) {
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
    getDefaultsAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.getDefaults = false;
    },
    getDefaultsCompleted(taskContext, taskResult) {
      const defaults = taskResult.output;
      this.vpnEndpointAddress = defaults.vpn.host;
      this.vpnEndpointPort = defaults.vpn.port.toString();
      this.network = defaults.vpn.network;
      this.loading.getDefaults = false;

      if (this.q.page === "create") {
        this.focusElement("network");
      }
    },
    selectCreateCluster() {
      this.$router.push({
        path: "/init",
        query: { page: "fqdn", action: "create" },
      });

      if (this.isPasswordChangeNeeded) {
        this.focusElement("currentPassword");
      }
    },
    selectJoinCluster() {
      this.$router.push({
        path: "/init",
        query: { page: "fqdn", action: "join" },
      });
    },
    selectRestoreCluster() {
      this.$router.push({ path: "/init", query: { page: "restore" } });

      if (this.isPasswordChangeNeeded) {
        this.focusElement("currentPassword");
      }
    },
    goToRestoreType() {
      this.restore.step = "type";
      this.restore.type = "";
    },
    goToRestoreFromFile() {
      this.restore.step = "file";
      this.restore.type = "file";
    },
    goToRestoreFromUrl() {
      this.restore.step = "url";
      this.restore.type = "url";
    },
    goToRestoreFromFileOrUrl() {
      if (this.restore.type == "file") {
        this.goToRestoreFromFile();
      } else {
        this.goToRestoreFromUrl();
      }
    },
    goToWelcomePage() {
      this.$router.push({ path: "/init", query: { page: "welcome" } });
    },
    goToFqdnPage() {
      const action = this.$route.query.action; // "create" or "join"

      this.$router.push({
        path: "/init",
        query: { page: "fqdn", action: action },
      });
    },
    onFqdnSet() {
      const action = this.$route.query.action; // "create" or "join"

      this.$router.push({
        path: "/init",
        query: { page: action, action: action },
      });

      if (action === "join") {
        this.focusElement("joinCode");
      }
    },
    async getClusterStatus() {
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
            toastTimeout: 0, // persistent notification
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response && err.response.status == 403) {
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
            toastTimeout: 0, // persistent notification
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response && err.response.status == 403) {
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
    isFqdn(fqdn) {
      // Regular expression for a valid FQDN
      const fqdnRegex = /^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$/;

      // Check if the provided FQDN matches the regular expression
      return fqdnRegex.test(fqdn);
    },
    validateCreateCluster() {
      this.clearErrors(this);
      let isValidationOk = true;

      if (!this.network) {
        this.error.network = "common.required";

        if (isValidationOk) {
          this.focusElement("network");
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
      const eventId = this.getUuid();
      this.createClusterProgress = 0;

      // register to task validation
      this.$root.$once(
        `${taskAction}-validation-failed-${eventId}`,
        this.createClusterValidationFailed
      );

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.createClusterAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.createClusterCompleted
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.createClusterProgressUpdated
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: {
            network: this.network,
            endpoint: this.vpnEndpointAddress + ":" + this.vpnEndpointPort,
          },
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            toastTimeout: 0, // persistent notification
            isProgressNotified: true,
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        if (err.response && err.response.status == 403) {
          this.isMaster = false;
        } else {
          console.error(`error creating task ${taskAction}`, err);
          this.error.createCluster = this.getErrorMessage(err);
          this.loading.getDefaults = false;
          return;
        }
      }
      this.isCreatingCluster = true;
    },
    createClusterValidationFailed(validationErrors) {
      this.loading.getDefaults = false;
      this.isCreatingCluster = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        // set i18n error message
        this.error[param] = this.getI18nStringWithFallback(
          "init." + validationError.error,
          "error." + validationError.error
        );

        if (!focusAlreadySet) {
          this.focusElement(param);
          focusAlreadySet = true;
        }
      }
    },
    createClusterCompleted(taskContext) {
      this.setClusterInitializedInStore(true);
      this.$root.$emit("clusterInitialized");
      this.$router.replace("/status");
      this.isCreatingCluster = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );
    },
    createClusterAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.error.createCluster = this.$t("error.generic_error");
      this.isCreatingCluster = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );
    },
    createClusterProgressUpdated(progress) {
      this.createClusterProgress = progress;
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
          this.error.joinCode = "init.the_join_code_can_not_be_decoded";

          if (isValidationOk) {
            this.focusElement("joinCode");
            isValidationOk = false;
          }
        }

        if (!decoded) {
          this.error.joinCode = "init.the_join_code_is_not_correctly_encoded";

          if (isValidationOk) {
            this.focusElement("joinCode");
            isValidationOk = false;
          }
        } else {
          let [endpoint, token] = decoded.split("|");

          if (!(endpoint && token)) {
            this.error.joinCode = "init.the_join_code_cannot_be_parsed";

            if (isValidationOk) {
              this.focusElement("joinCode");
              isValidationOk = false;
            }
          } else {
            this.joinEndpoint = endpoint;
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
        if (err.response && err.response.status == 403) {
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
      this.isJoiningCluster = false;
      this.$router.push("/init?page=redirect");

      const nodeId = taskResult.output.nodeId;

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
      this.error.joinCode = this.$t("init." + validationErrors[0].error, {
        value: validationErrors[0].value,
      });
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
      const result = await this.fileToBase64(file).catch((e) => Error(e));

      if (result instanceof Error) {
        console.log("error converting file to base64:", result.message);
        return;
      } else {
        this.restore.base64FileUploaded = result.split(";base64,")[1];
      }
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
      this.loading.retrieveClusterBackup = false;
      let focusAlreadySet = false;

      for (const validationError of validationErrors) {
        const param = validationError.parameter;

        if (param == "backup_password") {
          this.restore.isBackupPasswordAccordionOpenAndDisabled = true;
        }

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
      this.restore.summary = taskResult.output;
      this.restore.step = "summary";
      this.loading.retrieveClusterBackup = false;
      this.restore.backupPassword = "";
    },
    async restoreCluster() {
      this.error.restoreCluster = "";
      this.loading.restoreCluster = true;
      this.restore.step = "restoringCluster";
      const taskAction = "restore-cluster";
      const eventId = this.getUuid();
      this.restore.progress = 0;

      // register to task error
      this.$root.$once(
        `${taskAction}-aborted-${eventId}`,
        this.restoreClusterAborted
      );

      // register to task completion
      this.$root.$once(
        `${taskAction}-completed-${eventId}`,
        this.restoreClusterCompleted
      );

      // register to task progress to update progress bar
      this.$root.$on(
        `${taskAction}-progress-${eventId}`,
        this.restoreClusterProgress
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          extra: {
            title: this.$t("action." + taskAction),
            isNotificationHidden: true,
            isProgressNotified: true,
            toastTimeout: 0, // persistent notification
            eventId,
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.restoreCluster = this.getErrorMessage(err);
        return;
      }
    },
    restoreClusterProgress(progress) {
      this.restore.progress = progress;
    },
    restoreClusterAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.restoreCluster = false;
      this.error.restoreCluster = this.$t("error.generic_error");

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );
    },
    restoreClusterCompleted(taskContext) {
      this.loading.restoreCluster = false;

      // unregister to task progress
      this.$root.$off(
        `${taskContext.action}-progress-${taskContext.extra.eventId}`
      );

      if (this.restore.summary.single_node_cluster) {
        // restore apps
        this.readBackupRepositories();
      } else {
        // skip restore apps
        this.getClusterStatus();
      }
    },
    async readBackupRepositories() {
      this.error.readBackupRepositories = "";
      this.loading.readBackupRepositories = true;
      const taskAction = "read-backup-repositories";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(
        taskAction + "-aborted",
        this.readBackupRepositoriesAborted
      );

      // register to task completion
      this.$root.$off(taskAction + "-completed");
      this.$root.$once(
        taskAction + "-completed",
        this.readBackupRepositoriesCompleted
      );

      const res = await to(
        this.createClusterTask({
          action: taskAction,
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
        this.error.readBackupRepositories = this.getErrorMessage(err);
        return;
      }
    },
    readBackupRepositoriesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.readBackupRepositories = false;
    },
    readBackupRepositoriesCompleted(taskContext, taskResult) {
      let instances = taskResult.output;

      // sort instances by timestamp
      instances.sort(this.sortByProperty("timestamp")).reverse();
      this.restore.instances = instances;

      if (this.restore.instances.length) {
        this.restore.step = "apps";
      } else {
        // there is no app to restore, go to cluster status
        this.getClusterStatus();
      }
      this.loading.readBackupRepositories = false;
    },
    onSelectInstances(selectedInstances) {
      this.restore.selectedInstances = selectedInstances;
    },
    showSkipRestoreAppsModal() {
      this.restore.isShownSkipRestoreAppsModal = true;
    },
    hideSkipRestoreAppsModal() {
      this.restore.isShownSkipRestoreAppsModal = false;
    },
    prepareRestoreModulesData() {
      const data = [];

      for (const i of this.restore.selectedInstances) {
        const instance = {
          repository: i.repository_id,
          path: i.path,
          snapshot: "",
          node: 1,
        };
        data.push(instance);
      }
      return data;
    },
    async restoreModules() {
      const inputData = this.prepareRestoreModulesData();

      this.error.restoreModules = "";
      this.loading.restoreModules = true;
      const taskAction = "restore-modules";

      // register to task error
      this.$root.$off(taskAction + "-aborted");
      this.$root.$once(taskAction + "-aborted", this.restoreModulesAborted);

      // register to task validation
      this.$root.$off(taskAction + "-validation-ok");
      this.$root.$once(
        taskAction + "-validation-ok",
        this.restoreModulesValidationOk
      );
      this.$root.$off(taskAction + "-validation-failed");
      this.$root.$once(
        taskAction + "-validation-failed",
        this.restoreModulesValidationFailed
      );

      // register to task completion
      this.$root.$once(taskAction + "-completed", this.restoreModulesCompleted);

      const res = await to(
        this.createClusterTask({
          action: taskAction,
          data: inputData,
          extra: {
            title: this.$t("action." + taskAction),
            description: this.$tc("init.restoring_apps_c", inputData.length, {
              num: inputData.length,
            }),
          },
        })
      );
      const err = res[0];

      if (err) {
        console.error(`error creating task ${taskAction}`, err);
        this.error.restoreModules = this.getErrorMessage(err);
        return;
      }
    },
    restoreModulesAborted(taskResult, taskContext) {
      console.error(`${taskContext.action} aborted`, taskResult);
      this.loading.restoreModules = false;
    },
    restoreModulesValidationOk() {
      // go to cluster status page
      this.getClusterStatus();
    },
    restoreModulesValidationFailed() {
      this.loading.restoreModules = false;
    },
    restoreModulesCompleted() {
      // update app drawer to show restored instances
      this.$root.$emit("reloadAppDrawer");
    },
    skipRestoreApps() {
      // go to cluster status page
      this.getClusterStatus();
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

.welcome-grid {
  max-width: 56rem;
}

.bx--form .bx--form-item {
  margin-bottom: $spacing-06;
}

.file-uploader {
  margin-bottom: 0 !important;
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

.initialize-cluster .welcome {
  margin-top: 2rem;
  margin-bottom: 4rem;
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
}

.initialize-cluster .tile-description {
  color: $text-02;
}

.initialize-cluster .footer-buttons {
  display: flex;
  justify-content: flex-end;
  margin-top: $spacing-07;

  .cv-button {
    position: relative;
    top: $spacing-05;
    left: $spacing-05;
  }
}
</style>
