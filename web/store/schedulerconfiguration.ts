import { reactive } from "@nuxtjs/composition-api";
import {
  applySchedulerConfiguration,
  getSchedulerConfiguration,
} from "~/api/v1/schedulerconfiguration";
import { SchedulerConfiguration } from "~/api/v1/types";

type stateType = {
  selectedConfig: selectedConfig | null;
  schedulerconfigurations: SchedulerConfiguration[];
};

type selectedConfig = {
  // isNew represents whether this Config is a new one or not.
  isNew: boolean;
  item: SchedulerConfiguration;
  resourceKind: string;
};

export default function schedulerconfigurationStore() {
  const state: stateType = reactive({
    selectedConfig: null,
    schedulerconfigurations: [],
  });

  return {
    get selected() {
      return state.selectedConfig;
    },

    resetSelected() {
      state.selectedConfig = null;
    },

    select() {
      this.fetchSelected();
    },

    async fetchSelected() {
      const c = await getSchedulerConfiguration();
      if (c) {
        state.selectedConfig = {
          isNew: true,
          item: c,
          resourceKind: "SchedulerConfiguration",
        };
      }
    },

    async fetch() {
      const c = await getSchedulerConfiguration();
      if (c) {
        state.schedulerconfigurations = [c];
      }
    },

    async apply(cfg: SchedulerConfiguration, onError: (_: string) => void) {
      await applySchedulerConfiguration(cfg, onError);
    },

    async delete(_: string) {
      // This function do nothing, but exist to satisfy interface on ResourceBar.vue.
    },
  };
}

export type SchedulerConfigurationStore = ReturnType<
  typeof schedulerconfigurationStore
>;
