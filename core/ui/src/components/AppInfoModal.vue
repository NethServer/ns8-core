<template>
  <div v-if="app">
    <cv-modal
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
          <div>
            <span class="label"
              >{{ $t("software_center.documentation") }}
            </span>
            <span class="value">
              <cv-link :href="app.docs.documentation_url" target="_blank">
                {{ app.docs.documentation_url }}
              </cv-link>
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <span class="label">{{ $t("software_center.bugs") }}</span>
          <span class="value">
            <cv-link :href="app.docs.bug_url" target="_blank">
              {{ app.docs.bug_url }}
            </cv-link>
          </span>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label">{{ $t("software_center.source_code") }}</span>
            <span class="value">
              <cv-link :href="app.docs.code_url" target="_blank">
                {{ app.docs.code_url }}
              </cv-link>
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <div>
            <span class="label"
              >{{ $t("software_center.source_package") }}
            </span>
            <span class="value">
              {{ app.source }}
            </span>
          </div>
        </div>
        <div class="key-value-setting">
          <span class="label"
            >{{ $tc("software_center.authors", app.authors.length) }}
          </span>
          <span class="value">
            <span v-if="app.authors.length == 1"
              >{{ app.authors[0].name }}
              <cv-link :href="'mailto:' + app.authors[0].email" target="_blank">
                {{ app.authors[0].email }}
              </cv-link>
            </span>
            <ul v-else>
              <li
                v-for="(author, index) in app.authors"
                :key="index"
                class="author"
              >
                {{ author.name }}
                <cv-link :href="'mailto:' + author.email" target="_blank">
                  {{ author.email }}
                </cv-link>
              </li>
            </ul>
          </span>
        </div>
      </template>
      <template slot="secondary-button">{{ $t("common.close") }}</template>
    </cv-modal>
  </div>
</template>

<script>
import { UtilService } from "@nethserver/ns8-ui-lib";
import ImageGallery from "../components/ImageGallery";

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
    },
  },
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";

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

.author {
  margin-left: $spacing-05;
}
</style>
