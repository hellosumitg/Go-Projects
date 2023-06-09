package main

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/gorilla/websocket"
)

const wsendpoint = "wss://fstream.binance.com/stream?streams=btcusdt@markPrice/btcusdt@depth"

// Global variables
var (
	WIDTH         = 0
	HEIGHT        = 0
	currMarkPrice = 0.0
	prevMarkPrice = 0.0
	fundingRate   = "n/a"
	ARROW_UP      = "↑"
	ARROW_DOWN    = "↓"
)

// Structure representing an orderbook entry
type OrderbookEntry struct {
	Price  float64
	Volume float64
}

// Custom sorting interface for orderbook entries by best ask
type byBestAsk []OrderbookEntry

func (a byBestAsk) Len() int           { return len(a) }
func (a byBestAsk) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byBestAsk) Less(i, j int) bool { return a[i].Price < a[j].Price }

// Custom sorting interface for orderbook entries by best bid
type byBestBid []OrderbookEntry

func (a byBestBid) Len() int           { return len(a) }
func (a byBestBid) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byBestBid) Less(i, j int) bool { return a[i].Price > a[j].Price }

// Structure representing the orderbook
type Orderbook struct {
	Asks map[float64]float64
	Bids map[float64]float64
}

// `NewOrderbook()` is a constructor function for creating a new insatnce of `Orderbook` struct
func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks: make(map[float64]float64),
		Bids: make(map[float64]float64),
	}
}

// Function to handle depth response for the received from the Binance API so as to maintain the orderbook
// It takes two slices, asks and bids, as input. It iterates over the asks and bids slices,
// parses the price and volume values from each entry,
// and calls the addAsk() and addBid() methods to update the order book accordingly.
func (ob *Orderbook) handleDepthResponse(asks, bids []interface{}) {
	for _, v := range asks {
		ask := v.([]interface{})
		price, _ := strconv.ParseFloat(ask[0].(string), 64)
		volume, _ := strconv.ParseFloat(ask[1].(string), 64)
		ob.addAsk(price, volume)
	}
	for _, v := range bids {
		bid := v.([]interface{})
		price, _ := strconv.ParseFloat(bid[0].(string), 64)
		volume, _ := strconv.ParseFloat(bid[1].(string), 64)
		ob.addBid(price, volume)
	}
}

// Function to add a bid to the orderbook
func (ob *Orderbook) addBid(price, volume float64) {
	if _, ok := ob.Bids[price]; ok {
		// If the volume is zero, the entry is deleted from the respective map.
		// Otherwise, it adds or updates the volume for the given price.
		if volume == 0.0 {
			delete(ob.Bids, price)
			return
		}
	}
	ob.Bids[price] = volume
}

// Function to add an ask to the orderbook
func (ob *Orderbook) addAsk(price, volume float64) {
	if volume == 0.0 {
		delete(ob.Asks, price)
		return
	}
	ob.Asks[price] = volume
}

// The two methods `getBids()` and `getAsks()`, Function to get the best `bids` and `asks` entries from the orderbook
func (ob *Orderbook) getBids() []OrderbookEntry {
	depth := 10
	entries := make(byBestBid, len(ob.Bids)) // creates a slice entries of type byBestBid (custom sorting interface) with a length equal to the number of bids in the Bids map.
	i := 0
	// It iterates over the Bids map and filters out entries with zero volume.
	for price, volume := range ob.Bids {
		if volume == 0 {
			continue
		}
		// Then, it assigns the price and volume values to the corresponding OrderbookEntry struct and adds it to the entries slice.
		entries[i] = OrderbookEntry{
			Price:  price,
			Volume: volume,
		}
		i++
	}
	//  Finally, it sorts the entries slice and returns a portion of it based on the depth value (in this case, 10).
	sort.Sort(entries)
	if len(entries) >= depth {
		return entries[:depth]
	}
	return entries
}

// Function to get the best asks from the orderbook
func (ob *Orderbook) getAsks() []OrderbookEntry {
	depth := 10
	entries := make(byBestAsk, len(ob.Asks))
	i := 0
	for price, volume := range ob.Asks {
		entries[i] = OrderbookEntry{
			Price:  price,
			Volume: volume,
		}
		i++
	}
	sort.Sort(entries)
	if len(entries) >= depth {
		return entries[:depth]
	}
	return entries
}

