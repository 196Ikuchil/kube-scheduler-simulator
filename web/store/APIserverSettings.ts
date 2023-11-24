import { reactive, Ref } from "@nuxtjs/composition-api";
import {APIServerConfigs } from "~/types/kubeconfig";
import { InjectionKey } from "@nuxtjs/composition-api";

type stateType = {
    cfgs: APIServerConfigs | null;
  };

export default function apiServerConfigsStore() {
    const initialState: stateType = {
        cfgs: null,
    };

    const state: stateType = reactive({ ...initialState });
    return {
        get cfgs(): APIServerConfigs | null {
            return state.cfgs
        },
        setNewConfigs(cfgs: APIServerConfigs) {
            state.cfgs = cfgs
        }
    }
}

export type APIServerConfigsStore = ReturnType<typeof apiServerConfigsStore>;

export const APIServerSettingsKey: InjectionKey<APIServerConfigsStore> = Symbol("APIServerSettings")
