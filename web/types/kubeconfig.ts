// APIServerConfigs represents API server settings to connecting and authenticating.
export type APIServerConfigs = {
  url: string;
  clusterCa: string;
  userCert: string;
  userKey: string;
};
