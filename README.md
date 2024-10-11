# Pinata Files CLI

![cover](https://dweb.mypinata.cloud/ipfs/QmNcdx9t48z7RQUXUZZHmuc4zBfyBxKLjDfEgmfhiop7j7?img-format=webp)
The official CLI for the Files API written in Go

## Installation

> [!NOTE]
> If you are on Windows please use WSL when installing. If you get an error that it was not able to resolve the github host run `git config --global --unset http.proxy`

### Install Script

The easiest way to install is to copy and paste this script into your terminal

```bash
curl -fsSL https://cli.pinata.cloud/install | bash
```

### Homebrew

If you are on MacOS and have homebrew installed you can run the command below to install the CLI

```
brew install PinataCloud/files-cli/files-cli
```

### Building from Source

To build and instal from source make sure you have [Go](https://go.dev/) installed on your computer and the following command returns a version:

```
go version
```

Then paste and run the following into your terminal:

```
git clone https://github.com/PinataCloud/files-cli && cd files-cli && go install .
```

### Linux Binary

As versions are released you can visit the [Releases](https://github.com/PinataCloud/files-cli/releases) page and download the appropriate binary for your system, them move it into your bin folder.

For example, this is how I install the CLI for my Raspberry Pi

```
wget https://github.com/PinataCloud/files-cli/releases/download/v0.1.0/files-cli_Linux_arm64.tar.gz

tar -xzf files-cli_Linux_arm64.tar.gz

sudo mv pinata /usr/bin
```

## Usage

The Pinata CLI is equipped with the majortiry of features on the Pinata API.

### `auth`

With the CLI installed you will first need to authenticate it with your [Pinata JWT](https://docs.pinata.cloud/account-management/api-keys). Run this command and follow the steps to setup the CLI!

```
pinata auth
```

### `upload`

```
NAME:
   pinata upload - Upload a file to Pinata

USAGE:
   pinata upload [command options] [path to file]

OPTIONS:
   --group value, -g value  Upload a file to a specific group by passing in the groupId
   --name value, -n value   Add a name for the file you are uploading. By default it will use the filename on your system. (default: "nil")
   --verbose                Show upload progress (default: false)
   --help, -h               show help
```

### `files`

```
NAME:
   pinata files - Interact with your files on Pinata

USAGE:
   pinata files command [command options] [arguments...]

COMMANDS:
   delete, d  Delete a file by ID
   get, g     Get file info by ID
   update, u  Update a file by ID
   list, l    List most recent files
   help, h    Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

#### `get`

```
NAME:
   pinata files get - Get file info by ID

USAGE:
   pinata files get [command options] [ID of file]

OPTIONS:
   --help, -h  show help
```

#### `list`

```
NAME:
   pinata files list - List most recent files

USAGE:
   pinata files list [command options] [arguments...]

OPTIONS:
   --name value, -n value    Filter by name of the target file
   --cid value, -c value     Filter results by CID
   --group value, -g value   Filter results by group ID
   --mime value, -m value    Filter results by file mime type
   --amount value, -a value  The number of files you would like to return
   --token value, -t value   Paginate through file results using the pageToken
   --cidPending              Filter results based on whether or not the CID is pending (default: false)
   --help, -h                show help
```

#### `update`

```
NAME:
   pinata files update - Update a file by ID

USAGE:
   pinata files update [command options] [ID of file]

OPTIONS:
   --name value, -n value  Update the name of a file
   --help, -h              show help
```

#### `delete`

```
NAME:
   pinata files delete - Delete a file by ID

USAGE:
   pinata files delete [command options] [ID of file]

OPTIONS:
   --help, -h  show help
```

### `groups`

```
NAME:
   pinata groups - Interact with file groups

USAGE:
   pinata groups command [command options] [arguments...]

COMMANDS:
   create, c  Create a new group
   list, l    List groups on your account
   update, u  Update a group
   delete, d  Delete a group by ID
   get, g     Get group info by ID
   help, h    Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

#### `create`

```
NAME:
   pinata groups create - Create a new group

USAGE:
   pinata groups create [command options] [name of group]

OPTIONS:
   --public, -p  Determine if the group should be public (default: false)
   --help, -h    show help
```

#### `get`

```
NAME:
   pinata groups get - Get group info by ID

USAGE:
   pinata groups get [command options] [ID of group]

OPTIONS:
   --help, -h  show help
```

#### `list`

```
NAME:
   pinata groups list - List groups on your account

USAGE:
   pinata groups list [command options] [arguments...]

OPTIONS:
   --public, -p              List only public groups (default: false)
   --amount value, -a value  The number of groups you would like to return (default: "10")
   --name value, -n value    Filter groups by name
   --token value, -t value   Paginate through results using the pageToken
   --help, -h                show help
```

### `gateways`

```
USAGE:
   pinata gateways command [command options] [arguments...]

COMMANDS:
   set, s   Set your default gateway to be used by the CLI
   sign, s  Get a signed URL for a file by CID
   help, h  Shows a list of commands or help for one command
```

#### `set`

> [!TIP]
> Pass no arguments and get a list of your gateways to choose from!

```
NAME:
   pinata gateways set - Set your default gateway to be used by the CLI

USAGE:
   pinata gateways set [command options] [domain of the gateway]

OPTIONS:
   --help, -h  show help
```

#### `sign`

```
NAME:
   pinata gateways sign - Get a signed URL for a file by CID

USAGE:
   pinata gateways sign [command options] [cid of the file, number of seconds the url is valid for]
   example: pinata gateways sign bafkreih5aznjvttude6c3wbvqeebb6rlx5wkbzyppv7garjiubll2ceym4 30

OPTIONS:
   --help, -h  show help
```

## Contact

If you have any questions please feel free to reach out to us!

[team@pinata.cloud](mailto:team@pinata.cloud)
