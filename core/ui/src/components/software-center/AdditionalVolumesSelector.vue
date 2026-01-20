<!--
  Copyright (C) 2026 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="volume-selector">
    <div class="volume-list">
      <!-- skeleton -->
      <div v-if="loading">
        <cv-tile
          v-for="index in 3"
          :key="index"
          :light="light"
          class="volume-tile"
        >
          <cv-skeleton-text
            :paragraph="true"
            :line-count="2"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
      <!-- no volume to restore -->
      <NsEmptyState
        v-else-if="!additionalVolumes.length"
        :title="$t('software_center.no_volume_to_restore')"
      />
      <!-- volume list -->
      <NsTile
        v-else
        v-for="(volume, index) of additionalVolumesLoaded"
        :key="index"
        :light="light"
        kind="selectable"
        v-model="volume.selected"
        value="snapshotValue"
        @click="deselectOtherVolumes(volume)"
        class="volume-tile"
      >
        <div class="mg-bottom-lg">
          <strong v-if="volume.default">
            {{ $t("software_center.default_storage") }}
          </strong>
          <strong v-else>
            {{ volume.label ? volume.label : volume.ui_name }}
          </strong>
        </div>
        <div>{{ volume.path }}</div>
        <div class="mg-bottom-sm">
          {{ (volume.used / (1024 * 1024 * 1024)).toFixed(2) }} GB
          {{ $t("software_center.used") }}
          {{ $t("software_center.of") }}
          {{ (volume.size / (1024 * 1024 * 1024)).toFixed(2) }} GB ({{
            ((volume.used / volume.available) * 100).toFixed(2)
          }}% {{ $t("software_center.used") }})
        </div>
        <NsProgressBar
          :value="loading ? 0 : (volume.used / volume.available) * 100"
          :loading="loading"
          :warningThreshold="70"
          :dangerThreshold="90"
          :height="'10px'"
          :useHealthyColor="false"
          class="mg-bottom-lg"
        />
      </NsTile>
      <infinite-loading
        :identifier="infiniteId"
        @infinite="infiniteScrollHandler"
      ></infinite-loading>
    </div>
  </div>
</template>

<script>
import { UtilService, DateTimeService } from "@nethserver/ns8-ui-lib";
import _cloneDeep from "lodash/cloneDeep";

export default {
  name: "AdditionalVolumesSelector",
  components: {},
  mixins: [UtilService, DateTimeService],
  props: {
    volumes: {
      type: Array,
      required: true,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    light: Boolean,
  },
  data() {
    return {
      additionalVolumes: [],
      // infinite scroll
      additionalVolumesLoaded: [],
      pageNum: 0,
      pageSize: 20,
      infiniteId: +new Date(),
    };
  },
  computed: {
    selectedSnapshot() {
      return this.additionalVolumes.find((i) => i.selected);
    },
  },
  watch: {
    selectedSnapshot: function () {
      this.$emit("selectVolume", this.selectedSnapshot);
    },
    volumes: function () {
      this.updateInternalVolumes();
    },
    additionalVolumes: function () {
      this.additionalVolumesLoaded = [];
      this.pageNum = 0;
      this.infiniteId += 1;
      this.infiniteScrollHandler();
    },
  },
  created() {
    this.updateInternalVolumes();
  },
  methods: {
    updateInternalVolumes() {
      // deep copy (needed to avoid reactivity issues)
      let additionalVolumes = _cloneDeep(this.volumes);
      // select the first additionnal volume by default
      additionalVolumes.forEach((volume, index) => {
        volume.selected = index === 0;
      });
      this.additionalVolumes = additionalVolumes;
    },
    deselectOtherVolumes(volume) {
      for (let s of this.additionalVolumes) {
        if (s.path !== volume.path) {
          s.selected = false;
        }
      }
    },
    infiniteScrollHandler($state) {
      const pageItems = this.additionalVolumes.slice(
        this.pageNum * this.pageSize,
        (this.pageNum + 1) * this.pageSize
      );

      if (pageItems.length) {
        this.pageNum++;
        this.additionalVolumesLoaded.push(...pageItems);

        if ($state) {
          $state.loaded();
        }
      } else {
        if ($state) {
          $state.complete();
        }
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.volume-selector {
  display: flex;
  flex-direction: column;
}

.volume-list {
  overflow-y: auto;
  max-height: 21rem;
}

.ns-tile.volume-tile,
.cv-tile.volume-tile {
  margin-bottom: $spacing-03;
}

.secondary-row {
  margin-top: $spacing-02;
  color: $ui-04;
}
</style>
