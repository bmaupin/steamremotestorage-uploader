package main

import (
	"./steam"
	"./util"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var fAppId uint32 = 480
var fPublishedFileID uint64 = 0
var fItemTitle string = "New item"
var fItemDescription string = ""
var fFile string = ""
var fPreviewFile string = ""
var fChangeNote string = ""
var fTags string = ""

func parseArgs() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s: [options]\n", os.Args[0])
		util.PrintDefaults()
	}

	flag.CommandLine.Init("", flag.ExitOnError)

	flag.Var((*util.Uint32Flag)(&fAppId), "a", "0:`AppID` of the game")
	flag.Uint64Var(&fPublishedFileID, "i", 0, "1:Existing workshop item `id` or 0 for new item")
	flag.StringVar(&fItemTitle, "t", "New item", "2:Item `title`")
	flag.StringVar(&fItemDescription, "d", "", "3:Item `description`")
	flag.StringVar(&fFile, "f", "", "4:Path to `file` for upload")
	flag.StringVar(&fPreviewFile, "p", "", "5:Path to preview `image` for upload")
	flag.StringVar(&fChangeNote, "n", "", "6:Change `note`")
	flag.StringVar(&fTags, "tags", "", "7:Comma-separated list of item `tags`")
	flag.Parse()

	if fPublishedFileID == 0 && len(fFile) == 0 {
		fmt.Println("Error: Either an existing item id or a file to upload is required\n")
		flag.Usage()
		os.Exit(1)
	}

	if fPublishedFileID == 0 && len(fChangeNote) > 0 {
		fmt.Println("Error: Change notes can only be set when updating an existing item\n")
		flag.Usage()
		os.Exit(1)
	}
}

