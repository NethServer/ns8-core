<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <h1 class="page-title">{{ $t("software_center.title") }}</h1>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-md-5">
        <cv-tile :light="true" class="content-tile">
          <h2>Main content</h2>
          <p>This is some tile content</p>
        </cv-tile>
      </div>
      <div class="bx--col-md-3">
        <cv-tile :light="true" class="content-tile">
          <h2>Side content</h2>
          <p>This is some tile content</p>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "SoftwareCenter",
};
</script>

<style scoped lang="scss">
@import "../styles/carbon-utils";
</style>
