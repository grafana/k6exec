Run k6 with a seamless extension user experience.

`k6exec` is a [k6] launcher that automatically provides [k6] with the [extensions] used by the test. In order to do this, it analyzes the script arguments of the `run` and `archive` subcommands, detects the extensions to be used and their version constraints.

The launcher acts as a drop-in replacement for the `k6` command. For more convenient use, it is advisable to create an alias or shell script called `k6` for the launcher. The alias can be used in exactly the same way as the `k6` command, with the difference that it generates the real `k6` on the fly based on the extensions you want to use.

Any k6 command can be used. Use the `help` command to list the available k6 commands.

Since k6exec tries to emulate the `k6` command line, the `help` command or the `--help` flag cannot be used to display help from `k6exec` command itself. The `k6exec` help can be displayed using the `--usage` flag:

    k6exec --usage

### Prerequisites

k6exec tries to provide the appropriate k6 executable after detecting the extension dependencies. This can be done using a build service or a native builder.

#### Build Service

No additional installation is required to use the build service, just provide the build service URL.

The build service URL can be specified in the `K6_BUILD_SERVICE_URL` environment variable or by using the `--build-service-url` flag.

If the build service requires authentication, you can specify the authentication token using the `K6_BUILD_SERVICE_AUTH` environment variable.

If the `k6_BUILD_SERVICE_URL` is not specified, `k6exec` tries to use the build service provided by Grafana Cloud K6 using the credential obtained from the [k6 cloud login](https://grafana.com/docs/grafana-cloud/testing/k6/author-run/tokens-and-cli-authentication/) command. You can also provide this credentials using the `K6_CLOUD_TOKEN` environment variable.

### Dependencies

Dependencies can come from three sources: k6 test script, manifest file, `K6_DEPENDENCIES` environment variable. Instead of these three sources, a k6 archive can also be specified, which can contain all three sources.

#### Pragma

Version constraints can be specified using the JavaScript `"use ..."` pragma syntax for k6 and extensions. Put the following lines at the beginning of the test script:

```js
"use k6 >= v0.52";
"use k6 with k6/x/faker > 0.2";
```

Any number of `"use k6"` pragmas can be used.

> **Note**
> The use of pragmas is completely optional for JavaScript type extensions, it is only necessary if you want to specify version constraints.

The pragma syntax can also be used to specify an extension dependency that is not referenced in an import expression. A typical example of this is the Output type extension such as [xk6-top]:

```js
"use k6 with top >= 0.1";
```

Read the version constraints syntax in the [Version Constraints](#version-constraints) section

#### Environment

The extensions to be used and optionally their version constraints can also be specified in the `K6_DEPENDENCIES` environment variable. The value of the environment variable K6_DEPENDENCIES is a list of elements separated by semicolons. Each element specifies an extension (or k6 itself) and optionally its version constraint.

```
k6>=0.52;k6/x/faker>=0.3;k6/x/sql>=0.4
```

#### Manifest

The manifest file is a JSON file, the `dependencies` property of which can specify extension dependencies and version constraints. The value of the `dependencies` property is a JSON object. The property names of this object are the extension names (or k6) and the values ​​are the version constraints.

```json
{
  "dependencies": {
    "k6": ">=0.52",
    "k6/x/faker": ">=0.3",
    "k6/x/sql": ">=0.4"
  }
}
```

The manifest file is a file named `package.json`, which is located closest to the k6 test script or the current directory, depending on whether the given subcommand has a test script argument (e.g. run, archive) or not (e.g. version). The `package.json` file is searched for up to the root of the directory hierarchy.

### Limitations

Version constraints can be specified in several sources ([pragma](#pragma), [environment](#environment), [manifest](#manifest)) but cannot be overwritten. That is, for a given extension, the version constraints from different sources must either be equal, or only one source can contain a version constraint.

[k6]: https://k6.io
[extensions]: https://grafana.com/docs/k6/latest/extensions/
[xk6-top]: https://github.com/szkiba/xk6-top
[Masterminds/semver]: https://github.com/Masterminds/semver
