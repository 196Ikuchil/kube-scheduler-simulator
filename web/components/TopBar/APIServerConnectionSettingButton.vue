<template v-if="$config.alphaKubeConfigSetting">
  <v-dialog v-model="data.dialog" width="500">
    <template #activator="{ on }">
      <div>
        <v-btn color="ma-2" :retain-focus-on-click="false" v-on="on">
        API Server
      </v-btn>
      </div>
    </template>

    <v-card>
      <v-card-title class="text-h5 grey lighten-2">
        API server connection setting
      </v-card-title>

      <v-card-text>
        You can configure the API server to connect to and its credentials.
        NOTE: In most cases, we assume that you are using the Kwok. Please paste
        the following commands into this. `kwokctl get kubeconfig | base64`
        <v-text-field
          v-model="data.base64Kubeconfig"
          :rules="[rules.required]"
          label="kubeconfig encoded in base64"
          placeholder="YXBpVmVyc2lvbjogdjEKY2x1c3R......"
        />
      </v-card-text>

      <v-divider></v-divider>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="green darken-1" text @click="data.dialog = false">
          Cancel
        </v-btn>
        <v-btn color="green darken-1" text @click="ApplyConfiguration()">
          Apply
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { defineComponent, inject, reactive, useContext } from "@nuxtjs/composition-api";
import { KubeConfig } from "~/plugins/configs/KubeConfig";
import SnackBarStoreKey from "../StoreKey/SnackBarStoreKey";
import { ExportAPIKey } from "~/api/APIProviderKeys";
import { Buffer } from "buffer";
import { APIServerConfigs } from "~/types/kubeconfig";
import { APIServerSettingsKey } from "~/store/APIserverSettings";

export default defineComponent({
  setup() {
    const apiServerSettingsStore = inject(APIServerSettingsKey);
    const data = reactive({
      dialog: false,
      base64Kubeconfig: "",
    });

    const exportAPI = inject(ExportAPIKey);
    if (!exportAPI) {
      throw new Error(`${ExportAPIKey.description} is not provided`);
    }

    const snackbarstore = inject(SnackBarStoreKey);
    if (!snackbarstore) {
      throw new Error(`${SnackBarStoreKey.description} is not provided`);
    }
    const setServerErrorMessage = (error: string) => {
      snackbarstore.setServerErrorMessage(error);
    };

    const ApplyConfiguration = async () => {
      const strCfg = Buffer.from(data.base64Kubeconfig, "base64").toString();
      const kubeconfig = new KubeConfig();
      kubeconfig.loadFromString(strCfg);
      try {
        const serverCfgs = getAPIServerConfigs(kubeconfig);
        if (!apiServerSettingsStore) {
          throw new Error("apiServerSettingsStore is not injected.")
        }
        apiServerSettingsStore.setNewConfigs(serverCfgs)
      } catch (e: any) {
        setServerErrorMessage(e);
      } finally {
        data.dialog = false;
      }
    };

    // getAPIServerConfigs returns APIServerConfigs which has information of API server for connecting and authenticating.
    // This function is inspired from https://github.com/godaddy/kubernetes-client/blob/2f0676b8b35914fad90365671e9d508e8fb417aa/backends/request/config.js
    const getAPIServerConfigs = (cfg: KubeConfig): APIServerConfigs => {
      const cluster = cfg.getCurrentCluster();
      const user = cfg.getCurrentUser();

      let ca, cert, key, url: string;
      if (!cluster || !cluster || !user) {
        throw new Error(
          "Not enough information in kubeconfig. One or more of the following values is missing. `contexts`, `clusters`, `users`"
        );
      }
      // format: `https://localhost:3000`
      url = cluster.server;
      if (cluster.caData) {
        ca = Buffer.from(cluster.caData, "base64").toString();
      } else {
        console.log(
          "cluster.caData is missing. Currently, the feature to load from file is not implemented."
        );
        throw new Error("cluster.caData is missing.");
      }

      if (user.certData) {
        cert = Buffer.from(user.certData, "base64").toString();
      } else {
        console.log(
          "user.certData is missing. Currently, the feature to load from file is not implemented."
        );
        throw new Error("user.certData is missing.");
      }

      if (user.keyData) {
        key = Buffer.from(user.keyData, "base64").toString();
      } else {
        console.log(
          "user.keyData is missing. Currently, the feature to load from file is not implemented."
        );
        throw new Error("user.keyData is missing.");
      }
      return {
        url: url,
        clusterCa: ca,
        userCert: cert,
        userKey: key,
      } as APIServerConfigs;
    };

    // text form validator
    const required = (v: string) => !!v || "fill out this form";

    return {
      data,
      ApplyConfiguration,
      rules: {
        required,
      },
    };
  },
});
</script>
