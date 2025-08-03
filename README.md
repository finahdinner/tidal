# <img width="32" height="32" alt="icon" src="https://github.com/user-attachments/assets/3330f2cc-729b-40a7-a585-0d32f93a7053" /> Tidal

### **Tidal** is a Twitch stream title manager that allows you to dynamically update your stream titles in a context-aware, creative, and often humorous way.

Tidal has access to real-time data from your Twitch stream and channel, known as **Stream Variables**. These include values such as your current viewer count, stream category, and follower count.

You can use these **Stream Variables** in two ways:

-   Inject them directly into your live stream title.
-   Feed them into custom prompts sent to a Large Language Model, which returns dynamic responses which you assign to **AI-Generated Variables**.

Both **Stream Variables** and **AI-Generated Variables** can be used in customisable **Title Templates**.

Templates automatically substitute variable names with their current values, and are configured to update on a user-defined regular schedule.

## Installation

### Downloadable binaries

[See Latest Releases](https://github.com/finahdinner/tidal/releases/latest)

### Compile from source

1. Install Go: https://go.dev/doc/install — ensure that `go` is on your system `PATH`.
2. Install Git: https://git-scm.com/downloads — ensure that `git` is on your system `PATH`.

#### Windows

3. Install GCC: [https://code.visualstudio.com/docs/cpp/config-mingw](https://code.visualstudio.com/docs/cpp/config-mingw)
4. Compile the binary, `tidal.exe`:
- **Powershell**:

    ```
    git clone https://github.com/finahdinner/tidal.git tmp-tidal; Set-Location tmp-tidal; go build -o ../tidal.exe; Set-Location ..; Remove-Item tmp-tidal -Recurse -Force
    ```
- **CMD**:

    ```
    git clone https://github.com/finahdinner/tidal.git tmp-tidal && cd tmp-tidal && go build -o ..\tidal.exe && cd .. && rmdir /s /q tmp-tidal
    ```

#### Linux

3. Install GCC (consult your specific distro's instructions)
4. Compile the binary, `tidal`:

    ```
    git clone https://github.com/finahdinner/tidal.git tmp-tidal && cd tmp-tidal && go build -o ../tidal && cd .. && rm -rf tmp-tidal
    ```

## Post-Installation Setup

#### _A basic overview of the usage of the application is outlined in the `Help` section found within the main window, and more detailed instructions are found by clicking the `?` icon in the top right corner of most subsections within the application._

#### _You may also find the following helpful for constructing LLM prompts that produce reliable responses: https://www.promptingguide.ai/introduction/tips_

#### In order for Tidal to dynamically update your Twitch stream title, you must first configure the following:

1. **Twitch Channel & Application Credentials**

-   Navigate to the `Stream Variables` section and click on the settings cog in the top left corner to input these credentials - this subsection includes detailed instructions on how to fill in each field and authenticate your Twitch account.

2. **LLM Credentials**

-   Navigate to the `AI-generated Variables` section and click on the settings cog in the top left corner to input these credentials - this subsection includes detailed instructions on how to fill in each field.

## Example Tidal Usage

1. Define an **AI-Generated Variable** called `GameJoke`, which instructs an LLM with the following:

```
I am currently livestreaming the following on Twitch: $$StreamCategory
Write me a very short, family-friendly joke about what I am streaming.
Do not exceed more than 60 words, and ensure that you respond only with the joke - no additional text.
You may use emojis too if applicable.
```

2. Create a title template with the following text:

```
Streaming $$StreamCategory to $$NumViewers - $$GameJoke
```

3. Set Tidal to update the title (using the above title template) every 3 minutes.

4. Example of a generated title (with some terrible humour):

```
Streaming Old School RuneScape to 14 viewers - Why did the pickpocket get kicked out of Ardougne? He kept stealing the punchlines! 🥷😂
```
