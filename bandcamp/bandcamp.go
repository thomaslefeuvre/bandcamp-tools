package bandcamp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ItemURLHints struct {
	CustomDomain         *string `json:"custom_domain"`
	CustomDomainVerified any     `json:"custom_domain_verified"`
	ItemType             string  `json:"item_type"`
	Slug                 string  `json:"slug"`
	Subdomain            string  `json:"subdomain"`
}

type ItemArt struct {
	ArtID    int64  `json:"art_id"`
	ThumbURL string `json:"thumb_url"`
	URL      string `json:"url"`
}

type Item struct {
	// Added                    time.Time    `json:"added"`
	AlbumID                  int64    `json:"album_id"`
	AlbumTitle               string   `json:"album_title"`
	AlsoCollectedCount       int      `json:"also_collected_count"`
	BandID                   int64    `json:"band_id"`
	BandImageID              *int64   `json:"band_image_id"`
	BandLocation             *string  `json:"band_location"`
	BandName                 *string  `json:"band_name"`
	BandURL                  *string  `json:"band_url"`
	Currency                 *string  `json:"currency"`
	Discount                 *float64 `json:"discount"`
	DownloadAvailable        *bool    `json:"download_available"`
	FanID                    int64    `json:"fan_id"`
	FeaturedTrack            *int64   `json:"featured_track"`
	FeaturedTrackDuration    *float64 `json:"featured_track_duration"`
	FeaturedTrackEncodingsID *int64   `json:"featured_track_encodings_id"`
	FeaturedTrackIsCustom    *bool    `json:"featured_track_is_custom"`
	FeaturedTrackLicenseID   *int64   `json:"featured_track_license_id"`
	FeaturedTrackNumber      *int     `json:"featured_track_number"`
	FeaturedTrackTitle       *string  `json:"featured_track_title"`
	FeaturedTrackURL         *string  `json:"featured_track_url"`
	GenreID                  int      `json:"genre_id"`
	GiftID                   *int64   `json:"gift_id"`
	GiftRecipientName        *string  `json:"gift_recipient_name"`
	GiftSenderName           *string  `json:"gift_sender_name"`
	GiftSenderNote           *string  `json:"gift_sender_note"`
	HasDigitalDownload       *bool    `json:"has_digital_download"`
	Hidden                   *bool    `json:"hidden"`
	Index                    *int     `json:"index"`
	IsGiftable               bool     `json:"is_giftable"`
	IsPreorder               bool     `json:"is_preorder"`
	IsPrivate                bool     `json:"is_private"`
	IsPurchasable            bool     `json:"is_purchasable"`
	IsSetPrice               bool     `json:"is_set_price"`
	IsSubscriberOnly         bool     `json:"is_subscriber_only"`
	IsSubscriptionItem       bool     `json:"is_subscription_item"`
	ItemArt                  ItemArt  `json:"item_art"`
	ItemArtID                int64    `json:"item_art_id"`
	ItemArtIDs               *string  `json:"item_art_ids"`
	ItemArtURL               string   `json:"item_art_url"`
	ItemID                   int64    `json:"item_id"`
	ItemTitle                string   `json:"item_title"`
	ItemType                 string   `json:"item_type"`
	ItemURL                  string   `json:"item_url"`
	Label                    *string  `json:"label"`
	LabelID                  *int64   `json:"label_id"`
	LicensedItem             *bool    `json:"licensed_item"`
	ListenInAppURL           *string  `json:"listen_in_app_url"`
	MerchIDs                 any      `json:"merch_ids"`
	MerchSnapshot            *string  `json:"merch_snapshot"`
	MerchSoldOut             *bool    `json:"merch_sold_out"`
	MessageCount             *int     `json:"message_count"`
	NumStreamableTracks      int      `json:"num_streamable_tracks"`
	PackageDetails           *string  `json:"package_details"`
	Price                    float64  `json:"price"`
	Purchased                any      `json:"purchased"`
	ReleaseCount             *int     `json:"release_count"`
	Releases                 *string  `json:"releases"`
	RequireEmail             *bool    `json:"require_email"`
	SaleItemID               *int64   `json:"sale_item_id"`
	SaleItemType             *string  `json:"sale_item_type"`
	ServiceName              *string  `json:"service_name"`
	ServiceURLFragment       *string  `json:"service_url_fragment"`
	Token                    string   `json:"token"`
	TralbumID                int64    `json:"tralbum_id"`
	TralbumType              string   `json:"tralbum_type"`
	// Updated                  time.Time    `json:"updated"`
	URLHints  ItemURLHints `json:"url_hints"`
	VariantID *int64       `json:"variant_id"`
	Why       *string      `json:"why"`
}

type Data struct {
	Items         []Item                 `json:"items"`
	MoreAvailable bool                   `json:"more_available"`
	ItemLookup    map[string]interface{} `json:"item_lookup"`
	LastToken     string                 `json:"last_token"`
	Tracklists    map[string][]Track     `json:"tracklists"`
	PurchaseInfos map[string]interface{} `json:"purchase_infos"`
	Collectors    map[string]interface{} `json:"collectors"`
}

