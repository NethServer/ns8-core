<!--
  Copyright (C) 2024 Nethesis S.r.l.
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
        v-else-if="!additionnalVolumes.length"
        :title="$t('software_center.no_volume_to_restore')"
      />
      <!-- volume list -->
      <NsTile
        v-else
        v-for="(volume, index) of additionnalVolumesLoaded"
        :key="index"
        :light="light"
        kind="selectable"
        v-model="volume.selected"
        value="snapshotValue"
        @click="deselectOtherVolumes(volume)"
        class="volume-tile"
      >
        <div>
          {{ volume.path }} - {{ (volume.size / (1024 * 1024 * 1024)).toFixed(2) }} GB
        </div>
        <!-- <div v-if="index == 0" class="secondary-row">
          {{ $t("backup.most_recent") }}
        </div> -->
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
  name: "AdditionnalVolumesSelector",
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
      additionnalVolumes: [],
      // infinite scroll
      additionnalVolumesLoaded: [],
      pageNum: 0,
      pageSize: 20,
      infiniteId: +new Date(),
    };
  },
  computed: {
    selectedSnapshot() {
      return this.additionnalVolumes.find((i) => i.selected);
    },
  },
  watch: {
    selectedSnapshot: function () {
      this.$emit("select", this.selectedSnapshot);
    },
    volumes: function () {
      this.updateInternalVolumes();
    },
    additionnalVolumes: function () {
      this.additionnalVolumesLoaded = [];
      this.pageNum = 0;
      this.infiniteId += 1;
      this.infiniteScrollHandler();
    },
  },
  created() {
    this.updateInternalVolumes();
  },
  methods: {
    // formatSnapshotTimestamp(timestamp) {
    //   return new Intl.DateTimeFormat(navigator.language, {
    //     dateStyle: "long",
    //     timeStyle: "short",
    //   }).format(new Date(timestamp * 1000));
    // },
    updateInternalVolumes() {
      // deep copy (needed to avoid reactivity issues)
      let additionnalVolumes = _cloneDeep(this.volumes);

      for (const volume of additionnalVolumes) {
        volume.selected = false;
      }
      this.additionnalVolumes = additionnalVolumes;
    },
    deselectOtherVolumes(volume) {
      for (let s of this.additionnalVolumes) {
        if (s.path !== volume.path) {
          s.selected = false;
        }
      }
    },
    infiniteScrollHandler($state) {
      const pageItems = this.additionnalVolumes.slice(
        this.pageNum * this.pageSize,
        (this.pageNum + 1) * this.pageSize
      );

      if (pageItems.length) {
        this.pageNum++;
        this.additionnalVolumesLoaded.push(...pageItems);

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
