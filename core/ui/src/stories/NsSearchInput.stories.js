import NsSearchInput from "../components/NsSearchInput.vue";

//// delete file?

export default {
  title: "Components/NsSearchInput",
  component: NsSearchInput,
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsSearchInput },
  template: '<NsSearchInput v-bind="$props"></NsSearchInput>',
});

export const Default = Template.bind({});
Default.args = {
  placeholder: "Search...",
  label: "",
  value: "",
  helperText: "",
  invalidMessage: "",
};
