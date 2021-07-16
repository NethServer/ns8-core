<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("about.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">Dokuwiki about</cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
let nethserver = window.nethserver;

export default {
  name: "About",
  components: {},
  data() {
    return {
      q: {
        page: "about",
      },
      currentUrl: "",
      urlCheckInterval: null,
    };
  },
  beforeRouteEnter(to, from, next) {
    console.log(
      "about beforeRouteEnter, parent url",
      window.parent.location.href
    ); ////
    next((vm) => {
      nethserver.watchQueryData(vm);
      vm.urlCheckInterval = nethserver.initUrlBinding(vm, vm.q.page); ////
    });
  },
  beforeRouteLeave(to, from, next) {
    console.log(
      "about beforeRouteLeave, clear urlCheckInterval",
      this.urlCheckInterval
    ); ////
    clearInterval(this.urlCheckInterval);
    next();
  },
};
</script>
