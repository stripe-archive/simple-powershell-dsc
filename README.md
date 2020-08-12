# simple-powershell-dsc

This is a simple implementation of the Powershell DSC[¹][1] pull server[²][2]
in Go. It is intended to allow applying simple to moderately-complex
configuration to Windows in a reliable, cross-platform manner.

[1]: https://docs.microsoft.com/en-us/powershell/scripting/dsc/overview/overview?view=powershell-7
[2]: https://docs.microsoft.com/en-us/powershell/scripting/dsc/pull-server/enactingconfigurations?view=powershell-7

## Brief overview

There are three major interfaces that define a pull server:

- `ConfigurationRepository`, which is responsible for allowing a client to
  fetch a configuration or module from the server.
- `ReportServer`, which is responsible for receiving or returning "reports" on
  what the Powershell DSC agent is doing.
- `NodeStatus`, which is responsible for storing and reconciling the status of
  a Powershell DSC agent on a node.

These three interfaces are defined in the file `dsc/interfaces.go`, and are
implemented by various packages in the subdirectories `dsc/config`,
`dsc/report`, and `dsc/status`.

For testing, there is a simple HTTP server package under `cmd/http/` that uses
local filesystem-backed storage for all three; you can put configuration under
`test/config` and it will be served to clients that request it. For example:

```
$ tree test/config
test/config
|-- config
|   |-- HelloWorld
|   `-- WindowsHardening
`-- modules
    |-- AuditPolicyDsc
    |   `-- 1.3.0.0.zip
    `-- SecurityPolicyDsc
        `-- 2.6.0.0.zip
```

For actual deployment, the directory `cmd/lambda/` contains a AWS Lambda
package that serves configuration, stores registration, and saves reports in a
S3 bucket.

## Using a pull server

The following DSC configuration can be used to instruct a Windows client to
pull configuration from a given URL:

```powershell
[DSCLocalConfigurationManager()]
configuration DscPull
{
    $server = 'https://URL-HERE'

    Node localhost
    {
        Settings
        {
            RefreshMode = 'Pull'
            RefreshFrequencyMins = 30
            RebootNodeIfNeeded = $true
            DebugMode = 'All'
        }

        ReportServerWeb TestReportServer
        {
            ServerURL = "$server/"
            RegistrationKey = 'faee15cf-3403-41e7-8006-d7f2f86afc72'
        }

        ConfigurationRepositoryWeb TestPullServer
        {
            ServerURL = "$server/"
            RegistrationKey = 'faee15cf-3403-41e7-8006-d7f2f86afc72'
            ConfigurationNames = @('HelloWorld')
        }
    }
}

DscPull -OutputPath .\
```

When running this Powershell file, an output MOF will be generated in the
current directory for each `Node`, in this case, `localhost.meta.mof`. You can
then apply that MOF file (called a "metaconfiguration MOF" by Microsoft) to a
Windows node with the following command:

```powershell
Set-DSCLocalConfigurationManager –ComputerName localhost –Path .\path\to\directory\with\the-mof\ –Verbose -Wait
```

Some useful commands are:
- `Get-DscLocalConfigurationManager` to print information about the local agent
- `Get-DscConfiguration` to print the current configuration for this node
- `Update-DscConfiguration -Wait` to force the local agent to check in with the server immediately
- `Get-DscConfigurationStatus [-All]` to print historical information about DSC agent runs

## Useful links

- [Powershell DSC official protocol specification](https://msdn.microsoft.com/library/dn393548.aspx)
- [Powershell DSC official protocol specification errata](https://msdn.microsoft.com/library/mt612824.aspx)
