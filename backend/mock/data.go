package mock

import "simplicity/items"

func GenerateMockData() []items.ItemData {
	return []items.ItemData{
		{
			Title:       "First Item",
			Description: "This is the first item",
			Images:      []string{"image1.jpg", "image2.jpg"},
			Tags:        []string{"tag1", "tag2"},
		},
		{
			Title:       "Second Item",
			Description: "This is the second item",
			Images:      []string{"image3.jpg", "image4.jpg"},
			Tags:        []string{"tag3", "tag4"},
		},
		{
			Title:       "Third Item",
			Description: "This is the third item",
			Images:      []string{"image5.jpg", "image6.jpg"},
			Tags:        []string{"tag5", "tag6"},
		},
		{
			Title:       "Fourth Item",
			Description: "This is the fourth item",
			Images:      []string{"image7.jpg", "image8.jpg"},
			Tags:        []string{"tag7", "tag8"},
		},
		{
			Title:       "Fifth Item",
			Description: "This is the fifth item",
			Images:      []string{"image9.jpg", "image10.jpg"},
			Tags:        []string{"tag9", "tag10"},
		},
	}

}
