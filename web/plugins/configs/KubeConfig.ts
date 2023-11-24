import {
  newContexts,
  newClusters,
  newUsers,
  Context,
  Cluster,
  User,
  ConfigOptions,
} from "@kubernetes/client-node/dist/config_types";
import yaml from "js-yaml";

// KubeConfig
export class KubeConfig {
  /**
   * The list of all known clusters
   */
  public "clusters": Cluster[];

  /**
   * The list of all known users
   */
  public "users": User[];

  /**
   * The list of all known contexts
   */
  public "contexts": Context[];

  /**
   * The name of the current context
   */
  public "currentContext": string;

  constructor() {
    this.contexts = [];
    this.clusters = [];
    this.users = [];
  }

  public getContextObject(name: string): Context | null {
    if (!this.contexts) {
      return null;
    }
    return findObject(this.contexts, name);
  }
  public getCurrentCluster(): Cluster | null {
    const context = this.getCurrentContextObject();
    if (!context) {
      return null;
    }
    return this.getCluster(context.cluster);
  }

  public getCluster(name: string): Cluster | null {
    return findObject(this.clusters, name);
  }

  public getCurrentUser(): User | null {
    const ctx = this.getCurrentContextObject();
    if (!ctx) {
      return null;
    }
    return this.getUser(ctx.user);
  }

  public getUser(name: string): User | null {
    return findObject(this.users, name);
  }

  private getCurrentContextObject(): Context | null {
    return this.getContextObject(this.currentContext);
  }

  public getCurrentContext(): string {
    return this.currentContext;
  }

  public loadFromString(config: string, opts?: Partial<ConfigOptions>): void {
    const obj = yaml.load(config) as any;
    this.clusters = newClusters(obj.clusters, opts);
    this.contexts = newContexts(obj.contexts, opts);
    this.users = newUsers(obj.users, opts);
    this.currentContext = obj["current-context"];
  }
}

export interface Named {
  name: string;
}

function findObject<T extends Named>(list: T[], name: string): T | null {
  if (!list) {
    return null;
  }
  for (const obj of list) {
    if (obj.name === name) {
      return obj;
    }
  }
  return null;
}
