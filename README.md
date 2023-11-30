# go-jel-redeploy

Simple tool to redeploy your [jelastic](https://www.virtuozzo.com/application-platform-partners/) node groups. 

You need to provide a `config.json` file using `--config` argument of tool, with the following structure. The tool will redeploy nodes with the specified `tag` sequentially.

    {
        "host": "your_jcloud_host",
        "session": "your_api_token",
        "deployments": [
            {
                "envName": "environment_name",
                "nodeGroup": {
                    "name": "nodegroup_name",
                    "tag": "container_tag"
                }
            }
        ]
    }

### Windows example

You can download executable of the tool in [releases](https://github.com/faridprogrammer/go-jel-redeploy/releases/) section of the repository

    [Path to your executable]\go-jel-redeploy.exe -config path/to/config.json