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
          <span class="setting-value">{{
            repository.parameters.b2_account_key
          }}</span>
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
          <span class="setting-label">{{
            $t("backup.aws_secret_access_key")
          }}</span>
          <span class="setting-value">{{
            repository.parameters.aws_secret_access_key
          }}</span>
        </div>
      </template>
      <!-- //// handle all providers -->
      <!-- password -->
      <div class="mg-bottom-sm">
        <span class="setting-label">{{
          $t("backup.repository_password")
        }}</span>
        <cv-link @click="togglePassword" class="toggle-password">
          {{ isPasswordShown ? $t("common.hide") : $t("common.show") }}
        </cv-link>
        <NsCodeSnippet
          v-if="isPasswordShown"
          :copyTooltip="$t('common.copy_to_clipboard')"
          :copy-feedback="$t('common.copied_to_clipboard')"
          :feedback-aria-label="$t('common.copied_to_clipboard')"
          :wrap-text="true"
          :moreText="$t('common.show_more')"
          :lessText="$t('common.show_less')"
          light
          hideExpandButton
          >{{ repository.password }}</NsCodeSnippet
        >
      </div>
    </template>
    <template slot="primary-button">{{ $t("common.close") }}</template>
  </cv-modal>
</template>

<script>
export default {
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
      isPasswordShown: false,
    };
  },
  watch: {
    isShown: function () {
      if (this.isShown) {
        // hide password
        this.isPasswordShown = false;
      }
    },
  },
  methods: {
    togglePassword() {
      this.isPasswordShown = !this.isPasswordShown;
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.setting-label {
  display: inline-block;
  margin-right: $spacing-03;
  margin-bottom: $spacing-02;
  font-weight: bold;
}

.setting-value {
  word-wrap: break-word;
}
</style>