func createItem() {
	var steamError bool = false
	var fb []byte

	fmt.Println("Creating new item...")

	fb, _ = ioutil.ReadFile(util.ArgSlice(filepath.Abs(fFile))[0].(string))
	if !steam.SteamRemoteStorage().FileWrite(filepath.Base(fFile), util.GoStringToCString(string(fb)), len(fb)) {
		fmt.Println("Writing file to cloud failed")
		os.Exit(1)
	}
	defer steam.SteamRemoteStorage().FileDelete(filepath.Base(fFile))
	util.PtrFree(&fb)

	if len(fPreviewFile) > 0 {
		fb, _ = ioutil.ReadFile(util.ArgSlice(filepath.Abs(fPreviewFile))[0].(string))
		if !steam.SteamRemoteStorage().FileWrite(filepath.Base(fPreviewFile), util.GoStringToCString(string(fb)), len(fb)) {
			fmt.Println("Writing preview file to cloud failed")
			os.Exit(1)
		}
		defer steam.SteamRemoteStorage().FileDelete(filepath.Base(fPreviewFile))
		util.PtrFree(&fb)
	}

	var steamTags steam.SteamParamStringArray_t
	if len(fTags) > 0 {
		tags := strings.Split(fTags, ",")
		_steamTags, cleanupSteamTags := util.GoStringArrayToSteamStringArray(tags)
		steamTags = *_steamTags
		defer cleanupSteamTags()
	} else {
		steamTags = steam.NewSteamParamStringArray_t()
	}

	hSteamAPICall := steam.SteamRemoteStorage().PublishWorkshopFile(
		filepath.Base(fFile),
		(func() string {
			if len(fPreviewFile) > 0 {
				return filepath.Base(fPreviewFile)
			} else {
				return ""
			}
		})(),
		uint(fAppId),
		fItemTitle,
		fItemDescription,
		steam.K_ERemoteStoragePublishedFileVisibilityPrivate,
		steamTags,
		steam.K_EWorkshopFileTypeCommunity,
	)

	RemoteStoragePublishFileResult := steam.NewRemoteStoragePublishFileResult_t()

	for true {
		if steam.SteamUtils().IsAPICallCompleted(hSteamAPICall, &steamError) {
			steam.SteamUtils().GetAPICallResult(
				hSteamAPICall,
				RemoteStoragePublishFileResult.Swigcptr(),
				steam.Sizeof_RemoteStoragePublishFileResult_t,
				steam.RemoteStoragePublishFileResult_tK_iCallback,
				&steamError,
			)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if steamError {
		fmt.Println("Steam api call PublishWorkshopFile() failed:", steam.SteamUtils().GetAPICallFailureReason(hSteamAPICall))
		os.Exit(1)
	}

	if RemoteStoragePublishFileResult.GetM_eResult() != 1 {
		fmt.Println("Item upload failed:", RemoteStoragePublishFileResult.GetM_eResult())
		os.Exit(1)
	}

	fmt.Println("Your new workshop item id is:", RemoteStoragePublishFileResult.GetM_nPublishedFileId())
}

func updateItem() {
	var steamError bool = false
	var fb []byte

	fmt.Println("Updating item...")

	if len(fFile) > 0 {
		fb, _ = ioutil.ReadFile(util.ArgSlice(filepath.Abs(fFile))[0].(string))
		if !steam.SteamRemoteStorage().FileWrite(filepath.Base(fFile), util.GoStringToCString(string(fb)), len(fb)) {
			fmt.Println("Writing file to cloud failed")
			os.Exit(1)
		}
		defer steam.SteamRemoteStorage().FileDelete(filepath.Base(fFile))
		util.PtrFree(&fb)
	}

	if len(fPreviewFile) > 0 {
		fb, _ = ioutil.ReadFile(util.ArgSlice(filepath.Abs(fPreviewFile))[0].(string))
		if !steam.SteamRemoteStorage().FileWrite(filepath.Base(fPreviewFile), util.GoStringToCString(string(fb)), len(fb)) {
			fmt.Println("Writing preview file to cloud failed")
			os.Exit(1)
		}
		defer steam.SteamRemoteStorage().FileDelete(filepath.Base(fPreviewFile))
		util.PtrFree(&fb)
	}

	PublishedFileUpdateHandle := steam.SteamRemoteStorage().CreatePublishedFileUpdateRequest(fPublishedFileID)

	if !util.IsFlagDefault("t") {
		steam.SteamRemoteStorage().UpdatePublishedFileTitle(PublishedFileUpdateHandle, fItemTitle)
	}

	if len(fItemDescription) > 0 {
		steam.SteamRemoteStorage().UpdatePublishedFileDescription(PublishedFileUpdateHandle, fItemDescription)
	}

	if len(fFile) > 0 {
		steam.SteamRemoteStorage().UpdatePublishedFileFile(PublishedFileUpdateHandle, filepath.Base(fFile))
	}

	if len(fPreviewFile) > 0 {
		steam.SteamRemoteStorage().UpdatePublishedFilePreviewFile(PublishedFileUpdateHandle, filepath.Base(fPreviewFile))
	}

	if len(fChangeNote) > 0 {
		steam.SteamRemoteStorage().UpdatePublishedFileSetChangeDescription(PublishedFileUpdateHandle, fChangeNote)
	}

	//steam.SteamRemoteStorage().UpdatePublishedFileVisibility(PublishedFileUpdateHandle, steam.K_ERemoteStoragePublishedFileVisibilityPrivate)

	if len(fTags) > 0 {
		tags := strings.Split(fTags, ",")
		steamTags, cleanupSteamTags := util.GoStringArrayToSteamStringArray(tags)
		defer cleanupSteamTags()
		steam.SteamRemoteStorage().UpdatePublishedFileTags(PublishedFileUpdateHandle, *steamTags)
	}

	hSteamAPICall := steam.SteamRemoteStorage().CommitPublishedFileUpdate(PublishedFileUpdateHandle)

	RemoteStorageUpdatePublishedFileResult := steam.NewRemoteStorageUpdatePublishedFileResult_t()

	for true {
		if steam.SteamUtils().IsAPICallCompleted(hSteamAPICall, &steamError) {
			steam.SteamUtils().GetAPICallResult(
				hSteamAPICall,
				RemoteStorageUpdatePublishedFileResult.Swigcptr(),
				steam.Sizeof_RemoteStorageUpdatePublishedFileResult_t,
				steam.RemoteStorageUpdatePublishedFileResult_tK_iCallback,
				&steamError,
			)
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	if steamError {
		fmt.Println("Steam api call CommitPublishedFileUpdate() failed:", steam.SteamUtils().GetAPICallFailureReason(hSteamAPICall))
		os.Exit(1)
	}

	if RemoteStorageUpdatePublishedFileResult.GetM_eResult() != 1 {
		fmt.Println("Item upload failed:", RemoteStorageUpdatePublishedFileResult.GetM_eResult())
		os.Exit(1)
	}

	fmt.Println("Update complete")
}

func main() {
	parseArgs()

	ep, err := os.Executable()
	if err != nil {
		fmt.Println("Could not get executable path")
		os.Exit(1)
	}
	if ioutil.WriteFile(
		filepath.Join(filepath.Dir(ep), "steam_appid.txt"),
		[]byte(strconv.FormatUint(uint64(fAppId), 10)),
		0644,
	) != nil {
		fmt.Println("Failed to write to steam_appid.txt")
	}

	if !steam.SteamAPI_Init() {
		fmt.Println("Failed to initialize steam api")
		os.Exit(1)
	}
	defer steam.SteamAPI_Shutdown()

	if fPublishedFileID > 0 {
		updateItem()
	} else {
		createItem()
	}
}
