<template>
  <div v-if="1 < 10">sdf</div>
  <div v-else>sdf</div>
</template>

<script>

//// use JSX?

export default {
  name: "Icon",
  props: {
    name: {
      type: String,
      required: true,
    },
    width: {
      type: [Number, String],
      default: 64,
    },
    height: {
      type: [Number, String],
      default: 64,
    },
    color: {
      type: String,
      default: "currentColor",
    },
  },
};
</script>
