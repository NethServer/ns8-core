import NsEmptyState from "../components/NsEmptyState.vue";
import Gear from "../components/pictograms/Gear";

export default {
  title: "Components/NsEmptyState",
  component: NsEmptyState,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsEmptyState },
  template: '<NsEmptyState v-bind="$props"></NsEmptyState>',
});

export const Default = Template.bind({});
Default.args = {
  title: "No search results",
  description: "Try changing your search criteria",
};

const CustomPictogramTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsEmptyState, Gear },
  template: '<NsEmptyState v-bind="$props"><Gear /></NsEmptyState>',
});

export const CustomPictogram = CustomPictogramTemplate.bind({}); ////
CustomPictogram.args = {
  title: "No search results",
  description: "Try changing your search criteria",
};
