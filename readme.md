# Podcast-Host
Podcast-Host is a golang webserver designed to generate and host your RSS feed for your podcast. This includes:
  - Favicon
  - Generating the XML Feed from the yaml configuration
  - Ability to host multiple podcasts based off folder structure

# New Features!

  - Dynamically generated XML Files based off yaml
  - Hosting multiple podcasts off the same server

# Configuration Values
  - |  Static.EpisodeLocation | Static directory for hosting your podcast episodes.  |
  - |  Static.Favicon | Icon for your website  |
  - |  Static.Images | Static directory for Images in your website and episodes  |
  - |  RSS.SearchFolder | Directory that the program will look for to find folders. This Directory name will be included in your URL Path.  |
  - |  RSS.PodcastFilename | Name of the Yaml configuration that defines your podcast.  |
# Getting Started
  - Create a folder under the Config defined SearchFolder, this is the podcast folder by default. This will be part of your URL Path, so make sure it complies with URL Structure formats.
  - Modify your podcast.yaml file to contain all the details about your podcast and episodes. See the podcast.yaml for a detailed explaination, or example.podcast.yaml for an example generation
  - Place your static files in the appropriate configuration file locations