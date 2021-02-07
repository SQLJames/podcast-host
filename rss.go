package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"io"
)

var rsstemplate = `
	<rss version="2.0" xmlns:itunes="http://www.itunes.com/dtds/podcast-1.0.dtd" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom">
	<channel>
	<title>{{.Show.Required.Title}}</title>
	{{if .Show.Recommended.Website}}<link>{{.Show.Recommended.Website}}</link>{{ end }}
	<language>{{.Show.Required.Language}}</language>
	{{if .Show.Situational.Copyright}}<copyright>{{.Show.Situational.Copyright}}</copyright>{{end}}
	{{if .Show.Recommended.Author}}<itunes:author>{{.Show.Recommended.Author}}</itunes:author> {{end}}
	<itunes:summary>{{.Show.Required.Description}}</itunes:summary>
	<description>{{.Show.Required.Description}}</description>
	{{if .Show.Situational.Type}}<itunes:type>{{.Show.Situational.Type}}</itunes:type>{{end}}
	{{if .Show.Recommended.Owner}} 
	<itunes:owner>
	{{if .Show.Recommended.Owner.Name}} <itunes:name>{{.Show.Recommended.Owner.Name}}</itunes:name>{{ end }}
	{{if .Show.Recommended.Owner.Email}} <itunes:email>{{.Show.Recommended.Owner.Email}}</itunes:email>{{ end }}
	</itunes:owner>
	{{ end }}
	<itunes:explicit>{{.Show.Required.Explicit}}</itunes:explicit>
	<itunes:image href="{{.Show.Required.Image}}" />
	{{range .Show.Required.Categories}}
	<itunes:category text="{{.Category}}">
	{{if .Subcategory}} <itunes:category text="{{.Subcategory}}" />{{end}}
	</itunes:category>
	{{end}}
	{{if .Show.Recommended.HostURL}} <atom:link href="{{.Show.Recommended.HostURL}}" rel="self" type="application/rss+xml" />{{end}}
	{{if .Show.Situational.Block}}<itunes:block>{{.Show.Situational.Block}}</itunes:block>{{end}}
	{{if .Show.Situational.NewURLFeed}}
	<itunes:new-feed-url>
	{{.Show.Situational.NewURLFeed}}
	</itunes:new-feed-url>
	{{end}}
	{{range .Episodes}}
	<item>
	{{if .Situational.EpisodeType}}<itunes:episodeType>{{.Situational.EpisodeType}}</itunes:episodeType>{{end}}
	{{if .Situational.Block}}<itunes:block>{{.Situational.Block}}</itunes:block>{{end}}
	{{if .Situational.EpisodeNumber}}<itunes:episode>{{.Situational.EpisodeNumber}}</itunes:episode>{{end}}
	{{if .Situational.Season}}<itunes:season>{{.Situational.Season}}</itunes:season>{{end}}
	<title>{{.Required.Title}}</title>
	{{if .Recommended.Description}}
	<description>
	{{.Recommended.Description}}
	</description>
	{{end}}
	{{if .Recommended.Link}}<link>{{.Recommended.Link}}</link>{{end}}
	{{if .Recommended.Image}}
	<itunes:image href="{{.Recommended.Image}}" />
	{{end}}
	<enclosure 
	length="{{.Required.Enclosure.Length}}" 
	type="{{.Required.Enclosure.Type}}" 
	url="{{.Required.Enclosure.URL}}"
	/>
	{{if .Recommended.GUID}}<guid isPermaLink="false">{{.Recommended.GUID}}</guid>{{end}}
	{{if .Recommended.PublishDate}}<pubDate>{{.Recommended.PublishDate}}</pubDate>{{end}}
	{{if .Recommended.Duration}}<itunes:duration>{{.Recommended.Duration}}</itunes:duration>{{end}}
	{{if .Recommended.Explicit}}<itunes:explicit>{{.Recommended.Explicit}}</itunes:explicit>{{end}}
	</item>
	{{end}}
	</channel>
	</rss>
`

// Podcast single podcast detail
type Podcast struct {
	Show     ShowDetail `yaml:"Show"`
	Web      WebDetail  `yaml:"Website"`
	Episodes []Episode  `yaml:"Episodes"`
}

// ShowDetail is the details about you podcast
type ShowDetail struct {
	Required    RequiredShowDetails    `yaml:"Required"`
	Recommended RecommendedShowDetails `yaml:"Recommended"`
	Situational SituationalShowDetails `yaml:"Situational"`
}

