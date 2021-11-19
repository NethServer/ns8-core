import { NsLottieAnimation } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsLottieAnimation",
  component: NsLottieAnimation,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsLottieAnimation },
  template: '<NsLottieAnimation v-bind="$props" style="width: 64px"/>',
});

export const Default = Template.bind({});
Default.args = {
  animateOnHover: true,
  loop: false,
  autoPlay: true,
  refName: "rocket",
  animationData: require("./assets/rocket.json"),
};
