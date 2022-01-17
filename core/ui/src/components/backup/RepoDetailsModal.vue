<template>
  <cv-modal size="default" :visible="isShown" @modal-hidden="$emit('hide')">
    <template slot="title">{{ $t("backup.repository_details") }}</template>
    <template slot="content">
      <div class="mg-bottom-sm">
        <span class="setting-label">{{ $t("backup.name") }}</span>
        <span class="setting-value">{{ repository.name }}</span>
      </div>
      <div class="mg-bottom-sm">
        <span class="setting-label">{{ $t("backup.provider") }}</span>
        <span class="setting-value">{{
          $t("backup." + repository.provider)
        }}</span>
      </div>
      <div class="mg-bottom-sm">
        <span class="setting-label">{{ $t("backup.url") }}</span>
        <span class="setting-value">{{ repository.url }}</span>
      </div>
      <!-- backblaze -->
      <template v-if="repository.provider == 'backblaze'">
        <div class="mg-bottom-sm">
          <span class="setting-label">{{ $t("backup.b2_account_id") }}</span>
          <span class="setting-value">{{
            repository.parameters.b2_account_id
          }}</span>
        </div>
        <div class="mg-bottom-sm">
          <span class="setting-label">{{ $t("backup.b2_account_key") }}</span>
          <cv-link @click="toggleBackblazeAccountKey">
            {{
              isShownBackblazeAccountKey ? $t("common.hide") : $t("common.show")
            }}
          </cv-link>
          <NsCodeSnippet
            v-if="isShownBackblazeAccountKey"
            :copyTooltip="$t('common.copy_to_clipboard')"
            :copy-feedback="$t('common.copied_to_clipboard')"
            :feedback-aria-label="$t('common.copied_to_clipboard')"
            :wrap-text="true"
            :moreText="$t('common.show_more')"
            :lessText="$t('common.show_less')"
            light
            hideExpandButton
            class="password-snippet"
            >{{ repository.parameters.b2_account_key }}</NsCodeSnippet
          >
        </div>
      </template>
      <!-- amazon s3 -->
      <template v-if="repository.provider == 'aws'">
        <div class="mg-bottom-sm">
          <span class="setting-label">{{
            $t("backup.aws_access_key_id")
          }}</span>
          <span class="setting-value">{{
            repository.parameters.aws_access_key_id
          }}</span>
        </div>
        <div class="mg-bottom-sm">
          <span class="setting-label">{{
            $t("backup.aws_default_region")
          }}</span>
          <span class="setting-value">{{
            repository.parameters.aws_default_region
          }}</span>
        </div>
        <div class="mg-bottom-sm">
          <span class="setting-label">
            {{ $t("backup.aws_secret_access_key") }}
          </span>
          <cv-link @click="toggleAwsSecretAccessKey">
            {{
              isShownAwsSecretAccessKey ? $t("common.hide") : $t("common.show")
            }}
          </cv-link>
          <NsCodeSnippet
            v-if="isShownAwsSecretAccessKey"
            :copyTooltip="$t('common.copy_to_clipboard')"
            :copy-feedback="$t('common.copied_to_clipboard')"
            :feedback-aria-label="$t('common.copied_to_clipboard')"
            :wrap-text="true"
            :moreText="$t('common.show_more')"
            :lessText="$t('common.show_less')"
            light
            hideExpandButton
            class="password-snippet"
            >{{ repository.parameters.aws_secret_access_key }}</NsCodeSnippet
          >
        </div>
      </template>
      <!-- generic s3 -->
      <template v-if="repository.provider == 'genericS3'">
        generic s3 ////
      </template>
      <!-- azure -->
      <template v-if="repository.provider == 'azure'"> azure //// </template>
      <!-- //// handle all providers -->
      <!-- password -->
      <div class="mg-bottom-sm">
        <span class="setting-label">
          {{ $t("backup.repository_password") }}
          <cv-interactive-tooltip
            alignment="center"
            direction="top"
            class="info"
          >
            <template slot="trigger">
              <Information16 />
            </template>
            <template slot="content">
              <div>{{ $t("backup.repo_password_tooltip") }}</div>
            </template>
          </cv-interactive-tooltip>
        </span>
        <cv-link @click="togglePassword">
          {{ isShownRepoPassword ? $t("common.hide") : $t("common.show") }}
        </cv-link>
        <NsCodeSnippet
          v-if="isShownRepoPassword"
          :copyTooltip="$t('common.copy_to_clipboard')"
          :copy-feedback="$t('common.copied_to_clipboard')"
          :feedback-aria-label="$t('common.copied_to_clipboard')"
          :wrap-text="true"
          :moreText="$t('common.show_more')"
          :lessText="$t('common.show_less')"
          light
          hideExpandButton
          class="password-snippet"
          >{{ repository.password }}</NsCodeSnippet
        >
      </div>
    </template>
    <template slot="primary-button">{{ $t("common.close") }}</template>
  </cv-modal>
</template>

<script>
import Information16 from "@carbon/icons-vue/es/information/16";

export default {
  components: { Information16 },
  name: "RepoDetailsModal",
  props: {
    isShown: {
      type: Boolean,
      default: true,
    },
    repository: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isShownRepoPassword: false,
      isShownAwsSecretAccessKey: false,
      isShownBackblazeAccountKey: false,
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // hide secrets
        this.isShownRepoPassword = false;
        this.isShownAwsSecretAccessKey = false;
        this.isShownBackblazeAccountKey = false;
      }
    },
  },
  methods: {
    togglePassword() {
      this.isShownRepoPassword = !this.isShownRepoPassword;
    },
    toggleAwsSecretAccessKey() {
      this.isShownAwsSecretAccessKey = !this.isShownAwsSecretAccessKey;
    },
    toggleBackblazeAccountKey() {
      this.isShownBackblazeAccountKey = !this.isShownBackblazeAccountKey;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.setting-label {
  display: inline-block;
  margin-right: $spacing-04;
  margin-bottom: $spacing-02;
  font-weight: bold;
}

.setting-value {
  word-wrap: break-word;
}

.password-snippet {
  margin-bottom: $spacing-07;
}
</style>
