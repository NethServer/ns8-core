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
  isErrorShown: false,
  errorTitle: "",
  errorDescription: "",
};

const DescriptionSlotTemplate = (args, { argTypes }) => ({
  props: Object.keys(argTypes),
  components: { NsDangerDeleteModal },
  template: `<NsDangerDeleteModal v-bind="$props" @confirmDelete="testConfirmDelete('test')" @hide="testHide('test')">
      <template slot="description">
        Description slot
      </template>
      </NsDangerDeleteModal>`,
});

export const DescriptionSlot = DescriptionSlotTemplate.bind({});
DescriptionSlot.args = {
  isShown: true,
  name: "importantObject",
  title: "Confirm deletion",
  warning: "Please read carefully",
  description: "",
  typeToConfirm: "",
  isErrorShown: false,
  errorTitle: "",
  errorDescription: "",
};
