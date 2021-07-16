<template>
  <div class="bx--grid bx--grid--full-width">
    <div class="bx--row">
      <div class="bx--col-lg-16 page-title">
        <h2>{{ $t("status.title") }}</h2>
      </div>
    </div>
    <div class="bx--row">
      <div class="bx--col-lg-16">
        <cv-tile :light="true" class="content-tile">
          <div>Dokuwiki status</div>
        </cv-tile>
      </div>
    </div>
  </div>
</template>

<script>
let nethserver = window.nethserver;

export default {
  name: "Status",
  components: {},
  data() {
    return {
      q: {
        page: "status",
      },
      urlCheckInterval: null,
      isRedirectChecked: false,
      redirectTimeout: 0,
    };
  },
  beforeRouteEnter(to, from, next) {
    console.log(
      "status beforeRouteEnter, parent url",
      window.parent.location.href
    ); ////
    next((vm) => {
      nethserver.watchQueryData(vm);
      vm.urlCheckInterval = nethserver.initUrlBinding(vm, vm.q.page); ////
    });
  },
  beforeRouteLeave(to, from, next) {
    console.log(
      "status beforeRouteLeave, clear urlCheckInterval",
      this.urlCheckInterval
    ); ////
    clearInterval(this.urlCheckInterval);
    next();
  },
  mounted() {
    // show status page after a little delay to avoid page flickering when user directly access a page different from status
    this.redirectTimeout = setTimeout(
      () => (this.isRedirectChecked = true),
      200
    );
  },
  beforeUnmount() {
    clearTimeout(this.redirectTimeout);
    console.log("timeout cleared"); ////
  },
};
</script>
