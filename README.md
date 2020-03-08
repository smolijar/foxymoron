<div align="center">

<img src="https://media.giphy.com/media/xsE65jaPsUKUo/source.gif" />

# 🦊 Foxymoron

[![science](https://forthebadge.com/images/badges/built-with-science.svg)](https://xkcd.com/136/)
[![forthebadge](https://forthebadge.com/images/badges/uses-badges.svg)](https://forthebadge.com)
[![gopher](https://forthebadge.com/images/badges/made-with-go.svg)](https://www.reddit.com/r/golang/comments/e4x7x/does_the_go_gopher_mascot_have_a_name_if_not_i/)
[![forthebadge](https://forthebadge.com/images/badges/pretty-risque.svg)](#💀-pretty-risque-business)
[![water](https://forthebadge.com/images/badges/powered-by-water.svg)](https://www.quora.com/What-is-water-3)
[![forthebadge](https://forthebadge.com/images/badges/does-not-contain-treenuts.svg)](https://giphy.com/gifs/satisfying-squirrel-nut-x9RBTJVPqZy1y)

REST API server proxy for GitLab to allow extended commit queries accross multiple repositories. Build your [Commits Logs From Last Night](http://www.commitlogsfromlastnight.com/) over your GitLab instance!

</div>

## 🚀 Demo

There is a demo server running at [foxymoron.appspot.com](https://foxymoron.appspot.com), and API docs available [there](https://foxymoron.appspot.com/docs/index.html) as well with working demo client.

## ✨ Features
  - 🔑 Authenticate via GitLab API token and provide GitLab instance URL
  - 🌐 Fetch commits across all available projects
  - 📈 Fetch commit stats for given timespan, group by project / namespace

## 💀 Pretty risque business
This section should be probably the first thing in the readme. Whoops.

Querying all commits within a certain time frame is hard work with the current GitLab API. You can only list commits for a given project (where you can select a time frame). Foxymoron fetches all available projects (you cannot filter by date here), filters them for upper/lower bounds using _last activity_ and _created at_ and queries only those. It's not bad, but can result in tens, hundereds of requests with no delays or liming.Even listing projects could generate a considerable load on your GitLab. Paging size is at maximum of 100, so you should be safe with hundreds to thousand of visible projects, which seems reasoanble for an average GitLab instance.

PS.: OMG, please don't run it on [gitlab.com](https://gitlab.com), we might get in trouble.

PPS.: I guess it's too late. You don't know me and you have never been here. Also Go is a boardgame with black and white beans as far as you know.

PPPS.: This is a doom machine. It's Denial of Service as a Service. GitLab DoSAAS. Don't tell anyone about this project, please.


## ⚖️ License

This project is licensed under [MIT](./LICENSE).