// RequiredShowDetails information required by itunes RSS feed
type RequiredShowDetails struct {
	Title       string         `yaml:"Title"`
	Description string         `yaml:"Description"`
	Language    string         `yaml:"Language"`
	Image       string         `yaml:"Image"`
	Categories  []ShowCategory `yaml:"Categories"`
	Explicit    string         `yaml:"Explicit"`
}

// ShowCategory Categories defined by Apple
type ShowCategory struct {
	Category    string `yaml:"Category"`
	Subcategory string `yaml:"Subcategory"`
}

// RecommendedShowDetails information recommended by itunes RSS feed
type RecommendedShowDetails struct {
	Website string           `yaml:"Website"`
	Author  string           `yaml:"Author"`
	Owner   ShowOwnerDetails `yaml:"Owner"`
	HostURL string           `yaml:"HostURL"`
}

// ShowOwnerDetails information about the owner of the podcast
type ShowOwnerDetails struct {
	Name  string `yaml:"Name"`
	Email string `yaml:"Email"`
}

// SituationalShowDetails the situational details about your
type SituationalShowDetails struct {
	Type       string `yaml:"Type"`
	Copyright  string `yaml:"Copyright"`
	NewURLFeed string `yaml:"NewURLFeed"`
	Block      string `yaml:"Block"`
	Complete   string `yaml:"Complete"`
}

// WebDetail the details about how the url RSS will be hosted
type WebDetail struct {
	URL string `yaml:"URL"`
}

// Episode Holds the Detail about an Episode
type Episode struct {
	Required    RequiredEpisodeDetails    `yaml:"Required"`
	Recommended RecommendedEpisodeDetails `yaml:"Recommended"`
	Situational SituationalEpisodeDetails `yaml:"Situational"`
}

// RequiredEpisodeDetails Required details about a show
type RequiredEpisodeDetails struct {
	Title     string                   `yaml:"Title"`
	Enclosure RequiredEpisodeEnclosure `yaml:"Enclosure"`
}

// RequiredEpisodeEnclosure Required details about a show length URL and format
type RequiredEpisodeEnclosure struct {
	URL    string `yaml:"URL"`
	Length string `yaml:"Length"`
	Type   string `yaml:"Type"`
}

// RecommendedEpisodeDetails Recommended details about the episode
type RecommendedEpisodeDetails struct {
	GUID        string `yaml:"GUID"`
	PublishDate string `yaml:"PublishDate"`
	Duration    string `yaml:"Duration"`
	Link        string `yaml:"Link"`
	Image       string `yaml:"Image"`
	Explicit    string `yaml:"Explicit"`
	Description string `yaml:"Description"`
}

// SituationalEpisodeDetails Situational details about the episode
type SituationalEpisodeDetails struct {
	EpisodeNumber string `yaml:"EpisodeNumber"`
	Season        string `yaml:"Season"`
	EpisodeType   string `yaml:"EpisodeType"`
	Block         string `yaml:"Block"`
}

//CreateFeed creates the RSS feed to a file
func CreateFeed(cast Podcast, f io.Writer) {
	tmpl, err := template.New(cast.Web.URL).Parse(rsstemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(f, cast)
	if err != nil {
		panic(err)
	}

}

func getShowDetails(filename string) (Podcast, error) {
	var PC Podcast
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading YAML file: %s\n", err)
		return PC, err
	}
	err = yaml.Unmarshal(yamlFile, &PC)
	if err != nil {
		log.Fatalf("Error parsing YAML file: %s\n", err)
		return PC, err
	}

	return PC, nil
}

//CreateFeeds Searches the provided path for all podcasts
//and looks for the yaml file to generate the RSS feed.
func (pc *ProgramConfig) CreateFeeds() {
	//flag.StringVar(&Config, "f", "./config/config.yaml", "YAML file to parse.")
	//flag.Parse()
	log.Println("Looking for Folders under the provided Search Folder")
	log.Println(pc.RSS.SearchFolder)
	files, err := ioutil.ReadDir(pc.RSS.SearchFolder)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {

		filename := filepath.Join(pc.RSS.SearchFolder, f.Name())
		//log.Println("Found: " + filename)
		show, _ := getShowDetails(filepath.Join(filename, pc.RSS.PodcastFilename))
		log.Printf("%+v\n", show)
		rssfile, err := os.Create(filepath.Join(filename, "rss.xml"))
		if err != nil {
			log.Fatal(err)
		}
		CreateFeed(show, rssfile)
	}
}