// DataBlob is a rendered JSON blob that's returned on the initial page fetch. Subsequent fetches
// for data are done using POST requests to their API; urls look like
// https://bandcamp.com/api/fancollection/1/wishlist_items with the fan_id, and older_than_token.
// older_than_token refers to the last "last_token" that was given to the client. If you're going to
// be paginating the wishlist, you should use the "last_token" in the Wishlist struct.
type DataBlob struct {
	// TrackList is a list of tracks in your collection. The page blob contains roughly 40 or so
	// tracks.
	TrackList []BlobTrack `json:"track_list"`

	// ItemCache stores maps of id->item. Sequences use the ids and they should be fetched out of
	// these ItemCaches.
	ItemCache struct {
		Collection map[string]Item `json:"collection"`
		Wishlist   map[string]Item `json:"wishlist"`
	} `json:"item_cache"`

	CollectionData ItemData `json:"collection_data"`
	WishlistData   ItemData `json:"wishlist_data"`

	FanData struct {
		ID int `json:"fan_id"`
	} `json:"fan_data"`
}

type Track struct {
	Artist      string                 `json:"artist"`
	Duration    float64                `json:"duration"`
	File        map[string]interface{} `json:"file"`
	ID          *int64                 `json:"id"`
	Title       *string                `json:"title"`
	TrackNumber *int                   `json:"track_number"`
}

// func LoadWishlist(user string) ([]Item, error) {
// 	return nil, errors.New("unimplemented function")
// }

// ItemData describes the items and their order for a particular kind of item. Item IDs are listed
// in `Sequence` and `PendingSequence`. Data for each item can be found by matching the Item ID to
// an ID in the appropriate `ItemCache`.
type ItemData struct {
	// LastToken is used as a pagination cursor to fetch the next batch of items
	LastToken string `json:"last_token"`

	// Sequence is the order of the items on the current page to be rendered in. See
	// `PendingSequence` for an explanation of what these lists mean.
	Sequence []string `json:"sequence"`

	// PendingSequence is a sequence that needs to be shown first. This is the first batch of items
	// for a given item category ("followers", "wishlist", etc) that isn't baked into the current
	// page.
	//
	// Visiting the wishlist page will result in an empted "pending_sequence" list, and items in the
	// "sequence" list, whereas "colleciton_data" will have a full "pending_sequence" list and an
	// empty "sequence" list.
	PendingSequence []string `json:"pending_sequence"`
}

type BlobTrack struct {
	BandName string `json:"band_namp"`
	Title    string `json:"title"`
	AlbumID  int    `json:"album_id"`
}

type APIItemsResponse struct {
	Items         []Item `json:"items"`
	MoreAvailable bool   `json:"more_available"`
	LastToken     string `json:"last_token"`
}

// Item represents an API return item. Ambiguous while exploring API. "purchased" will always be
// false for items in the wishlist
type Item2 struct {
	// Date added to wishlist
	Added string `json:"added"`

	// URL of the item in the wishlist
	ItemURL string `json:"item_url"`

	// Type of the item. Could be "album", or "track" (maybe "merch")
	ItemType string `json:"item_type"`
}

func aux(fanID, lastpageToken string) (APIItemsResponse, error) {
	// fmt.Println(fanID, lastpageToken)
	request := map[string]string{
		"fan_id":           fanID,
		"older_than_token": lastpageToken,
	}
	requestJSON, _ := json.Marshal(request)

	// fmt.Println("here")

	resp, err := http.Post("https://bandcamp.com/api/fancollection/1/wishlist_items", "application/json", strings.NewReader(string(requestJSON)))
	if err != nil {
		return APIItemsResponse{}, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return APIItemsResponse{}, err
	}

	// fmt.Println(string(body))

	var response APIItemsResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return APIItemsResponse{}, err
	}

	return response, nil
}

func FetchWishlist(user string) ([]Item, error) {
	resp, err := http.Get(fmt.Sprintf("https://bandcamp.com/%s/wishlist", user))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract the baked-in datablob in the HTML
	datablobExp := regexp.MustCompile("id=\"pagedata\".*?data-blob=\"(.*?)\">")
	datablobMatch := datablobExp.FindStringSubmatch(string(body))
	pagedata := strings.ReplaceAll(datablobMatch[1], "&quot;", "\"")

	// Unmarshal the datablob
	var datablob DataBlob
	if err := json.Unmarshal([]byte(pagedata), &datablob); err != nil {
		return nil, err
	}

	var items []Item

	for _, trackID := range datablob.WishlistData.Sequence {
		for itemcacheID, item := range datablob.ItemCache.Wishlist {
			if trackID == itemcacheID {
				items = append(items, item)
				break
			}
		}
	}

	fanID := strconv.Itoa(datablob.FanData.ID)
	lastToken := datablob.WishlistData.LastToken

	var auxErr error
	var nextResp APIItemsResponse
	moreAvailable := true
	for moreAvailable {
		// fmt.Println(fanID, lastToken)
		nextResp, auxErr = aux(fanID, lastToken)
		if auxErr != nil {
			break
		}

		items = append(items, nextResp.Items...)
		moreAvailable = nextResp.MoreAvailable
		lastToken = nextResp.LastToken
	}

	if auxErr != nil {
		return nil, fmt.Errorf("error in auxiliary fn: %w", auxErr)
	}

	return items, nil
}

func SyncWishlist(user string, datafile string) error {
	items, err := FetchWishlist(user)
	if err != nil {
		return err
	}

	f, err := os.Create(datafile)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
