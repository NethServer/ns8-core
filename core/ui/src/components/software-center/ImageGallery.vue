<!--
  Copyright (C) 2023 Nethesis S.r.l.
  SPDX-License-Identifier: GPL-3.0-or-later
-->
<template>
  <div>
    <VueGallery
      :images="images"
      :index="index"
      @close="index = null"
      :options="{ fullScreen: fullScreen }"
    ></VueGallery>
    <div class="previews">
      <div
        class="image"
        v-for="(image, imageIndex) in images"
        :key="imageIndex"
        @click="index = imageIndex"
        :style="{
          backgroundImage: 'url(' + image + ')',
        }"
      ></div>
    </div>
  </div>
</template>

<script>
import VueGallery from "vue-gallery";

export default {
  name: "ImageGallery",
  components: { VueGallery },
  props: {
    images: Array,
    fullScreen: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      index: null,
    };
  },
  created() {
    document.addEventListener("fullscreenchange", this.fullScreenChange, false);
    document.addEventListener(
      "mozfullscreenchange",
      this.fullScreenChange,
      false
    );
    document.addEventListener(
      "MSFullscreenChange",
      this.fullScreenChange,
      false
    );
    document.addEventListener(
      "webkitfullscreenchange",
      this.fullScreenChange,
      false
    );
  },
  beforeDestroy() {
    document.removeEventListener("fullscreenchange", this.fullScreenChange);
    document.removeEventListener("mozfullscreenchange", this.fullScreenChange);
    document.removeEventListener("MSFullscreenChange", this.fullScreenChange);
    document.removeEventListener(
      "webkitfullscreenchange",
      this.fullScreenChange
    );
  },
  methods: {
    fullScreenChange() {
      if (
        document.webkitIsFullScreen === false ||
        document.mozFullScreen === false ||
        document.msFullscreenElement === false
      ) {
        // fullscreen exit, close gallery
        this.index = null;
      }
    },
  },
};
</script>

<style scoped lang="scss">
@import "../../styles/carbon-utils";

.previews {
  display: flex;
  overflow-x: auto;
}

.image {
  width: 300px;
  height: 200px;
  flex-shrink: 0;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center center;
  border: 1px solid $ui-03;
  margin: 0 $spacing-03 $spacing-03 0;
  cursor: pointer;
}

.image:last-child {
  margin-right: 0;
}
</style>
