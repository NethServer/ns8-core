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
              <h3>{{ app.name }}</h3>
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
              {{ app.versions.length ? app.versions[0].tag : "-" }}
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <span class="label"
            >{{ $tc("software_center.authors", app.authors.length) }}
          </span>
          <span class="value">
            <span v-if="app.authors.length == 1"
              >
              <cv-link  :href="'mailto:' + app.authors[0].email" target="_blank">
                {{ app.authors[0].name }}
              </cv-link>
            </span>
            <ul v-else>
              <li
                v-for="(author, index) in app.authors"
                :key="index"
              >
                <cv-link  :href="'mailto:' + author.email" target="_blank">
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
            <span v-if="app.versions[0]['labels']['org.nethserver.images'].split(' ').length == 1"
              >{{ app.versions[0]['labels']['org.nethserver.images'] }}
            </span>
            <ul class="image-list" v-else>
              <li
                v-for="(image, index) in app.versions[0]['labels']['org.nethserver.images'].split(' ')"
                :key="index"
              >
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
            </span> &bull;
            <span>
              <cv-link :href="app.docs.bug_url" target="_blank">
                {{ $t("software_center.bugs") }}
              </cv-link>
            </span> &bull;
            <span>
              <cv-link :href="app.source" target="_blank">
                {{ $t("software_center.source_package") }}
              </cv-link>
            </span> &bull;
            <span>
              <cv-link :href="app.docs.code_url" target="_blank">
                {{ $t("software_center.source_code") }}
              </cv-link>
            </span>
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

export default {
  name: "AppInfoModal",
  components: { ImageGallery },
  mixins: [UtilService],
  props: {
    isShown: Boolean,
    app: { type: [Object, null] },
  },
  methods: {
    getApplicationDescription(app) {
      return this.getAppDescription(app, this);
    },
    getApplicationCategories(app) {
      return this.getAppCategories(app, this);
    }
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
  margin-left: $spacing-04;
}
</style>
