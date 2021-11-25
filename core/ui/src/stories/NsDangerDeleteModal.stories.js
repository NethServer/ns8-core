import { NsDangerDeleteModal } from "@nethserver/ns8-ui-lib";

export default {
  title: "Components/NsDangerDeleteModal",
  component: NsDangerDeleteModal,
  argTypes: {
    testConfirmDelete: { action: "confirmDeleteAction" },
    testHide: { action: "hideAction" },
  },
};

const Template = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsDangerDeleteModal },
  template:
    '<NsDangerDeleteModal v-bind="$props" @confirmDelete="testConfirmDelete(\'test\')" @hide="testHide(\'test\')" />',
});

export const Default = Template.bind({});
Default.args = {
  isShown: true,
  name: "importantObject",
  title: "Confirm deletion",
  warning: "Please read carefully",
  description:
    "Do you really want to delete this object? This action is irreversible",
  typeToConfirm: "",
};
