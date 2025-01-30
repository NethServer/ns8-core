<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div v-if="app">
    <NsModal
      size="default"
      :visible="isShown"
      @modal-hidden="$emit('close')"
      class="no-pad-modal"
    >
      <template slot="title">{{ $t("software_center.app_info") }}</template>
      <!-- v-if="isShown" is needed to reset modal scroll to top -->
      <template v-if="isShown" slot="content">
        <div class="mg-bottom-md">
          <div class="logo-and-name">
            <div class="app-logo">
              <img
                :src="
                  app.logo
                    ? app.logo
                    : require('@/assets/module_default_logo.png')
                "
                :alt="app.name + ' logo'"
              />
            </div>
            <div class="app-name">
              <h3 data-modal-primary-focus>{{ app.name || app.id }}</h3>
            </div>
          </div>
        </div>
        <div v-if="app.screenshots.length" class="mg-bottom-md">
          <div class="screenshots">
            <ImageGallery :fullScreen="true" :images="app.screenshots" />
          </div>
        </div>
        <div class="description">
          {{ getApplicationDescription(app) }}
        </div>
        <!-- warning for rootfull app -->
        <NsInlineNotification
          v-if="app.rootfull && app.certification_level < 3"
          kind="warning"
          :title="$t('software_center.rootfull_app_warning_title')"
          :description="
            $t('software_center.rootfull_app_warning_description', {
              appName: app.name,
            })
          "
          :showCloseButton="false"
        />
        <div class="key-value-setting">
          <span class="label">
            {{ $t("software_center.certification") }}
          </span>
          <span class="value">
            <CertificationLevelBadge :level="app.certification_level" />
          </span>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">{{
              $tc("software_center.categories", app.categories.length)
            }}</span>
            <span class="value">
              {{ getApplicationCategories(app) }}
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">
              {{ $t("software_center.latest_version") }}
            </span>
            <span class="value">
              {{
                app.versions.length
                  ? app.versions[0].tag
                  : $t("common.not_available")
              }}
              <span v-if="app.upstream_name"> ({{ app.upstream_name }}) </span>
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">{{ $t("software_center.repository") }}</span>
            <span class="value">
              {{ app.repository || "-" }}
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <span class="label"
            >{{ $tc("software_center.authors", app.authors.length) }}
          </span>
          <span class="value">
            <span v-if="app.authors.length == 1">
              <cv-link :href="'mailto:' + app.authors[0].email" target="_blank">
                {{ app.authors[0].name }}
              </cv-link>
            </span>
            <ul v-else>
              <li v-for="(author, index) in app.authors" :key="index">
                <cv-link :href="'mailto:' + author.email" target="_blank">
                  {{ author.name }}
                </cv-link>
              </li>
            </ul>
          </span>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">
              {{ $t("software_center.images_label") }}
            </span>
            <span class="value">
              <ul class="image-list unordered-list">
                <li>
                  {{ app.source }}:{{
                    app.versions.length ? app.versions[0].tag : "-"
                  }}
                </li>
                <li v-for="(image, index) in appImages" :key="index">
                  {{ image }}
                </li>
              </ul>
            </span>
          </div>
        </div>
        <div class="info-divider"></div>
        <div class="key-value-setting">
          <div>
            <span>
              <cv-link :href="app.docs.documentation_url" target="_blank">
                {{ $t("software_center.documentation") }}
              </cv-link>
            </span>
            &bull;
            <span>
              <cv-link :href="app.docs.bug_url" target="_blank">
                {{ $t("software_center.bugs") }}
              </cv-link>
            </span>
            &bull;
            <span>
              <cv-link :href="app.docs.code_url" target="_blank">
                {{ $t("software_center.source_code") }}
              </cv-link>
            </span>
            <template v-if="app.docs.terms_url">
              &bull;
              <span>
                <cv-link :href="app.docs.terms_url" target="_blank">
                  {{ $t("common.terms_and_conditions") }}
                </cv-link>
              </span>
            </template>
            <template v-if="app.docs.relnotes_url">
              &bull;
              <span>
                <cv-link :href="app.docs.relnotes_url" target="_blank">
                  {{ $t("common.release_notes") }}
                </cv-link>
              </span>
            </template>
          </div>
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
    </NsModal>
  </div>
</template>

<script>
import { UtilService } from "@nethserver/ns8-ui-lib";
import ImageGallery from "./ImageGallery";
import CertificationLevelBadge from "./CertificationLevelBadge.vue";

export default {
  name: "AppInfoModal",
  components: { ImageGallery, CertificationLevelBadge },
  mixins: [UtilService],
  props: {
    isShown: Boolean,
    app: { type: [Object, null] },
  },
  computed: {
    appImages() {
      if (
        this.app &&
        this.app.versions.length &&
        this.app.versions[0]["labels"]["org.nethserver.images"]
      ) {
        return this.app.versions[0]["labels"]["org.nethserver.images"]
          .trim()
          .replace(/\s+/g, " ")
          .split(" ");
      } else {
        return [];
      }
    },
  },
  methods: {
    getApplicationDescription(app) {
      return this.getAppDescription(app, this);
    },
    getApplicationCategories(app) {
      return this.getAppCategories(app, this);
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.logo-and-name {
  display: flex;
  align-items: center;
  margin-top: $spacing-05;
  margin-bottom: $spacing-05;
}

.app-logo {
  width: 4rem;
  height: 4rem;
  margin-right: $spacing-05;
  flex-shrink: 0;
}

.app-logo img {
  width: 100%;
  height: 100%;
}

.description {
  margin-bottom: $spacing-06;
}

.info-divider {
  margin-top: $spacing-02;
  margin-bottom: $spacing-03;
  border-bottom: 1px solid $ui-04;
}

.image-list {
  margin-left: $spacing-03;
}
</style>
