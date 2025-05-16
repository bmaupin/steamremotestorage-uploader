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

#### Create a new item

When creating an item, you must provide a file to publish and all else is optional, e.g.

```
steamremotestorage-uploader -a 12345 -t "My mod name" -f "mymod.file" -p assets/preview.jpg
```

ⓘ Change notes do not work for new items and tag functionality has not yet been implemented. To add change notes and tags, update the item after publishing

#### Update an item

When updating an item, you must provide the item ID and all else is optional, e.g.

```
steamremotestorage-uploader -a 12345 -i 123456798 -t "My mod name" -tags "Other,Gameplay,Leaders" -f "mymod.file" -p assets/preview.jpg -n "v2: Version summary"
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

## Troubleshooting

#### `Item upload failed: 21`

Updating an existing mod normally goes quickly, but sometimes it takes a while and throws this error. If it happens, close and reopen Steam and try again.

#### `IPC function call IClientRemoteStorage::FileWrite took too long`

This can be ignored

#### `panic: runtime error: index out of range [0] with length 0`

Make sure all of the filenames (mod file, preview image, etc) are correct and try again
