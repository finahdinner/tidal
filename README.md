# Tidal

A streaming client for Twitch that allows you to dynamically and creatively update your current stream title using real-time stream/channel data alongside the power of LLMs.

---

## Installation

#### _(Downloadable binaries coming soon...)_

1. Install Go: https://go.dev/doc/install — ensure that `go` is on your system `PATH`.
2. Install Git: https://git-scm.com/downloads — ensure that `git` is on your system `PATH`.
3. Compile the binary:

#### Windows

```
git clone https://github.com/finahdinner/tidal.git tmp-tidal && cd tmp-tidal && go build -o ../myapp.exe && cd .. && rmdir /s /q tmp-tidal
```

#### macOS & Linux

```
git clone https://github.com/finahdinner/tidal.git tmp-tidal && cd tmp-tidal && go build -o ../myapp && cd .. && rm -rf tmp-tidal
```

## Basic Instructions

#### A basic overview of the usage of the application is outlined in the `Help` section found within the main window, and more detailed instructions are found by clicking the `?` icon in the top right corner of most subsections within the application.

#### In order to use Tidal to dynamically update your Twitch stream titles, you will first need to configure the following:

1. **Twitch Channel & Application Credentials**

-   Navigate to the `Stream Variables` section and click on the settings cog in the top left corner to input these credentials - this subsection includes detailed instructions on how to fill in each field and authenticate your Twitch account.

2. **LLM Credentials**

-   Navigate to the `AI-generated` section and click on the settings cog in the top left corner to input these credentials - this subsection includes detailed instructions on how to fill in each field and configure LLM credentials.
