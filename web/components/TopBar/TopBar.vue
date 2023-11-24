<template>
  <v-sheet class="transparent">
    <APIServerConnectionSettingButton />
    <SchedulerConfigurationEditButton />
    <ExportButton />
    <ImportButton />
    <ResetButton />
  </v-sheet>
</template>

<script lang="ts">
import { defineComponent } from "@nuxtjs/composition-api";
import APIServerConnectionSettingButton from "./APIServerConnectionSettingButton.vue";
import SchedulerConfigurationEditButton from "./SchedulerConfigurationEditButton.vue";
import ExportButton from "./ExportButton.vue";
import ImportButton from "./ImportButton.vue";
import ResetButton from "./ResetButton.vue";

export default defineComponent({
  components: {
    APIServerConnectionSettingButton,
    ImportButton,
    ExportButton,
    SchedulerConfigurationEditButton,
    ResetButton,
  },
});
</script>
