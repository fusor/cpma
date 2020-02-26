## Generating pre-migration report with CPMA

---

### 1. Prerequisites

Prior to working with CPMA you need to deploy an OCP 3.7+ cluster(OCP 3.7, 3.9, 3.10, 3.11 are supported).
In order to generate a report, CPMA interacts with OCP using it's API. This means KUBECONFIG is required as a well as user used to talk with OCP cluster API needs to have privileges to list nodes, projects, pods etc. Recommended cluster role is `cluster-admin`. It can be configured with following command: `oc adm policy add-cluster-role-to-user cluster-admin <username>`.

---

### 2. General CPMA configuration

CPMA can be configured using either:

1. Interactive prompts. This is simillar to `openshift-install`. The tool can be run with no configuration, all required values will be prompted. You can see an example below. This is the most recommended way, because prompts will guide you through needed values and it can generate a configuration based on prompted values that will be used later.

![prompt](https://user-images.githubusercontent.com/20123872/60581251-c0f57100-9d86-11e9-9ab3-7681b840731a.gif)


2. CLI parameters. All configuration values can be passed using CLI parameters. For example: `./cpma --source cluster.example.com --work-dir ./dir` Refer to CPMA's [README.md](https://github.com/konveyor/cpma#usage) for full list of parameters.

3. Predefined configuration file. You can manually create a yaml configuration based on this [example](https://github.com/konveyor/cpma/blob/master/examples/cpma-config.example.yaml). Configuration file path can be passed using `--config` parameter, or place `cpma.yaml` in your home directory.

4. Environmental variables. It is also possible to pass all configuration values as environmental variables. List of variables can be found in [README.md](https://github.com/konveyor/cpma#e2e-tests)

---

### 3. Using CPMA to generate report

Once the configuration has been provided, by either prompt, CLI or ENV parameters, or configuration file, the report will 
be generated.

---

## 4. Reading report.json

Generated report will be placed inside specified working directory in format of a json file. We are planning to visualize report in one of future mig-ui versions.

You can find example report.json in this scenario directory.