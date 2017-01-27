package cmd

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"

	"github.com/Shopify/themekit/kit"
)

const settingsDataKey = "config/settings_data.json"

var uploadCmd = &cobra.Command{
	Use:   "upload <filenames>",
	Short: "Upload theme file(s) to shopify",
	Long: `Upload will upload specific files to shopify servers if provided file names.
If no filenames are provided then upload will upload every file in the project
to shopify.

For more documentation please see http://shopify.github.io/themekit/commands/#upload
`,
	RunE: forEachClient(upload, uploadSettingsData),
}

func upload(client kit.ThemeClient, filenames []string, wg *sync.WaitGroup) {
	defer wg.Done()

	if client.Config.ReadOnly {
		kit.LogErrorf("[%s]environment is reaonly", kit.GreenText(client.Config.Environment))
		return
	}

	var err error
	localAssets := []kit.Asset{}

	if len(filenames) == 0 {
		localAssets, err = client.LocalAssets()
		if err != nil {
			kit.LogError(err)
			return
		}
	} else {
		for _, filename := range filenames {
			asset, err := client.LocalAsset(filename)
			if err != nil {
				kit.LogError(err)
				return
			}
			localAssets = append(localAssets, asset)
		}
	}

	bar := newProgressBar(len(localAssets)-2, client.Config.Environment)
	for _, asset := range localAssets {
		if asset.Key == settingsDataKey {
			continue
		}
		wg.Add(1)
		go performUpload(client, asset, bar, wg)
	}
}

func performUpload(client kit.ThemeClient, asset kit.Asset, bar *mpb.Bar, wg *sync.WaitGroup) {
	defer wg.Done()
	defer incBar(bar)

	resp, err := client.UpdateAsset(asset)
	if err != nil {
		kit.LogErrorf("[%s]%s", kit.GreenText(client.Config.Environment), err)
	} else if verbose {
		kit.Printf(
			"[%s] Successfully performed Update on file %s from %s",
			kit.GreenText(client.Config.Environment),
			kit.GreenText(asset.Key),
			kit.YellowText(resp.Host),
		)
	}
}

func uploadSettingsData(client kit.ThemeClient, filenames []string, wg *sync.WaitGroup) {
	if client.Config.ReadOnly {
		return
	}

	doupload := func() {
		asset, err := client.LocalAsset(settingsDataKey)
		if err != nil {
			kit.LogError(err)
			return
		}
		wg.Add(1)
		go performUpload(client, asset, nil, wg)
	}

	if len(filenames) == 0 {
		doupload()
	} else {
		for _, filename := range filenames {
			if filename == settingsDataKey {
				doupload()
				return
			}
		}
	}
}
