# Steam workshop uploader for Linux using old SteamRemoteStorage API

For uploading using the newer UGC API: see here: https://github.com/nnnn20430/steamugc-uploader

## Usage

1. Download the binary from [Releases](../../releases) and extract it

1. Run the binary to see usage, e.g.

   ```
   $ cd steamremotestorage-uploader/amd64
   $ ./steamremotestorage-uploader
   ./steamremotestorage-uploader: -f file [options]
     -a AppID
           AppID of the game (default 480)
     -i id
           Existing workshop item id or 0 for new item (default 0)
     -t title
           Item title (default "New item")
     -d description
           Item description (default "")
     -f file
           Path to file for upload (default "")
     -p image
           Path to preview image for upload (default "")
     -n note
           Change note (default "")
     -tags tags
           Comma-separated list of item tags (default "")
   ```

#### Run from path

If you would like to put the binary in your path so you can run it from anywhere:

1. Copy `steamremotestorage-uploader` and `libsteam_api.so` to a directory in your path (e.g. `~/.bin`)

1. Run this command in that directory to update the path where `libsteam_api.so` will be searched for:

   ```
   patchelf --set-rpath "$PWD" steamremotestorage-uploader
   ```

## Development

#### Prerequisites

1. Clone this project

1. Go to https://partner.steamgames.com/downloads/list

1. Download Steamworks SDK 1.42

1. Extract it to the directory the project was cloned to, e.g.

   ```
   cd steamremotestorage-uploader
   unzip /path/to/steamworks_sdk_142.zip
   ```

#### Run without build

1. Make changes to code as desired

1. Set required environment variables

   ```
   export CGO_LDFLAGS="-L/${PWD}/sdk/redistributable_bin/linux64 -lsteam_api"
   export GO111MODULE=off
   ```

1. Run the code

   ```
   LD_LIBRARY_PATH=sdk/redistributable_bin/linux64/ go run main.go
   ```

#### Build

1. Set required environment variables

   ```
   export CGO_LDFLAGS="-L/${PWD}/sdk/redistributable_bin/linux64 -lsteam_api -s -w"
   export GO111MODULE=off
   ```

1. Build the project

   ```
   go build -o steamremotestorage-uploader main.go
   ```

1. Run the binary, e.g.

   ```
   LD_LIBRARY_PATH=sdk/redistributable_bin/linux64/ ./steamremotestorage-uploader
   ```
