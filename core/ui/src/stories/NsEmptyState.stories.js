import { NsEmptyState, GearPictogram } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsEmptyState",
  component: NsEmptyState,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsEmptyState },
  template:
    '<NsEmptyState v-bind="$props">\
      <template v-slot:description>\
        Try changing your search criteria\
      </template>\
    </NsEmptyState>',
});

export const Default = Template.bind({});
Default.args = {
  title: "No search results",
};

const CustomPictogramTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsEmptyState, GearPictogram },
  template:
    '<NsEmptyState v-bind="$props">\
      <template v-slot:pictogram>\
        <GearPictogram />\
      </template>\
      <template v-slot:description>\
        Try changing your search criteria\
      </template>\
    </NsEmptyState>',
});

export const CustomPictogram = CustomPictogramTemplate.bind({});
CustomPictogram.args = {
  title: "No search results",
};

export const LottieAnimation = Template.bind({});
LottieAnimation.args = {
  title: "No search results",
  animationTitle: "ghost",
  loop: 1,
  animationData: require("./assets/ghost.json"),
};
