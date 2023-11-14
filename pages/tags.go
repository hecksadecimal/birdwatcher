package pages

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"cherp.chat/birdwatcher/cherp_api"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type search struct {
	Title    string `json:"title"`
	Query    string `json:"query"`
	LatestId int    `json:"latestID"`
}

var tagsPage *container.TabItem
var maxSearchItems = 5
var searchItems []search
var searchItemsContainer *fyne.Container
var searchItemsLabelContent binding.String
var addSearchForm widget.Form

func deleteAtIndex(slice []search, index int) []search {
	return append(slice[:index], slice[index+1:]...)
}

func InitTagSearches() {
	if _, err := os.Stat("./searchtags.json"); err == nil {
		cherp_api.Load("./searchtags.json", &searchItems)
	}
	if len(searchItems) > maxSearchItems {
		searchItems = searchItems[:maxSearchItems]
	}
	SaveTagSearches()
}

func SaveTagSearches() {
	err := cherp_api.Save("./searchtags.json", &searchItems)
	if err != nil {
		fmt.Println(err)
	}
}

func RegenerateSearches() {
	searchItemsContainer.RemoveAll()
	updateSearchCount()
	for index, element := range searchItems {
		url, _ := url.Parse(cherp_api.BaseUrl + "directory/search/" + element.Query)
		hyperlinkWidget := widget.NewHyperlink("Search Page", url)
		hyperlinkWidget.Alignment = fyne.TextAlignTrailing

		queryLabelWidget := widget.NewLabel(element.Query)
		queryLabelWidget.Wrapping = fyne.TextWrapBreak

		cardWidget := widget.NewCard(
			element.Title,
			"",
			container.NewVBox(
				queryLabelWidget,
				container.NewHBox(
					widget.NewButton("Delete", func() {
						searchItems = deleteAtIndex(searchItems, index)
						SaveTagSearches()
						updateSearchCount()
						RegenerateSearches()
					}),
					hyperlinkWidget,
				),
			),
		)

		searchItemsContainer.Add(
			cardWidget,
		)
	}
}

func updateSearchCount() {
	searchItemsLabelContent.Set(fmt.Sprintf("%d/%d tag notifiers", len(searchItems), maxSearchItems))
	SaveTagSearches()
}

func trimQueryTerm(query string) string {
	return strings.Replace(query, cherp_api.BaseUrl+"directory/search/", "", 1)
}

func addSearchTerm(title string, query string) bool {
	if title == "" || query == "" {
		return false
	}
	if len(searchItems) >= maxSearchItems {
		return false
	}

	newSearch := search{title, trimQueryTerm(query), 0}
	searchItems = append(searchItems, newSearch)

	RegenerateSearches()
	return true
}

func GenerateTagsPage() {
	InitTagSearches()

	searchItemsLabelContent = binding.NewString()
	searchItemsLabelContent.Set("")
	updateSearchCount()

	searchItemsContainer = container.NewVBox()

	newTagTitle := binding.NewString()
	newTagTitle.Set("")
	newTagTitleTextEntry := widget.NewEntry()

	newTagTerms := binding.NewString()
	newTagTerms.Set("")
	newTagTermsTextEntry := widget.NewEntry()

	addSearchForm = widget.Form{
		Items: []*widget.FormItem{
			{Text: "Tagset Name", Widget: newTagTitleTextEntry},
			{Text: "Search Terms", Widget: newTagTermsTextEntry},
		},
		OnSubmit: func() {
			newTagTitle.Set(newTagTitleTextEntry.Text)
			newTagTerms.Set(newTagTermsTextEntry.Text)

			if addSearchTerm(newTagTitleTextEntry.Text, newTagTermsTextEntry.Text) {
				newTagTitleTextEntry.SetText("")
				newTagTermsTextEntry.SetText("")
			}

		},
		SubmitText: "Add Search Term",
	}

	advancedSearchUrl, _ := url.Parse(cherp_api.BaseUrl + "directory/advsearch")

	searchControlsContainer := container.NewVBox(
		widget.NewLabelWithData(searchItemsLabelContent),
		widget.NewLabel("Use advanced search on "+cherp_api.BaseUrl+" and copy the resulting url to 'search terms'."),
		widget.NewHyperlink("Advanced Search", advancedSearchUrl),
		&addSearchForm,
	)

	tagsPageContent := container.NewVBox(
		searchControlsContainer,
		container.NewHScroll(searchItemsContainer),
	)

	tagsPage = container.NewTabItem("Tags", tagsPageContent)
	RegenerateSearches()
}

func TickTags() {
	if (CurrentTick+12)%chatsTickRate == 0 && cherp_api.LoggedIn {
		for index, item := range searchItems {
			results := cherp_api.GetSearchResults(item.Query)
			if results["prompts"] != nil {
				prompts := results["prompts"].([]interface{})
				lastPrompt := prompts[0].(map[string]interface{})
				latestPromptID, _ := strconv.Atoi(lastPrompt["ID"].(string))
				if latestPromptID != item.LatestId {
					cherp_api.SendNotification("Cherp Tags", "Your TagSet '"+item.Title+"' has new activity!")
					searchItems[index].LatestId = latestPromptID
					SaveTagSearches()
				}
			}
			time.Sleep(time.Second / 2) // Lets not hammer the server harder than we need to.
		}
	}
}
