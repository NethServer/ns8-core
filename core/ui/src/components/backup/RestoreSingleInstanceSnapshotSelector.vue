<!--
  Copyright (C) 2024 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div class="snapshot-selector">
    <div class="snapshot-list">
      <!-- skeleton -->
      <div v-if="loading">
        <cv-tile
          v-for="index in 3"
          :key="index"
          :light="light"
          class="snapshot-tile"
        >
          <cv-skeleton-text
            :paragraph="true"
            :line-count="2"
          ></cv-skeleton-text>
        </cv-tile>
      </div>
      <!-- no snapshot to restore -->
      <NsEmptyState
        v-else-if="!internalSnapshots.length"
        :title="$t('backup.no_snapshot_to_restore')"
      />
      <!-- snapshot list -->
      <NsTile
        v-else
        v-for="(snapshot, index) of internalSnapshotsLoaded"
        :key="index"
        :light="light"
        kind="selectable"
        v-model="snapshot.selected"
        value="snapshotValue"
        @click="deselectOtherSnapshots(snapshot)"
        class="snapshot-tile"
      >
        <div>
          {{ formatDate(new Date(snapshot.timestamp * 1000), "PPpp") }}
        </div>
        <div v-if="index == 0" class="secondary-row">
          {{ $t("backup.most_recent") }}
        </div>
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
  name: "RestoreSingleInstanceSnapshotSelector",
  components: {},
  mixins: [UtilService, DateTimeService],
  props: {
    snapshots: {
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
      internalSnapshots: [],
      // infinite scroll
      internalSnapshotsLoaded: [],
      pageNum: 0,
      pageSize: 20,
      infiniteId: +new Date(),
    };
  },
  computed: {
    selectedSnapshot() {
      return this.internalSnapshots.find((i) => i.selected);
    },
  },
  watch: {
    selectedSnapshot: function () {
      this.$emit("select", this.selectedSnapshot);
    },
    snapshots: function () {
      this.updateInternalSnapshots();
    },
    internalSnapshots: function () {
      this.internalSnapshotsLoaded = [];
      this.pageNum = 0;
      this.infiniteId += 1;
      this.infiniteScrollHandler();
    },
  },
  created() {
    this.updateInternalSnapshots();
  },
  methods: {
    updateInternalSnapshots() {
      // deep copy (needed to avoid reactivity issues)
      let internalSnapshots = _cloneDeep(this.snapshots);

      for (const snapshot of internalSnapshots) {
        snapshot.selected = false;
      }
      this.internalSnapshots = internalSnapshots;
    },
    deselectOtherSnapshots(snapshot) {
      for (let s of this.internalSnapshots) {
        if (s.id !== snapshot.id) {
          s.selected = false;
        }
      }
    },
    infiniteScrollHandler($state) {
      const pageItems = this.internalSnapshots.slice(
        this.pageNum * this.pageSize,
        (this.pageNum + 1) * this.pageSize
      );

      if (pageItems.length) {
        this.pageNum++;
        this.internalSnapshotsLoaded.push(...pageItems);

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

.snapshot-selector {
  display: flex;
  flex-direction: column;
}

.snapshot-list {
  overflow-y: auto;
  max-height: 21rem;
}

.ns-tile.snapshot-tile,
.cv-tile.snapshot-tile {
  margin-bottom: $spacing-03;
}

.secondary-row {
  margin-top: $spacing-02;
  color: $ui-04;
}
</style>
