package tasks


type Categories struct {
	CatMap       map[int][]Item
	CatList      []Category
	CurrentIndex int
	Size	     int
}


type Category struct {
	Title string
	Id    string
}


type Item struct {
	Id    		string `json:"Id"`
	Title 		string `json:"Title"`
	Due   		string `json:"Due"`
	Link  		string `json:"SelfLink"`
	Status		string `json:"Status"`// "needsAction" or "completed"
	Notes 		string `json:"Notes"`
}
