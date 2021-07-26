<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :width="width"
    :height="height"
    viewBox="0 0 32 32"
    :aria-labelledby="title"
    role="presentation"
  >
    <title :id="title" lang="en">{{ title }} icon</title>
    <g :fill="color">
      <slot />
    </g>
  </svg>
</template>

<script>
//// FILE COPIED FROM CORE

export default {
  name: "Pictogram",
  props: {
    title: {
      type: String,
      default: "",
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
