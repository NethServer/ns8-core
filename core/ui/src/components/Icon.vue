<template>
  <div>icon test</div>
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