// Structure representing the trade result from Binance
type BinanceTradeResult struct {
	Data struct {
		Price string `json:"p"`
	} `json:"data"`
}

// Structure representing the depth result from Binance API
type BinanceDepthResult struct {
	Asks [][]string `json:"a"`
	Bids [][]string `json:"b"`
}

// Structure representing the depth response from Binance
type BinanceDepthResponse struct {
	Stream string             `json:"stream"`
	Data   BinanceDepthResult `json:"data"`
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}

	// This code establishes a WebSocket connection to the Binance API using the github.com/gorilla/websocket
	conn, _, err := websocket.DefaultDialer.Dial(wsendpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	var (
		ob     = NewOrderbook()
		result map[string]interface{} // `result` variable will be used to store the JSON response from the WebSocket connection.
	)

	go func() {
		// goroutine to continuously read JSON messages from the WebSocket connection using the `conn.ReadJSON()` method.
		// It parses the received JSON into the `result` map
		for {
			if err := conn.ReadJSON(&result); err != nil {
				log.Fatal(err)
			}
			// below `stream` variable denotes the `stream` field in the received JSON to determine if it is a depth update or a mark price update.
			stream := result["stream"].(string) 
			// If it is a depth update, it extracts the asks and bids data from the result map and calls the ob.handleDepthResponse() method to update the order book.
			if stream == "btcusdt@depth" {
				data := result["data"].(map[string]interface{})
				asks := data["a"].([]interface{})
				bids := data["b"].([]interface{})
				ob.handleDepthResponse(asks, bids)
			}
			// If it is a marketprice update, it updates the previous mark price, current mark price, and funding rate variables.
			if stream == "btcusdt@markPrice" {
				prevMarkPrice = currMarkPrice
				data := result["data"].(map[string]interface{})
				priceStr := data["p"].(string)
				fundingRate = data["r"].(string)
				currMarkPrice, _ = strconv.ParseFloat(priceStr, 64)
			}
		}
	}()

	isrunning := true

	margin := 2
	pheight := 3

	pticker := widgets.NewParagraph()
	pticker.Title = "Binancef"
	pticker.Text = "[BTCUSDT](fg:cyan)"
	pticker.SetRect(0, 0, 14, pheight)

	pprice := widgets.NewParagraph()
	pprice.Title = "Market price"
	ppriceOffset := 14 + 14 + margin + 2
	pprice.SetRect(14+margin, 0, ppriceOffset, pheight)

	pfund := widgets.NewParagraph()
	pfund.Title = "Funding rate"
	pfund.SetRect(ppriceOffset+margin, 0, ppriceOffset+margin+16, 3)

	tob := widgets.NewTable()
	out := make([][]string, 20)
	for i := 0; i < 20; i++ {
		out[i] = []string{"n/a", "n/a"}
	}
	tob.TextStyle = ui.NewStyle(ui.ColorWhite)
	tob.SetRect(0, pheight+2, 30, 22+pheight+2)
	tob.PaddingBottom = 0
	tob.PaddingTop = 0
	tob.RowSeparator = false
	tob.TextAlignment = ui.AlignCenter
	for isrunning {
		var (
			asks = ob.getAsks()
			bids = ob.getBids()
		)
		if len(asks) >= 10 {
			for i := 0; i < 10; i++ {
				out[i] = []string{fmt.Sprintf("[%.2f](fg:red)", asks[i].Price), fmt.Sprintf("[%.2f](fg:cyan)", asks[i].Volume)}
			}
		}
		if len(bids) >= 10 {
			for i := 0; i < 10; i++ {
				out[i+10] = []string{fmt.Sprintf("[%.2f](fg:green)", bids[i].Price), fmt.Sprintf("[%.2f](fg:cyan)", bids[i].Volume)}
			}
		}
		tob.Rows = out

		pprice.Text = getMarketPrice()
		pfund.Text = fmt.Sprintf("[%s](fg:yellow)", fundingRate)
		ui.Render(pticker, pprice, pfund, tob)
		time.Sleep(time.Millisecond * 20)
	}
}

// Function to get the market price with arrow indicator
func getMarketPrice() string {
	price := fmt.Sprintf("[%s %.2f](fg:green)", ARROW_UP, currMarkPrice)
	if prevMarkPrice > currMarkPrice {
		price = fmt.Sprintf("[%s %.2f](fg:red)", ARROW_DOWN, currMarkPrice)
	}
	return price
}
